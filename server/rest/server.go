// Package rest includes rest api related utilities.
package rest

import (
	"context"
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/agukrapo/ports/parser"
	"github.com/agukrapo/ports/service"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	glog "github.com/labstack/gommon/log"
	"github.com/rs/zerolog/log"
)

const (
	shutdownTimeout = 3 * time.Second
)

type Server struct {
	address string
	e       *echo.Echo
	service *service.Service
}

// New instantiates a new Server.
func New(port string, service *service.Service) *Server {
	e := echo.New()
	e.HideBanner = true
	e.HidePort = true
	e.Logger.SetLevel(glog.OFF)
	e.Use(middleware.Recover(), loggingMW)

	s := &Server{
		e:       e,
		address: fmt.Sprintf(":%s", port),
		service: service,
	}

	e.PUT("/upload", s.upload)

	return s
}

// Start starts a Server.
func (s *Server) Start() {
	go func() {
		if err := s.e.Start(s.address); err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.Error().Err(err).Msg("REST server start failed")

			_ = syscall.Kill(syscall.Getpid(), syscall.SIGINT)
		}
	}()

	log.Info().Msgf("REST server started at %s", s.address)
}

// Listen blocks until an os.Interrupt occurs.
func (s *Server) Listen() {
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt)
	<-quit

	ctx, cancel := context.WithTimeout(context.Background(), shutdownTimeout)
	defer cancel()

	if err := s.e.Shutdown(ctx); err != nil {
		log.Error().Err(err).Msg("REST server shutdown failed")
	}

	log.Info().Msg("REST server closed")
}

func loggingMW(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) (err error) {
		if err = next(c); err != nil {
			c.Error(err)
		}

		req := c.Request()
		res := c.Response()

		event := log.Debug()
		event.Str("request", fmt.Sprintf("%s %s", req.Method, req.URL))
		event.Str("status", fmt.Sprintf("%d %s", res.Status, http.StatusText(res.Status)))

		if err != nil {
			event.Err(err)
		}

		event.Msg("Endpoint call")

		return err
	}
}

func (s *Server) upload(c echo.Context) error {
	file, err := c.FormFile("file")
	if err != nil {
		return c.JSON(http.StatusBadRequest, err.Error())
	}

	src, err := file.Open()
	if err != nil {
		return err
	}
	defer safeClose(src)

	p, err := parser.New(src)
	if err != nil {
		return err
	}

	s.service.Process(c.Request().Context(), p)

	return c.NoContent(http.StatusOK)
}

func safeClose(c io.Closer) {
	if err := c.Close(); err != nil {
		log.Error().Err(err).Msg("Close failed")
	}
}
