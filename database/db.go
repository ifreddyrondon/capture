package database

import (
	"context"
	"net/http"

	"gopkg.in/mgo.v2"
)

type DB struct {
	Session *mgo.Session
}

// Close terminates the session. It's a runtime error to use a session
// after it has been closed.
func (db *DB) Close() {
	if db.Session != nil {
		db.Session.Close()
	}
}

// CopySession return a session with the same parameters as the original
// and preserves the exact authentication information from the original session.
func (db *DB) CopySession() *mgo.Session {
	if db.Session != nil {
		return db.Session.Copy()
	}
	return nil
}

// Ctx set into the context request a DB value with a session to perform actions.
func (db *DB) Ctx(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		session := db.CopySession()
		defer session.Close()
		ctx := context.WithValue(r.Context(), "DB", session.DB(""))
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

// Open establishes a new session with the mongod server, it returns a *DB
// [mongodb://][user:pass@]host1[:port1][,host2[:port2],...][/database][?options]
func Open(url string) (*DB, error) {
	var err error
	db := new(DB)
	db.Session, err = mgo.Dial(url)
	return db, err
}
