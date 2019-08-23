/*
 * Copyright © 2019. TIBCO Software Inc.
 * This file is subject to the license terms contained
 * in the license file that is distributed with this file.
 */
package neo4jupsert

import (
	//	b64 "encoding/base64"
	"fmt"
	//	"strings"
	"sync"

	"github.com/TIBCOSoftware/flogo-lib/core/activity"
	"github.com/TIBCOSoftware/flogo-lib/core/data"
	"github.com/TIBCOSoftware/flogo-lib/logger"
	"github.com/TIBCOSoftware/labs-graphbuilder-lib/dbservice/neo4j"
	"github.com/TIBCOSoftware/labs-graphbuilder-lib/model"
	"github.com/TIBCOSoftware/labs-graphbuilder-lib/util"
)

const (
	Connection     = "neo4jConnection"
	attrWithPrefix = "attrWithPrefix"
)

var log = logger.GetLogger("tibco-activity-neo4jupsert")

type Neo4jUpsertActivity struct {
	metadata            *activity.Metadata
	neo4jServices       map[string]*neo4j.Neo4jService
	activityToConnector map[string]string
	mux                 sync.Mutex
}

func NewActivity(metadata *activity.Metadata) activity.Activity {
	return &Neo4jUpsertActivity{
		metadata:            metadata,
		neo4jServices:       make(map[string]*neo4j.Neo4jService),
		activityToConnector: make(map[string]string),
	}
}

func (a *Neo4jUpsertActivity) Metadata() *activity.Metadata {
	return a.metadata
}

func (a *Neo4jUpsertActivity) Eval(context activity.Context) (done bool, err error) {

	neo4jService, err := a.getNeo4jService(context)

	if nil != err {
		return false, err
	}

	a.mux.Lock()
	defer a.mux.Unlock()

	var graph model.Graph
	graph, err = GetGraph(context)

	if nil != err {
		return false, err
	}

	err = neo4jService.UpsertGraph(graph)

	if nil != err {
		return false, err
	}

	return true, nil
}

func (a *Neo4jUpsertActivity) getNeo4jService(context activity.Context) (*neo4j.Neo4jService, error) {
	myId := util.ActivityId(context)

	log.Info("(getNeo4jService) entering - myId = ", myId)

	neo4jService := a.neo4jServices[a.activityToConnector[myId]]
	if nil == neo4jService {
		a.mux.Lock()
		defer a.mux.Unlock()

		neo4jService = a.neo4jServices[a.activityToConnector[myId]]
		if nil == neo4jService {

			log.Info("(getNeo4jService) Initializing Neo4j Service start ...")

			connection, exist := context.GetSetting(Connection)
			if !exist {
				return nil, activity.NewError("Neo4j connection is not configured", "Neo4j-UPSERT-4001", nil)
			}

			connectionInfo, _ := data.CoerceToObject(connection)
			if connectionInfo == nil {
				return nil, activity.NewError("Neo4j connection not able to be parsed", "Neo4j-UPSERT-4002", nil)
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

			neo4jService, err = neo4j.NewNeo4jServiceFactory().GetService(properties)

			if nil != err {
				return nil, err
			}
			a.neo4jServices[connectorName] = neo4jService
			a.activityToConnector[myId] = connectorName

			log.Info("(getNeo4jService) Initializing Dgraph Service end ...")
		}
	}

	log.Info("(getNeo4jService) return - myId = ", myId, ", dgraphService = ", neo4jService)

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
