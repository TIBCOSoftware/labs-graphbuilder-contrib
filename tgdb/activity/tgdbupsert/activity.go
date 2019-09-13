/*
 * Copyright © 2019. TIBCO Software Inc.
 * This file is subject to the license terms contained
 * in the license file that is distributed with this file.
 */
package tgdbupsert

import (
	"sync"

	"github.com/TIBCOSoftware/flogo-lib/core/activity"
	"github.com/TIBCOSoftware/flogo-lib/core/data"
	"github.com/TIBCOSoftware/flogo-lib/logger"
	"github.com/TIBCOSoftware/labs-graphbuilder-lib/dbservice/tgdb"
	"github.com/TIBCOSoftware/labs-graphbuilder-lib/util"
)

const (
	Connection = "tgdbConnection"
	KeepAlive  = "keepAlive"
)

var log = logger.GetLogger("tibco-activity-tgdbupsert")

type TGDBUpsertActivity struct {
	metadata            *activity.Metadata
	activityToConnector map[string]string
	mux                 sync.Mutex
}

func NewActivity(metadata *activity.Metadata) activity.Activity {
	return &TGDBUpsertActivity{
		metadata:            metadata,
		activityToConnector: make(map[string]string),
	}
}

func (a *TGDBUpsertActivity) Metadata() *activity.Metadata {
	return a.metadata
}

func (a *TGDBUpsertActivity) Eval(context activity.Context) (done bool, err error) {

	log.Info("(TGDBUpsertActivity) entering ......")

	tgdbService, err := a.getTGDBService(context)

	if nil != err {
		return false, err
	}

	graph, _ := context.GetInput("Graph").(map[string]interface{})["graph"].(map[string]interface{})

	err = tgdbService.UpsertGraph(graph)

	if nil != err {
		return false, err
	}

	log.Info("(TGDBUpsertActivity) exit normally ......")

	return true, nil
}

func (a *TGDBUpsertActivity) getTGDBService(context activity.Context) (*tgdb.TGDBService, error) {
	myId := util.ActivityId(context)

	tgdbService := tgdb.GetFactory().GetService(a.activityToConnector[myId])
	if nil == tgdbService {
		a.mux.Lock()
		defer a.mux.Unlock()
		tgdbService = tgdb.GetFactory().GetService(a.activityToConnector[myId])
		if nil == tgdbService {
			log.Info("Initializing TGDB Service start ...")
			connection, exist := context.GetSetting(Connection)
			if !exist {
				return nil, activity.NewError("TGDB connection is not configured", "TGDB-UPSERT-4001", nil)
			}

			connectionInfo, _ := data.CoerceToObject(connection)
			if connectionInfo == nil {
				return nil, activity.NewError("TGDB connection not able to be parsed", "TGDB-UPSERT-4002", nil)
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

			allowEmptyStringKey, exist := context.GetSetting("allowEmptyStringKey")
			if exist {
				properties["allowEmptyStringKey"] = allowEmptyStringKey
			} else {
				log.Warn("allowEmptyStringKey configuration is not configured, will make type defininated implicit!")
			}

			log.Info("(getTGDBService) - properties = ", properties)

			tgdbService, _ = tgdb.GetFactory().CreateService(connectorName, properties)

			log.Info("Initializing TGDB Service end, tgdbService = ", tgdbService)
		}
	}

	return tgdbService, nil
}
