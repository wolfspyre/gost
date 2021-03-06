package postgis

import (
	"fmt"
	"github.com/geodan/gost/src/sensorthings/entities"
)

// tables as defined in postgis
var (
	thingTable                        = "thing"
	locationTable                     = "location"
	historicalLocationTable           = "historicallocation"
	sensorTable                       = "sensor"
	observedPropertyTable             = "observedproperty"
	datastreamTable                   = "datastream"
	observationTable                  = "observation"
	featureOfInterestTable            = "featureofinterest"
	thingToLocationTable              = "thing_to_location"
	locationToHistoricalLocationTable = "location_to_historicallocation"
)

// thing fields
var (
	thingID          = "id"
	thingName        = "name"
	thingDescription = "description"
	thingProperties  = "properties"
)

// location fields
var (
	locationID           = "id"
	locationName         = "name"
	locationDescription  = "description"
	locationEncodingType = "encodingtype"
	locationLocation     = "location"
)

// thingToLocationTable fields
var (
	thingToLocationThingID    = "thing_id"
	thingToLocationLocationID = "location_id"
)

// historical location fields
var (
	historicalLocationID         = "id"
	historicalLocationTime       = "time"
	historicalLocationThingID    = "thing_id"
	historicalLocationLocationID = "location_id"
)

// locationToHistoricalLocation fields
var (
	locationToHistoricalLocationLocationID           = "location_id "
	locationToHistoricalLocationHistoricalLocationID = "historicallocation_id "
)

var (
	asSeparator = "_"
	idField     = "id"
)

// sensor fields
var (
	sensorID           = idField
	sensorName         = "name"
	sensorDescription  = "description"
	sensorEncodingType = "encodingtype"
	sensorMetadata     = "metadata"
)

// observed property fields
var (
	observedPropertyID          = idField
	observedPropertyName        = "name"
	observedPropertyDescription = "description"
	observedPropertyDefinition  = "definition"
)

// datastream fields
var (
	datastreamID                 = idField
	datastreamName               = "name"
	datastreamDescription        = "description"
	datastreamUnitOfMeasurement  = "unitofmeasurement"
	datastreamObservationType    = "observationtype"
	datastreamObservedArea       = "observedarea"
	datastreamPhenomenonTime     = "phenomenontime"
	datastreamResultTime         = "resulttime"
	datastreamThingID            = "thing_id"
	datastreamSensorID           = "sensor_id"
	datastreamObservedPropertyID = "observedproperty_id"
)

// observation fields
var (
	observationID                  = idField
	observationData                = "data"
	observationPhenomenonTime      = "phenomenontime"
	observationResultTime          = "resulttime"
	observationResult              = "result"
	observationValidTime           = "validtime"
	observationResultQuality       = "resultquality"
	observationParameters          = "parameters"
	observationStreamID            = "stream_id"
	observationFeatureOfInterestID = "featureofinterest_id"
)

// feature of interest fields
var (
	foiID                 = idField
	foiName               = "name"
	foiDescription        = "description"
	foiEncodingType       = "encodingtype"
	foiFeature            = "feature"
	foiOriginalLocationID = "original_location_id"
)

type ParamFactory func(values map[string]interface{}) (entities.Entity, error)

type QueryParseInfo struct {
	QueryIndex   int
	ParamFactory ParamFactory
	EntityType   entities.EntityType
	Entity       entities.Entity
	SubEntities  []*QueryParseInfo
}

func (q *QueryParseInfo) Init(entityType entities.EntityType, queryIndex int) {
	q.QueryIndex = queryIndex
	q.EntityType = entityType
	switch e := entityType; e {
	case entities.EntityTypeThing:
		q.Entity = &entities.Thing{}
		q.ParamFactory = thingParamFactory
		break
	case entities.EntityTypeFeatureOfInterest:
		q.Entity = &entities.FeatureOfInterest{}
		q.ParamFactory = featureOfInterestParamFactory
		break
	case entities.EntityTypeLocation:
		q.Entity = &entities.Location{}
		q.ParamFactory = locationParamFactory
		break
	case entities.EntityTypeObservation:
		q.Entity = &entities.Observation{}
		q.ParamFactory = observationParamFactory
		break
	case entities.EntityTypeObservedProperty:
		q.Entity = &entities.ObservedProperty{}
		q.ParamFactory = observedPropertyParamFactory
		break
	case entities.EntityTypeDatastream:
		q.Entity = &entities.Datastream{}
		q.ParamFactory = datastreamParamFactory
		break
	case entities.EntityTypeHistoricalLocation:
		q.Entity = &entities.HistoricalLocation{}
		q.ParamFactory = historicalLocationParamFactory
		break
	case entities.EntityTypeSensor:
		q.Entity = &entities.Sensor{}
		q.ParamFactory = sensorParamFactory
		break
	}
}

func (q *QueryParseInfo) GetQueryParseInfoByQueryIndex(id int) *QueryParseInfo {
	if q.QueryIndex == id {
		return q
	}

	for _, qpi := range q.SubEntities {
		t := qpi.GetQueryParseInfoByQueryIndex(id)
		if t != nil {
			return t
		}
	}

	return nil
}

// GetNextQueryIndex returns the next query index number based on the added entities/sub entities
func (q *QueryParseInfo) GetNextQueryIndex() int {
	qi := q.QueryIndex
	if len(q.SubEntities) > 0 {
		lastSub := q.SubEntities[len(q.SubEntities)-1]
		qi = lastSub.GetNextQueryIndex()
	}

	return qi + 1
}

// GetQueryIDRelationMap returns the query index relations, ie QueryParseInfo with sub entity datastream thing qid = 0, datastream qid = 1
// example: returns [1]0 - datastream (1) relates to thing (0)
// returns nil if no relation exists
func (q *QueryParseInfo) GetQueryIDRelationMap(relationMap map[int]int) map[int]int {
	if relationMap == nil {
		relationMap = map[int]int{}
	}

	if len(q.SubEntities) == 0 {
		return relationMap
	}

	for _, qpi := range q.SubEntities {
		relationMap[qpi.QueryIndex] = q.QueryIndex
		relationMap = qpi.GetQueryIDRelationMap(relationMap)
	}

	return relationMap
}

func (q *QueryParseInfo) Parse(values map[string]interface{}) (entities.Entity, error) {
	return q.ParamFactory(values)
}

var asMappings = map[entities.EntityType]map[string]string{
	entities.EntityTypeThing: {
		thingID:          constructAs(thingTable, thingID),
		thingName:        constructAs(thingTable, thingName),
		thingDescription: constructAs(thingTable, thingDescription),
		thingProperties:  constructAs(thingTable, thingProperties),
	},
	entities.EntityTypeLocation: {
		locationID:           constructAs(locationTable, locationID),
		locationName:         constructAs(locationTable, locationName),
		locationDescription:  constructAs(locationTable, locationDescription),
		locationEncodingType: constructAs(locationTable, locationEncodingType),
		locationLocation:     constructAs(locationTable, locationLocation),
	},
	entities.EntityTypeThingToLocation: {
		thingToLocationThingID:    constructAs(thingToLocationTable, thingToLocationThingID),
		thingToLocationLocationID: constructAs(thingToLocationTable, thingToLocationLocationID),
	},
	entities.EntityTypeLocationToHistoricalLocation: {
		locationToHistoricalLocationLocationID:           constructAs(locationToHistoricalLocationTable, locationToHistoricalLocationLocationID),
		locationToHistoricalLocationHistoricalLocationID: constructAs(locationToHistoricalLocationTable, locationToHistoricalLocationHistoricalLocationID),
	},
	entities.EntityTypeHistoricalLocation: {
		historicalLocationID:         constructAs(historicalLocationTable, historicalLocationID),
		historicalLocationTime:       constructAs(historicalLocationTable, historicalLocationTime),
		historicalLocationThingID:    constructAs(historicalLocationTable, historicalLocationThingID),
		historicalLocationLocationID: constructAs(historicalLocationTable, historicalLocationLocationID),
	},
	entities.EntityTypeSensor: {
		sensorID:           constructAs(sensorTable, sensorID),
		sensorName:         constructAs(sensorTable, sensorName),
		sensorDescription:  constructAs(sensorTable, sensorDescription),
		sensorEncodingType: constructAs(sensorTable, sensorEncodingType),
		sensorMetadata:     constructAs(sensorTable, sensorMetadata),
	},
	entities.EntityTypeObservedProperty: {
		observedPropertyID:          constructAs(observedPropertyTable, observedPropertyID),
		observedPropertyName:        constructAs(observedPropertyTable, observedPropertyName),
		observedPropertyDescription: constructAs(observedPropertyTable, observedPropertyDescription),
		observedPropertyDefinition:  constructAs(observedPropertyTable, observedPropertyDefinition),
	},
	entities.EntityTypeObservation: {
		observationID:                  constructAs(observationTable, observationID),
		observationData:                constructAs(observationTable, observationData),
		observationPhenomenonTime:      constructAs(observationTable, observationPhenomenonTime),
		observationResultTime:          constructAs(observationTable, observationResultTime),
		observationResult:              constructAs(observationTable, observationResult),
		observationValidTime:           constructAs(observationTable, observationValidTime),
		observationResultQuality:       constructAs(observationTable, observationResultQuality),
		observationParameters:          constructAs(observationTable, observationParameters),
		observationStreamID:            constructAs(observationTable, observationStreamID),
		observationFeatureOfInterestID: constructAs(observationTable, observationFeatureOfInterestID),
	},
	entities.EntityTypeFeatureOfInterest: {
		foiID:                 constructAs(featureOfInterestTable, foiID),
		foiName:               constructAs(featureOfInterestTable, foiName),
		foiDescription:        constructAs(featureOfInterestTable, foiDescription),
		foiEncodingType:       constructAs(featureOfInterestTable, foiEncodingType),
		foiFeature:            constructAs(featureOfInterestTable, foiFeature),
		foiOriginalLocationID: constructAs(featureOfInterestTable, foiOriginalLocationID),
	},
	entities.EntityTypeDatastream: {
		datastreamID:                 constructAs(datastreamTable, datastreamID),
		datastreamName:               constructAs(datastreamTable, datastreamName),
		datastreamDescription:        constructAs(datastreamTable, datastreamDescription),
		datastreamUnitOfMeasurement:  constructAs(datastreamTable, datastreamUnitOfMeasurement),
		datastreamObservationType:    constructAs(datastreamTable, datastreamObservationType),
		datastreamObservedArea:       constructAs(datastreamTable, datastreamObservedArea),
		datastreamPhenomenonTime:     constructAs(datastreamTable, datastreamPhenomenonTime),
		datastreamResultTime:         constructAs(datastreamTable, datastreamResultTime),
		datastreamThingID:            constructAs(datastreamTable, datastreamThingID),
		datastreamSensorID:           constructAs(datastreamTable, datastreamSensorID),
		datastreamObservedPropertyID: constructAs(datastreamTable, datastreamObservedPropertyID),
	},
}

func constructAs(table, field string) string {
	return fmt.Sprintf("%s%s%s", table, asSeparator, field)
}

var tableMappings = map[entities.EntityType]string{
	entities.EntityTypeThing:              thingTable,
	entities.EntityTypeLocation:           locationTable,
	entities.EntityTypeThingToLocation:    thingToLocationTable,
	entities.EntityTypeHistoricalLocation: historicalLocationTable,
	entities.EntityTypeSensor:             sensorTable,
	entities.EntityTypeObservedProperty:   observedPropertyTable,
	entities.EntityTypeObservation:        observationTable,
	entities.EntityTypeFeatureOfInterest:  featureOfInterestTable,
	entities.EntityTypeDatastream:         datastreamTable,
}

// maps an entity property name to the right field
var selectMappings = map[entities.EntityType]map[string]string{
	entities.EntityTypeThing: {
		thingID:          fmt.Sprintf("%s.%s", thingTable, thingID),
		thingName:        fmt.Sprintf("%s.%s", thingTable, thingName),
		thingDescription: fmt.Sprintf("%s.%s", thingTable, thingDescription),
		thingProperties:  fmt.Sprintf("%s.%s", thingTable, thingProperties),
	},
	entities.EntityTypeLocation: {
		locationID:           fmt.Sprintf("%s.%s", locationTable, locationID),
		locationName:         fmt.Sprintf("%s.%s", locationTable, locationName),
		locationDescription:  fmt.Sprintf("%s.%s", locationTable, locationDescription),
		locationEncodingType: fmt.Sprintf("%s.%s", locationTable, locationEncodingType),
		locationLocation:     fmt.Sprintf("public.ST_AsGeoJSON(%s.%s)", locationTable, locationLocation),
	},
	entities.EntityTypeThingToLocation: {
		thingToLocationThingID:    fmt.Sprintf("%s.%s", thingToLocationTable, thingToLocationThingID),
		thingToLocationLocationID: fmt.Sprintf("%s.%s", thingToLocationTable, thingToLocationLocationID),
	},
	entities.EntityTypeHistoricalLocation: {
		historicalLocationID:         fmt.Sprintf("%s.%s", historicalLocationTable, historicalLocationID),
		historicalLocationTime:       fmt.Sprintf("%s.%s", historicalLocationTable, historicalLocationTime),
		historicalLocationThingID:    fmt.Sprintf("%s.%s", historicalLocationTable, historicalLocationThingID),
		historicalLocationLocationID: fmt.Sprintf("%s.%s", historicalLocationTable, historicalLocationLocationID),
	},
	entities.EntityTypeSensor: {
		sensorID:           fmt.Sprintf("%s.%s", sensorTable, sensorID),
		sensorName:         fmt.Sprintf("%s.%s", sensorTable, sensorName),
		sensorDescription:  fmt.Sprintf("%s.%s", sensorTable, sensorDescription),
		sensorEncodingType: fmt.Sprintf("%s.%s", sensorTable, sensorEncodingType),
		sensorMetadata:     fmt.Sprintf("%s.%s", sensorTable, sensorMetadata),
	},
	entities.EntityTypeObservedProperty: {
		observedPropertyID:          fmt.Sprintf("%s.%s", observedPropertyTable, observedPropertyID),
		observedPropertyName:        fmt.Sprintf("%s.%s", observedPropertyTable, observedPropertyName),
		observedPropertyDescription: fmt.Sprintf("%s.%s", observedPropertyTable, observedPropertyDescription),
		observedPropertyDefinition:  fmt.Sprintf("%s.%s", observedPropertyTable, observedPropertyDefinition),
	},
	entities.EntityTypeObservation: {
		observationID:                  fmt.Sprintf("%s.%s", observationTable, observationID),
		observationData:                fmt.Sprintf("%s.%s", observationTable, observationData),
		observationPhenomenonTime:      fmt.Sprintf("%s.%s -> '%s'", observationTable, observationData, "phenomenonTime"),
		observationResultTime:          fmt.Sprintf("%s.%s -> '%s'", observationTable, observationData, "resultTime"),
		observationResult:              fmt.Sprintf("%s.%s -> '%s'", observationTable, observationData, observationResult),
		observationValidTime:           fmt.Sprintf("%s.%s -> '%s'", observationTable, observationData, "validTime"),
		observationResultQuality:       fmt.Sprintf("%s.%s -> '%s'", observationTable, observationData, "resultQuality"),
		observationParameters:          fmt.Sprintf("%s.%s -> '%s'", observationTable, observationData, observationParameters),
		observationStreamID:            fmt.Sprintf("%s.%s", observationTable, observationStreamID),
		observationFeatureOfInterestID: fmt.Sprintf("%s.%s", observationTable, observationFeatureOfInterestID),
	},
	entities.EntityTypeFeatureOfInterest: {
		foiID:                 fmt.Sprintf("%s.%s", featureOfInterestTable, foiID),
		foiName:               fmt.Sprintf("%s.%s", featureOfInterestTable, foiName),
		foiDescription:        fmt.Sprintf("%s.%s", featureOfInterestTable, foiDescription),
		foiEncodingType:       fmt.Sprintf("%s.%s", featureOfInterestTable, foiEncodingType),
		foiFeature:            fmt.Sprintf("public.ST_AsGeoJSON(%s.%s)", featureOfInterestTable, foiFeature),
		foiOriginalLocationID: fmt.Sprintf("%s.%s", featureOfInterestTable, foiOriginalLocationID),
	},
	entities.EntityTypeDatastream: {
		datastreamID:                 fmt.Sprintf("%s.%s", datastreamTable, datastreamID),
		datastreamName:               fmt.Sprintf("%s.%s", datastreamTable, datastreamName),
		datastreamDescription:        fmt.Sprintf("%s.%s", datastreamTable, datastreamDescription),
		datastreamUnitOfMeasurement:  fmt.Sprintf("%s.%s", datastreamTable, datastreamUnitOfMeasurement),
		datastreamObservationType:    fmt.Sprintf("%s.%s", datastreamTable, datastreamObservationType),
		datastreamObservedArea:       fmt.Sprintf("public.ST_AsGeoJSON(%s.%s)", datastreamTable, datastreamObservedArea),
		datastreamPhenomenonTime:     fmt.Sprintf("%s.%s", datastreamTable, datastreamPhenomenonTime),
		datastreamResultTime:         fmt.Sprintf("%s.%s", datastreamTable, datastreamResultTime),
		datastreamThingID:            fmt.Sprintf("%s.%s", datastreamTable, datastreamThingID),
		datastreamSensorID:           fmt.Sprintf("%s.%s", datastreamTable, datastreamSensorID),
		datastreamObservedPropertyID: fmt.Sprintf("%s.%s", datastreamTable, datastreamObservedPropertyID),
	},
}

func createJoinMappings(tableMappings map[entities.EntityType]string) map[entities.EntityType]map[entities.EntityType]string {
	joinMappings := map[entities.EntityType]map[entities.EntityType]string{
		entities.EntityTypeThing: { // get thing by ...
			entities.EntityTypeDatastream:         fmt.Sprintf("WHERE %s = %s", selectMappings[entities.EntityTypeThing][thingID], selectMappings[entities.EntityTypeDatastream][datastreamThingID]),
			entities.EntityTypeHistoricalLocation: fmt.Sprintf("WHERE %s = %s", selectMappings[entities.EntityTypeThing][thingID], selectMappings[entities.EntityTypeHistoricalLocation][historicalLocationThingID]),
			entities.EntityTypeLocation: fmt.Sprintf("INNER JOIN %s ON %s = %s AND %s = %s",
				tableMappings[entities.EntityTypeThingToLocation],
				selectMappings[entities.EntityTypeThing][thingID],
				selectMappings[entities.EntityTypeThingToLocation][thingToLocationThingID],
				selectMappings[entities.EntityTypeLocation][thingID],
				selectMappings[entities.EntityTypeThingToLocation][thingToLocationLocationID]),
		},
		entities.EntityTypeLocation: { // get Location by ...
			entities.EntityTypeHistoricalLocation: fmt.Sprintf("INNER JOIN %s ON %s = %s AND %s = %s",
				tableMappings[entities.EntityTypeLocationToHistoricalLocation],
				selectMappings[entities.EntityTypeLocationToHistoricalLocation][locationToHistoricalLocationLocationID],
				selectMappings[entities.EntityTypeLocation][locationID],
				selectMappings[entities.EntityTypeLocationToHistoricalLocation][locationToHistoricalLocationHistoricalLocationID],
				selectMappings[entities.EntityTypeHistoricalLocation][historicalLocationID]),
			entities.EntityTypeThing: fmt.Sprintf("INNER JOIN %s ON %s = %s AND %s = %s",
				tableMappings[entities.EntityTypeThingToLocation],
				selectMappings[entities.EntityTypeThingToLocation][thingToLocationLocationID],
				selectMappings[entities.EntityTypeLocation][locationID],
				selectMappings[entities.EntityTypeThingToLocation][thingToLocationThingID],
				selectMappings[entities.EntityTypeThing][thingID]),
		},
		entities.EntityTypeHistoricalLocation: { // get HistoricalLocation by ...
			entities.EntityTypeLocation: fmt.Sprintf("INNER JOIN %s ON %s = %s AND %s = %s",
				tableMappings[entities.EntityTypeLocationToHistoricalLocation],
				selectMappings[entities.EntityTypeLocationToHistoricalLocation][locationToHistoricalLocationHistoricalLocationID],
				selectMappings[entities.EntityTypeHistoricalLocation][historicalLocationID],
				selectMappings[entities.EntityTypeLocationToHistoricalLocation][locationToHistoricalLocationLocationID],
				selectMappings[entities.EntityTypeLocation][locationID]),
			entities.EntityTypeThing: fmt.Sprintf("WHERE %s = %s", selectMappings[entities.EntityTypeHistoricalLocation][historicalLocationThingID], selectMappings[entities.EntityTypeThing][thingID]),
		},
		entities.EntityTypeSensor: { // get sensor by ...
			entities.EntityTypeDatastream: fmt.Sprintf("WHERE %s = %s", selectMappings[entities.EntityTypeSensor][sensorID], selectMappings[entities.EntityTypeDatastream][datastreamSensorID]),
		},
		entities.EntityTypeObservedProperty: { // get observed property by ...
			entities.EntityTypeDatastream: fmt.Sprintf("WHERE %s = %s", selectMappings[entities.EntityTypeObservedProperty][observedPropertyID], selectMappings[entities.EntityTypeDatastream][datastreamObservedPropertyID]),
		},
		entities.EntityTypeObservation: { // get observation by ...
			entities.EntityTypeDatastream:        fmt.Sprintf("WHERE %s = %s", selectMappings[entities.EntityTypeObservation][observationStreamID], selectMappings[entities.EntityTypeDatastream][datastreamID]),
			entities.EntityTypeFeatureOfInterest: fmt.Sprintf("WHERE %s = %s", selectMappings[entities.EntityTypeObservation][observationFeatureOfInterestID], selectMappings[entities.EntityTypeFeatureOfInterest][foiID]),
		},
		entities.EntityTypeFeatureOfInterest: { // get feature of interest by ...
			entities.EntityTypeObservation: fmt.Sprintf("WHERE %s = %s", selectMappings[entities.EntityTypeFeatureOfInterest][foiID], selectMappings[entities.EntityTypeObservation][observationFeatureOfInterestID]),
		},
		entities.EntityTypeDatastream: { // get Datastream by ...
			entities.EntityTypeThing:            fmt.Sprintf("WHERE %s = %s", selectMappings[entities.EntityTypeDatastream][datastreamThingID], selectMappings[entities.EntityTypeThing][thingID]),
			entities.EntityTypeSensor:           fmt.Sprintf("WHERE %s = %s", selectMappings[entities.EntityTypeDatastream][datastreamSensorID], selectMappings[entities.EntityTypeSensor][sensorID]),
			entities.EntityTypeObservedProperty: fmt.Sprintf("WHERE %s = %s", selectMappings[entities.EntityTypeDatastream][datastreamObservedPropertyID], selectMappings[entities.EntityTypeObservedProperty][observedPropertyID]),
			entities.EntityTypeObservation:      fmt.Sprintf("WHERE %s = %s", selectMappings[entities.EntityTypeDatastream][datastreamID], selectMappings[entities.EntityTypeObservation][observationStreamID]),
		},
	}

	return joinMappings
}

func createTableMappings(schema string) map[entities.EntityType]string {
	if len(schema) > 0 {
		schema = fmt.Sprintf("%s.", schema)
	}

	tables := map[entities.EntityType]string{
		entities.EntityTypeThing:                        fmt.Sprintf("%s%s", schema, thingTable),
		entities.EntityTypeLocation:                     fmt.Sprintf("%s%s", schema, locationTable),
		entities.EntityTypeHistoricalLocation:           fmt.Sprintf("%s%s", schema, historicalLocationTable),
		entities.EntityTypeSensor:                       fmt.Sprintf("%s%s", schema, sensorTable),
		entities.EntityTypeObservedProperty:             fmt.Sprintf("%s%s", schema, observedPropertyTable),
		entities.EntityTypeDatastream:                   fmt.Sprintf("%s%s", schema, datastreamTable),
		entities.EntityTypeObservation:                  fmt.Sprintf("%s%s", schema, observationTable),
		entities.EntityTypeFeatureOfInterest:            fmt.Sprintf("%s%s", schema, featureOfInterestTable),
		entities.EntityTypeThingToLocation:              fmt.Sprintf("%s%s", schema, thingToLocationTable),
		entities.EntityTypeLocationToHistoricalLocation: fmt.Sprintf("%s%s", schema, locationToHistoricalLocationTable),
	}

	return tables
}
