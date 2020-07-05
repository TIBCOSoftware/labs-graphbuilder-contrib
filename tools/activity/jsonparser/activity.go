/*
 * Copyright © 2020. TIBCO Software Inc.
 * This file is subject to the license terms contained
 * in the license file that is distributed with this file.
 */
package jsonparser

import (
	"fmt"
	"strings"
	"sync"

	"github.com/TIBCOSoftware/flogo-lib/core/activity"
	"github.com/TIBCOSoftware/flogo-lib/core/data"
	"github.com/TIBCOSoftware/flogo-lib/logger"
	"github.com/TIBCOSoftware/labs-graphbuilder-lib/parser/json"
	"github.com/TIBCOSoftware/labs-graphbuilder-lib/util"
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

	in, ok := ctx.GetInput(input).(string)

	if !ok {
		log.Warn("Illegal input data : ", ctx.GetInput(input))
		return false, nil
	}

	tupleArray := parser.Parse([]byte(in))

	if nil == tupleArray {
		log.Warn("Input data can not be parsed : ", in)
		return false, nil
	}

	log.Info("Valid data : ", tupleArray)

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
			mandatoryAttrs := make(map[string]bool)

			outputFieldnames, _ := ctx.GetSetting("OutputFieldnames")
			log.Info("Processing handlers : outputFieldnames = ", outputFieldnames)

			var multiInstancePathPrefix string
			for _, outputFieldname := range outputFieldnames.([]interface{}) {
				outputFieldnameInfo := outputFieldname.(map[string]interface{})
				attribute := &json.Attribute{}
				attribute.SetName(outputFieldnameInfo["AttributeName"].(string))
				attribute.SetType(outputFieldnameInfo["Type"].(string))
				if nil != outputFieldnameInfo["Optional"] && "no" == outputFieldnameInfo["Optional"].(string) {
					attribute.SetOptional(false)
					mandatoryAttrs[attribute.GetName()] = true
				} else {
					attribute.SetOptional(true)
				}
				if nil != outputFieldnameInfo["Default"] && "" != outputFieldnameInfo["Default"].(string) {
					//attribute.SetDValue()
				}

				jsonPath := outputFieldnameInfo["JSONPath"].(string)
				pos := strings.Index(jsonPath, "[]")
				log.Debug("[JSONParserActivity::Eval] jsonPath : ", jsonPath, " ,pos : ", pos)
				if 0 <= pos {
					log.Debug("[JSONParserActivity::Eval] jsonPath[0:pos+1] : ", jsonPath[0:pos+1])
					if "" != multiInstancePathPrefix {
						if jsonPath[0:pos+1] != multiInstancePathPrefix {
							return nil, fmt.Errorf("multiInstancePathPrefix confixt !!")
						}
					} else {
						multiInstancePathPrefix = jsonPath[0 : pos+1]
					}
					attribute.SetMultiInstance(true)
				}

				attributeMap[jsonPath] = attribute
				log.Info("[JSONParserActivity::Eval] Mapping : ", outputFieldnameInfo["JSONPath"].(string), " -> ", outputFieldnameInfo["AttributeName"].(string))
			}

			parser = json.NewJSONParser(attributeMap, mandatoryAttrs, multiInstancePathPrefix)
			dateFormat, _ := ctx.GetSetting(setting_DateFormat)
			if nil != dateFormat {
				parser.SetDateFromatString(dateFormat.(string))
			}

			a.parsers[myId] = parser
		}
	}
	return parser, nil
}
