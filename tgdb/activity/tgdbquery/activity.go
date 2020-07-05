/*
 * Copyright © 2020. TIBCO Software Inc.
 * This file is subject to the license terms contained
 * in the license file that is distributed with this file.
 */
package tgdbquery

import (
	"errors"
	"reflect"
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
	QueryType_Match     = "match"

	input_QueryParams    = "params"
	input_QueryType      = "queryType"
	output_Data          = "queryResult"
	output_DataContent   = "content"
	output_DataSuccess   = "success"
	output_DataError     = "error"
	output_DataErrorCode = "code"
	output_DataErrorMsg  = "message"

	Error_LoadDBService = 0
	Error_ConnectServer = 1
	Error_FindQueryType = 2
	Error_ExecuteQuery  = 3
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
		sendOutput(context, a.buildQueryResult(nil, false, Error_LoadDBService, err.Error()))
		return true, err
	}

	a.mux.Lock()
	defer a.mux.Unlock()

	queryType := context.GetInput(input_QueryType).(string)

	queryResult := make(map[string]interface{})
	metadata, err := tgdbService.GetMetadata()
	if nil != err {
		log.Error(err.Error())
		sendOutput(context, a.buildQueryResult(nil, false, Error_ConnectServer, err.Error()))
		return true, err
	}

	switch queryType {
	case QueryType_Metadata:
		queryResult[output_DataContent] = tgdb.BuildMetadata(metadata)
	case QueryType_NodeTypes:
		queryResult[output_DataContent] = tgdb.BuildMetadata(metadata)["nodeTypes"]
	case QueryType_EdgeTypes:
		queryResult[output_DataContent] = tgdb.BuildMetadata(metadata)["edgeTypes"]
	case QueryType_Search, QueryType_Match:
		query := context.GetInput(input_QueryParams).(*data.ComplexObject).Value.(map[string]interface{})
		_, language, parameters := a.buildQueryParams(query)
		switch language {
		case QueryLanguage_Gremlin:
			queryResult, err = a.handleGremlin(tgdbService, parameters)
		default:
			queryResult, err = a.handleTGQL(tgdbService, parameters)
		}

		if nil != err {
			log.Error(err.Error())
			queryResult = a.buildQueryResult(nil, false, Error_ExecuteQuery, err.Error())
		}
	default:
		queryResult = a.buildQueryResult(nil, false, Error_FindQueryType, errors.New("Query type not found! "))
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

	var language string
	if nil != parameters[tgdb.Query_Language] {
		language = parameters[tgdb.Query_Language].(string)
	}

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

func (a *TGDBQueryActivity) handleGremlin(
	tgdbService *tgdb.TGDBService,
	parameters map[string]interface{}) (map[string]interface{}, error) {
	resultSet, tgErr := tgdbService.GremlinQuery(parameters)
	if nil == tgErr {
		content := make(map[string]interface{})
		if nil != resultSet {
			content["nodes"] = make([]interface{}, 0)
			content["result"] = make([]interface{}, 0)
			for resultSet.HasNext() {
				result := resultSet.Next()
				log.Info("######################## handleGremlin #########################")
				log.Info(reflect.TypeOf(result).String())
				log.Info("################################################################")
				switch reflect.TypeOf(result).String() {
				case "*model.Node":
					node := result.(types.TGNode)
					content["nodes"] = append(content["nodes"].([]interface{}), tgdb.BuildNode(tgdbService, node))
				case "[]interface {}":
					nodes := make([]interface{}, 0)
					for _, value := range result.([]interface{}) {
						log.Info("+++++++++++++++++++++++++++++++++++++++++++++++++")
						log.Info(reflect.TypeOf(value).String())
						log.Info("+++++++++++++++++++++++++++++++++++++++++++++++++")
						node := value.(types.TGNode)
						nodes = append(nodes, tgdb.BuildNode(tgdbService, node))
					}
					content["nodes"] = append(content["nodes"].([]interface{}), nodes)
				default:
					content["result"] = append(content["result"].([]interface{}), result)
				}
			}
			if 0 == len(content["nodes"].([]interface{})) {
				delete(content, "nodes")
			}
		}
		return a.buildQueryResult(content, true, nil, nil), nil
	}
	return nil, tgErr
}

func (a *TGDBQueryActivity) handleTGQL(
	tgdbService *tgdb.TGDBService,
	parameters map[string]interface{}) (map[string]interface{}, error) {

	resultSet, tgErr := tgdbService.TGQLQuery(parameters)
	if nil == tgErr {
		if nil != resultSet {
			tgResult := make(map[string]interface{})
			nodes := make(map[int64]types.TGEntity)
			tgResult["nodes"] = nodes
			edges := make(map[int64]types.TGEntity)
			tgResult["edges"] = edges
			for resultSet.HasNext() {
				result := resultSet.Next()
				log.Info("######################## handleTGQL #########################")
				log.Info(reflect.TypeOf(result).String())
				log.Info("#############################################################")
				entity, ok := result.(types.TGEntity)
				if ok {
					switch entity.GetEntityKind() {
					case types.EntityKindEdge:
						edges[entity.GetVirtualId()] = entity
					case types.EntityKindNode:
						nodes[entity.GetVirtualId()] = entity
					}
				} else {
					return nil, errors.New("Unexpected result type : " + reflect.TypeOf(result).String())
				}
			}
			return a.buildQueryResult(tgdb.BuildResult(tgdbService, tgResult), true, nil, nil), nil
		} else {
			return a.buildQueryResult(tgdb.BuildResult(tgdbService, nil), true, nil, nil), nil
		}
	}
	return nil, tgErr
}

func (a *TGDBQueryActivity) buildQueryResult(
	content interface{},
	success bool,
	errorCode interface{},
	errorMsg interface{}) map[string]interface{} {

	log.Info("%%%%%%%%%%%%%%%%%%%%%% queryResult %%%%%%%%%%%%%%%%%%%%%%")
	log.Info("content   : ", content)
	log.Info("success   : ", success)
	log.Info("errorCode : ", errorCode)
	log.Info("errorMsg  : ", errorMsg)
	log.Info("%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%")

	queryResult := make(map[string]interface{})

	if success {
		queryResult[output_DataContent] = content
		queryResult[output_DataSuccess] = true
	} else {
		error := make(map[string]interface{})
		error[output_DataErrorCode] = errorCode
		error[output_DataErrorMsg] = errorMsg
		queryResult[output_DataError] = error
		queryResult[output_DataSuccess] = false
	}
	return queryResult
}
