/*
 * Copyright © 2019. TIBCO Software Inc.
 * This file is subject to the license terms contained
 * in the license file that is distributed with this file.
 */
package jsondeserializer

import (
	"encoding/json"

	"github.com/TIBCOSoftware/flogo-lib/core/activity"
	"github.com/TIBCOSoftware/flogo-lib/core/data"
	"github.com/TIBCOSoftware/flogo-lib/logger"
)

// activityLogger is the default logger for the Filter Activity
var log = logger.GetLogger("activity-jsonparser")

const (
	setting_DateFormat = "DateFormat"
	input              = "JSONString"
	output             = "Data"
)

// JSONDeserializerActivity is an Activity that is used to Filter a message to the console
type JSONDeserializerActivity struct {
	metadata *activity.Metadata
}

// NewActivity creates a new AppActivity
func NewActivity(metadata *activity.Metadata) activity.Activity {
	aJSONDeserializerActivity := &JSONDeserializerActivity{
		metadata: metadata,
	}
	return aJSONDeserializerActivity
}

// Metadata returns the activity's metadata
func (a *JSONDeserializerActivity) Metadata() *activity.Metadata {
	return a.metadata
}

// Eval implements api.Activity.Eval - Filters the Message
func (a *JSONDeserializerActivity) Eval(ctx activity.Context) (done bool, err error) {

	in := ctx.GetInput(input).(string)

	var rootObject interface{}

	err = json.Unmarshal([]byte(in), &rootObject)
	if nil != err {
		return false, err
	}

	jsondata := &data.ComplexObject{Metadata: "Data", Value: rootObject}

	ctx.SetOutput(output, jsondata)
	return true, nil
}
