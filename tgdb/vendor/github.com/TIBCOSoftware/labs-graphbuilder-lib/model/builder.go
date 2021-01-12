/*
 * Copyright © 2020. TIBCO Software Inc.
 * This file is subject to the license terms contained
 * in the license file that is distributed with this file.
 */
package model

import (
	"strings"
	"sync"

	"github.com/TIBCOSoftware/labs-graphbuilder-lib/util"
)

//-============================-//
//     Define GraphBuilder
//-============================-//

type GraphBuilder struct {
	mux sync.Mutex
}

func NewGraphBuilder() *GraphBuilder {
	builder := &GraphBuilder{}
	return builder
}

func (builder *GraphBuilder) CreateGraph(graphId string, model *GraphDefinition) Graph {
	graph := NewGraphImpl(model.GetId(), graphId)
	graph.SetModel(model.Export())
	return graph
}

func (builder *GraphBuilder) CreateUndefinedGraph(modelId string, graphId string) Graph {
	graph := NewGraphImpl(modelId, graphId)
	return graph
}

func (builder *GraphBuilder) BuildGraph(
	graph *Graph,
	model *GraphDefinition,
	nodes interface{},
	edges interface{},
	allowNullKey bool,
) {

	//fmt.Println("[GraphBuilder:BuildGraph] entering ........ ")
	//	fmt.Println("[GraphBuilder:BuildGraph] nodes : ", nodes)
	//	fmt.Println("[GraphBuilder:BuildGraph] edges : ", edges)

	nodesWrapper, nodesWrapperValid := nodes.([]interface{})
	nodeMap := make(map[string]*Node)
	if nodesWrapperValid {
		if nil != nodesWrapper[0] {
			nodes, nodesValid := nodesWrapper[0].(map[string]interface{})
			if nodesValid {
				for nodeConfKey, node := range nodes {

					//					fmt.Println("nodekey before = ", nodekey)
					var nodeType string
					checkPos := strings.LastIndex(nodeConfKey, "_")
					if 0 <= checkPos && util.IsInteger(nodeConfKey[checkPos+1:len(nodeConfKey)]) {
						nodeType = string(nodeConfKey[0:checkPos])
					} else {
						nodeType = nodeConfKey
					}
					//					fmt.Println("nodekey after = " + nodekey)

					var attrs map[string]interface{}
					attrWrapper, attrWrapperValid := node.([]interface{})
					if attrWrapperValid {
						if nil != attrWrapper[0] {
							attrs = attrWrapper[0].(map[string]interface{})
						}
					}
					if nil == attrs {
						attrs = make(map[string]interface{})
					}

					_skipCondition := false
					if nil != attrs["_skipCondition"] {
						_skipCondition = attrs["_skipCondition"].(bool)
						delete(attrs, "_skipCondition")
					}

					if !_skipCondition {
						node := builder.BuildNode(graph, model, nodeType, attrs, allowNullKey)
						if nil != node {
							nodeMap[nodeConfKey] = node
						}
					}
				}
			}
		}
	}

	edgesWrapper, edgesWrapperValid := edges.([]interface{})
	//fmt.Println("edgesWrapper = ", edgesWrapper, ", edgesWrapperValid = ", edgesWrapperValid)
	if edgesWrapperValid {
		if nil != edgesWrapper[0] {
			edges, edgesValid := edgesWrapper[0].(map[string]interface{})
			//fmt.Println("edges = ", edges, ", edgesValid = ", edgesValid)
			if edgesValid {
				for edgeConfKey, edge := range edges {
					var attrs map[string]interface{}
					attrWrapper, attrWrapperValid := edge.([]interface{})
					//fmt.Println("attrWrapper = ", attrWrapper, ", attrWrapperValid = ", attrWrapperValid)
					if attrWrapperValid && 0 < len(attrWrapper) && nil != attrWrapper[0] {
						attrs = attrWrapper[0].(map[string]interface{})
					}
					if nil == attrs {
						attrs = make(map[string]interface{})
					}

					//fmt.Println("attrs = ", attrs)

					_skipCondition := false
					if nil != attrs["_skipCondition"] {
						_skipCondition = attrs["_skipCondition"].(bool)
						delete(attrs, "_skipCondition")
					}

					/* attribute could include "vertices" */
					if !_skipCondition {
						builder.BuildEdge(graph, model, nodeMap, edgeConfKey, attrs, allowNullKey)
					}
				}
			}
		}
	}
	//fmt.Println("[GraphBuilder:BuildGraph] exit ........ ")

}

func (builder *GraphBuilder) BuildNode(
	graph *Graph,
	model *GraphDefinition,
	nodeType string,
	attributesInfo map[string]interface{},
	allowNullKey bool,
) *Node {

	//fmt.Println("[GraphBuilder:BuildNode] entering ........ ")

	builder.mux.Lock()
	defer builder.mux.Unlock()

	nodeModel := model.GetNodeDefinition(nodeType)
	keyDefinition := nodeModel._keyDefinition
	key := make([]interface{}, len(keyDefinition))
	for i := 0; i < len(keyDefinition); i++ {
		attr := attributesInfo[keyDefinition[i]]
		if !allowNullKey && nil == attr {
			/* null key is not allowed */
			return nil
		}
		key[i] = attr
	}

	node := (*graph).UpsertNode(nodeType, key)
	for attrKey, attrVal := range attributesInfo {
		attribute := NewAttribute(nodeModel._attributes[attrKey], attrVal)
		node._attributes[attribute.GetName()] = attribute
	}

	//fmt.Println("[GraphBuilder:BuildNode] exit ........ ")

	return node
}

func (builder *GraphBuilder) BuildEdge(
	graph *Graph,
	model *GraphDefinition,
	nodeMap map[string]*Node,
	edgeType string,
	attributesInfo map[string]interface{},
	allowNullKey bool, /* Not in use */
) {

	//fmt.Println("[GraphBuilder:BuildEdge] entering, edgeType = ", edgeType)

	builder.mux.Lock()
	defer builder.mux.Unlock()

	edgeModel := model.GetEdgeDefinition(edgeType)
	keyDefinition := edgeModel._keyDefinition

	var verticesValues map[string]interface{}
	if nil != attributesInfo["vertices"] {
		rawVerticesArray := attributesInfo["vertices"].([]interface{})
		if 0 < len(rawVerticesArray) {
			verticesValues = rawVerticesArray[0].(map[string]interface{})
		}

		delete(attributesInfo, "vertices")
	}

	fromNodes, toNodes := builder.buildVerexes(graph, model, nodeMap, edgeType, verticesValues)

	var fromNode *Node
	var toNode *Node

	for _, fromNode = range fromNodes {
		for _, toNode = range toNodes {
			/* Allow duplicate? */
			key := make([]interface{}, len(keyDefinition))
			for i := 0; i < len(keyDefinition); i++ {
				key[i] = attributesInfo[keyDefinition[i]]
			}
			//fmt.Println("[GraphBuilder:BuildEdge] graph.UpsertEdge : ", edgeType, ", key = ", key)
			edge := (*graph).UpsertEdge(edgeType, key, fromNode, toNode)
			for attrKey, attrVal := range attributesInfo {
				attribute := NewAttribute(edgeModel._attributes[attrKey], attrVal)
				edge._attributes[attribute.GetName()] = attribute
			}
		}
	}

	//fmt.Println("[GraphBuilder:BuildEdge] exit ........ ")

}

func (builder *GraphBuilder) buildVerexes(
	graph *Graph,
	model *GraphDefinition,
	nodeMap map[string]*Node,
	edgeType string,
	verticesConfKey map[string]interface{}) (map[NodeId]*Node, map[NodeId]*Node) {

	//fmt.Println("[GraphBuilder:BuildVertices] entering ........ ")

	edgeModel := model.GetEdgeDefinition(edgeType)

	var fromNodes map[NodeId]*Node
	var toNodes map[NodeId]*Node

	/* verticesConfKey = map[from:airport_0 to:airport_0] */
	fromNodeConfKey := verticesConfKey["from"]
	toNodeConfKey := verticesConfKey["to"]

	//fmt.Println("[GraphBuilder:BuildVertices] fromNodeConfKey = ", fromNodeConfKey, ", toNodeConfKey = ", toNodeConfKey)
	//fmt.Println("[GraphBuilder:BuildVertices] graph = ", graph)
	//fmt.Println("[GraphBuilder:BuildVertices] edgeModel._fromNodeType = ", edgeModel._fromNodeType, ", edgeModel._toNodeType = ", edgeModel._toNodeType)

	if nil != fromNodeConfKey {
		fromNodes = make(map[NodeId]*Node)
		fromNode := nodeMap[fromNodeConfKey.(string)]
		fromNodes[fromNode.NodeId] = fromNode
	} else {
		fromNodes = (*graph).GetNodesByType(edgeModel._fromNodeType)
	}

	if nil != toNodeConfKey {
		toNodes = make(map[NodeId]*Node)
		toNode := nodeMap[toNodeConfKey.(string)]
		toNodes[toNode.NodeId] = toNode
	} else {
		toNodes = (*graph).GetNodesByType(edgeModel._toNodeType)
	}

	//fmt.Println("[GraphBuilder:BuildVertices] exit, fromNodes = ", fromNodes, ", toNodes = ", toNodes)

	return fromNodes, toNodes
}

func (builder *GraphBuilder) Export(g *Graph, graphModel *GraphDefinition) map[string]interface{} {

	nodeDefinitions := graphModel._nodeDefinitions
	edgeDefinitions := graphModel._edgeDefinitions

	data := make(map[string]interface{})
	data["id"] = (*g).GetId()
	data["modelId"] = (*g).GetModelId()
	data["model"] = graphModel.Export()

	nodesData := make(map[string]interface{})
	for nodeId, node := range (*g).GetNodes() {
		nodeData := make(map[string]interface{})
		attrsData := make(map[string]interface{})
		for attrName, attribute := range node._attributes {
			attrData := make(map[string]interface{})
			attrData["name"] = attribute._name
			attrData["value"] = attribute._value
			attrData["type"] = attribute._type.String()
			attrsData[attrName] = attrData
		}
		nodeData["type"] = node._type
		nodeData["keyAttributeName"] = nodeDefinitions[node._type]._keyDefinition
		nodeData["key"] = node._key
		nodeData["attributes"] = attrsData
		nodesData[nodeId.ToString()] = nodeData
	}
	data["nodes"] = nodesData

	edgesData := make(map[string]interface{})
	for edgeId, edge := range (*g).GetEdges() {
		edgeData := make(map[string]interface{})
		attrsData := make(map[string]interface{})
		for attrName, attribute := range edge._attributes {
			attrData := make(map[string]interface{})
			attrData["name"] = attribute._name
			attrData["value"] = attribute._value
			attrData["type"] = attribute._type.String()
			attrsData[attrName] = attrData
		}
		edgeData["type"] = edge._type
		edgeData["from"] = edge._fromNodeId.ToString()
		edgeData["to"] = edge._toNodeId.ToString()
		edgeData["keyAttributeName"] = edgeDefinitions[edge._type]._keyDefinition
		edgeData["key"] = edge._key
		edgeData["attributes"] = attrsData
		edgesData[edgeId.ToString()] = edgeData
	}
	data["edges"] = edgesData

	//log.Debug("[GraphBuilder::Export] graph : ", data)

	return data
}

func ReconstructGraph(graphData map[string]interface{}) Graph {

	graph := NewGraphImpl(graphData["modelId"].(string), graphData["id"].(string))
	graph.SetModel(graphData["model"].(map[string]interface{}))

	nodes := util.CastGenMap(graphData["nodes"])
	for _, value := range nodes {
		nodeData := util.CastGenMap(value)
		node := NewNode(
			util.CastString(nodeData["type"]),
			util.CastGenArray(nodeData["key"]),
		)

		attributes := util.CastGenMap(nodeData["attributes"])
		for attrName, value := range attributes {
			attrData := util.CastGenMap(value)
			dataType, ok := ToTypeEnum(util.CastString(attrData["type"]))
			if !ok {
				dataType = TypeString
			}
			attribute := Attribute{
				_name:  util.CastString(attrData["name"]),
				_value: attrData["value"],
				_type:  dataType,
			}
			node.SetAttribute(attrName, &attribute)
		}
		graph.SetNode(node.NodeId, node)
	}

	//fmt.Println("Graph : ", graph.GetNodes())

	edges := util.CastGenMap(graphData["edges"])
	for _, value := range edges {
		edgeData := util.CastGenMap(value)
		fromId := *(&NodeId{}).FromString(util.CastString(edgeData["from"]))
		toId := *(&NodeId{}).FromString(util.CastString(edgeData["to"]))
		//fmt.Println("EdgeType : ", edgeData["type"], ", fromID = ", fromId, ", toID = ", toId)
		edge := NewEdge(
			util.CastString(edgeData["type"]),
			util.CastGenArray(edgeData["key"]),
			graph.GetNode(fromId),
			graph.GetNode(toId),
		)

		attributes := util.CastGenMap(edgeData["attributes"])
		for attrName, value := range attributes {
			attrData := util.CastGenMap(value)
			dataType, ok := ToTypeEnum(util.CastString(attrData["type"]))
			if !ok {
				dataType = TypeString
			}
			attribute := Attribute{
				_name:  util.CastString(attrData["name"]),
				_value: attrData["value"],
				_type:  dataType,
			}
			edge.SetAttribute(attrName, &attribute)
		}
		graph.SetEdge(edge.EdgeId, edge)
	}

	return graph
}
