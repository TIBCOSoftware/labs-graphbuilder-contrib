/*
 * Copyright © 2019. TIBCO Software Inc.
 * This file is subject to the license terms contained
 * in the license file that is distributed with this file.
 */
package dummybuilder

import (
	b64 "encoding/base64"
	"strings"
	"sync"

	"github.com/TIBCOSoftware/flogo-lib/core/activity"
	"github.com/TIBCOSoftware/flogo-lib/core/data"
	"github.com/TIBCOSoftware/flogo-lib/logger"
)

var log = logger.GetLogger("tibco-activity-dummy-graphbuilder")

var initialized bool = false

const (
	GraphModel = "GraphModel"
	Nodes      = "Nodes"
	Edges      = "Edges"
)

type DummyBuilderActivity struct {
	metadata *activity.Metadata
	mux      sync.Mutex
}

func NewActivity(metadata *activity.Metadata) activity.Activity {
	aBuilderActivity := &DummyBuilderActivity{
		metadata: metadata,
	}

	return aBuilderActivity
}

func (a *DummyBuilderActivity) Metadata() *activity.Metadata {
	return a.metadata
}

func (a *DummyBuilderActivity) Eval(context activity.Context) (done bool, err error) {

	log.Info("[DummyBuilderActivity:Eval] entering ........ ")

	err = a.getGraphModel(context)

	if nil != err {
		return false, err
	}

	log.Info("[DummyBuilderActivity:Eval] Exit ........ ")

	return true, nil
}

func (a *DummyBuilderActivity) getGraphModel(context activity.Context) error {
	graphmodel, exist := context.GetSetting(GraphModel)
	if !exist {
		return activity.NewError("GraphModel is not configured", "DUMMY_GRAPHBUILDER-4002", nil)
	}
	log.Info("Got model = ", graphmodel)

	//Read graph model details
	connectionInfo, _ := data.CoerceToObject(graphmodel)
	if connectionInfo == nil {
		return activity.NewError("Unable extract model", "DUMMY_GRAPHBUILDER-4001", nil)
	}
	log.Info("Got connectionInfo = ", connectionInfo)

	var jsonmodel []byte
	var modelName string
	connectionSettings, _ := connectionInfo["settings"].([]interface{})
	if connectionSettings != nil {
		for _, v := range connectionSettings {
			setting, _ := data.CoerceToObject(v)

			if nil != setting {
				if setting["name"] == "model" {
					modelcontent, _ := data.CoerceToObject(setting["value"])
					jsonmodel, _ = b64.StdEncoding.DecodeString(strings.Split(modelcontent["content"].(string), ",")[1])
				} else if setting["name"] == "name" {
					modelName = setting["value"].(string)
				}
			}
		}
	}

	if "" == modelName {
		return activity.NewError("Unable to get builder name", "DUMMY_GRAPHBUILDER-4003", nil)
	}
	log.Info("Got modelName = ", modelName)

	if nil == jsonmodel {
		return activity.NewError("Unable to get model string", "DUMMY_GRAPHBUILDER-4004", nil)
	}
	log.Info("Got jsonmodel = ", jsonmodel)

	log.Info("******* Graph Model loaded correctly ********")

	return nil
}
