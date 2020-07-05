/*
 * Copyright © 2020. TIBCO Software Inc.
 * This file is subject to the license terms contained
 * in the license file that is distributed with this file.
 */
package aggregate

import (
	"sync"

	"github.com/TIBCOSoftware/flogo-lib/core/activity"
	"github.com/TIBCOSoftware/flogo-lib/logger"
	"github.com/TIBCOSoftware/labs-graphbuilder-lib/tools"
)

// activityLogger is the default logger for the Filter Activity
var log = logger.GetLogger("activity-aggregate")

const (
	input  = "Data"
	output = "AggregatedData"
)

// Aggregate is an Activity that is used to Aggregate input data
type Aggregate struct {
	metadata    *activity.Metadata
	initialized bool
	aggStates   map[string]map[tools.Index]map[string]tools.DataState
	mux         sync.Mutex
}

// NewActivity creates a new AppActivity
func NewActivity(metadata *activity.Metadata) activity.Activity {
	aCSVParserActivity := &Aggregate{metadata: metadata}
	return aCSVParserActivity
}

// Metadata returns the activity's metadata
func (a *Aggregate) Metadata() *activity.Metadata {
	return a.metadata
}

// Eval implements api.Activity.Eval - Filters the Message
func (a *Aggregate) Eval(ctx activity.Context) (done bool, err error) {

	log.Info("(Aggregated Eval entering .....")

	err = a.init(ctx)

	if nil != err {
		return false, err
	}

	dataTuples := ctx.GetInput(input).(map[string]interface{})

	log.Info("(Eval) input data : ", dataTuples)

	aggregatedDataTuples := make(map[string]interface{})
	if nil != dataTuples {
		for queryId, dataTuple := range dataTuples {
			log.Info("queryResult id = ", queryId, ", dataTuple = ", dataTuple)

			aggStateOfQuery := a.aggStates[queryId]
			if nil == aggStateOfQuery {
				aggStateOfQuery = make(map[tools.Index]map[string]tools.DataState)
				a.aggStates[queryId] = aggStateOfQuery
			}

			data := dataTuple.(map[string]interface{})
			if nil == data["parameter"] || 0 == len(data["parameter"].(map[string]interface{})) {
				aggregatedDataTuples[queryId] = dataTuples[queryId]
			} else {
				aggregatedDataTuple := tools.Agg(aggStateOfQuery, data)
				log.Info("queryResult id = ", queryId, ", aggregatedDataTuple = ", aggregatedDataTuple)
				if 0 != len(aggregatedDataTuple) {
					aggregatedDataTuples[queryId] = aggregatedDataTuple
				}

			}
		}
		if 0 == len(aggregatedDataTuples) {
			return false, nil
		}
	} else {
		log.Warn("Nill tuples won't be processed, dataTuples = ", dataTuples)
		return false, nil
	}

	log.Info("(Eval) output data = ", aggregatedDataTuples)

	ctx.SetOutput(output, aggregatedDataTuples)

	return true, nil
}

func (a *Aggregate) init(context activity.Context) error {

	if !a.initialized {
		a.mux.Lock()
		defer a.mux.Unlock()
		if !a.initialized {
			a.initialized = true
			a.aggStates = make(map[string]map[tools.Index]map[string]tools.DataState)
		}
	}

	return nil
}
