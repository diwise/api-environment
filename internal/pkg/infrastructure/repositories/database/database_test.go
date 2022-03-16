package database

import (
	"fmt"
	"testing"
	"time"

	"github.com/matryer/is"
	"github.com/rs/zerolog/log"
)

func TestThatStoreAirQualityObservedStoresStuffCorrectly(t *testing.T) {
	is, db := setupTest(t)

	aqo, err := db.StoreAirQualityObserved("entityId", "deviceId", 15.0, 20.0, 25.0, time.Now().UTC())
	is.NoErr(err) // error when storing new air quality observed...
	is.Equal(aqo.DeviceId, "deviceId")
}

func TestThatGetEntitiesReturnsAllStoredAirQualityObserveds(t *testing.T) {
	is, db := setupTest(t)

	createAirQualityObserveds(db, 3)

	aqos, err := db.GetAirQualityObserveds("", time.Time{}, time.Time{}, 1000)
	is.NoErr(err)
	is.Equal(len(aqos), 3)
}

func setupTest(t *testing.T) (*is.I, Datastore) {
	is := is.New(t)
	log := log.Logger
	db, err := NewDatabaseConnection(NewSQLiteConnector(log))
	is.NoErr(err) // error when creating new database connection

	return is, db
}

func createAirQualityObserveds(db Datastore, times int) {
	i := 0

	for i < times {
		db.StoreAirQualityObserved(fmt.Sprintf("entityId%d", i), fmt.Sprintf("entityId%d", i), 15.0, 20.0, 25.0, time.Now().UTC())
		i++
	}
}
