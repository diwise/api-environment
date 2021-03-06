package database

import (
	"fmt"
	"os"
	"time"

	"github.com/diwise/api-environment/internal/pkg/infrastructure/repositories/models"
	"github.com/rs/zerolog"
	"gorm.io/driver/postgres"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

type Datastore interface {
	GetAirQualityObserveds(deviceId string, from, to time.Time, limit uint64) ([]models.AirQualityObserved, error)
	StoreAirQualityObserved(entityId, deviceId string, co2, humidity, temperature float64, timestamp time.Time) (*models.AirQualityObserved, error)
}

type myDB struct {
	impl *gorm.DB
	log  zerolog.Logger
}

func getEnv(key, fallback string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}
	return fallback
}

//ConnectorFunc is used to inject a database connection method into NewDatabaseConnection
type ConnectorFunc func() (*gorm.DB, zerolog.Logger, error)

//NewSQLiteConnector opens a connection to a local sqlite database
func NewSQLiteConnector(log zerolog.Logger) ConnectorFunc {
	return func() (*gorm.DB, zerolog.Logger, error) {
		db, err := gorm.Open(sqlite.Open("file::memory:"), &gorm.Config{
			Logger: logger.Default.LogMode(logger.Silent),
		})

		if err == nil {
			db.Exec("PRAGMA foreign_keys = ON")
		}

		return db, log, err
	}
}

//NewPostgreSQLConnector opens a connection to a postgresql database
func NewPostgreSQLConnector(log zerolog.Logger) ConnectorFunc {
	dbHost := os.Getenv("DIWISE_SQLDB_HOST")
	username := os.Getenv("DIWISE_SQLDB_USER")
	dbName := os.Getenv("DIWISE_SQLDB_NAME")
	password := os.Getenv("DIWISE_SQLDB_PASSWORD")
	sslMode := getEnv("DIWISE_SQLDB_SSLMODE", "disable")

	dbURI := fmt.Sprintf("host=%s user=%s dbname=%s sslmode=%s password=%s", dbHost, username, dbName, sslMode, password)

	return func() (*gorm.DB, zerolog.Logger, error) {
		sublogger := log.With().Str("host", dbHost).Str("database", dbName).Logger()

		for {
			sublogger.Info().Msg("connecting to database host")
			db, err := gorm.Open(postgres.Open(dbURI), &gorm.Config{
				Logger: logger.New(
					&sublogger,
					logger.Config{
						SlowThreshold:             time.Second,
						LogLevel:                  logger.Info,
						IgnoreRecordNotFoundError: false,
						Colorful:                  false,
					},
				),
			})

			if err != nil {
				sublogger.Fatal().Msg("failed to connect to database")
				time.Sleep(3 * time.Second)
			} else {
				return db, sublogger, nil
			}
		}
	}
}

//NewDatabaseConnection initializes a new connection to the database and wraps it in a Datastore
func NewDatabaseConnection(connect ConnectorFunc) (Datastore, error) {
	impl, log, err := connect()
	if err != nil {
		return nil, err
	}

	db := &myDB{
		impl: impl.Debug(),
		log:  log,
	}

	db.impl.AutoMigrate(
		&models.AirQualityObserved{},
	)

	return db, nil
}

func (db *myDB) StoreAirQualityObserved(entityId, deviceId string, co2, humidity, temperature float64, timestamp time.Time) (*models.AirQualityObserved, error) {
	aqo := models.AirQualityObserved{
		EntityId:    entityId,
		DeviceId:    deviceId,
		CO2:         co2,
		Humidity:    humidity,
		Temperature: temperature,
		Timestamp:   timestamp,
	}

	result := db.impl.Create(&aqo)
	if result.Error != nil {
		return nil, result.Error
	}

	return &aqo, nil
}

func (db *myDB) GetAirQualityObserveds(deviceId string, from, to time.Time, limit uint64) ([]models.AirQualityObserved, error) {
	aqos := []models.AirQualityObserved{}
	gorm := db.impl.Order("timestamp DESC")

	if deviceId != "" {
		gorm = gorm.Where("device = ?", deviceId)
	}

	if !from.IsZero() || !to.IsZero() {
		gorm = insertTemporalSQL(gorm, "timestamp", from, to)
		if gorm.Error != nil {
			return nil, gorm.Error
		}
	}

	result := gorm.Limit(int(limit)).Find(&aqos)
	if result.Error != nil {
		return nil, result.Error
	}

	return aqos, nil
}

func insertTemporalSQL(gorm *gorm.DB, property string, from, to time.Time) *gorm.DB {
	if !from.IsZero() {
		gorm = gorm.Where(fmt.Sprintf("%s >= ?", property), from)
		if gorm.Error != nil {
			return gorm
		}
	}

	if !to.IsZero() {
		gorm = gorm.Where(fmt.Sprintf("%s < ?", property), to)
	}

	return gorm
}
