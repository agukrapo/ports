// Package parser includes json parsing utilities.
package parser

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
)

// Port represents a port json object.
type Port struct {
	Key         string
	Timezone    string    `json:"timezone"`
	Coordinates []float64 `json:"coordinates"`
	Name        string    `json:"name"`
	City        string    `json:"city"`
	Province    string    `json:"province"`
	Country     string    `json:"country"`
	Alias       []string  `json:"alias"`
	Unlocs      []string  `json:"unlocs"`
	Code        string    `json:"code"`
}

// Iterator represents a port json stream iterator.
type Iterator struct {
	dec *json.Decoder
}

// New instantiates an Iterator.
func New(r io.Reader) (*Iterator, error) {
	d := json.NewDecoder(r)

	if _, err := d.Token(); errors.Is(err, io.EOF) {
		return nil, errors.New("empty input")
	} else if err != nil {
		return nil, err
	}

	return &Iterator{
		dec: d,
	}, nil
}

// More tells if there is a Port in the Iterator.
func (i *Iterator) More() bool {
	return i.dec.More()
}

// Next populates the input port with the next Port in the Iterator.
func (i *Iterator) Next(port *Port) error {
	token, err := i.dec.Token()
	if err != nil {
		return err
	}

	k, ok := token.(string)
	if !ok {
		return fmt.Errorf("invalid key: %v", token)
	}

	if err := i.dec.Decode(&port); err != nil {
		return err
	}

	port.Key = k
	return nil
}

// Packet represents either a parsed Port or an error.
type Packet struct {
	Port *Port
	Err  error
}

// Stream return a Packet unbuffered channel.
func (i *Iterator) Stream(ctx context.Context) chan Packet {
	out := make(chan Packet)

	go func() {
		defer close(out)

		for i.More() {
			var p Port
			if err := i.Next(&p); err != nil {
				out <- Packet{Err: err}
				return
			}

			select {
			case out <- Packet{Port: &p}:
			case <-ctx.Done():
				return
			}
		}
	}()

	return out
}
