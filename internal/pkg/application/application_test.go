package application

import (
	"testing"
	"time"

	"github.com/diwise/api-environment/internal/pkg/infrastructure/repositories/database"
	"github.com/diwise/api-environment/internal/pkg/infrastructure/repositories/models"
	"github.com/matryer/is"
	"github.com/rs/zerolog/log"
)

func newAppForTesting() (*database.DatastoreMock, EnvironmentApp) {
	db := &database.DatastoreMock{
		StoreAirQualityObservedFunc: func(entityId, deviceId string, co2, humidity, temperature float64, timestamp time.Time) (*models.AirQualityObserved, error) {
			return nil, nil
		},
	}

	return db, NewEnvironmentApp(db, log.Logger)
}

func TestStoreAirQuality(t *testing.T) {
	is := is.New(t)
	db, app := newAppForTesting()

	err := app.StoreAirQualityObserved("aqoID", "refDeviceId", 0.0, 0.0, 0.0, time.Now().UTC())
	is.NoErr(err)
	is.Equal(len(db.StoreAirQualityObservedCalls()), 1)
}
