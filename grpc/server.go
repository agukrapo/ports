// Package grpc includes grpc api server definitions.
package grpc

import (
	"errors"
	"fmt"
	"io"
	"net"
	"os"
	"os/signal"
	"syscall"

	"github.com/agukrapo/ports/parser"
	"github.com/agukrapo/ports/service"
	"github.com/rs/zerolog/log"
	"google.golang.org/grpc"
)

// Server represents am upload gRPC server.
type Server struct {
	UnimplementedUploadServer

	s       *grpc.Server
	address string
	service *service.Service
}

// NewServer instantiates a new Server.
func NewServer(port string, service *service.Service) *Server {
	s := grpc.NewServer()
	out := &Server{
		s:       s,
		address: fmt.Sprintf(":%s", port),
		service: service,
	}

	RegisterUploadServer(s, out)

	return out
}

// Start starts a Server.
func (s *Server) Start() {
	go func() {
		if err := s.start(); err != nil {
			log.Error().Err(err).Msg("GRPC server start failed")

			_ = syscall.Kill(syscall.Getpid(), syscall.SIGINT)
		}
	}()
	log.Info().Msgf("GRPC server started at %s", s.address)
}

func (s *Server) start() error {
	lis, err := net.Listen("tcp", s.address)
	if err != nil {
		return err
	}

	if err := s.s.Serve(lis); err != nil {
		return err
	}

	return nil
}

// Listen blocks until an os.Interrupt occurs.
func (s *Server) Listen() {
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt)
	<-quit

	s.s.GracefulStop()

	log.Info().Msg("GRPC server closed")
}

func (s *Server) Upload(stream Upload_UploadServer) error {
	tmp, err := os.CreateTemp(os.TempDir(), "upload")
	if err != nil {
		return err
	}
	defer safeClose(tmp)

	for {
		req, err := stream.Recv()
		if errors.Is(err, io.EOF) {
			break
		}
		if err != nil {
			return err
		}

		if _, err := tmp.Write(req.Chunk); err != nil {
			return err
		}
	}

	if _, err := tmp.Seek(0, io.SeekStart); err != nil {
		return err
	}

	p, err := parser.New(tmp)
	if err != nil {
		return err
	}

	s.service.Process(stream.Context(), p)

	return stream.SendAndClose(&Response{Result: "ok"})
}

func safeClose(c io.Closer) {
	if err := c.Close(); err != nil {
		log.Error().Err(err).Msg("Close failed")
	}
}
