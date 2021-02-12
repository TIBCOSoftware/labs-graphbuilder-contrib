/*
 * Copyright © 2020. TIBCO Software Inc.
 * This file is subject to the license terms contained
 * in the license file that is distributed with this file.
 */
package dgraphquery

import (
	b64 "encoding/base64"
	"fmt"
	"strings"
	"sync"

	"github.com/TIBCOSoftware/flogo-lib/core/activity"
	"github.com/TIBCOSoftware/flogo-lib/core/data"
	"github.com/TIBCOSoftware/flogo-lib/logger"
	"github.com/TIBCOSoftware/labs-graphbuilder-lib/dbservice/dgraph"
	"github.com/TIBCOSoftware/labs-graphbuilder-lib/util"
)

const (
	Setting_Connection = "dgraphConnection"
	Setting_typeTag    = "typeTag"

	Setting_QueryServiceType = "queryServiceType"
	input_QueryParams        = "queryParams"
	input_QueryString        = "queryString"
	input_PathParams         = "pathParams"
	input_QueryType          = "queryType"
	input_EntityType         = "entityType"
	output_Data              = "queryResult"
	QueryType_Search         = "search"
)

var log = logger.GetLogger("tibco-activity-tgdbquery")

type DgraphQueryActivity struct {
	metadata            *activity.Metadata
	activityToConnector map[string]string
	mux                 sync.Mutex
}

func NewActivity(metadata *activity.Metadata) activity.Activity {
	return &DgraphQueryActivity{
		metadata:            metadata,
		activityToConnector: make(map[string]string),
	}
}

func (a *DgraphQueryActivity) Metadata() *activity.Metadata {
	return a.metadata
}

func (a *DgraphQueryActivity) Eval(context activity.Context) (done bool, err error) {

	dgraphService, err := a.getDgraphService(context)

	if nil != err {
		return false, err
	}

	a.mux.Lock()
	defer a.mux.Unlock()

	pathParams := context.GetInput(input_PathParams).(*data.ComplexObject).Value.(map[string]interface{})
	queryType := pathParams[input_QueryType]
	log.Info("query type = ", queryType)

	queryResult := make(map[string]interface{})

	switch queryType {
	case QueryType_Search:
		log.Info("context.GetInput(input_QueryParams) = ", context.GetInput(input_QueryParams))
		query := context.GetInput(input_QueryParams).(*data.ComplexObject).Value.(map[string]interface{})
		log.Info("query = ", query)
		log.Info("dgraphService = ", dgraphService)
		result, error := dgraphService.Query(query[input_QueryString].(string))
		if nil != error {
			queryResult["data"] = "{}"
		} else {
			queryResult["data"] = result
		}
		break
	default:

	}

	queryResult["success"] = true

	log.Info("query result = ", queryResult)

	complexdata := &data.ComplexObject{Metadata: output_Data, Value: queryResult}
	log.Info("complexdata = ", complexdata)
	context.SetOutput(output_Data, complexdata)

	return true, nil
}

func (a *DgraphQueryActivity) getDgraphService(context activity.Context) (*dgraph.DgraphService, error) {
	myId := util.ActivityId(context)

	log.Debug("(getDgraphService) entering - myId = ", myId)

	dgraphService := dgraph.GetFactory().GetService(a.activityToConnector[myId])
	if nil == dgraphService {
		a.mux.Lock()
		defer a.mux.Unlock()

		dgraphService = dgraph.GetFactory().GetService(a.activityToConnector[myId])
		if nil == dgraphService {
			var err error
			log.Info("(getDgraphService) Initializing DGraph Service start ...")

			connection, exist := context.GetSetting(Setting_Connection)
			if !exist {
				return nil, activity.NewError("Dgraph connection is not configured", "Degraph-UPSERT-4001", nil)
			}

			connectionInfo, _ := data.CoerceToObject(connection)
			if connectionInfo == nil {
				return nil, activity.NewError("Dgraph connection not able to be parsed", "Degraph-UPSERT-4002", nil)
			}

			properties := make(map[string]interface{})
			connectionSettings, _ := connectionInfo["settings"].([]interface{})

			if nil == connectionSettings {
				return nil, fmt.Errorf("Unable to get connection setting!")
			}

			var connectorName string
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
					} else if setting["name"] == "schema" {
						if nil != setting["value"] {
							content, err := data.CoerceToObject(setting["value"])
							if nil != err {
								break
							}
							schemaBytes, err := b64.StdEncoding.DecodeString(strings.Split(content["content"].(string), ",")[1])
							if nil != err {
								break
							}
							properties["schema"] = strings.Fields(string(schemaBytes))
						}
					}
				}
			}

			//typeName, exist := context.GetSetting(Setting_typeTag)
			//if !exist {
			//	return nil, activity.NewError("Type tag is not configured", "Degraph-UPSERT-4003", nil)
			//}

			//properties["typeName"] = typeName

			//addPrefixToAttr, _ := context.GetSetting(attrWithPrefix)

			//properties["addPrefixToAttr"] = addPrefixToAttr

			//if nil != err {
			//	return nil, err
			//}
			//properties["graphModel"] = graph.GetModel()

			log.Debug("properties : ", properties)

			dgraphService, err = dgraph.GetFactory().CreateService(connectorName, properties)

			if nil != err {
				return nil, err
			}
			a.activityToConnector[myId] = connectorName

			log.Info("(getDgraphService) Initializing Dgraph Service end ...")
		}
	}

	log.Info("(getDgraphService) return - myId = ", myId, ", dgraphService = ", dgraphService)

	return dgraphService, nil
}
