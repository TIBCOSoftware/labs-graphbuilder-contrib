/*
 * Copyright © 2019. TIBCO Software Inc.
 * This file is subject to the license terms contained
 * in the license file that is distributed with this file.
 */
package graphquerier

import (
	"encoding/json"
	"fmt"
	"sync"

	"github.com/TIBCOSoftware/flogo-lib/core/activity"
	"github.com/TIBCOSoftware/flogo-lib/logger"

	"github.com/TIBCOSoftware/labs-graphbuilder-lib/model"
	"github.com/TIBCOSoftware/labs-graphbuilder-lib/query"
	"github.com/TIBCOSoftware/labs-graphbuilder-lib/util"
)

const (
	FILTER_FILE = "testFilter"
)

var log = logger.GetLogger("tibco-activity-graphquerier")

type GraphQuerierActivity struct {
	metadata *activity.Metadata
	mux      sync.Mutex
}

func NewActivity(metadata *activity.Metadata) activity.Activity {
	return &GraphQuerierActivity{
		metadata: metadata,
	}
}

func (a *GraphQuerierActivity) Metadata() *activity.Metadata {
	return a.metadata
}

func (a *GraphQuerierActivity) Eval(context activity.Context) (done bool, err error) {
	log.Info("(Eval) GraphQuerier entering ......... ")

	a.mux.Lock()
	defer a.mux.Unlock()

	graphId := context.GetInput("graphId")
	if nil == graphId {
		return false, fmt.Errorf("Unable to get graphId from input!!")
	}
	/* UpsertGraph */
	graph := model.GetGraphManager().GetGraph(graphId.(string), graphId.(string))

	queryId := context.GetInput("queryId")
	if nil == queryId {
		return false, fmt.Errorf("Unable to get queryId from input!!")
	}

	rawQuery, exists := context.GetInput("Query").(map[string]interface{})
	if !exists {
		return false, fmt.Errorf("Unable to get data from input!!")
	}
	/* Register Query */
	query, err := query.GetQueryManager().CreateQuery(
		graphId.(string),
		queryId.(string),
		rawQuery,
	)

	if nil != err {
		return false, err
	}

	log.Info("(GraphQuerierActivity::Eval) queryId = ", queryId, ", graphId = ", graphId)

	matchedData := make(map[string]interface{})
	matchedData[queryId.(string)] = query.Search(graph)

	context.SetOutput("MatchedData", matchedData)

	log.Info("(Eval) GraphQuerier exit ......... ")
	return true, nil
}

func (a *GraphQuerierActivity) loadTest(graphId string, context activity.Context) (map[string]*query.Query, error) {
	queries := query.GetQueryManager().GetQueries(graphId)
	if nil == queries {
		a.mux.Lock()
		defer a.mux.Unlock()
		queries := query.GetQueryManager().GetQueries(graphId)
		if nil == queries {
			filterFile, exists := context.GetSetting(FILTER_FILE)
			if !exists {
				return nil, fmt.Errorf("No default filter defined!")
			}
			filterString, err := util.ReadFile(filterFile.(string))
			if nil != err {
				return nil, err
			}
			var queryObject interface{}

			err = json.Unmarshal([]byte(filterString), &queryObject)
			if nil != err {
				return nil, err
			}

			query.GetQueryManager().CreateQuery(graphId, "DEFAULT_QUERY", queryObject.(map[string]interface{}))
			queries = query.GetQueryManager().GetQueries(graphId)
		}

	}

	return queries, nil
}
