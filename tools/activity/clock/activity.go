/*
 * Copyright © 2020. TIBCO Software Inc.
 * This file is subject to the license terms contained
 * in the license file that is distributed with this file.
 */
package clock

import (
	"fmt"
	"reflect"
	"sync"
	"time"

	"github.com/TIBCOSoftware/flogo-lib/core/activity"
	"github.com/TIBCOSoftware/flogo-lib/core/data"
	"github.com/TIBCOSoftware/flogo-lib/logger"
	"github.com/TIBCOSoftware/labs-graphbuilder-lib/tools"
	"github.com/TIBCOSoftware/labs-graphbuilder-lib/util"
)

// activityLogger is the default logger for the Filter Activity
var log = logger.GetLogger("activity-clock")

const (
	input      = "iCurrentTime"
	output     = "oCurrentTime"
	time_field = "Datetime"
)

type Clock struct {
	metadata             *activity.Metadata
	initialized          bool
	InputDatetimeType    string
	OutputDatetimeType   string
	InputDatetimeFormat  string
	OutputDatetimeFormat string
	mux                  sync.Mutex
}

// NewActivity creates a new AppActivity
func NewActivity(metadata *activity.Metadata) activity.Activity {
	aCSVParserActivity := &Clock{metadata: metadata}
	return aCSVParserActivity
}

// Metadata returns the activity's metadata
func (a *Clock) Metadata() *activity.Metadata {
	return a.metadata
}

// Eval implements api.Activity.Eval - Filters the Message
func (a *Clock) Eval(ctx activity.Context) (done bool, err error) {
	log.Info("(Clock.eval) Entering .........")

	err = a.init(ctx)

	if nil != err {
		return false, err
	}

	a.mux.Lock()
	defer a.mux.Unlock()

	inputTuple := ctx.GetInput(input).(*data.ComplexObject)

	iCurrentTime := inputTuple.Value.(map[string]interface{})[time_field]
	log.Info("iCurrentTime = ", iCurrentTime, ", type = ", reflect.TypeOf(iCurrentTime).String())

	var currentTime time.Time
	if "String" == a.InputDatetimeType {
		log.Info("String type input .........")
		currentTime, err = time.Parse(a.InputDatetimeFormat, iCurrentTime.(string))
		if nil != err {
			return false, err
		}
	} else {
		log.Info("Date type input .........")
		longCurrentLong, err := util.ConvertToLong(iCurrentTime)
		//	iCurrentTime, err := coerce.ToInt64(inputTuple.Value.(map[string]interface{})[time_field])
		if nil != err {
			return false, err
		}
		currentTime = time.Unix(longCurrentLong.(int64), 0)
	}

	log.Info("currentTime = ", currentTime)

	tools.GetClock().SetCurrentTime(currentTime.Unix())

	var oCurrentTime interface{}
	if "String" == a.OutputDatetimeType {
		log.Info("String type output .........")
		oCurrentTime = currentTime.Format(a.OutputDatetimeFormat)
	} else {
		log.Info("Date type output .........")
		oCurrentTime = currentTime.Unix()
	}

	outputTupple := make(map[string]interface{})
	outputTupple[time_field] = oCurrentTime

	complexdata := &data.ComplexObject{Metadata: output, Value: outputTupple}
	ctx.SetOutput(output, complexdata)

	return true, nil
}

func (a *Clock) init(context activity.Context) error {

	if !a.initialized {
		a.mux.Lock()
		defer a.mux.Unlock()
		if !a.initialized {

			InputDatetimeType, _ := context.GetSetting("InputDatetimeType")
			InputDatetimeFormat, _ := context.GetSetting("InputDatetimeFormat")
			if "string" == InputDatetimeType && "" == InputDatetimeFormat {
				return fmt.Errorf("InputDatetimeFormat not set !!!")
			}

			OutputDatetimeType, _ := context.GetSetting("OutputDatetimeType")
			OutputDatetimeFormat, _ := context.GetSetting("OutputDatetimeFormat")
			if "string" == OutputDatetimeType && "" == OutputDatetimeFormat {
				return fmt.Errorf("OutputDatetimeFormat not set !!!")
			}

			fmt.Println("InputDatetimeType = ", InputDatetimeType)
			fmt.Println("InputDatetimeFormat = ", InputDatetimeFormat)
			fmt.Println("OutputDatetimeType = ", OutputDatetimeType)
			fmt.Println("OutputDatetimeFormat = ", OutputDatetimeFormat)

			a.InputDatetimeType = InputDatetimeType.(string)
			a.OutputDatetimeType = OutputDatetimeType.(string)
			a.InputDatetimeFormat = InputDatetimeFormat.(string)
			a.OutputDatetimeFormat = OutputDatetimeFormat.(string)
			a.initialized = true
		}
	}

	return nil
}
