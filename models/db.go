package models

import (
	"database/sql"
	"log"
	"os"
)

type Datastore interface {
}

type DB struct {
	*sql.DB
}

// OpenDB opens a connection and returns a DB struct with an active handle to the sqlite db.
func OpenDB(dbPath string) (*DB, error) {
	// sqlite setup and verification
	db, err := sql.Open("sqlite3", dbPath)
	if err != nil {
		log.Printf("Error when opening sqlite3: %s\n", err)
		return nil, err
	}
	if db == nil {
		log.Printf("db nil when opened")
		return nil, err
	}

	// sql.open may just validate its arguments without creating a connection to the database.
	// To verify that the data source name is valid, call Ping.
	err = db.Ping()
	if err != nil {
		log.Printf("Error when pinging db: %s\n", err)
		return nil, err
	}

	return &DB{db}, nil
}

// NewDB returns a connection to new database.
//
// If the given name exists, rename it to *.old, overwrites any existing *.old db.
func NewDB(dbPath string) (*DB, error) {
	if _, err := os.Stat(dbPath); err == nil {
		log.Println(dbPath, "existed, renaming to .old.")
		err := os.Rename(dbPath, dbPath+".old")
		if err != nil {
			log.Printf("Error when renaming db: %s\n", err)
			return nil, err
		}
	}
	return OpenDB(dbPath)
}
