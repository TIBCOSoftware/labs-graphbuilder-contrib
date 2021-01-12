/*
 * Copyright © 2020. TIBCO Software Inc.
 * This file is subject to the license terms contained
 * in the license file that is distributed with this file.
 */
package model

//-====================-//
// Define TraversalNode
//-====================-//

type TraversalNode struct {
	Node
	outEdges map[string](map[string]*TraversalEdge)
	inEdges  map[string](map[string]*TraversalEdge)
}

func (this *TraversalNode) SetEdge(edge *TraversalEdge, outbound bool) bool {
	if outbound {
		setEdgeInternal(this.outEdges, edge)
	} else {
		setEdgeInternal(this.inEdges, edge)
	}
	return true
}

func (this *TraversalNode) GetAllEdges() map[string]*TraversalEdge {
	idMap := make(map[string]*TraversalEdge)
	for edgeId, edge := range this.GetEdges(true) {
		idMap[edgeId] = edge
	}
	for edgeId, edge := range this.GetEdges(false) {
		idMap[edgeId] = edge
	}
	return idMap
}

func (this *TraversalNode) GetEdges(outbound bool) map[string]*TraversalEdge {
	if outbound {
		return idToEdgeMap(this.outEdges)
	} else {
		return idToEdgeMap(this.inEdges)
	}
}

func idToEdgeMap(edges map[string](map[string]*TraversalEdge)) map[string]*TraversalEdge {
	idMap := make(map[string]*TraversalEdge)
	for _, typedMap := range edges {
		for edgeId, edge := range typedMap {
			idMap[edgeId] = edge
		}
	}
	return idMap
}

func (this *TraversalNode) GetEdgeByType(edgeType string, outbound bool) map[string]*TraversalEdge {
	if outbound {
		return this.outEdges[edgeType]
	} else {
		return this.inEdges[edgeType]
	}
}

func setEdgeInternal(edgeGroup map[string](map[string]*TraversalEdge), edge *TraversalEdge) {
	edgeType := edge.GetType()
	typeEdges := edgeGroup[edgeType]
	if nil == typeEdges {
		typeEdges = make(map[string]*TraversalEdge)
		edgeGroup[edgeType] = typeEdges
	}
	typeEdges[edge.EdgeId.ToString()] = edge
}

func (this *TraversalNode) Update(node *Node) {
	for attrName, attr := range node.GetAttributes() {
		attrValue := attr.GetValue()
		if nil != attrValue {
			this.GetAttribute(attrName).SetValue(attrValue)
		}
	}
}

func NewTraversalNode(node *Node) *TraversalNode {
	var n TraversalNode
	n.Node = *node
	n.outEdges = make(map[string](map[string]*TraversalEdge))
	n.inEdges = make(map[string](map[string]*TraversalEdge))

	return &n
}

//-====================-//
//     Define Edge
//-====================-//

type TraversalEdge struct {
	Edge
	_fromNode *TraversalNode
	_toNode   *TraversalNode
}

func (this *TraversalEdge) GetAllNodes() []*TraversalNode {
	allNodes := make([]*TraversalNode, 2)
	allNodes[0] = this._fromNode
	allNodes[1] = this._toNode
	return allNodes
}

func (this *TraversalEdge) GetFromNode() *TraversalNode {
	return this._fromNode
}

func (this *TraversalEdge) GetToNode() *TraversalNode {
	return this._toNode
}

func (this *TraversalEdge) Update(edge *Edge) {
	for attrName, attr := range edge.GetAttributes() {
		attrValue := attr.GetValue()
		if nil != attrValue {
			this.GetAttribute(attrName).SetValue(attrValue)
		}
	}
}

func NewTraversalEdge(edge *Edge, fromNode *TraversalNode, toNode *TraversalNode) *TraversalEdge {
	e := &TraversalEdge{}

	e.Edge = *edge

	e._fromNode = fromNode
	fromNode.SetEdge(e, true)
	e._fromNodeId = &fromNode.NodeId
	e._fromNodeKey = fromNode.GetKey()
	e._fromNodeKeyHash = fromNode.GetKeyHash()

	e._toNode = toNode
	toNode.SetEdge(e, false)
	e._toNodeId = &toNode.NodeId
	e._toNodeKey = toNode.GetKey()
	e._toNodeKeyHash = toNode.GetKeyHash()

	return e
}

//-=========================-//
//     Define GraphImpl
//-=========================-//

type TraversalGraph struct {
	id          string
	modelId     string
	model       map[string]interface{}
	edges       map[EdgeId]*TraversalEdge
	nodes       map[NodeId]*TraversalNode
	edgesByType map[string](map[EdgeId]*TraversalEdge)
	nodesByType map[string](map[NodeId]*TraversalNode)
}

func (g *TraversalGraph) GetId() string {
	return g.id
}

func (g *TraversalGraph) GetModelId() string {
	return g.modelId
}

func (g *TraversalGraph) SetModel(model map[string]interface{}) {
	g.model = model
}

func (g *TraversalGraph) GetModel() map[string]interface{} {
	return g.model
}

func (g *TraversalGraph) UpsertGraph(graph Graph) {
	nodeMap := graph.GetNodes()
	if nil != nodeMap {
		for _, node := range graph.GetNodes() {
			g.UpsertNode(node)
		}
	}

	if nil != graph.GetEdges() {
		for _, edge := range graph.GetEdges() {
			g.UpsertEdge(edge, nodeMap[*edge.GetFromId()], nodeMap[*edge.GetToId()])
		}
	}
}

func (g *TraversalGraph) UpsertNode(node *Node) *TraversalNode {
	//fmt.Println("node = ", node)
	var tNode *TraversalNode
	if nil != g.nodes[node.NodeId] {
		tNode = g.nodes[node.NodeId]
		tNode.Update(node)
	} else {
		tNode := NewTraversalNode(node)
		g.SetNode(tNode.NodeId, tNode)
	}
	//fmt.Println("graph01 = ", g)
	return tNode
}

func (g *TraversalGraph) UpsertEdge(edge *Edge, from *Node, to *Node) {
	//fmt.Println("edge = ", edge)
	if nil != g.edges[edge.EdgeId] {
		g.edges[edge.EdgeId].Update(edge)
	} else {
		fromNode := g.UpsertNode(from)
		toNode := g.UpsertNode(to)
		newEdge := NewTraversalEdge(edge, fromNode, toNode)
		g.SetEdge(newEdge.EdgeId, newEdge)
	}
	//fmt.Println("graph02 = ", g)
}

func (g *TraversalGraph) GetNode(id NodeId) *TraversalNode {
	return g.nodes[id]
}

func (g *TraversalGraph) GetNodeByTypeByKey(nodeType string, nodeId NodeId) *TraversalNode {
	typedNodeById := g.nodesByType[nodeType]
	return typedNodeById[nodeId]
}

func (g *TraversalGraph) GetNodesByType(nodeType string) map[NodeId]*TraversalNode {
	return g.nodesByType[nodeType]
}

func (g *TraversalGraph) GetNodes() map[NodeId]*TraversalNode {
	return g.nodes
}

func (g *TraversalGraph) SetNode(id NodeId, node *TraversalNode) {
	g.nodes[id] = node
	nodeMap := g.nodesByType[node._type]
	if nil == nodeMap {
		nodeMap = make(map[NodeId]*TraversalNode)
		g.nodesByType[node._type] = nodeMap
	}
	nodeMap[id] = node
}

func (g *TraversalGraph) GetEdge(id EdgeId) *TraversalEdge {
	return g.edges[id]
}

func (g *TraversalGraph) SetEdge(id EdgeId, edge *TraversalEdge) {
	g.edges[id] = edge
	edgeMap := g.edgesByType[edge._type]
	if nil == edgeMap {
		edgeMap = make(map[EdgeId]*TraversalEdge)
		g.edgesByType[edge._type] = edgeMap
	}
	edgeMap[id] = edge
}

func (g *TraversalGraph) GetEdges() map[EdgeId]*TraversalEdge {
	return g.edges
}

func (g *TraversalGraph) GetEntityKeyNamesForNode(nodeType string) []string {
	nodeKeyMap := g.model["nodes"].(map[string]interface{})["keyMap"].(map[string][]string)
	return nodeKeyMap[nodeType]
}

func (g *TraversalGraph) GetEntityKeyNamesForEdge(edgeType string) []string {
	edgeKeyMap := g.model["edges"].(map[string]interface{})["keyMap"].(map[string][]string)
	return edgeKeyMap[edgeType]
}

func (g *TraversalGraph) Clear() {
}

func NewTraversalGraph(modelId string, id string) *TraversalGraph {
	var g TraversalGraph
	g.id = id
	g.modelId = modelId
	g.edges = make(map[EdgeId]*TraversalEdge)
	g.nodes = make(map[NodeId]*TraversalNode)
	g.edgesByType = make(map[string](map[EdgeId]*TraversalEdge))
	g.nodesByType = make(map[string](map[NodeId]*TraversalNode))
	return &g
}
