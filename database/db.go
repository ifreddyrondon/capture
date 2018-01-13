package database

import (
	"context"
	"net/http"

	"log"

	"gopkg.in/mgo.v2"
)

type DataSource struct {
	session *mgo.Session
}

// Finalize implements the Finalizer interface from bastion to be executed as graceful shutdown.
func (ds *DataSource) Finalize() error {
	log.Printf("[finalizer:data source] closing the main session")
	ds.session.Close()
	return nil
}

// GetCtx implements the DataSource interface.
func (ds *DataSource) Ctx() func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if ds.session == nil {
				return
			}
			session := ds.session.Copy()
			defer session.Close()
			ctx := context.WithValue(r.Context(), "DataSource", session.DB(""))
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

// Finalize implements the Finalizer interface from bastion to be executed as graceful shutdown.
func (ds *DataSource) DB() *mgo.Database {
	return ds.session.DB("")
}

// Open establishes a new session with the mongod server, it returns a *DataSource
// [mongodb://][user:pass@]host1[:port1][,host2[:port2],...][/database][?options]
func Open(url string) (*DataSource, error) {
	var err error
	db := new(DataSource)
	db.session, err = mgo.Dial(url)
	return db, err
}
