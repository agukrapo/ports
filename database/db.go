// Package database includes data storage related utilities.
package database

import (
	"github.com/lib/pq"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
	"gorm.io/gorm/logger"
)

// Database represents a Postgres storage.
type Database struct {
	db *gorm.DB
}

// New instantiates a new Database.
func New(dsn string) (*Database, error) {
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Silent),
	})
	if err != nil {
		return nil, err
	}

	if err := db.AutoMigrate(&Port{}); err != nil {
		return nil, err
	}

	return &Database{
		db: db,
	}, nil
}

// Close releases Database resources.
func (db *Database) Close() error {
	sqlDB, err := db.db.DB()
	if err != nil {
		return err
	}

	return sqlDB.Close()
}

// Port represents a ports database table.
type Port struct {
	Key       string `gorm:"primarykey"`
	Code      string
	Name      string
	City      string
	Province  string
	Country   string
	Timezone  string
	Latitude  float64
	Longitude float64
	Unlocs    pq.StringArray `gorm:"type:text[]"`
	Alias     pq.StringArray `gorm:"type:text[]"`
}

// Upsert inserts a new Port, or updates it if already present.
func (db *Database) Upsert(port *Port) error {
	tx := db.db.Clauses(clause.OnConflict{
		Columns:   []clause.Column{{Name: "key"}},
		UpdateAll: true,
	}).Create(&port)

	return tx.Error
}
