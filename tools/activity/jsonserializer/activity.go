/*
 * Copyright © 2020. TIBCO Software Inc.
 * This file is subject to the license terms contained
 * in the license file that is distributed with this file.
 */
package jsonserializer

import (
	"encoding/json"

	"github.com/TIBCOSoftware/flogo-lib/core/activity"
	"github.com/TIBCOSoftware/flogo-lib/core/data"
	"github.com/TIBCOSoftware/flogo-lib/logger"
)

var log = logger.GetLogger("activity-jsondeserializer")

const (
	iData       = "Data"
	oJSONString = "JSONString"
)

type JSONSerializerActivity struct {
	metadata *activity.Metadata
}

func NewActivity(metadata *activity.Metadata) activity.Activity {
	aJSONSerializerActivity := &JSONSerializerActivity{
		metadata: metadata,
	}
	return aJSONSerializerActivity
}

func (a *JSONSerializerActivity) Metadata() *activity.Metadata {
	return a.metadata
}

func (a *JSONSerializerActivity) Eval(ctx activity.Context) (done bool, err error) {

	data, ok := ctx.GetInput(iData).(*data.ComplexObject).Value.(map[string]interface{})
	if !ok {
		log.Warn("No valid data ... ")
	}

	jsondata, err := json.Marshal(data)
	if nil != err {
		logger.Warn("Unable to serialize object, reason : ", err.Error())
		return false, nil
	}

	ctx.SetOutput(oJSONString, string(jsondata))

	return true, nil
}
