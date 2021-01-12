/*
 * Copyright © 2020. TIBCO Software Inc.
 * This file is subject to the license terms contained
 * in the license file that is distributed with this file.
 */
package tgdb

import (
	"fmt"
	"strings"
	"tgdb"
)

type TGEntityWrapper struct {
	id    string
	isNew bool
}

func (this *TGEntityWrapper) String() string {
	return fmt.Sprintf("Node {id : %s, isNew : %s}", this.id, this.isNew)
}

type TGNodeWrapper struct {
	TGEntityWrapper
	node *tgdb.TGNode
}

func NewTGNodeWrapper(node tgdb.TGNode, isNew bool) *TGNodeWrapper {
	nodeWrapper := TGNodeWrapper{}
	nodeWrapper.isNew = isNew
	nodeWrapper.SetNode(&node)
	return &nodeWrapper
}

func (this *TGNodeWrapper) SetNode(node *tgdb.TGNode) {
	this.node = node
}

func (this *TGNodeWrapper) GetNode() *tgdb.TGNode {
	return this.node
}

func (this *TGNodeWrapper) Info() string {
	node := *this.node
	return fmt.Sprintf("%s : a node", node.GetEntityType().GetName())
}

type TGEdgeWrapper struct {
	TGEntityWrapper
	edge *tgdb.TGEdge
}

func NewTGEdgeWrapper(edge tgdb.TGEdge, isNew bool) *TGEdgeWrapper {
	edgeWrapper := TGEdgeWrapper{}
	edgeWrapper.isNew = isNew
	edgeWrapper.SetEdge(&edge)
	return &edgeWrapper
}

func (this *TGEdgeWrapper) SetEdge(edge *tgdb.TGEdge) {
	this.edge = edge
}

func (this *TGEdgeWrapper) GetEdge() *tgdb.TGEdge {
	return this.edge
}

func (this *TGEdgeWrapper) Info() string {
	edge := *this.edge
	return fmt.Sprintf(
		"%s : %s - %s",
		edge.GetEntityType().GetName(),
		edge.GetVertices()[0].GetEntityType().GetName(),
		edge.GetVertices()[1].GetEntityType().GetName())
}

type EntityKeeper interface {
	Populate(tgConnection tgdb.TGConnection)
	AddNode(nodeTypeStr string, key interface{}, tgnode *TGNodeWrapper)
	AddEdge(edgeTypeStr string, key interface{}, tgedge *TGEdgeWrapper)
	GetNode(nodeTypeStr string, key interface{}) *TGNodeWrapper
	GetEdge(edgeTypeStr string, key interface{}) *TGEdgeWrapper
	GetNodes() []*TGNodeWrapper
	GetEdges() []*TGEdgeWrapper
	Clear()
	Print()
}

func NewUpsertEntityKeeper() EntityKeeper {
	keeper := &UpsertEntityKeeper{}
	keeper.nodes = make(map[string](map[interface{}]*TGNodeWrapper))
	keeper.edges = make(map[string](map[interface{}]*TGEdgeWrapper))

	return keeper
}

type UpsertEntityKeeper struct {
	ReadyEntityKeeper
}

func (this *UpsertEntityKeeper) Populate(tgConnection tgdb.TGConnection) {
	nodeWrappers := this.GetNodes()
	for _, nodeWrapper := range nodeWrappers {
		node := *nodeWrapper.node
		if nodeWrapper.isNew {
			logger.Info(fmt.Sprintf("Insert node -> %s", nodeWrapper.Info()))
			tgConnection.InsertEntity(node)
		} else {
			logger.Info(fmt.Sprintf("Update node -> %s", nodeWrapper.Info()))
			tgConnection.UpdateEntity(node)
		}
	}

	edgeWrappers := this.GetEdges()
	for _, edgeWrapper := range edgeWrappers {
		edge := *edgeWrapper.edge
		if edgeWrapper.isNew {
			logger.Info(fmt.Sprintf("Insert edge -> %s", edgeWrapper.Info()))
			tgConnection.InsertEntity(edge)
		} else {
			if !strings.HasSuffix(edge.GetEntityType().GetName(), "_event") {
				logger.Info(fmt.Sprintf("Update edge -> %s", edgeWrapper.Info()))
				tgConnection.UpdateEntity(edge)
			} else {
				logger.Info(fmt.Sprintf("Blocking update edge -> %s", edgeWrapper.Info()))
			}
		}
	}
	this.Clear()

}

func (this *UpsertEntityKeeper) AddNode(nodeTypeStr string, key interface{}, tgnode *TGNodeWrapper) {
	logger.Debug(fmt.Sprintf("(UpsertEntityKeeper.AddNode) %s, %s", nodeTypeStr, key))

	cachedNodePerType := this.nodes[nodeTypeStr]
	if nil == cachedNodePerType {
		cachedNodePerType = make(map[interface{}]*TGNodeWrapper)
		this.nodes[nodeTypeStr] = cachedNodePerType
	}
	cachedNodePerType[key] = tgnode
}

func (this *UpsertEntityKeeper) AddEdge(edgeTypeStr string, key interface{}, tgedge *TGEdgeWrapper) {
	logger.Debug(fmt.Sprintf("(UpsertEntityKeeper.AddEdge) %s, %s", edgeTypeStr, key))

	cachedEdgePerType := this.edges[edgeTypeStr]
	if nil == cachedEdgePerType {
		cachedEdgePerType = make(map[interface{}]*TGEdgeWrapper)
		this.edges[edgeTypeStr] = cachedEdgePerType
	}
	cachedEdgePerType[key] = tgedge
}

func NewDeleteEntityKeeper(deleteNode bool, deleteEdge bool) EntityKeeper {
	keeper := &DeleteEntityKeeper{}
	keeper.nodes = make(map[string](map[interface{}]*TGNodeWrapper))
	keeper.edges = make(map[string](map[interface{}]*TGEdgeWrapper))
	keeper.deleteNode = deleteNode
	keeper.deleteEdge = deleteEdge

	return keeper
}

type DeleteEntityKeeper struct {
	ReadyEntityKeeper
	deleteNode bool
	deleteEdge bool
}

func (this *DeleteEntityKeeper) Populate(tgConnection tgdb.TGConnection) {
	logger.Debug("(DeleteEntityKeeper.Populate) enter ...... ")
	defer logger.Debug("(DeleteEntityKeeper.Populate) exit ...... ")

	nodeWrappers := this.GetNodes()
	if this.deleteNode {
		for _, nodeWrapper := range nodeWrappers {
			node := *nodeWrapper.node
			if nodeWrapper.isNew {
				logger.Warn(fmt.Sprintf("Why delete new node -> %s, %s", node.GetEntityType().GetName(), node.String()))
				/* do nothing */
			} else {
				logger.Info(fmt.Sprintf("Delete node -> %s, %s", node.GetEntityType().GetName(), node.String()))
				tgConnection.DeleteEntity(node)
			}
		}
	}

	if this.deleteEdge {
		edgeWrappers := this.GetEdges()
		for _, edgeWrapper := range edgeWrappers {
			edge := *edgeWrapper.edge
			if edgeWrapper.isNew {
				logger.Warn(fmt.Sprintf("Why delete new edge -> %s, %s", edge.GetEntityType().GetName(), edge.String()))
				/* do nothing */
			} else {
				logger.Info(fmt.Sprintf("Delete edge -> %s, %s", edge.GetEntityType().GetName(), edge.String()))
				tgConnection.DeleteEntity(edge)
			}
		}
	}
	this.Clear()

}

func (this *DeleteEntityKeeper) AddNode(nodeTypeStr string, key interface{}, tgnode *TGNodeWrapper) {
	if !tgnode.isNew {
		logger.Debug(fmt.Sprintf("(DeleteEntityKeeper.AddNode) %s, %s", nodeTypeStr, key))
		cachedNodePerType := this.nodes[nodeTypeStr]
		if nil == cachedNodePerType {
			cachedNodePerType = make(map[interface{}]*TGNodeWrapper)
			this.nodes[nodeTypeStr] = cachedNodePerType
		}
		cachedNodePerType[key] = tgnode
	}
}

func (this *DeleteEntityKeeper) AddEdge(edgeTypeStr string, key interface{}, tgedge *TGEdgeWrapper) {
	if !tgedge.isNew {
		logger.Debug(fmt.Sprintf("(DeleteEntityKeeper.AddEdge) %s, %s", edgeTypeStr, key))
		cachedEdgePerType := this.edges[edgeTypeStr]
		if nil == cachedEdgePerType {
			cachedEdgePerType = make(map[interface{}]*TGEdgeWrapper)
			this.edges[edgeTypeStr] = cachedEdgePerType
		}
		cachedEdgePerType[key] = tgedge
	}
}

type ReadyEntityKeeper struct {
	nodes map[string](map[interface{}]*TGNodeWrapper)
	edges map[string](map[interface{}]*TGEdgeWrapper)
}

func (r *ReadyEntityKeeper) GetNode(nodeTypeStr string, key interface{}) *TGNodeWrapper {

	logger.Debug(fmt.Sprintf("[RedayEntityKeeper:getNode] Target node : %s, %s", nodeTypeStr, key))
	logger.Debug(fmt.Sprintf("[RedayEntityKeeper:getNode] Available ready node : %s", r.nodes))
	if val, ok := r.nodes[nodeTypeStr]; ok {
		logger.Debug(fmt.Sprintf("[RedayEntityKeeper:getNode] Available target type ready node : %s ", r.nodes[nodeTypeStr]))
		return val[key]
	}
	return nil
}

func (r *ReadyEntityKeeper) GetEdge(edgeTypeStr string, key interface{}) *TGEdgeWrapper {
	if val, ok := r.edges[edgeTypeStr]; ok {
		return val[key]
	}
	return nil
}

func (r *ReadyEntityKeeper) GetNodes() []*TGNodeWrapper {
	allNodes := make([]*TGNodeWrapper, 0)

	for _, nodesInType := range r.nodes {
		for _, node := range nodesInType {
			allNodes = append(allNodes, node)
		}
	}
	return allNodes
}

func (r *ReadyEntityKeeper) GetEdges() []*TGEdgeWrapper {
	allEdges := make([]*TGEdgeWrapper, 0)
	for _, edgesInType := range r.edges {
		for _, edge := range edgesInType {
			allEdges = append(allEdges, edge)
		}
	}
	return allEdges
}

func (r *ReadyEntityKeeper) Clear() {
	for nodesByType := range r.nodes {
		delete(r.nodes, nodesByType)
	}
	for edgesByType := range r.edges {
		delete(r.edges, edgesByType)
	}
}

func (r *ReadyEntityKeeper) Print() {
	logger.Debug("+++++++ ready nodes +++++++")
	logger.Debug(fmt.Sprintf("%s", r.nodes))
	logger.Debug("+++++++ reday edges +++++++")
	logger.Debug(fmt.Sprintf("%s", r.edges))
	logger.Debug("++++++++++++++++++++++++++++")
}
