package main

import (
	"context"
	"errors"
	"fmt"
	"io"
	"os"
	"os/signal"
	"time"

	"github.com/agukrapo/ports/database"
	"github.com/agukrapo/ports/grpc"
	"github.com/agukrapo/ports/parser"
	"github.com/agukrapo/ports/rest"
	"github.com/agukrapo/ports/service"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

type usageError string

func (e usageError) Error() string {
	return string(e) + `

usage:
  ports COMMAND ARG

available COMMANDS
  cli: command line interface, requires an extra file path argument
  rest: REST server
`
}

func main() {
	if err := exec(); err != nil {
		_, _ = fmt.Fprintf(os.Stderr, "%v\n", err)
		os.Exit(1)
	}
}

func exec() error {
	log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr, TimeFormat: time.Kitchen})

	if len(os.Args) < 2 {
		return usageError("missing COMMAND argument")
	}

	switch os.Args[1] {
	case "cli":
		return runCLI()
	case "rest":
		return runREST()
	case "grpc-server":
		return runGRPCServer()
	case "grpc-client":
		return runGRPCClient()
	default:
		return usageError("invalid COMMAND argument: " + os.Args[1])
	}
}

func runCLI() error {
	if len(os.Args) < 3 {
		return errors.New("file path argument missing")
	}

	file, err := os.Open(os.Args[2])
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

	svc := service.New(db)

	ctx, cancel := context.WithCancel(context.Background())

	go func() {
		quit := make(chan os.Signal, 1)
		signal.Notify(quit, os.Interrupt)
		<-quit

		log.Info().Msg("Gracefully stopping...")
		cancel()
	}()

	log.Info().Msg("Ports cli started")
	svc.Process(ctx, src)
	log.Info().Msg("Ports cli finished")

	return nil
}

func runREST() error {
	db, err := openDB()
	if err != nil {
		return err
	}
	defer safeClose(db)

	port, ok := os.LookupEnv("PORT")
	if !ok {
		port = "8080"
	}

	server := rest.New(port, service.New(db))

	server.Start()
	server.Listen()

	return nil
}

func runGRPCServer() error {
	db, err := openDB()
	if err != nil {
		return err
	}
	defer safeClose(db)

	port, ok := os.LookupEnv("PORT")
	if !ok {
		port = "8080"
	}

	server := grpc.NewServer(port, service.New(db))

	server.Start()
	server.Listen()
	return nil
}

func runGRPCClient() error {
	if len(os.Args) < 3 {
		return errors.New("server address argument missing")
	}
	if len(os.Args) < 4 {
		return errors.New("file path argument missing")
	}

	client, err := grpc.NewClient(os.Args[2])
	if err != nil {
		return err
	}

	file, err := os.Open(os.Args[3])
	if err != nil {
		return err
	}
	defer safeClose(file)

	ctx, cancel := context.WithCancel(context.Background())
	go func() {
		quit := make(chan os.Signal, 1)
		signal.Notify(quit, os.Interrupt)
		<-quit

		log.Info().Msg("Gracefully stopping...")
		cancel()
	}()

	if err := client.Upload(ctx, file); err != nil {
		return err
	}

	return nil
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
