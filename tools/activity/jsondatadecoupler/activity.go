/*
 * Copyright © 2020. TIBCO Software Inc.
 * This file is subject to the license terms contained
 * in the license file that is distributed with this file.
 */
package jsondatadecoupler

import (
	"fmt"
	"strings"
	"sync"

	"github.com/TIBCOSoftware/flogo-lib/core/activity"
	"github.com/TIBCOSoftware/flogo-lib/core/data"
	"github.com/TIBCOSoftware/flogo-lib/logger"
	"github.com/TIBCOSoftware/labs-graphbuilder-lib/util"
)

// activityLogger is the default logger for the Filter Activity
var log = logger.GetLogger("activity-jsondatadecoupler")

const (
	setting_DateFormat = "DateFormat"
	input              = "JSONString"
	output             = "Data"
)

// JSONDataDecouplerActivity is an Activity that is used to Filter a message to the console
type JSONDataDecouplerActivity struct {
	metadata          *activity.Metadata
	mux               sync.Mutex
	targetMap         map[string]string
	targetElementsMap map[string][]string
}

// NewActivity creates a new AppActivity
func NewActivity(metadata *activity.Metadata) activity.Activity {
	aJSONDataDecouplerActivity := &JSONDataDecouplerActivity{
		metadata: metadata,
	}
	return aJSONDataDecouplerActivity
}

// Metadata returns the activity's metadata
func (a *JSONDataDecouplerActivity) Metadata() *activity.Metadata {
	return a.metadata
}

// Eval implements api.Activity.Eval - Filters the Message
func (a *JSONDataDecouplerActivity) Eval(ctx activity.Context) (done bool, err error) {

	originJSONObject := ctx.GetInput(input)
	targetPath, targetPathElements, _ := a.getTarget(ctx)
	target := originJSONObject
	for _, targetPathElement := range targetPathElements {
		target = target.(map[string]interface{})[targetPathElement]
	}

	targetArray := target.([]interface{})
	outputArray := make([]interface{}, len(targetArray))
	for index, targetElement := range targetArray {
		outputElement := make(map[string]interface{})
		outputElement["originJSONObject"] = originJSONObject
		outputElement[fmt.Sprintf("%s.%s", targetPath, "Index")] = index
		outputElement[fmt.Sprintf("%s.%s", targetPath, "Element")] = targetElement
		if len(targetArray)-1 == index {
			outputElement["LastElement"] = true
		} else {
			outputElement["LastElement"] = false
		}
	}

	jsondata := &data.ComplexObject{Metadata: "Data", Value: outputArray}

	ctx.SetOutput(output, jsondata)

	return true, nil
}

func (a *JSONDataDecouplerActivity) getTarget(ctx activity.Context) (string, []string, error) {
	myId := util.ActivityId(ctx)
	target := a.targetMap[myId]
	targetElements := a.targetElementsMap[myId]

	if nil == targetElements {
		a.mux.Lock()
		defer a.mux.Unlock()
		target = a.targetMap[myId]
		targetElements := a.targetElementsMap[myId]
		if nil == targetElements {
			decoupleTarget, _ := ctx.GetSetting("decoupleTarget")
			a.targetMap[myId] = decoupleTarget.(string)
			a.targetElementsMap[myId] = strings.Split(decoupleTarget.(string), ".")
		}
	}
	return target, targetElements, nil
}
