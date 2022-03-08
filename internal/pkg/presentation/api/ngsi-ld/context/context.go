package context

import (
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/diwise/api-environment/internal/pkg/application"
	"github.com/diwise/ngsi-ld-golang/pkg/datamodels/fiware"
	"github.com/diwise/ngsi-ld-golang/pkg/ngsi-ld"
	"github.com/diwise/ngsi-ld-golang/pkg/ngsi-ld/types"
	"github.com/rs/zerolog"
)

type contextSource struct {
	app application.EnvironmentApp
	log zerolog.Logger
}

//CreateSource instantiates and returns a Fiware ContextSource that wraps the provided application interface
func CreateSource(app application.EnvironmentApp, log zerolog.Logger) ngsi.ContextSource {
	return &contextSource{
		app: app,
		log: log,
	}
}

func (cs contextSource) CreateEntity(typeName, entityID string, req ngsi.Request) error {
	if typeName != fiware.AirQualityObservedTypeName {
		errorMessage := fmt.Sprintf("entity type %s not supported", typeName)
		cs.log.Error().Msg(errorMessage)
		return errors.New(errorMessage)
	}

	aqo := &fiware.AirQualityObserved{}
	err := req.DecodeBodyInto(aqo)
	if err != nil {
		fmt.Print(err.Error())
		return err
	}

	dateObserved, err := time.Parse(time.RFC3339, aqo.DateObserved.Value.Value)
	if err != nil {
		return err
	}

	entity := strings.TrimPrefix(aqo.ID, fiware.AirQualityObservedIDPrefix)

	refDevice := ""
	if aqo.RefDevice != nil {
		refDevice = strings.TrimPrefix(aqo.RefDevice.Object, fiware.DeviceIDPrefix)
	}

	co2 := 0.0
	if aqo.CO2 != nil {
		co2 = aqo.CO2.Value
	}

	humidity := 0.0
	if aqo.RelativeHumidity != nil {
		humidity = aqo.RelativeHumidity.Value
	}

	temp := 0.0
	if aqo.Temperature != nil {
		temp = aqo.Temperature.Value
	}

	err = cs.app.StoreAirQualityObserved(entity, refDevice, co2, humidity, temp, dateObserved)

	return err
}

func (cs contextSource) GetEntities(query ngsi.Query, callback ngsi.QueryEntitiesCallback) error {
	var err error

	if query == nil {
		return errors.New("GetEntities: query may not be nil")
	}

	deviceId := ""
	if query.HasDeviceReference() {
		deviceId = strings.TrimPrefix(query.Device(), fiware.DeviceIDPrefix)
	}

	from := time.Time{}
	to := time.Time{}
	if query.IsTemporalQuery() {
		from, to = query.Temporal().TimeSpan()
	}

	limit := query.PaginationLimit()

	aqos, err := cs.app.RetrieveAirQualityObserveds(deviceId, from, to, limit)
	if err != nil {
		return err
	}

	for _, a := range aqos {
		entity := fiware.NewAirQualityObserved(a.EntityId, a.Latitude, a.Longitude, a.Timestamp.Format(time.RFC3339)).WithCO2(a.CO2).WithRelativeHumidity(a.Humidity).WithTemperature(a.Temperature)
		if a.DeviceId != "" {
			entity.RefDevice = types.NewSingleObjectRelationship(a.DeviceId)
		}

		err = callback(entity)
		if err != nil {
			break
		}
	}

	return err
}

func (cs contextSource) GetProvidedTypeFromID(entityID string) (string, error) {
	if cs.ProvidesEntitiesWithMatchingID(entityID) {
		return fiware.AirQualityObservedTypeName, nil
	}

	return "", errors.New("no entities found with matching type")
}

func (cs contextSource) ProvidesAttribute(attributeName string) bool {
	return attributeName == "airquality"
}

func (cs contextSource) ProvidesEntitiesWithMatchingID(entityID string) bool {
	return strings.HasPrefix(entityID, fiware.AirQualityObservedIDPrefix)
}

func (cs contextSource) ProvidesType(typeName string) bool {
	return typeName == fiware.AirQualityObservedTypeName
}

func (cs contextSource) RetrieveEntity(entityID string, request ngsi.Request) (ngsi.Entity, error) {
	return nil, errors.New("retrieve entity not implemented")
}

func (cs contextSource) UpdateEntityAttributes(entityID string, req ngsi.Request) error {
	return errors.New("UpdateEntityAttributes is not supported by this service")
}
