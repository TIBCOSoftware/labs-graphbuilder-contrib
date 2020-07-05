/*
 * Copyright © 2020. TIBCO Software Inc.
 * This file is subject to the license terms contained
 * in the license file that is distributed with this file.
 */
package tgdbdelete

import (
	"errors"
	"sync"

	"github.com/TIBCOSoftware/flogo-lib/core/activity"
	"github.com/TIBCOSoftware/flogo-lib/core/data"
	"github.com/TIBCOSoftware/flogo-lib/logger"
	"github.com/TIBCOSoftware/labs-graphbuilder-lib/dbservice"
	"github.com/TIBCOSoftware/labs-graphbuilder-lib/dbservice/factory"
	"github.com/TIBCOSoftware/labs-graphbuilder-lib/util"
)

const (
	Connection = "tgdbConnection"
	KeepAlive  = "keepAlive"
)

var log = logger.GetLogger("tibco-activity-tgdbdelete")

type TGDBDeleteActivity struct {
	metadata            *activity.Metadata
	activityToConnector map[string]string
	entityFilter        map[string]interface{}
	mux                 sync.Mutex
}

func NewActivity(metadata *activity.Metadata) activity.Activity {
	return &TGDBDeleteActivity{
		metadata:            metadata,
		activityToConnector: make(map[string]string),
		entityFilter:        make(map[string]interface{}),
	}
}

func (a *TGDBDeleteActivity) Metadata() *activity.Metadata {
	return a.metadata
}

func (a *TGDBDeleteActivity) Eval(context activity.Context) (done bool, err error) {

	log.Info("(TGDBDeleteActivity) entering ......")
	defer log.Info("(TGDBDeleteActivity) exit ......")

	tgdbService, entityFilter, err := a.getTGDBService(context)
	if nil != err {
		return false, err
	}
	log.Info("(TGDBDeleteActivity) entityFilter = ", entityFilter)

	iInputData := context.GetInput("Graph")
	if nil == iInputData {
		return false, errors.New("Illegal nil graph data")
	}
	log.Info("(TGDBDeleteActivity) iInputData = ", iInputData)

	inputData, ok := iInputData.(map[string]interface{})
	if !ok {
		return false, errors.New("Illegal graph data type, should be map[string]interface{}.")
	}

	iGraph := inputData["graph"]
	if nil == iGraph {
		return false, errors.New("Illegal nil graph content")
	}

	graph, ok := iGraph.(map[string]interface{})

	if !ok {
		return false, errors.New("Illegal graph content, should be map[string]interface{}.")
	}

	log.Info("(TGDBDeleteActivity) tgdbService = ", tgdbService, ", entityFilter = ", entityFilter, ", graph = ", graph)
	err = tgdbService.DeleteGraph(entityFilter, graph)
	if nil != err {
		return false, err
	}

	return true, nil
}

func (a *TGDBDeleteActivity) getTGDBService(context activity.Context) (dbservice.UpsertService, int, error) {
	myId := util.ActivityId(context)

	tgdbService := factory.GetFactory(dbservice.TGDB).GetUpsertService(a.activityToConnector[myId])
	entityFilter := a.entityFilter[myId]
	if nil == tgdbService {
		a.mux.Lock()
		defer a.mux.Unlock()
		tgdbService = factory.GetFactory(dbservice.TGDB).GetUpsertService(a.activityToConnector[myId])
		entityFilter = a.entityFilter[myId]
		if nil == tgdbService {
			log.Info("Initializing TGDB Service start ...")
			defer log.Info("Initializing TGDB Service done ...")

			connection, exist := context.GetSetting(Connection)
			if !exist {
				return nil, 0, activity.NewError("TGDB connection is not configured", "TGDB-DELETE-4001", nil)
			}

			connectionInfo, _ := data.CoerceToObject(connection)
			if connectionInfo == nil {
				return nil, 0, activity.NewError("TGDB connection not able to be parsed", "TGDB-DELETE-4002", nil)
			}

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
						} else if setting["name"] == KeepAlive {
							properties[KeepAlive], _ = data.CoerceToBoolean(setting["value"])
						} else if setting["name"] == "name" {
							connectorName, _ = data.CoerceToString(setting["value"])
						}
					}
				}

				a.activityToConnector[myId] = connectorName
			}

			filterString, _ := context.GetSetting("filter")
			switch filterString {
			case "Node":
				entityFilter = 0
			case "Edge":
				entityFilter = 1
			case "Both":
				entityFilter = 2
			}

			a.entityFilter[myId] = entityFilter

			log.Info("(getTGDBService) - properties = ", properties, ", entityFilter = ", entityFilter)

			var err error
			tgdbService, err = factory.GetFactory(dbservice.TGDB).CreateUpsertService(connectorName, properties)

			log.Info("(getTGDBService) - tgdbService = ", tgdbService)

			//tgdb.GetFactory().CreateService(connectorName, properties)
			if nil != err {
				return nil, 0, err
			}
		}
	}

	return tgdbService, entityFilter.(int), nil
}
