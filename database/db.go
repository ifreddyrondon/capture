package database

import (
	"context"
	"log"
	"net/http"

	"gopkg.in/mgo.v2"
)

var DB *mgo.Session

func MongoCtx(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		session := DB.Copy()
		defer session.Close()
		ctx := context.WithValue(r.Context(), "DB", session.DB(""))
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

// CreateConnection establishes a new session with the mongod server
// [mongodb://][user:pass@]host1[:port1][,host2[:port2],...][/database][?options]
func CreateConnection(url string) {
	var err error
	DB, err = mgo.Dial(url)
	if err != nil {
		log.Fatal(err)
	}
}
