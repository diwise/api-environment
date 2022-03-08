package context

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/diwise/api-environment/internal/pkg/application"
	"github.com/diwise/api-environment/internal/pkg/infrastructure/repositories/models"
	"github.com/diwise/ngsi-ld-golang/pkg/ngsi-ld"
	"github.com/matryer/is"
	"github.com/rs/zerolog/log"
)

func TestStoreAirQualityObserved(t *testing.T) {
	req, _ := http.NewRequest("POST", "/ngsi-ld/v1/entities", bytes.NewBuffer([]byte(aqoJson)))
	w := httptest.NewRecorder()

	is, app, ctxReg := testSetup(t)

	ngsi.NewCreateEntityHandler(ctxReg).ServeHTTP(w, req)

	is.Equal(w.Code, http.StatusCreated)
	is.Equal(len(app.StoreAirQualityObservedCalls()), 1)
}

func TestRetrieveAirQualityObserveds(t *testing.T) {
	req, _ := http.NewRequest("GET", "/ngsi-ld/v1/entities?type=AirQualityObserved", nil)
	w := httptest.NewRecorder()

	is, app, ctxReg := testSetup(t)

	ngsi.NewQueryEntitiesHandler(ctxReg).ServeHTTP(w, req)

	is.Equal(w.Code, http.StatusOK)
	is.Equal(len(app.RetrieveAirQualityObservedsCalls()), 1)
}

func testSetup(t *testing.T) (*is.I, *application.EnvironmentAppMock, ngsi.ContextRegistry) {
	is := is.New(t)

	log := log.Logger
	app := &application.EnvironmentAppMock{
		StoreAirQualityObservedFunc: func(entityId, deviceId string, co2, humidity, temperature float64, timestamp time.Time) error {
			return nil
		},
		RetrieveAirQualityObservedsFunc: func() ([]models.AirQualityObserved, error) {
			return nil, nil
		},
	}

	ctxReg := ngsi.NewContextRegistry()
	ctxSource := CreateSource(app, log)
	ctxReg.Register(ctxSource)

	return is, app, ctxReg
}

const aqoJson string = `{
    "id": "urn:ngsi-ld:AirQualityObserved:Madrid-AmbientObserved-28079004-2016-03-15T11:00:00",
    "type": "AirQualityObserved",
    "dateObserved": {
        "value": {
			"@type": "Property",
        	"@value": "2016-03-15T11:00:00Z"
		}
    },
    "temperature": {
        "type": "Property",
        "value": 12.2
    },
    "location": {
        "type": "GeoProperty",
        "value": {
            "type": "Point",
            "coordinates": [-3.712247222222222, 40.423852777777775]
        }
    },
    "relativeHumidity": {
        "type": "Property",
        "value": 0.54
    },
    "CO2": {
        "type": "Property",
        "value": 500,
        "unitCode": "GP"
    },
    "@context": [
        "https://schema.lab.fiware.org/ld/context",
        "https://uri.etsi.org/ngsi-ld/v1/ngsi-ld-core-context.jsonld"
    ]
}`
