/*
 * Copyright © 2019. TIBCO Software Inc.
 * This file is subject to the license terms contained
 * in the license file that is distributed with this file.
 */
package tgdbquery

import (
	"errors"
	"sync"

	"github.com/TIBCOSoftware/flogo-lib/core/activity"
	"github.com/TIBCOSoftware/flogo-lib/core/data"
	"github.com/TIBCOSoftware/flogo-lib/logger"
	"github.com/TIBCOSoftware/labs-graphbuilder-lib/dbservice"
	"github.com/TIBCOSoftware/labs-graphbuilder-lib/dbservice/factory"
	"github.com/TIBCOSoftware/labs-graphbuilder-lib/dbservice/tgdb"
	"github.com/TIBCOSoftware/labs-graphbuilder-lib/util"
	"github.com/TIBCOSoftware/tgdb-client/client/goAPI/types"
)

const (
	Setting_Connection = "tgdbConnection"

	QueryLanguage         = "language"
	QueryLanguage_Gremlin = "gremlin"
	QueryLanguage_TGQL    = "tgql"

	QueryType_Metadata  = "metadata"
	QueryType_NodeTypes = "nodetypes"
	QueryType_EdgeTypes = "edgetypes"
	QueryType_Node      = "node"
	QueryType_Search    = "search"

	input_QueryParams = "params"
	input_QueryType   = "queryType"
	output_Data       = "queryResult"

	ErrorLoadDBService = 0
	ErrorConnectServer = 1
	ErrorFindQueryType = 2
	ErrorExecuteQuery  = 3
)

var log = logger.GetLogger("tibco-activity-tgdbquery")

type TGDBQueryActivity struct {
	metadata            *activity.Metadata
	activityToConnector map[string]string
	mux                 sync.Mutex
}

func NewActivity(metadata *activity.Metadata) activity.Activity {
	return &TGDBQueryActivity{
		metadata:            metadata,
		activityToConnector: make(map[string]string),
	}
}

func (a *TGDBQueryActivity) Metadata() *activity.Metadata {
	return a.metadata
}

func (a *TGDBQueryActivity) Eval(context activity.Context) (done bool, err error) {

	tgdbService, err := a.getTGDBService(context)

	if nil != err {
		log.Error(err.Error())
		sendOutput(context, a.buildQueryResult(nil, false, ErrorLoadDBService, err.Error()))
		return true, err
	}

	a.mux.Lock()
	defer a.mux.Unlock()

	queryType := context.GetInput(input_QueryType).(string)
	//log.Info("query type => ", queryType)

	queryResult := make(map[string]interface{})
	metadata, err := tgdbService.GetMetadata()
	if nil != err {
		log.Error(err.Error())
		sendOutput(context, a.buildQueryResult(nil, false, ErrorConnectServer, err.Error()))
		return true, err
	}

	switch queryType {
	case QueryType_Metadata:
		queryResult["data"] = tgdb.BuildMetadata(metadata)
		break
	case QueryType_NodeTypes:
		queryResult["data"] = tgdb.BuildMetadata(metadata)["nodeTypes"]
		break
	case QueryType_EdgeTypes:
		queryResult["data"] = tgdb.BuildMetadata(metadata)["edgeTypes"]
		break
	case QueryType_Search:
		query := context.GetInput(input_QueryParams).(*data.ComplexObject).Value.(map[string]interface{})
		_, language, parameters := a.buildQueryParams(query)
		var resultSet types.TGResultSet
		var tgErr types.TGError
		switch language {
		case QueryLanguage_Gremlin:
			{
				resultSet, tgErr = tgdbService.GremlinQuery(parameters)
			}
		case QueryLanguage_TGQL:
			{
				resultSet, tgErr = tgdbService.TGQLQuery(parameters)
			}
		default:
			{
				resultSet, tgErr = tgdbService.TGQLQuery(parameters)
			}
		}

		if nil == tgErr {
			if nil != resultSet {
				tgResult := make(map[string]map[int64]types.TGEntity)
				tgResult["nodes"] = make(map[int64]types.TGEntity)
				tgResult["edges"] = make(map[int64]types.TGEntity)
				for resultSet.HasNext() {
					entity := resultSet.Next().(types.TGEntity)
					switch entity.GetEntityKind() {
					case types.EntityKindEdge:
						tgResult["edges"][entity.GetVirtualId()] = entity
					case types.EntityKindNode:
						tgResult["nodes"][entity.GetVirtualId()] = entity
					}
				}

				queryResult = a.buildQueryResult(tgdb.BuildResult(tgdbService, tgResult), true, nil, nil)
			}
		} else {
			log.Error(tgErr.Error())
			queryResult = a.buildQueryResult(nil, false, ErrorExecuteQuery, tgErr.Error())
		}
		break
	default:
		queryResult = a.buildQueryResult(nil, false, ErrorFindQueryType, errors.New("Query type not found! "))
	}

	sendOutput(context, queryResult)

	return true, nil
}

func sendOutput(
	context activity.Context,
	queryResult interface{}) {

	complexdata := &data.ComplexObject{
		Metadata: output_Data,
		Value:    queryResult,
	}
	context.SetOutput(output_Data, complexdata)
}

func (a *TGDBQueryActivity) getTGDBService(context activity.Context) (*tgdb.TGDBService, error) {
	myId := util.ActivityId(context)

	tgdbService := factory.GetFactory(dbservice.TGDB).GetUpsertService(a.activityToConnector[myId])
	if nil == tgdbService {
		a.mux.Lock()
		defer a.mux.Unlock()
		tgdbService = factory.GetFactory(dbservice.TGDB).GetUpsertService(a.activityToConnector[myId])
		if nil == tgdbService {
			log.Info("Initializing TGDB Service start ...")
			connection, exist := context.GetSetting(Setting_Connection)
			if !exist {
				return nil, activity.NewError("TGDB connection is not configured", "TGDB-QUERY-4001", nil)
			}

			connectionInfo, _ := data.CoerceToObject(connection)
			if connectionInfo == nil {
				return nil, activity.NewError("TGDB connection not able to be parsed", "TGDB-QUERY-4002", nil)
			}

			var connectorName string
			properties := make(map[string]interface{})
			connectionSettings, _ := connectionInfo["settings"].([]interface{})
			if connectionSettings != nil {
				for _, v := range connectionSettings {
					setting, _ := data.CoerceToObject(v)
					if setting != nil {
						if setting["name"] == "url" {
							url, _ := data.CoerceToString(setting["value"])
							properties["url"] = util.GetValue(url)
						} else if setting["name"] == "user" {
							properties["user"], _ = data.CoerceToString(setting["value"])
						} else if setting["name"] == "password" {
							properties["password"], _ = data.CoerceToString(setting["value"])
						} else if setting["name"] == "name" {
							connectorName, _ = data.CoerceToString(setting["value"])
						}
					}
				}
				log.Info(properties)

				tgdbService, _ = factory.GetFactory(dbservice.TGDB).CreateUpsertService(connectorName, properties)
				a.activityToConnector[myId] = connectorName
			}
			log.Info("Initializing TGDB Service end ...")
		}
	}

	return tgdbService.(*tgdb.TGDBService), nil
}

func (a *TGDBQueryActivity) buildQueryParams(parameters map[string]interface{}) (error, string, map[string]interface{}) {
	queryParams := make(map[string]interface{})
	query := make(map[string]interface{})
	queryParams[tgdb.Query] = query
	if nil != parameters[tgdb.Query_QueryString] {
		query[tgdb.Query_QueryString] = parameters[tgdb.Query_QueryString]
	}

	if nil != parameters[tgdb.Query_EdgeFilter] {
		query[tgdb.Query_EdgeFilter] = parameters[tgdb.Query_EdgeFilter]
	}

	if nil != parameters[tgdb.Query_TraversalCondition] {
		query[tgdb.Query_TraversalCondition] = parameters[tgdb.Query_TraversalCondition]
	}

	if nil != parameters[tgdb.Query_EndCondition] {
		query[tgdb.Query_EndCondition] = parameters[tgdb.Query_EndCondition]
	}

	if nil != parameters[tgdb.Query_OPT_PrefetchSize] {
		queryParams[tgdb.Query_OPT_PrefetchSize] = int(parameters[tgdb.Query_OPT_PrefetchSize].(float64))
	}

	if nil != parameters[tgdb.Query_OPT_TraversalDepth] {
		queryParams[tgdb.Query_OPT_TraversalDepth] = int(parameters[tgdb.Query_OPT_TraversalDepth].(float64))
	}

	if nil != parameters[tgdb.Query_OPT_EdgeLimit] {
		queryParams[tgdb.Query_OPT_EdgeLimit] = int(parameters[tgdb.Query_OPT_EdgeLimit].(float64))
	}

	var err error
	language, ok := parameters[QueryLanguage].(string)
	if !ok {
		err = errors.New("Language not defined")
	}
	return err, language, queryParams
}

func (a *TGDBQueryActivity) buildQueryResult(
	data interface{},
	success bool,
	errorCode interface{},
	errorMsg interface{}) map[string]interface{} {

	log.Debug("%%%%%%%%%%%%%%%%%%%%%% queryResult %%%%%%%%%%%%%%%%%%%%%%")
	log.Debug("data      : ", data)
	log.Debug("success   : ", success)
	log.Debug("errorCode : ", errorCode)
	log.Debug("errorMsg  : ", errorMsg)
	log.Debug("%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%")

	queryResult := make(map[string]interface{})

	if success {
		queryResult["data"] = data
		queryResult["success"] = true
	} else {
		error := make(map[string]interface{})
		error["code"] = errorCode
		error["message"] = errorMsg
		queryResult["error"] = error
		queryResult["success"] = false
	}
	return queryResult
}
