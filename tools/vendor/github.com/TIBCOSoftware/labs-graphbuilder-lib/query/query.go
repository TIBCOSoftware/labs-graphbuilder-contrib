/*
 * Copyright © 2020. TIBCO Software Inc.
 * This file is subject to the license terms contained
 * in the license file that is distributed with this file.
 */
package query

import (
	"reflect"
	"time"

	"github.com/TIBCOSoftware/labs-graphbuilder-lib/model"
	"github.com/TIBCOSoftware/labs-graphbuilder-lib/tools"
	"github.com/TIBCOSoftware/labs-graphbuilder-lib/util"
)

//-====================-//
//  Define COMPTEntity
//-====================-//

type COMPTEntity interface {
}

//-====================-//
//  Define COMPTResult
//-====================-//
type Scope struct {
	index    int
	maxIndex int
	name     string
	array    bool
}

type COMPTResult struct {
	inScope bool
	/* 0: out, 1: enter, 2: remain in, 3: leave */
	liveCycleFlag   int
	currentLocation []Scope
	path            []COMPTEntity
	rtns            map[string][]interface{}
	data            map[string]interface{}
}

func (this *COMPTResult) SetInScope(inScope bool) {
	this.inScope = inScope
}

func (this *COMPTResult) InScope() bool {
	return this.inScope
}

func (this *COMPTResult) SetLiveCycleFlag(liveCycleFlag int) {
	this.liveCycleFlag = liveCycleFlag
}

func (this *COMPTResult) GetLiveCycleFlag() int {
	return this.liveCycleFlag
}

func (this *COMPTResult) AddEntity(query COMPTEntity) {
	this.path = append(this.path, query)
}

func (this *COMPTResult) PickupData(entity model.Entity, compEntity BaseCOMPTEntity) {
	//	fmt.Println("Entity type -> ", compEntity.entityType, ", returns from down stream = ", this.rtns)
	if 0 != len(compEntity.returns) {
		for id, attrName := range compEntity.returns {
			//			fmt.Println("   id = ", id, ", attrName = ", attrName, ", attr = ", entity.GetAttribute(attrName))
			var attrValue interface{}
			if nil != entity.GetAttribute(attrName) {
				attrValue = entity.GetAttribute(attrName).GetValue()
			}
			this.SetRtnData(id, attrValue)
		}
	}
	//	fmt.Println("Entity type -> ", compEntity.entityType, ", this.returns = ", this.rtns)
}

func (this *COMPTResult) RecordData(id string, value interface{}) {
	//for index, scope := range this.currentLocation {

	//}

	if nil == this.rtns[id] {
		this.rtns[id] = make([]interface{}, 1)
		this.rtns[id][0] = value
	} else {
		this.rtns[id] = append(this.rtns[id], value)
	}
}

func (this *COMPTResult) SetRtnData(id string, value interface{}) {
	if nil == this.rtns[id] {
		this.rtns[id] = make([]interface{}, 1)
		this.rtns[id][0] = value
	} else {
		this.rtns[id] = append(this.rtns[id], value)
	}
}

func (this *COMPTResult) GetRtnDatas() map[string][]interface{} {
	this.rtns["_currentTime"][0] = time.Unix(tools.GetClock().GetCurrentTime(), 0)
	this.rtns["_liveCycleFlag"][0] = this.liveCycleFlag
	//	fmt.Println("Return : ", this.rtns)

	return this.rtns
}

func (this *COMPTResult) Visited(targetEntity COMPTEntity) bool {
	for _, entity := range this.path {
		if entity == targetEntity {
			return true
		}
	}
	return false
}

func (this *COMPTResult) EnterScope(scopename string, isArray bool) {
	this.currentLocation = append(this.currentLocation, Scope{name: scopename, array: isArray, index: -1, maxIndex: -1})
	//	fmt.Println("After enterScope : ", this.currentLocation) //, ", index : ", this.namespace[len(this.namespace)-1].index)
}

func (this *COMPTResult) LeaveScope(scopename string, isArray bool) {
	//	fmt.Println("Before leaveScope : ", this.currentLocation) //, ", index : ", this.namespace[len(this.namespace)-1].index)
	this.currentLocation = this.currentLocation[:len(this.currentLocation)-1]
}

func NewCOMPTResult() *COMPTResult {
	result := &COMPTResult{}
	result.inScope = false
	result.currentLocation = make([]Scope, 0)
	result.path = make([]COMPTEntity, 0)
	result.rtns = make(map[string][]interface{})
	result.rtns["_currentTime"] = make([]interface{}, 1)
	result.rtns["_liveCycleFlag"] = make([]interface{}, 1)
	result.data = make(map[string]interface{})
	return result
}

//-====================-//
//  Define COMPTEntity
//-====================-//

type BaseCOMPTEntity struct {
	Expression
	entityType string
	returns    map[string]string
}

func (this *BaseCOMPTEntity) GetReturn() map[string]string {
	return this.returns
}

//-====================-//
//   Define COMPTNode
//-====================-//

type COMPTNode struct {
	BaseCOMPTEntity
	outEdgesByType map[string]([]*COMPTEdge)
	inEdgesByType  map[string]([]*COMPTEdge)
}

func (this *COMPTNode) SetEdge(edge *COMPTEdge, outbound bool) {
	edgeType := edge.entityType
	if outbound {
		this.outEdgesByType[edgeType] = append(this.outEdgesByType[edgeType], edge)
	} else {
		this.inEdgesByType[edgeType] = append(this.inEdgesByType[edgeType], edge)
	}
}

func (this *COMPTNode) GetEdges() []*COMPTEdge {
	edges := make([]*COMPTEdge, 0)
	for _, outEdges := range this.outEdgesByType {
		for _, outEdge := range outEdges {
			edges = append(edges, outEdge)
		}
	}

	for _, inEdges := range this.inEdgesByType {
		for _, inEdge := range inEdges {
			edges = append(edges, inEdge)
		}
	}

	return edges
}

func (this *COMPTNode) Match(result COMPTResult, data *model.TraversalNode) bool {
	/* Match local */
	if this.entityType != data.GetType() {
		return false
	}

	//	fmt.Println("\n\n******** Check Node begin, Type -> ", this.entityType, ", Filter Expression -> ", this.Expression.ToString())
	//	defer fmt.Println("******** Check Node end, Type -> ", this.entityType, "\n\n")

	if !this.Eval(data.GetAttributes()) {
		//		fmt.Println("Node Expression Eval -> Failed!")
		return false
	}

	result.AddEntity(this)
	result.EnterScope(this.entityType, true)
	defer result.LeaveScope(this.entityType, true)

	/* Match downstream */
	if !matchDownStream(result, this.outEdgesByType, data, true) {
		return false
	}

	if !matchDownStream(result, this.inEdgesByType, data, false) {
		return false
	}

	/* pick up data only if all downstream condition matched */
	result.PickupData(data, this.BaseCOMPTEntity)

	return true
}

func matchDownStream(result COMPTResult, targetEdges map[string]([]*COMPTEdge), data *model.TraversalNode, outbound bool) bool {
	for targetEdgeType, targetEdgeGroup := range targetEdges {
		edges := data.GetEdgeByType(targetEdgeType, outbound)
		for _, targetEdge := range targetEdgeGroup {
			if result.Visited(targetEdge) {
				continue
			}

			atLeastOneMatched := false
			for _, edge := range edges {
				if targetEdge.Match(result, edge) {
					atLeastOneMatched = true
				}
			}
			if !atLeastOneMatched {
				return false
			}
		}
	}

	return true
}

func NewCOMPTNode(query map[string]interface{}) *COMPTNode {
	var node COMPTNode
	//node.BuildBasic(query)
	node.entityType = query["name"].(string)
	node.outEdgesByType = make(map[string]([]*COMPTEdge))
	node.inEdgesByType = make(map[string]([]*COMPTEdge))
	node.returns = make(map[string]string)

	node.returns = make(map[string]string)
	if nil != query["rtn"] {
		rtns := query["rtn"].(map[string]interface{})
		for id, attrName := range rtns {
			node.returns[id] = attrName.(string)[1:]
		}
	}
	//	fmt.Println("Type -> ", node.entityType, ", this.returns = ", node.returns)

	if nil != query["expr"] {
		expr := query["expr"].(map[string]interface{})
		if 0 == len(expr) {
			node.operator = CreateOperator("nop")
		} else {
			node.Build(expr)
		}
	} else {
		node.operator = CreateOperator("nop")
	}

	return &node
}

//-====================-//
//   Define COMPTEdge
//-====================-//

type COMPTEdge struct {
	BaseCOMPTEntity
	fromNode *COMPTNode
	toNode   *COMPTNode
}

func (this *COMPTEdge) GetNodes() []*COMPTNode {
	nodes := make([]*COMPTNode, 2)
	nodes[0] = this.fromNode
	nodes[1] = this.toNode
	return nodes
}

func (this *COMPTEdge) Match(result COMPTResult, edge *model.TraversalEdge) bool {
	/* Match local */
	if this.entityType != edge.GetType() {
		return false
	}

	//	fmt.Println("Edge Expression -> ", this.Expression.ToString())
	if !this.Eval(edge.GetAttributes()) {
		//		fmt.Println("Edge Expression Eval -> Failed!")
		return false
	}

	result.AddEntity(this)
	defer result.PickupData(edge, this.BaseCOMPTEntity)

	/* Match downstream */
	if !result.Visited(this.toNode) {
		return this.toNode.Match(result, edge.GetToNode())
	} else {
		return this.fromNode.Match(result, edge.GetFromNode())
	}
}

func NewCOMPTEdge(query map[string]interface{}, fromNode *COMPTNode, toNode *COMPTNode) *COMPTEdge {
	edge := &COMPTEdge{}
	//edge.BuildBasic(query)
	edge.entityType = query["name"].(string)
	edge.fromNode = fromNode
	fromNode.SetEdge(edge, true)
	edge.toNode = toNode
	toNode.SetEdge(edge, false)

	edge.returns = make(map[string]string)
	if nil != query["rtn"] {
		rtns := query["rtn"].(map[string]interface{})
		for id, attrName := range rtns {
			edge.returns[id] = attrName.(string)[1:]
		}
	}
	//	fmt.Println("Type -> ", edge.entityType, ", this.returns = ", edge.returns)

	if nil != query["expr"] {
		expr := query["expr"].(map[string]interface{})
		if 0 == len(expr) {
			edge.operator = CreateOperator("nop")
		} else {
			edge.Build(expr)
		}
	} else {
		edge.operator = CreateOperator("nop")
	}

	return edge
}

//-==================-//
//   Define Query
//-==================-//

type Query struct {
	manager   *QueryManager
	graphId   string
	id        string
	parameter map[string]interface{}
	planners  map[string]*QueryPlanner
	entryNode *COMPTNode
	results   map[model.NodeId]*COMPTResult
}

func (this *Query) Search(graph *model.TraversalGraph) map[string]interface{} {
	output := make(map[string]interface{})
	for nodeId, tNode := range graph.GetNodesByType(this.entryNode.entityType) {
		//		fmt.Println("(Query::Match) - nodeId = ", nodeId, ", tNode = ", tNode)
		//		printTraversal(tNode)
		result := this.results[nodeId]
		if nil == result {
			result = NewCOMPTResult()
		}
		//		fmt.Println("\n\n*************** Start matching ******************")
		result.SetInScope(this.entryNode.Match(*result, tNode))
		//		fmt.Println("*************** End matching ******************\n\n")
		if result.InScope() {
			result.SetLiveCycleFlag(1)
			output["parameter"] = this.parameter
			output[nodeId.ToString()] = result.GetRtnDatas()
		}
	}
	return output
}

func (this *Query) Match(sourceId string, graph *model.TraversalGraph, deltaGraph model.Graph) map[string]interface{} {
	output := make(map[string]interface{})
	//fmt.Println("(Query::Match) - this.entryNode.entityType = ", this.entryNode.entityType, ", this.entryNode = ", this.entryNode)
	planner := this.planners[sourceId]
	if nil == planner {
		planner = NewQueryPlanner()
		this.planners[sourceId] = planner
	}

	entryNode, tNodes := planner.GetEntry(deltaGraph, this.entryNode)
	//	fmt.Println("(Query::Match) - entryNode = ", entryNode, ", tNodes = ", tNodes)
	for nodeId, _ := range tNodes {
		tNode := graph.GetNode(nodeId)
		//		fmt.Println("(Query::Match) - nodeId = ", nodeId, ", tNode = ", tNode)
		//printTraversal(tNode)
		result := this.results[nodeId]
		if nil == result {
			result = NewCOMPTResult()
		}
		previousInScope := result.InScope()
		result.SetInScope(entryNode.Match(*result, tNode))
		switch previousInScope {
		case true:
			switch result.InScope() {
			case true:
				result.SetLiveCycleFlag(2)
			case false:
				result.SetLiveCycleFlag(3)
				delete(output, nodeId.ToString())
				output[nodeId.ToString()] = result.GetRtnDatas()
			}
		case false:
			switch result.InScope() {
			case true:
				result.SetLiveCycleFlag(1)
				output[nodeId.ToString()] = result.GetRtnDatas()
			case false:
				result.SetLiveCycleFlag(0)
			}
		}
		//		fmt.Println("(Query::Match) - Result lifecycle flag --> ", result.GetLiveCycleFlag())
	}

	if 0 != len(output) {
		output["parameter"] = this.parameter
	}

	return output
}

func NewQuery(manager *QueryManager, graphId string, targetMap map[string]interface{}) (*Query, error) {
	//	fmt.Println("targetMap = ", targetMap)
	query := &Query{
		manager: manager,
		graphId: graphId,
	}

	if nil != targetMap["agg"] {
		parameter := make(map[string]interface{})
		parameter["agg"] = targetMap["agg"].(map[string]interface{})
		if nil != targetMap["key"] {
			parameter["key"] = targetMap["key"].([]interface{})
		}
		if nil != targetMap["index"] {
			parameter["index"] = targetMap["index"].([]interface{})
		}
		query.parameter = parameter
	}

	//	fmt.Println("query = ", query)

	query.planners = make(map[string]*QueryPlanner)
	targetNodes := make(map[string]*COMPTNode)
	nodes := targetMap["nodes"].([]interface{})
	//fmt.Println("nodes = ", nodes)
	for index, node := range nodes {
		nodeInfo := node.(map[string]interface{})
		nodeType := nodeInfo["name"].(string)
		targetNodes[nodeType] = NewCOMPTNode(nodeInfo)
		if 0 == index {
			/* the first node is the entry node */
			query.entryNode = targetNodes[nodeType]
		}
	}

	edges := targetMap["edges"].([]interface{})
	//fmt.Println("edges = ", edges)
	for _, edge := range edges {
		edgeInfo := edge.(map[string]interface{})
		NewCOMPTEdge(edgeInfo, targetNodes[edgeInfo["from"].(string)], targetNodes[edgeInfo["to"].(string)])
	}
	//fmt.Println("targetNodes = ", targetNodes)

	return query, nil
}

//-=======================-//
//   Define QueryPlanner
//-=======================-//

type QueryPlanner struct {
	entryNodes []*COMPTNode
}

func (this *QueryPlanner) GetEntryNodes() []*COMPTNode {
	return this.entryNodes
}

func (this *QueryPlanner) GetEntry(deltaGraph model.Graph, definedEntryNode *COMPTNode) (*COMPTNode, map[model.NodeId]*model.Node) {
	if 0 == len(this.entryNodes) {
		this.findEntry(deltaGraph, make([]interface{}, 0), definedEntryNode)
	}

	for _, entryNode := range this.entryNodes {
		if nil != deltaGraph.GetNodesByType(entryNode.entityType) {
			return entryNode, deltaGraph.GetNodesByType(entryNode.entityType)
		}
	}
	return nil, nil
}

func (this *QueryPlanner) findEntry(deltaGraph model.Graph, path []interface{}, entity interface{}) {
	//	fmt.Println("Entity kind = ", reflect.TypeOf(entity).String(), ", path = ", path)
	if "*query.COMPTNode" == reflect.TypeOf(entity).String() {
		node := entity.(*COMPTNode)
		path = append(path, node)
		for _, edge := range node.GetEdges() {
			//	fmt.Println("Entity type = ", edge.entityType)
			if nil != edge && !util.Contains(path, edge) {
				this.findEntry(deltaGraph, path, edge)
			}
		}
	} else if "*query.COMPTEdge" == reflect.TypeOf(entity).String() {
		edge := entity.(*COMPTEdge)
		path = append(path, edge)
		if nil != edge.GetNodes() {
			for _, node := range edge.GetNodes() {
				//	fmt.Println("Entity type = ", node.entityType)
				if nil != node && !util.Contains(path, node) {
					//	fmt.Println("pass Entity type = ", node.entityType, ", delta graph = ", deltaGraph)
					entryNodes := deltaGraph.GetNodesByType(node.entityType)
					if nil != entryNodes {
						this.entryNodes = append(this.entryNodes, node)
					}
					this.findEntry(deltaGraph, path, node)
				}
			}
		}
	}
}

func NewQueryPlanner() *QueryPlanner {
	planner := &QueryPlanner{entryNodes: make([]*COMPTNode, 0)}
	return planner
}
