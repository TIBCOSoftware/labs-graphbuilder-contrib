/*
 * Copyright © 2019. TIBCO Software Inc.
 * This file is subject to the license terms contained
 * in the license file that is distributed with this file.
 */
package tgdbquery

import (
	"sync"

	"github.com/TIBCOSoftware/flogo-lib/core/activity"
	"github.com/TIBCOSoftware/flogo-lib/core/data"
	"github.com/TIBCOSoftware/flogo-lib/logger"
	"github.com/TIBCOSoftware/labs-graphbuilder-lib/dbservice/tgdb"
	"github.com/TIBCOSoftware/labs-graphbuilder-lib/util"
	"github.com/TIBCOSoftware/tgdb-client/client/goAPI/types"
)

const (
	Setting_Connection       = "tgdbConnection"
	Setting_QueryServiceType = "queryServiceType"
	input_QueryParams        = "queryParams"
	input_Get_KeyAttrNames   = "keyAttrNames"
	input_Get_KeyAttrValues  = "keyAttrValues"
	input_PathParams         = "pathParams"
	input_QueryType          = "queryType"
	input_EntityType         = "entityType"
	output_Data              = "queryResult"
	QueryType_Metadata       = "metadata"
	QueryType_NodeTypes      = "nodetypes"
	QueryType_EdgeTypes      = "edgetypes"
	QueryType_Node           = "node"
	QueryType_Search         = "search"
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
		sendOutput(context, nil, false, 0, err.Error())
		return true, err
	}

	a.mux.Lock()
	defer a.mux.Unlock()

	pathParams := context.GetInput(input_PathParams).(*data.ComplexObject).Value.(map[string]interface{})
	queryType := pathParams[input_QueryType]
	log.Info("query type => ", queryType)

	queryResult := make(map[string]interface{})
	metadata, err := tgdbService.GetMetadata()
	if nil != err {
		log.Error(err.Error())
		sendOutput(context, nil, false, 0, err.Error())
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
	case QueryType_Node:
		entityType := pathParams[input_EntityType].(string)
		//		log.Info("entity type = ", entityType)
		queryParams := context.GetInput(input_QueryParams).(*data.ComplexObject).Value.(map[string]interface{})
		//		log.Info("queryParams = ", queryParams)
		keyAttrValues := queryParams[input_Get_KeyAttrValues].([]interface{})
		keyAttrNames := queryParams[input_Get_KeyAttrNames].([]interface{})
		attributes := make(map[string]interface{})
		for index, keyAttrName := range keyAttrNames {
			log.Info("keyAttrName = ", keyAttrName, "keyAttrValue = ", keyAttrValues[index])
			attribute := make(map[string]interface{})
			attribute["name"] = keyAttrName.(string)
			attribute["value"] = keyAttrValues[index]
			attribute["type"] = "string"
			attributes[keyAttrName.(string)] = attribute
		}
		//		log.Info("attributes = ", attributes)
		pKey := make(map[string]interface{})
		pKey["attributes"] = attributes
		entity, _ := tgdbService.GetNode(entityType, pKey)
		if nil != entity {
			result := make(map[string]interface{})
			tgdb.BuildNode(tgdbService, entity.(types.TGNode), result)
			queryResult = a.buildQueryResult(result, true, nil, nil)
		} else {
			queryResult = a.buildQueryResult(nil, true, nil, nil)
		}
		break
	case QueryType_Search:
		query := context.GetInput(input_QueryParams).(*data.ComplexObject).Value.(map[string]interface{})
		language, parameters := a.buildQueryParams(query)
		var resultSet types.TGResultSet
		var tgErr types.TGError
		switch language {
		case QueryLanguage_Gremlin:
			{
				resultSet, tgErr = tgdbService.GremlinQuery(parameters)
			}
		default:
			{
				resultSet, tgErr = tgdbService.TGQLQuery(parameters)
			}
		}

		if nil == tgErr {
			if nil != resultSet {
				result := make(map[string]interface{})
				for resultSet.HasNext() {
					entity := resultSet.Next()
					tgdb.BuildNode(tgdbService, entity.(types.TGNode), result)
				}
				queryResult = a.buildQueryResult(result, true, nil, nil)
			}
		}
		break
	default:

	}

	sendOutput(context, queryResult, true, 0, "")

	return true, nil
}

func sendOutput(
	context activity.Context,
	content interface{},
	success bool,
	errorCode int,
	errorMessage string) {
	queryResult := make(map[string]interface{})
	queryResult["success"] = true
	queryResult["error"] = make(map[string]interface{})
	queryResult["data"] = content
	queryResult["success"] = success
	error := make(map[string]interface{})
	error["code"] = errorCode
	error["message"] = errorMessage
	queryResult["error"] = error
	complexdata := &data.ComplexObject{
		Metadata: output_Data,
		Value:    queryResult,
	}
	context.SetOutput(output_Data, complexdata)
}

func (a *TGDBQueryActivity) getTGDBService(context activity.Context) (*tgdb.TGDBService, error) {
	myId := util.ActivityId(context)

	tgdbService := tgdb.GetFactory().GetService(a.activityToConnector[myId])
	if nil == tgdbService {
		a.mux.Lock()
		defer a.mux.Unlock()
		tgdbService = tgdb.GetFactory().GetService(a.activityToConnector[myId])
		if nil == tgdbService {
			log.Info("Initializing TGDB Service start ...")
			connection, exist := context.GetSetting(Setting_Connection)
			if !exist {
				return nil, activity.NewError("TGDB connection is not configured", "TGDB-UPSERT-4001", nil)
			}

			connectionInfo, _ := data.CoerceToObject(connection)
			if connectionInfo == nil {
				return nil, activity.NewError("TGDB connection not able to be parsed", "TGDB-UPSERT-4002", nil)
			}

			//			queryType, exist := context.GetSetting(Setting_QueryServiceType)
			//			a.queryType = queryType.(string)

			var connectorName string
			properties := make(map[string]interface{})
			connectionSettings, _ := connectionInfo["settings"].([]interface{})
			if connectionSettings != nil {
				for _, v := range connectionSettings {
					setting, _ := data.CoerceToObject(v)
					if setting != nil {
						if setting["name"] == "url" {
							properties["url"], _ = data.CoerceToString(setting["value"])
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

				tgdbService, _ = tgdb.GetFactory().CreateService(connectorName, properties)
				a.activityToConnector[myId] = connectorName
			}
			log.Info("Initializing TGDB Service end ...")
		}
	}

	return tgdbService, nil
}

func (a *TGDBQueryActivity) buildQueryParams(parameters map[string]interface{}) map[string]interface{} {
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

	return language, queryParams
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
