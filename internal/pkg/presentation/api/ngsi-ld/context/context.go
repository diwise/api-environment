package context

import (
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/diwise/api-environment/internal/pkg/application"
	"github.com/diwise/ngsi-ld-golang/pkg/datamodels/fiware"
	"github.com/diwise/ngsi-ld-golang/pkg/ngsi-ld"
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
		return err
	}

	dateObserved, err := time.Parse(time.RFC3339, aqo.DateObserved.Value.Value)
	if err != nil {
		return err
	}

	err = cs.app.StoreAirQualityObserved(aqo.ID, aqo.RefDevice.Object, aqo.CO2.Value, aqo.RelativeHumidity.Value, aqo.Temperature.Value, dateObserved)

	return err
}

func (cs contextSource) GetEntities(query ngsi.Query, callback ngsi.QueryEntitiesCallback) error {
	return errors.New("not implemented yet")
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
