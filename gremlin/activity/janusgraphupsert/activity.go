/*
 * Copyright © 2019. TIBCO Software Inc.
 * This file is subject to the license terms contained
 * in the license file that is distributed with this file.
 */
package janusgraphupsert

import (
	"fmt"
	"sync"

	"github.com/TIBCOSoftware/flogo-lib/core/activity"
	"github.com/TIBCOSoftware/flogo-lib/core/data"
	"github.com/TIBCOSoftware/flogo-lib/logger"
	"github.com/TIBCOSoftware/labs-graphbuilder-lib/dbservice"
	"github.com/TIBCOSoftware/labs-graphbuilder-lib/dbservice/factory"
	"github.com/TIBCOSoftware/labs-graphbuilder-lib/model"
	"github.com/TIBCOSoftware/labs-graphbuilder-lib/util"
)

const (
	Connection     = "Connection"
	attrWithPrefix = "attrWithPrefix"
)

var log = logger.GetLogger("tibco-activity-gremlinupsert")

type JanusgraphUpsertActivity struct {
	metadata            *activity.Metadata
	activityToConnector map[string]string
	mux                 sync.Mutex
}

func NewActivity(metadata *activity.Metadata) activity.Activity {
	return &JanusgraphUpsertActivity{
		metadata:            metadata,
		activityToConnector: make(map[string]string),
	}
}

func (a *JanusgraphUpsertActivity) Metadata() *activity.Metadata {
	return a.metadata
}

func (a *JanusgraphUpsertActivity) Eval(context activity.Context) (done bool, err error) {
	log.Info("(JanusgraphUpsertActivity.Eval) entering ........")
	defer log.Info("(JanusgraphUpsertActivity.Eval) exit ........")

	neo4jService, err := a.getGremlinService(context)

	if nil != err {
		return false, err
	}

	var graph model.Graph
	graph, err = GetGraph(context)

	if nil != err {
		return false, err
	}

	err = neo4jService.UpsertGraph(graph, nil)

	if nil != err {
		return false, err
	}

	return true, nil
}

func (a *JanusgraphUpsertActivity) getGremlinService(context activity.Context) (dbservice.UpsertService, error) {
	myId := util.ActivityId(context)

	log.Info("(getGremlinService) entering - myId = ", myId)

	neo4jService := factory.GetFactory(dbservice.Gremlin).GetUpsertService(a.activityToConnector[myId])
	if nil == neo4jService {
		a.mux.Lock()
		defer a.mux.Unlock()

		neo4jService = factory.GetFactory(dbservice.Gremlin).GetUpsertService(a.activityToConnector[myId])
		if nil == neo4jService {

			log.Info("(getGremlinService) Initializing Gremlin Service start ...")

			connection, exist := context.GetSetting(Connection)
			if !exist {
				return nil, activity.NewError("Gremlin connection is not configured", "Gremlin-UPSERT-4001", nil)
			}

			connectionInfo, _ := data.CoerceToObject(connection)
			if connectionInfo == nil {
				return nil, activity.NewError("Gremlin connection not able to be parsed", "Gremlin-UPSERT-4002", nil)
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
					} /* else if setting["name"] == "schema" {
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
					}*/
				}
			}

			addPrefixToAttr, _ := context.GetSetting(attrWithPrefix)

			properties["addPrefixToAttr"] = addPrefixToAttr

			graph, err := GetGraph(context)

			log.Info(properties)

			if nil != err {
				return nil, err
			}
			properties["graphModel"] = graph.GetModel()

			log.Info("properties : ", properties)

			neo4jService, err = factory.GetFactory(dbservice.Gremlin).CreateUpsertService(connectorName, properties)

			if nil != err {
				return nil, err
			}
			a.activityToConnector[myId] = connectorName

			log.Info("(getGremlinService) Initializing Dgraph Service end ...")
		}
	}

	log.Info("(getGremlinService) return - myId = ", myId, ", dgraphService = ", neo4jService)

	return neo4jService, nil
}

func GetGraph(context activity.Context) (model.Graph, error) {
	graphData, ok := context.GetInput("Graph").(map[string]interface{})["graph"].(map[string]interface{})
	if !ok {
		return nil, fmt.Errorf("Fail to parse graph model !!")
	}

	graph := model.ReconstructGraph(graphData)

	return graph, nil
}
