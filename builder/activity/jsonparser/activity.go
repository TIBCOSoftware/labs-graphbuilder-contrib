/*
 * Copyright © 2019. TIBCO Software Inc.
 * This file is subject to the license terms contained
 * in the license file that is distributed with this file.
 */
package jsonparser

import (
	"sync"

	"github.com/TIBCOSoftware/flogo-lib/core/activity"
	"github.com/TIBCOSoftware/flogo-lib/core/data"
	"github.com/TIBCOSoftware/flogo-lib/logger"
	"github.com/TIBCOSoftware/graph-builder-flogo/lib/parser/json"
	"github.com/TIBCOSoftware/graph-builder-flogo/lib/util"
)

// activityLogger is the default logger for the Filter Activity
var log = logger.GetLogger("activity-jsonparser")

const (
	setting_DateFormat = "DateFormat"
	input              = "JSONString"
	output             = "Data"
)

// JSONParserActivity is an Activity that is used to Filter a message to the console
type JSONParserActivity struct {
	metadata *activity.Metadata
	parsers  map[string]*json.JSONParser
	mux      sync.Mutex
}

// NewActivity creates a new AppActivity
func NewActivity(metadata *activity.Metadata) activity.Activity {
	aJSONParserActivity := &JSONParserActivity{
		metadata: metadata,
		parsers:  make(map[string]*json.JSONParser),
	}
	return aJSONParserActivity
}

// Metadata returns the activity's metadata
func (a *JSONParserActivity) Metadata() *activity.Metadata {
	return a.metadata
}

// Eval implements api.Activity.Eval - Filters the Message
func (a *JSONParserActivity) Eval(ctx activity.Context) (done bool, err error) {

	parser, err := a.getParser(ctx)

	if nil != err {
		return false, err
	}

	in := ctx.GetInput(input).(string)

	tupleArray := parser.Parse([]byte(in))

	jsondata := &data.ComplexObject{Metadata: "Data", Value: tupleArray}
	ctx.SetOutput(output, jsondata)

	return true, nil
}

func (a *JSONParserActivity) getParser(ctx activity.Context) (*json.JSONParser, error) {
	myId := util.ActivityId(ctx)
	parser := a.parsers[myId]

	if nil == parser {
		a.mux.Lock()
		defer a.mux.Unlock()
		parser = a.parsers[myId]
		if nil == parser {

			attributeMap := make(map[string]*json.Attribute)

			outputFieldnames, _ := ctx.GetSetting("OutputFieldnames")
			log.Info("Processing handlers : outputFieldnames = ", outputFieldnames)

			for _, outputFieldname := range outputFieldnames.([]interface{}) {
				outputFieldnameInfo := outputFieldname.(map[string]interface{})
				attribute := &json.Attribute{}
				attribute.SetName(outputFieldnameInfo["AttributeName"].(string))
				attribute.SetType(outputFieldnameInfo["Type"].(string))
				attributeMap[outputFieldnameInfo["JSONPath"].(string)] = attribute
				log.Info("[JSONParserActivity::Eval] Mapping : ", outputFieldnameInfo["JSONPath"].(string), " -> ", outputFieldnameInfo["AttributeName"].(string))
			}

			parser = json.NewJSONParser(attributeMap)
			dateFormat, _ := ctx.GetSetting(setting_DateFormat)
			if nil != dateFormat {
				parser.SetDateFromatString(dateFormat.(string))
			}

			a.parsers[myId] = parser
		}
	}
	return parser, nil
}
