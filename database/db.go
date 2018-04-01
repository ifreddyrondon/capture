package database

import (
	"log"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
)

type Schemator interface {
	// Create (panic) runs schema migration.
	Create()
	// Drop (panic) delete schema.
	Drop()
}

type DataSource struct {
	DB *gorm.DB
}

// TODO: load this func into bastion RegisterOnShutdown
// Finalize implements the Finalizer interface from bastion to be executed as graceful shutdown.
func (ds *DataSource) OnShutdown() {
	log.Printf("[finalizer:data source] closing the main session")
	ds.DB.Close()
}

// Open establishes a connection with the database server and verify with a ping.
// it returns a *DataSource
func Open(url string) *DataSource {
	db, err := gorm.Open("postgres", url)
	if err != nil {
		log.Panic(err)
	}

	return &DataSource{DB: db}
}
