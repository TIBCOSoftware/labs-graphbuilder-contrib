/*
 * Copyright © 2019. TIBCO Software Inc.
 * This file is subject to the license terms contained
 * in the license file that is distributed with this file.
 */
package tgdbupsert

import (
	"errors"
	"sync"

	"git.tibco.com/git/product/ipaas/wi-contrib.git/connection/generic"
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
	defer log.Info("(TGDBUpsertActivity) exit ......")

	tgdbService, err := a.getTGDBService(context)
	if nil != err {
		return false, err
	}

	iInputData := context.GetInput("Graph")
	if nil == iInputData {
		return false, errors.New("Illegal nil graph data")
	}

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

	err = tgdbService.UpsertGraph(nil, graph)
	if nil != err {
		return false, err
	}

	return true, nil
}

func (a *TGDBUpsertActivity) getTGDBService(context activity.Context) (dbservice.UpsertService, error) {
	myId := util.ActivityId(context)

	tgdbService := factory.GetFactory(dbservice.TGDB).GetUpsertService(a.activityToConnector[myId])
	//tgdb.GetFactory().GetService(a.activityToConnector[myId])
	if nil == tgdbService {
		a.mux.Lock()
		defer a.mux.Unlock()
		tgdbService = factory.GetFactory(dbservice.TGDB).GetUpsertService(a.activityToConnector[myId])
		//tgdb.GetFactory().GetService(a.activityToConnector[myId])
		if nil == tgdbService {
			log.Info("Initializing TGDB Service start ...")
			defer log.Info("Initializing TGDB Service done ...")

			connection, exist := context.GetSetting(Connection)
			if !exist {
				return nil, activity.NewError("TGDB connection is not configured", "TGDB-UPSERT-4001", nil)
			}

			genericConn, err := generic.NewConnection(connection)
			if err != nil {
				return nil, err
			}

			var connectorName string
			properties := make(map[string]interface{})
			for name, value := range genericConn.Settings() {
				switch name {
				case "url":
					properties["url"], _ = data.CoerceToString(value)
				case "user":
					properties["user"], _ = data.CoerceToString(value)
				case "password":
					properties["password"], _ = data.CoerceToString(value)
				case KeepAlive:
					properties[KeepAlive], _ = data.CoerceToBoolean(value)
				case "name":
					connectorName, _ = data.CoerceToString(value)
				}
			}

			a.activityToConnector[myId] = connectorName
			allowEmptyStringKey, exist := context.GetSetting("allowEmptyStringKey")
			if exist {
				properties["allowEmptyStringKey"] = allowEmptyStringKey
			} else {
				log.Warn("allowEmptyStringKey configuration is not configured, will make type defininated implicit!")
			}

			log.Info("(getTGDBService) - properties = ", properties)

			tgdbService, err = factory.GetFactory(dbservice.TGDB).CreateUpsertService(connectorName, properties)
			//tgdb.GetFactory().CreateService(connectorName, properties)
			if nil != err {
				return nil, err
			}
		}
	}

	return tgdbService, nil
}
