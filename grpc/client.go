package grpc

import (
	"context"
	"errors"
	"io"

	"github.com/rs/zerolog/log"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

const defaultSize = 1024

// Client represents am upload gRPC client.
type Client struct {
	c UploadClient
}

// NewClient instantiates a new Client.
func NewClient(address string) (*Client, error) {
	conn, err := grpc.Dial(address, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, err
	}

	return &Client{
		c: NewUploadClient(conn),
	}, nil
}

// Upload sends a reader to the server.
func (c *Client) Upload(ctx context.Context, r io.Reader) error {
	stream, err := c.c.Upload(ctx)
	if err != nil {
		return err
	}

	buf := make([]byte, defaultSize)
	for {
		n, err := r.Read(buf)
		if errors.Is(err, io.EOF) {
			break
		}
		if err != nil {
			return err
		}

		if err := stream.Send(&Request{
			Chunk: buf[:n],
		}); err != nil {
			return err
		}
	}

	res, err := stream.CloseAndRecv()
	if err != nil {
		return err
	}

	log.Info().Msgf("Server response: %s", res.Result)

	return nil
}
