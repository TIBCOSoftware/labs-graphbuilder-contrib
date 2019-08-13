/*
 * Copyright © 2019. TIBCO Software Inc.
 * This file is subject to the license terms contained
 * in the license file that is distributed with this file.
 */
package dgraphupsert

import (
	b64 "encoding/base64"
	"fmt"
	"strings"
	"sync"

	"github.com/TIBCOSoftware/flogo-lib/core/activity"
	"github.com/TIBCOSoftware/flogo-lib/core/data"
	"github.com/TIBCOSoftware/flogo-lib/logger"
	"github.com/TIBCOSoftware/graph-builder-flogo/lib/dbservice/dgraph"
	"github.com/TIBCOSoftware/graph-builder-flogo/lib/model"
	"github.com/TIBCOSoftware/graph-builder-flogo/lib/util"
)

const (
	Connection         = "dgraphConnection"
	typeTag            = "typeTag"
	explicitType       = "explicitType"
	readableExternalId = "readableExternalId"
	attrWithPrefix     = "attrWithPrefix"
)

var log = logger.GetLogger("tibco-activity-dgraphupsert")

type DgraphUpsertActivity struct {
	metadata            *activity.Metadata
	activityToConnector map[string]string
	mux                 sync.Mutex
}

func NewActivity(metadata *activity.Metadata) activity.Activity {
	return &DgraphUpsertActivity{
		metadata:            metadata,
		activityToConnector: make(map[string]string),
	}
}

func (a *DgraphUpsertActivity) Metadata() *activity.Metadata {
	return a.metadata
}

func (a *DgraphUpsertActivity) Eval(context activity.Context) (done bool, err error) {

	log.Debug("(getDgraphService) entering ......")

	dgraphService, err := a.getDgraphService(context)

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

	err = dgraphService.UpsertGraph(graph)

	if nil != err {
		return false, err
	}

	log.Debug("(getDgraphService) exit ......")

	return true, nil
}

func (a *DgraphUpsertActivity) getDgraphService(context activity.Context) (*dgraph.DgraphService, error) {
	myId := util.ActivityId(context)

	log.Debug("(getDgraphService) entering - myId = ", myId)

	dgraphService := dgraph.GetFactory().GetService(a.activityToConnector[myId])
	if nil == dgraphService {
		a.mux.Lock()
		defer a.mux.Unlock()

		dgraphService = dgraph.GetFactory().GetService(a.activityToConnector[myId])
		if nil == dgraphService {

			log.Info("(getDgraphService) Initializing DGraph Service start ...")

			connection, exist := context.GetSetting(Connection)
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
					} else if setting["name"] == "tls" {
						if nil != setting["value"] {
							content, err := data.CoerceToObject(setting["value"])
							if nil != err {
								break
							}
							schemaBytes, err := b64.StdEncoding.DecodeString(strings.Split(content["content"].(string), ",")[1])
							if nil != err {
								break
							}
							properties["tls"] = strings.Fields(string(schemaBytes))
						}
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

			readableExternalId, exist := context.GetSetting(readableExternalId)
			if exist {
				properties["readableExternalId"] = readableExternalId
			} else {
				log.Info("readableExternalId configuration is not configured, will make readableExternalId true!")
			}

			explicitType, exist := context.GetSetting(explicitType)
			if exist {
				properties["explicitType"] = explicitType
			} else {
				log.Info("explicitType configuration is not configured, will make type defininated implicit!")
			}

			typeName, exist := context.GetSetting(typeTag)
			if exist {
				properties["typeName"] = typeName
			} else {
				log.Info("Type tag is not configured, will reate an predicate as type!")
			}

			addPrefixToAttr, _ := context.GetSetting(attrWithPrefix)

			properties["addPrefixToAttr"] = addPrefixToAttr

			graph, err := GetGraph(context)

			log.Debug(properties)

			if nil != err {
				return nil, err
			}
			properties["graphModel"] = graph.GetModel()

			log.Debug("properties : ", properties)

			dgraphService, err = dgraph.GetFactory().CreateService(connectorName, properties)

			if nil != err {
				return nil, err
			}
			a.activityToConnector[myId] = connectorName

			log.Info("(getDgraphService) Initializing Dgraph Service end ...")
		}
	}

	log.Debug("(getDgraphService) return - myId = ", myId, ", dgraphService = ", dgraphService)

	return dgraphService, nil
}

func GetGraph(context activity.Context) (model.Graph, error) {

	graphData, ok := context.GetInput("Graph").(map[string]interface{})["graph"].(map[string]interface{})
	if !ok {
		return nil, fmt.Errorf("Fail to parse graph model !!")
	}
	log.Debug("(getDgraph) json graph data = ", graphData)

	graph := model.ReconstructGraph(graphData)
	log.Debug("(getDgraph) graph obj = ", graph)

	return graph, nil
}
