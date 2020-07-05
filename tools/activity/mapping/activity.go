/*
 * Copyright © 2020. TIBCO Software Inc.
 * This file is subject to the license terms contained
 * in the license file that is distributed with this file.
 */
package mapping

import (
	"sync"

	"github.com/TIBCOSoftware/flogo-lib/core/activity"
	"github.com/TIBCOSoftware/flogo-lib/core/data"
	"github.com/TIBCOSoftware/flogo-lib/logger"
)

// activityLogger is the default logger for the Filter Activity
var log = logger.GetLogger("activity-mapping")

const (
	input  = "Mapping"
	output = "Mapped"
)

// Mapping is an Activity that is used to Filter a message to the console
type Mapping struct {
	metadata    *activity.Metadata
	initialized bool
	mux         sync.Mutex
}

// NewActivity creates a new AppActivity
func NewActivity(metadata *activity.Metadata) activity.Activity {
	aCSVParserActivity := &Mapping{metadata: metadata}
	return aCSVParserActivity
}

// Metadata returns the activity's metadata
func (a *Mapping) Metadata() *activity.Metadata {
	return a.metadata
}

// Eval implements api.Activity.Eval - Filters the Message
func (a *Mapping) Eval(ctx activity.Context) (done bool, err error) {

	err = a.init(ctx)

	if nil != err {
		return false, err
	}

	mappedTuple := ctx.GetInput(input).(*data.ComplexObject)
	log.Info("mapped data = ", mappedTuple.Value)

	//complexdata := &data.ComplexObject{Metadata: "Mapped", Value: mappedTuple}
	ctx.SetOutput("Mapped", mappedTuple)

	return true, nil
}

func (a *Mapping) init(context activity.Context) error {

	return nil
}
