package database

import (
	"testing"
	"time"

	"github.com/matryer/is"
	"github.com/rs/zerolog/log"
)

func TestThatCreateAirQualityObservedDoesNot(t *testing.T) {
	is, db := setupTest(t)

	_, err := db.StoreAirQualityObserved("", "", 0.0, 0.0, 0.0, time.Now().UTC())
	is.NoErr(err)
}

func setupTest(t *testing.T) (*is.I, Datastore) {
	is := is.New(t)
	log := log.Logger
	db, err := NewDatabaseConnection(NewSQLiteConnector(log))
	is.NoErr(err) // error when creating new database connection

	return is, db
}
