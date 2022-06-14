package main

import (
	"errors"
	"fmt"
	"io"
	"os"
	"os/signal"
	"time"

	"github.com/agukrapo/ports/database"
	"github.com/agukrapo/ports/parser"
	"github.com/agukrapo/ports/service"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

func main() {
	if err := exec(); err != nil {
		_, _ = fmt.Fprintf(os.Stderr, "%v\n", err)
		os.Exit(1)
	}
}

func exec() error {
	log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr, TimeFormat: time.Kitchen})

	file, err := openFile()
	if err != nil {
		return err
	}
	defer safeClose(file)

	src, err := parser.New(file)
	if err != nil {
		return err
	}

	db, err := openDB()
	if err != nil {
		return err
	}
	defer safeClose(db)

	svc := service.New(src, db)

	go func() {
		quit := make(chan os.Signal, 1)
		signal.Notify(quit, os.Interrupt)
		<-quit

		log.Info().Msg("Gracefully stopping...")
		svc.Stop()
	}()

	log.Info().Msg("Ports app started")
	svc.Start()
	log.Info().Msg("Ports app finished")

	return nil
}

func openFile() (io.ReadCloser, error) {
	if len(os.Args) == 1 {
		return nil, errors.New("file path argument missing")
	}

	return os.Open(os.Args[1])
}

func openDB() (*database.Database, error) {
	dsn, ok := os.LookupEnv("DATABASE_DSN")
	if !ok {
		return nil, errors.New("DATABASE_DSN environment missing")
	}

	return database.New(dsn)
}

func safeClose(c io.Closer) {
	if err := c.Close(); err != nil {
		log.Error().Err(err).Msg("Close failed")
	}
}
