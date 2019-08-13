/*
 * Copyright © 2019. TIBCO Software Inc.
 * This file is subject to the license terms contained
 * in the license file that is distributed with this file.
 */
package cloudconversation

import (
	"fmt"
	"sync"

	"github.com/TIBCOSoftware/flogo-lib/core/activity"
	"github.com/TIBCOSoftware/flogo-lib/core/data"
	"github.com/TIBCOSoftware/flogo-lib/logger"
	"github.com/TIBCOSoftware/graph-builder-flogo/lib/cloudconversation"
	"github.com/TIBCOSoftware/graph-builder-flogo/lib/util"
	//"github.com/TIBCOSoftware/tgdb-client/client/goAPI/types"
)

const (
	Setting_DatabaseType = "databaseType"
	Setting_Connection   = "databaseConnection"
	Setting_typeTag      = "typeTag"

	Setting_QueryServiceType = "queryServiceType"
	input_QueryParams        = "queryParams"
	input_QueryString        = "queryString"
	input_PathParams         = "pathParams"
	input_QueryType          = "queryType"
	input_EntityType         = "entityType"
	output_Data              = "queryResult"

	QueryType_FindEntityById    = "findEntityById"
	input_QueryParam_uuid       = "uuid"
	input_QueryParam_entityType = "entityType"
	input_QueryParam_info       = "info"
)

var log = logger.GetLogger("tibco-activity-cloudconversation")

type CloudConversationActivity struct {
	metadata *activity.Metadata
	mux      sync.Mutex
}

func NewActivity(metadata *activity.Metadata) activity.Activity {
	return &CloudConversationActivity{
		metadata: metadata,
	}
}

func (a *CloudConversationActivity) Metadata() *activity.Metadata {
	return a.metadata
}

func (a *CloudConversationActivity) Eval(context activity.Context) (done bool, err error) {

	queryService, err := a.getQueryService(context)

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
	case QueryType_FindEntityById:
		log.Info("context.GetInput(input_QueryParams) = ", context.GetInput(input_QueryParams))
		query := context.GetInput(input_QueryParams).(*data.ComplexObject).Value.(map[string]interface{})
		log.Info("query = ", query)
		log.Info("queryService = ", queryService)
		uuid := query[input_QueryParam_uuid].(string)
		entityType := query[input_QueryParam_entityType].(string)
		info := query[input_QueryParam_info].(string)

		result, error := queryService.FindEntityById(
			uuid, entityType, info,
		)

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

func (a *CloudConversationActivity) getQueryService(context activity.Context) (cloudconversation.ConvUIQuery, error) {
	myId := util.ActivityId(context)

	log.Debug("(getQueryService) entering - myId = ", myId)

	queryService := cloudconversation.GetQueryManager().GetQuery(myId)
	if nil == queryService {
		a.mux.Lock()
		defer a.mux.Unlock()

		queryService = cloudconversation.GetQueryManager().GetQuery(myId)
		if nil == queryService {
			log.Info("(getQueryService) Initializing DGraph Service start ...")
			var err error
			databaseType, _ := context.GetSetting(Setting_DatabaseType)

			connection, exist := context.GetSetting(Setting_Connection)
			if !exist {
				return nil, activity.NewError("DB connection is not configured", "Cloud-Conversation-4001", nil)
			}

			connectionInfo, _ := data.CoerceToObject(connection)
			if connectionInfo == nil {
				return nil, activity.NewError("DB connection not able to be parsed", "Cloud-Conversation-4002", nil)
			}

			connectionSettings, _ := connectionInfo["settings"].([]interface{})

			if nil == connectionSettings {
				return nil, fmt.Errorf("Unable to get connection setting!")
			}

			queryService, err = cloudconversation.GetQueryManager().CreateQuery(myId, databaseType.(string), connectionSettings)
			if nil != err {
				return nil, err
			}
			log.Info("(getQueryService) Initializing Dgraph Service end ...")
		}
	}

	log.Info("(getQueryService) return - myId = ", myId, ", queryService = ", queryService)

	return queryService, nil
}
