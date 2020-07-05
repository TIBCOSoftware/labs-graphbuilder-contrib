/*
 * Copyright © 2020. TIBCO Software Inc.
 * This file is subject to the license terms contained
 * in the license file that is distributed with this file.
 */
package xmlparser

import (
	"sync"

	"github.com/TIBCOSoftware/flogo-lib/core/activity"
	"github.com/TIBCOSoftware/flogo-lib/core/data"
	"github.com/TIBCOSoftware/flogo-lib/logger"
	"github.com/TIBCOSoftware/labs-graphbuilder-lib/parser/xml"
	"github.com/TIBCOSoftware/labs-graphbuilder-lib/util"
)

// activityLogger is the default logger for the Filter Activity
var log = logger.GetLogger("activity-xmlparser")

const (
	input  = "XMLString"
	output = "Data"
)

// JSONParserActivity is an Activity that is used to Filter a message to the console
type XMLParserActivity struct {
	metadata *activity.Metadata
	parsers  map[string]*xml.XMLParser
	mux      sync.Mutex
}

// NewActivity creates a new AppActivity
func NewActivity(metadata *activity.Metadata) activity.Activity {
	aXMLParserActivity := &XMLParserActivity{
		metadata: metadata,
		parsers:  make(map[string]*xml.XMLParser),
	}

	return aXMLParserActivity
}

// Metadata returns the activity's metadata
func (a *XMLParserActivity) Metadata() *activity.Metadata {
	return a.metadata
}

// Eval implements api.Activity.Eval - Filters the Message
func (a *XMLParserActivity) Eval(ctx activity.Context) (done bool, err error) {
	parser, _ := a.getParser(ctx)

	in := ctx.GetInput(input).(string)
	tuple := parser.Parse(in)

	xmldata := &data.ComplexObject{Metadata: "Data", Value: tuple}
	ctx.SetOutput("Data", xmldata)

	return true, nil
}

func (a *XMLParserActivity) getParser(ctx activity.Context) (*xml.XMLParser, error) {
	myId := util.ActivityId(ctx)
	parser := a.parsers[myId]

	if nil == parser {
		a.mux.Lock()
		defer a.mux.Unlock()
		parser = a.parsers[myId]
		if nil == parser {
			outputFieldnames, _ := ctx.GetSetting("OutputFieldnames")
			log.Info("Processing handlers : outputFieldnames = ", outputFieldnames)

			attributeMap := make(map[string]string)
			for _, outputFieldname := range outputFieldnames.([]interface{}) {
				outputFieldnameInfo := outputFieldname.(map[string]interface{})
				attributeMap[outputFieldnameInfo["XMLPath"].(string)] = outputFieldnameInfo["AttributeName"].(string)
				log.Info("Processing handlers : ", outputFieldnameInfo["XMLPath"].(string), " -> ", outputFieldnameInfo["AttributeName"].(string))
			}

			parser = xml.NewXMLParser(attributeMap)
			a.parsers[myId] = parser
		}
	}
	return parser, nil
}
