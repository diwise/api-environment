package models

import (
	"time"

	"gorm.io/gorm"
)

type AirQualityObserved struct {
	gorm.Model
	DeviceId    string
	CO2         float64
	Humidity    float64
	Temperature float64
	Latitude    float64
	Longitude   float64
	Timestamp   time.Time
}
