/*
 * Copyright © 2020. TIBCO Software Inc.
 * This file is subject to the license terms contained
 * in the license file that is distributed with this file.
 */
package simplefilter

import (
	"fmt"
	"sync"

	"github.com/TIBCOSoftware/flogo-lib/core/activity"
	"github.com/TIBCOSoftware/flogo-lib/core/data"
	"github.com/TIBCOSoftware/flogo-lib/logger"
)

// activityLogger is the default logger for the Filter Activity
var activityLogger = logger.GetLogger("activity-simple-filter")

const (
	passFilterFlag     = "filterExpression"
	sProceedOnlyOnEmit = "proceedOnlyOnEmit"
	ivValue            = "value"
	ovFiltered         = "filtered"
	ovValue            = "value"
)

// FilterActivity is an Activity that is used to Filter a message to the console
type FilterActivity struct {
	metadata *activity.Metadata
	mux      sync.Mutex
}

// NewActivity creates a new AppActivity
func NewActivity(md *activity.Metadata) activity.Activity {
	filterActivity := FilterActivity{metadata: md}

	return &filterActivity
}

// Metadata returns the activity's metadata
func (a *FilterActivity) Metadata() *activity.Metadata {
	return a.metadata
}

// Eval implements api.Activity.Eval - Filters the Message
func (a *FilterActivity) Eval(ctx activity.Context) (done bool, err error) {

	proceedOnlyOnEmit, err := a.getProceedOnlyOnEmit(ctx)

	if nil != err {
		return true, err
	}

	in := ctx.GetInput(ivValue)

	passFilter := ctx.GetInput(passFilterFlag).(bool)

	done = !(proceedOnlyOnEmit && !passFilter)

	activityLogger.Info("proceedOnlyOnEmit : ", proceedOnlyOnEmit, ", passFilter : ", passFilter, ", done : ", done)

	if done {
		ctx.SetOutput(ovFiltered, !passFilter)
		ctx.SetOutput(ovValue, in)
	}

	return done, nil
}

func (a *FilterActivity) getProceedOnlyOnEmit(context activity.Context) (bool, error) {
	setting, exists := context.GetSetting(sProceedOnlyOnEmit)
	if !exists {
		return true, fmt.Errorf("ProceedOnlyOnEmit flag not found!")
	}

	val, err := data.CoerceToBoolean(setting)
	if err != nil {
		return true, err
	}

	return val, nil
}
