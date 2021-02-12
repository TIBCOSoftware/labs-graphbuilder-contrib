/*
 * Copyright © 2020. TIBCO Software Inc.
 * This file is subject to the license terms contained
 * in the license file that is distributed with this file.
 */
package model

import (
	"sync"
)

type GraphType int

const (
	GRAPH  GraphType = 0
	TGRAPH GraphType = 1
)

func (this GraphType) int() int {
	index := [...]int{0, 1, 2}
	return index[this]
}

type ManagableGraph interface {
}

type GraphManager struct {
	graphs map[string]ManagableGraph
}

var (
	instance *GraphManager
	once     sync.Once
	mux      sync.Mutex
)

func GetGraphManager() *GraphManager {
	once.Do(func() {
		instance = &GraphManager{
			graphs: make(map[string]ManagableGraph),
		}
	})
	return instance
}

func (this *GraphManager) GetGraph(
	graphType GraphType,
	modelId string,
	graphId string,
) ManagableGraph {

	graph := this.graphs[graphId]
	if nil == graph {
		mux.Lock()
		defer mux.Unlock()
		graph = this.graphs[graphId]
		if nil == graph {
			graph = CreateGraph(graphType, modelId, graphId)
			this.graphs[graphId] = graph
		}
	}

	return graph
}

func CreateGraph(
	graphType GraphType,
	modelId string,
	graphId string,
) ManagableGraph {
	//fmt.Println("Create new graph, modelId = ", modelId, ", graphId = ", graphId)
	var graph ManagableGraph
	switch graphType {
	case GRAPH:
		graph = NewTraversalGraph(modelId, graphId)
	case TGRAPH:
		graph = NewGraphImpl(modelId, graphId)
	}
	return graph
}
