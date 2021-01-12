/*
 * Copyright © 2020. TIBCO Software Inc.
 * This file is subject to the license terms contained
 * in the license file that is distributed with this file.
 */
package query

import (
	"sync"
)

type QueryManager struct {
	queries map[string](map[string]*Query)
}

var (
	instance *QueryManager
	once     sync.Once
	mux      sync.Mutex
)

func GetQueryManager() *QueryManager {
	once.Do(func() {
		instance = &QueryManager{queries: make(map[string](map[string]*Query))}
	})
	return instance
}

func (this *QueryManager) GetQuery(graphId string, queryId string) *Query {
	if nil != this.queries[graphId] {
		return this.queries[graphId][queryId]
	}
	return nil
}

func (this *QueryManager) GetQueries(graphId string) map[string]*Query {
	return this.queries[graphId]
}

func (this *QueryManager) CreateQuery(
	graphId string,
	queryId string,
	queryDef map[string]interface{}) (*Query, error) {

	//	fmt.Println("(QueryManager::CreateQuery) - new query = ", queryDef)

	queriesByGraph := this.queries[graphId]
	if nil == queriesByGraph {
		queriesByGraph = make(map[string]*Query)
		this.queries[graphId] = queriesByGraph
	}
	query := queriesByGraph[queryId]
	if nil == query {
		mux.Lock()
		defer mux.Unlock()
		query = queriesByGraph[queryId]
		if nil == query {
			var err error
			query, err = NewQuery(this, graphId, queryDef)
			if nil != err {
				return nil, err
			}
			queriesByGraph[queryId] = query
		}
	}

	return query, nil
}
