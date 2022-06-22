// Package service includes the port service.
package service

import (
	"context"
	"fmt"

	"github.com/agukrapo/ports/database"
	"github.com/agukrapo/ports/parser"
	"github.com/rs/zerolog/log"
)

type source interface {
	Stream(context.Context) chan parser.Packet
}

type storage interface {
	Upsert(*database.Port) error
}

// Service represents a process that moves ports from a source to a destination.
type Service struct {
	storage storage
}

// New instantiates a new Service.
func New(storage storage) *Service {
	return &Service{
		storage: storage,
	}
}

// Process moves a Port from a source to the storage.
func (s *Service) Process(ctx context.Context, src source) {
	var out database.Port

	for in := range src.Stream(ctx) {
		if in.Err != nil {
			log.Error().Err(in.Err).Msg("Port parse failed")
			continue
		}

		if err := translate(in.Port, &out); err != nil {
			log.Error().Err(err).Msg("Port translation failed")
			continue
		}

		if err := s.storage.Upsert(&out); err != nil {
			log.Error().Err(err).Msg("Port upsert failed")
		}
	}
}

func translate(in *parser.Port, out *database.Port) error {
	if len(in.Coordinates) != 2 {
		return fmt.Errorf("key %s: invalid coordinates: %v", in.Key, in.Coordinates)
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
