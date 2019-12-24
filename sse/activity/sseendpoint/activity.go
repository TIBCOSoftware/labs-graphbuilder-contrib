/*
 * Copyright © 2019. TIBCO Software Inc.
 * This file is subject to the license terms contained
 * in the license file that is distributed with this file.
 */
package sseendpoint

import (
	"encoding/json"
	"sync"

	"github.com/TIBCOSoftware/flogo-lib/core/activity"
	"github.com/TIBCOSoftware/flogo-lib/core/data"
	"github.com/TIBCOSoftware/flogo-lib/logger"
	"github.com/TIBCOSoftware/labs-graphbuilder-lib/internet/sseserver"
	"github.com/TIBCOSoftware/labs-graphbuilder-lib/util"
)

var log = logger.GetLogger("activity-sse-endpoint")

const (
	cConnection     = "sseConnection"
	cConnectionName = "name"
	input           = "Data"
)

type SSEEndPoint struct {
	metadata            *activity.Metadata
	activityToConnector map[string]string
	mux                 sync.Mutex
}

func NewActivity(metadata *activity.Metadata) activity.Activity {
	aCSVParserActivity := &SSEEndPoint{
		metadata:            metadata,
		activityToConnector: make(map[string]string),
	}
	return aCSVParserActivity
}

func (a *SSEEndPoint) Metadata() *activity.Metadata {
	return a.metadata
}

func (a *SSEEndPoint) Eval(ctx activity.Context) (done bool, err error) {

	broker, err := a.getSSEBroker(ctx)

	if nil != err {
		return false, err
	}

	streamId, validId := ctx.GetInput("StreamId").(string)
	if !validId {
		log.Warn("Invalid stream id, expecting string type but get : ", ctx.GetInput("StreamId"))
		streamId = "*"
	}

	data, ok := ctx.GetInput("Data").(map[string]interface{})
	if !ok {
		log.Warn("Invalid data, expecting map[string]interface{} but get : ", ctx.GetInput("Data"))
	}

	if nil != data {
		bytes, err := json.Marshal(data)
		if err != nil {
			log.Warn("Unable to serialize data : ", data)
		}
		log.Debug("Serialized data : ", bytes)

		broker.SendData(streamId, bytes)
	}

	/*
		queryResults, ok := ctx.GetInput("Data").(map[string]interface{})
		if !ok {
			log.Warn("Invalid data, expecting map[string]interface{} but get : ", queryResults)
		}

		if nil != queryResults {
			for queryId, result := range queryResults {
				log.Info("queryResult id = ", queryId, ", result = ", result)
				bytes, err := json.Marshal(result.(map[string]interface{}))
				if err != nil {
					log.Warn("Unable to serialize data : ", result)
				}
				log.Debug("Serialized data : ", bytes)

				broker.SendData(queryId, bytes)
			}
		}
	*/
	return true, nil
}

func (a *SSEEndPoint) getSSEBroker(context activity.Context) (*sseserver.Server, error) {
	myId := util.ActivityId(context)

	sseBroker := sseserver.GetFactory().GetServer(a.activityToConnector[myId])
	if nil == sseBroker {
		log.Info("Look up SSE data broker start ...")
		connection, exist := context.GetSetting(cConnection)
		if !exist {
			return nil, activity.NewError("SSE connection is not configured", "TGDB-SSE-4001", nil)
		}

		connectionInfo, _ := data.CoerceToObject(connection)
		if connectionInfo == nil {
			return nil, activity.NewError("SSE connection not able to be parsed", "TGDB-SSE-4002", nil)
		}

		var connectorName string
		connectionSettings, _ := connectionInfo["settings"].([]interface{})
		if connectionSettings != nil {
			for _, v := range connectionSettings {
				setting, _ := data.CoerceToObject(v)
				if setting != nil {
					if setting["name"] == cConnectionName {
						connectorName, _ = data.CoerceToString(setting["value"])
					}
				}
			}
			a.activityToConnector[myId] = connectorName
			sseBroker = sseserver.GetFactory().GetServer(a.activityToConnector[myId])
		}
		log.Info("Look up SSE data broker end ...")
	}

	return sseBroker, nil
}
