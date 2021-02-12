/*
 * Copyright © 2020. TIBCO Software Inc.
 * This file is subject to the license terms contained
 * in the license file that is distributed with this file.
 */
package accumulator

import (
	"sync"

	"github.com/TIBCOSoftware/flogo-lib/core/activity"
	"github.com/TIBCOSoftware/flogo-lib/core/data"
	"github.com/TIBCOSoftware/flogo-lib/logger"
	"github.com/TIBCOSoftware/labs-graphbuilder-lib/util"
)

// activityLogger is the default logger for the Filter Activity
var log = logger.GetLogger("activity-accumulator")

const (
	arrayMode   = "ArrayMode"
	windowSize  = "WindowSize"
	lastElement = "LastElement"
	input       = "Input"
	output      = "Output"
)

// Mapping is an Activity that is used to Filter a message to the console
type Accumulator struct {
	metadata    *activity.Metadata
	initialized bool
	mux         sync.Mutex
	windows     map[string]*Window
}

// NewActivity creates a new AppActivity
func NewActivity(metadata *activity.Metadata) activity.Activity {
	aCSVParserActivity := &Accumulator{
		metadata: metadata,
		windows:  make(map[string]*Window),
	}
	return aCSVParserActivity
}

// Metadata returns the activity's metadata
func (a *Accumulator) Metadata() *activity.Metadata {
	return a.metadata
}

// Eval implements api.Activity.Eval - Filters the Message
func (a *Accumulator) Eval(ctx activity.Context) (done bool, err error) {
	log.Info("[Accumulator:Eval] entering ........ ")
	defer log.Info("[Accumulator:Eval] Exit ........ ")

	inputTuple := ctx.GetInput(input).(*data.ComplexObject)
	log.Info("Input data = ", inputTuple.Value)

	arrayMode, _ := ctx.GetSetting(arrayMode)
	var outputTuple []map[string]interface{}
	if nil != arrayMode && arrayMode.(bool) {
		rawOutputTuple := inputTuple.Value.([]interface{})
		outputTuple = make([]map[string]interface{}, len(rawOutputTuple))
		for index, tuple := range rawOutputTuple {
			outputTuple[index] = tuple.(map[string]interface{})
			if index < len(rawOutputTuple)-1 {
				outputTuple[index][lastElement] = false
			} else {
				outputTuple[index][lastElement] = true
			}
		}
		log.Info("(Array Mode) Output data = ", outputTuple)
	} else {
		window, err := a.getWindow(ctx)

		if nil != err {
			log.Error(err)
			return false, err
		}

		outputTuple, err = window.update(inputTuple.Value.(map[string]interface{}))
		if nil != err {
			log.Error(err)
			return false, err
		}
		log.Info("(Iterator Mode) Output data = ", outputTuple)
	}

	if nil != outputTuple {
		ctx.SetOutput(output, outputTuple)
		log.Debug("Output data = ", outputTuple)
	} else {
		return false, nil
	}

	return true, nil
}

func (a *Accumulator) getWindow(context activity.Context) (*Window, error) {
	myId := util.ActivityId(context)
	window := a.windows[myId]

	if nil == window {
		a.mux.Lock()
		defer a.mux.Unlock()
		window = a.windows[myId]
		if nil == window {
			windowSize, _ := context.GetSetting(windowSize)
			window = NewWindow(windowSize.(int))
			a.windows[myId] = window
		}
	}

	return window, nil
}

func NewWindow(maxSize int) *Window {
	if 0 >= maxSize {
		maxSize = 1
	}
	window := &Window{
		currentSize: 0,
		maxSize:     maxSize,
		tuples:      make([]map[string]interface{}, maxSize),
	}
	return window
}

type Window struct {
	currentSize int
	maxSize     int
	tuples      []map[string]interface{}
}

func (this *Window) update(tuple map[string]interface{}) ([]map[string]interface{}, error) {
	this.currentSize += 1
	this.tuples[this.currentSize-1] = tuple
	if this.currentSize >= this.maxSize {
		tuple[lastElement] = true
		this.currentSize = 0
		return this.tuples, nil
	} else {
		tuple[lastElement] = false
	}
	return nil, nil
}
