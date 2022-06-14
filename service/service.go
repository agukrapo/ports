// Package service includes the port service.
package service

import (
	"fmt"

	"github.com/agukrapo/ports/database"
	"github.com/agukrapo/ports/parser"
	"github.com/rs/zerolog/log"
)

type source interface {
	More() bool
	Next(*parser.Port) error
}

type destination interface {
	Upsert(*database.Port) error
}

// Service represents a process that moves ports from a source to a destination.
type Service struct {
	src  source
	dest destination

	closed bool
}

// New instantiates a new Service.
func New(src source, dest destination) *Service {
	return &Service{
		src:  src,
		dest: dest,
	}
}

// Start starts the Service.
func (s *Service) Start() {
	var in parser.Port
	var out database.Port

	for !s.closed && s.src.More() {
		if err := s.src.Next(&in); err != nil {
			log.Error().Err(err).Msg("Port parse failed")
			continue
		}

		if err := translate(&in, &out); err != nil {
			log.Error().Err(err).Msg("Port translation failed")
			continue
		}

		if err := s.dest.Upsert(&out); err != nil {
			log.Error().Err(err).Msg("Port upsert failed")
		}
	}
}

// Stop stops the Service.
func (s *Service) Stop() {
	s.closed = true
}

func translate(in *parser.Port, out *database.Port) error {
	if len(in.Coordinates) != 2 {
		return fmt.Errorf("invalid coordinates: %v", in.Coordinates)
	}

	out.Key = in.Key
	out.Code = in.Code
	out.Name = in.Name
	out.City = in.City
	out.Province = in.Province
	out.Country = in.Country
	out.Timezone = in.Timezone
	out.Latitude = in.Coordinates[0]
	out.Longitude = in.Coordinates[1]
	out.Alias = in.Alias
	out.Unlocs = in.Unlocs

	return nil
}
