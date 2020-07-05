/*
 * Copyright © 2020. TIBCO Software Inc.
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
	"github.com/TIBCOSoftware/labs-graphbuilder-lib/dbservice"
	"github.com/TIBCOSoftware/labs-graphbuilder-lib/dbservice/factory"
	"github.com/TIBCOSoftware/labs-graphbuilder-lib/model"
	"github.com/TIBCOSoftware/labs-graphbuilder-lib/util"

	"git.tibco.com/git/product/ipaas/wi-contrib.git/connection/generic"
)

const (
	Connection         = "dgraphConnection"
	cacheSize          = "cacheSize"
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

	log.Info("(DgraphUpsertActivity) entering ......")

	dgraphService, err := a.getDgraphService(context)

	if nil != err {
		log.Error("(DgraphUpsertActivity) exit after get service, with error = ", err.Error())
		return false, err
	}

	var graph model.Graph
	graph, err = GetGraph(context)

	if nil != err {
		log.Error("(DgraphUpsertActivity) exit after get graph data, with error = ", err.Error())
		return false, err
	}

	err = dgraphService.UpsertGraph(graph, nil)

	if nil != err {
		log.Error("(DgraphUpsertActivity) exit during upsert, with error = ", err.Error())
		return false, err
	}

	log.Info("(DgraphUpsertActivity) exit normally ......")

	return true, nil
}

func (a *DgraphUpsertActivity) getDgraphService(context activity.Context) (dbservice.UpsertService, error) {
	myId := util.ActivityId(context)

	log.Debug("(getDgraphService) entering - myId = ", myId)

	dgraphService := factory.GetFactory(dbservice.Dgraph).GetUpsertService(a.activityToConnector[myId])
	//dgraph.GetFactory().GetService(a.activityToConnector[myId])
	if nil == dgraphService {
		a.mux.Lock()
		defer a.mux.Unlock()

		dgraphService = factory.GetFactory(dbservice.Dgraph).GetUpsertService(a.activityToConnector[myId])
		if nil == dgraphService {

			log.Info("(getDgraphService) Initializing DGraph Service start ...")

			settingsMap, err := getConnectionSetting(context)
			if nil != err {
				return nil, err
			}

			properties := make(map[string]interface{})

			properties["version"], _ = data.CoerceToString(settingsMap["apiVersion"])
			properties["url"], _ = data.CoerceToString(settingsMap["url"])
			properties["user"], _ = data.CoerceToString(settingsMap["user"])
			properties["password"], _ = data.CoerceToString(settingsMap["password"])
			properties["tlsEnabled"], _ = data.CoerceToBoolean(settingsMap["tlsEnabled"])
			if properties["tlsEnabled"].(bool) {
				if nil != settingsMap["tls"] {
					content, err := data.CoerceToObject(settingsMap["tls"])
					if nil == err {
						tlsBytes, err := b64.StdEncoding.DecodeString(strings.Split(content["content"].(string), ",")[1])
						if nil == err {
							properties["tls"] = string(tlsBytes)
						}
					}
				}
			}

			properties["schemaGen"], _ = data.CoerceToString(settingsMap["schemaGen"])
			if "file" == properties["schemaGen"].(string) {
				if nil != settingsMap["schema"] {
					content, err := data.CoerceToObject(settingsMap["schema"])
					if nil == err {
						if nil != content["content"] {
							schemaBytes, err := b64.StdEncoding.DecodeString(strings.Split(content["content"].(string), ",")[1])
							if nil == err {
								properties["schema"] = string(schemaBytes)
							}
						}
					}
				}
			}

			connectorName, _ := data.CoerceToString(settingsMap["name"])

			cacheSize, exist := context.GetSetting(cacheSize)
			if exist {
				properties["cacheSize"] = cacheSize
			} else {
				log.Warn("cacheSize configuration is not configured, will turn off cache!")
			}

			readableExternalId, exist := context.GetSetting(readableExternalId)
			if exist {
				properties["readableExternalId"] = readableExternalId
			} else {
				log.Warn("readableExternalId configuration is not configured, will make readableExternalId true!")
			}

			explicitType, exist := context.GetSetting(explicitType)
			if exist {
				properties["explicitType"] = explicitType
			} else {
				log.Warn("explicitType configuration is not configured, will make type implicit!")
			}

			typeName, exist := context.GetSetting(typeTag)
			if exist {
				properties["typeName"] = typeName
			} else {
				log.Warn("Type tag is not configured, will reate an predicate as type!")
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

			dgraphService, err = factory.GetFactory(dbservice.Dgraph).CreateUpsertService(connectorName, properties)
			//dgraph.GetFactory().CreateService(connectorName, properties)

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

func getConnectionSetting(context activity.Context) (map[string]interface{}, error) {
	connection, exist := context.GetSetting(Connection)
	if !exist {
		return nil, activity.NewError("Connection is not configured", "Connection-4001", nil)
	}

	genericConn, err := generic.NewConnection(connection)
	if err != nil {
		return nil, err
	}

	settingsMap := make(map[string]interface{})
	for name, value := range genericConn.Settings() {
		if "" != name && nil != value {
			settingsMap[name] = value
		}
	}

	return settingsMap, nil
}
