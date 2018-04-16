package database

import (
	"log"

	"github.com/jinzhu/gorm"
	// register postgres drive
	_ "github.com/jinzhu/gorm/dialects/postgres"
)

// DataSource wrapper over the DB driver
type DataSource struct {
	*gorm.DB
}

// OnShutdown is executed as graceful shutdown.
func (ds *DataSource) OnShutdown() {
	log.Printf("[finalizer:data source] closing the main session")
	if err := ds.Close(); err != nil {
		log.Fatal(err)
	}
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
