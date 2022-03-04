package application

import (
	"time"

	"github.com/diwise/api-environment/internal/pkg/infrastructure/repositories/database"
	"github.com/rs/zerolog"
)

type EnvironmentApp interface {
	StoreAirQualityObserved(entityId, deviceId string, co2, humidity, temperature float64, timestamp time.Time) error
}

type app struct {
	db  database.Datastore
	log zerolog.Logger
}

func NewEnvironmentApp(db database.Datastore, log zerolog.Logger) EnvironmentApp {
	newApp := &app{
		db:  db,
		log: log,
	}

	return newApp
}

func (a *app) StoreAirQualityObserved(entityId, deviceId string, co2, humidity, temperature float64, timestamp time.Time) error {
	_, err := a.db.CreateAirQualityObserved(entityId, deviceId, co2, humidity, temperature, timestamp)
	if err != nil {
		return err
	}
	return nil
}
