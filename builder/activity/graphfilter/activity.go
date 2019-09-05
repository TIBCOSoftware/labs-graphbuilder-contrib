/*
 * Copyright © 2019. TIBCO Software Inc.
 * This file is subject to the license terms contained
 * in the license file that is distributed with this file.
 */
package graphfilter

import (
	"encoding/json"
	"fmt"
	"sync"
	"time"

	"github.com/TIBCOSoftware/flogo-lib/core/activity"
	"github.com/TIBCOSoftware/flogo-lib/logger"

	"github.com/TIBCOSoftware/labs-graphbuilder-lib/model"
	"github.com/TIBCOSoftware/labs-graphbuilder-lib/query"
	"github.com/TIBCOSoftware/labs-graphbuilder-lib/tools"
	"github.com/TIBCOSoftware/labs-graphbuilder-lib/util"
)

const (
	FILTER_FILE = "testFilter"
)

var log = logger.GetLogger("tibco-activity-graphfilter")

type GraphFilterActivity struct {
	metadata *activity.Metadata
	once     sync.Once
	mux      sync.Mutex
}

func NewActivity(metadata *activity.Metadata) activity.Activity {
	return &GraphFilterActivity{
		metadata: metadata,
	}
}

func (a *GraphFilterActivity) Metadata() *activity.Metadata {
	return a.metadata
}

func (a *GraphFilterActivity) Eval(context activity.Context) (done bool, err error) {
	log.Info("(Eval) GraphFilter entering, current time =  ", time.Unix(tools.GetClock().GetCurrentTime(), 0))

	a.mux.Lock()
	defer a.mux.Unlock()

	graphData, exists := context.GetInput("Graph").(map[string]interface{})["graph"].(map[string]interface{})
	if !exists {
		return false, fmt.Errorf("Unable to get data from input!!")
	}
	deltaGraph := model.ReconstructGraph(graphData)

	graphId := deltaGraph.GetId()
	a.init(graphId, context)

	log.Info("(GraphFilterActivity::Eval) graphId = ", graphId)

	/* UpsertGraph */
	graph := model.GetGraphManager().GetGraph(graphId, graphId)
	(*graph).UpsertGraph(deltaGraph)
	//	log.Info("(GraphFilterActivity::Eval) After upsert graph = ", graph)
	//	log.Info("(GraphFilterActivity::Eval) After upsert deltaGraph = ", deltaGraph)

	/* GetQuery */
	queries := query.GetQueryManager().GetQueries(graphId)

	if nil != queries {
		matchedDatas := make(map[string]interface{})
		for queryId, query := range queries {
			log.Info("(GraphFilterActivity::Eval) graphId = ", graphId, ", queryId = ", queryId)
			matchedData := query.Match(util.ActivityId(context), graph, deltaGraph)
			if 0 != len(matchedData) {
				matchedDatas[queryId] = matchedData
			}
		}

		log.Info("(GraphFilterActivity::Eval) matchedData = ", matchedDatas)

		if 0 != len(matchedDatas) {
			context.SetOutput("MatchedData", matchedDatas)
			return true, nil
		}
	}

	log.Info("(Eval) GraphFilter exit ... No matched data!!!")
	return false, nil
}

func (a *GraphFilterActivity) init(graphId string, context activity.Context) {

	a.once.Do(func() {
		filterFile, exists := context.GetSetting(FILTER_FILE)
		if exists {
			filterString, err := util.ReadFile(filterFile.(string))
			if nil != err {
				log.Info(err)
			}

			var queryObject interface{}
			err = json.Unmarshal([]byte(filterString), &queryObject)
			if nil != err {
				log.Info("No default filter defined!")
				return
			}
			query.GetQueryManager().CreateQuery(graphId, "DEFAULT_QUERY", queryObject.(map[string]interface{}))
		} else {
			log.Info("No default filter defined!")
		}
	})
}
