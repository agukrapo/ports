// Package parser includes json parsing utilities.
package parser

import (
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

// Parser represents a port json stream parser.
type Parser struct {
	dec *json.Decoder
}

// New instantiates a Parser.
func New(r io.Reader) (*Parser, error) {
	d := json.NewDecoder(r)

	if _, err := d.Token(); errors.Is(err, io.EOF) {
		return nil, errors.New("empty input")
	} else if err != nil {
		return nil, err
	}

	return &Parser{
		dec: d,
	}, nil
}

// More tells if there is a Port in the Parser.
func (p *Parser) More() bool {
	return p.dec.More()
}

// Next populates the input port with the next Port in the Parser.
func (p *Parser) Next(port *Port) error {
	token, err := p.dec.Token()
	if err != nil {
		return err
	}

	k, ok := token.(string)
	if !ok {
		return fmt.Errorf("invalid key: %v", token)
	}

	if err := p.dec.Decode(&port); err != nil {
		return err
	}

	port.Key = k
	return nil
}
