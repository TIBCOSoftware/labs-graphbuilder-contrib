/*
 * Copyright © 2020. TIBCO Software Inc.
 * This file is subject to the license terms contained
 * in the license file that is distributed with this file.
 */
package model

import (
	"strings"

	"github.com/TIBCOSoftware/labs-graphbuilder-lib/util"
)

//-====================-//
//   Define Attribute
//-====================-//

type Attribute struct {
	_name  string
	_value interface{}
	_type  DataType
}

func (a *Attribute) GetName() string {
	return a._name
}

func (a *Attribute) SetName(myName string) {
	a._name = myName
}

func (a *Attribute) GetValue() interface{} {
	return a._value
}

func (a *Attribute) SetValue(myValue interface{}) error {
	//fmt.Println("Attribute _name = ", a._name, ", _type = ", a._type, ", myValue = ", myValue)

	if nil == myValue {
		a._value = myValue
		return nil
	}

	var err error
	switch a._type {
	//case TypeString:
	//	break
	case TypeInteger:
		a._value, err = util.ConvertToInteger(myValue)
		break
	case TypeLong:
		a._value, err = util.ConvertToLong(myValue)
		break
	//case TypeDouble:
	//	break
	//case TypeBoolean:
	//	break
	//case TypeDate:
	//	break
	default:
		a._value = myValue
	}
	return err
}

func (a *Attribute) GetType() DataType {
	return a._type
}

func (a *Attribute) SetType(myType DataType) {
	a._type = myType
}

func NewAttribute(attributeModel *Attributefinition, value interface{}) *Attribute {
	attr := &Attribute{_name: attributeModel.name, _type: attributeModel.dataType}
	attr.SetValue(value)

	return attr
}

//-========================-//
//     Define EntityId
//-========================-//

type EntityId struct {
	_keyHash string
	_type    string
}

func (this *EntityId) GetType() string {
	return this._type
}

func (this *EntityId) GetKeyHash() string {
	return this._keyHash
}

func (this *EntityId) ToStringBuffer() *strings.Builder {
	var sb strings.Builder
	sb.WriteString(this._type)
	sb.WriteString("_")
	sb.WriteString(this._keyHash)
	return &sb
}

//-====================-//
//     Define Entity
//-====================-//

type Entity interface {
	SetKey(key []interface{})

	GetKey() []interface{}

	GetAttribute(name string) *Attribute

	SetAttribute(name string, attr *Attribute)

	GetAttributes() map[string]*Attribute
}

//-====================-//
//   Define BaseEntity
//-====================-//

type BaseEntity struct {
	_key        []interface{}
	_attributes map[string]*Attribute
}

func (e *BaseEntity) SetKey(key []interface{}) {
	e._key = key
}

func (e *BaseEntity) GetKey() []interface{} {
	return e._key
}

func (e *BaseEntity) GetAttribute(name string) *Attribute {
	return e._attributes[name]
}

func (e *BaseEntity) SetAttribute(name string, attr *Attribute) {
	e._attributes[name] = attr
}

func (e *BaseEntity) GetAttributes() map[string]*Attribute {
	return e._attributes
}

//-====================-//
//     Define NodeId
//-====================-//

type NodeId struct {
	EntityId
}

func (this *NodeId) ToString() string {
	return this.EntityId.ToStringBuffer().String()
}

func (this *NodeId) FromString(id string) *NodeId {
	pos := strings.LastIndexByte(id, '_')
	this._type = id[0:pos]
	this._keyHash = id[pos+1 : len(id)]
	return this
}

//-====================-//
//     Define Node
//-====================-//

type Node struct {
	NodeId
	BaseEntity
}

func NewNode(myType string, myKey []interface{}) *Node {
	var n Node
	n._key = myKey
	n._keyHash = Hash(myKey)
	n._type = myType
	n._attributes = make(map[string]*Attribute)
	return &n
}

//-====================-//
//     Define EdgeId
//-====================-//

type EdgeId struct {
	EntityId
	_fromNodeKeyHash string
	_toNodeKeyHash   string
}

func (eid *EdgeId) ToString() string {
	sb := eid.EntityId.ToStringBuffer()
	sb.WriteString("_")
	sb.WriteString(eid._fromNodeKeyHash)
	sb.WriteString("_")
	sb.WriteString(eid._toNodeKeyHash)
	return sb.String()
}

//-====================-//
//     Define Edge
//-====================-//

type Edge struct {
	EdgeId
	BaseEntity
	_fromNodeId  *NodeId
	_fromNodeKey []interface{}
	_toNodeId    *NodeId
	_toNodeKey   []interface{}
}

func (this *Edge) GetFromId() *NodeId {
	return this._fromNodeId
}

func (this *Edge) GetToId() *NodeId {
	return this._toNodeId
}

func (this *Edge) UpdateFrom(node Node) {
	this._fromNodeId = &node.NodeId
	this._fromNodeKey = node.GetKey()
	this._fromNodeKeyHash = node.GetKeyHash()
}

func (this *Edge) UpdateTo(node Node) {
	this._toNodeId = &node.NodeId
	this._toNodeKey = node.GetKey()
	this._toNodeKeyHash = node.GetKeyHash()
}

func NewEdge(myType string, myKey []interface{}, fromNode *Node, toNode *Node) *Edge {
	var e Edge

	e._fromNodeId = &fromNode.NodeId
	e._fromNodeKey = fromNode.GetKey()
	e._fromNodeKeyHash = fromNode.GetKeyHash()
	e._toNodeId = &toNode.NodeId
	e._toNodeKey = toNode.GetKey()
	e._toNodeKeyHash = toNode.GetKeyHash()
	e._type = myType
	e._key = myKey
	e._keyHash = Hash(myKey)
	e._attributes = make(map[string]*Attribute)

	return &e
}

//-====================-//
//     Define Graph
//-====================-//

type Graph interface {
	GetId() string
	GetModelId() string
	GetModel() map[string]interface{}
	GetNodes() map[NodeId]*Node
	GetEdges() map[EdgeId]*Edge
	UpsertGraph(graph map[string]interface{})
	UpsertNode(nodeType string, nodeKey []interface{}) *Node
	UpsertEdge(edgeType string, edgeKey []interface{}, fromNode *Node, toNode *Node) *Edge
	Merge(graph Graph)
	GetNodeByTypeByKey(nodeType string, nodeId NodeId) *Node
	GetNodesByType(nodeType string) map[NodeId]*Node
	GetEntityKeyNamesForNode(entityName string) []string
	GetEntityKeyNamesForEdge(entityName string) []string
	Clear()
}

//-=========================-//
//     Define GraphImpl
//-=========================-//

type GraphImpl struct {
	id          string
	modelId     string
	model       map[string]interface{}
	edges       map[EdgeId]*Edge
	nodes       map[NodeId]*Node
	edgesByType map[string](map[EdgeId]*Edge)
	nodesByType map[string](map[NodeId]*Node)
}

func (g *GraphImpl) GetId() string {
	return g.id
}

func (g *GraphImpl) GetModelId() string {
	return g.modelId
}

func (g *GraphImpl) SetModel(model map[string]interface{}) {
	g.model = model
}

func (g *GraphImpl) GetModel() map[string]interface{} {
	return g.model
}

func (g *GraphImpl) UpsertGraph(graph map[string]interface{}) {

}

func (g *GraphImpl) UpsertNode(nodeType string, nodeKey []interface{}) *Node {
	node := NewNode(nodeType, nodeKey)
	if nil != g.nodes[node.NodeId] {
		return g.nodes[node.NodeId]
	}
	g.SetNode(node.NodeId, node)
	return node
}

func (g *GraphImpl) UpsertNodeToo(node *Node) {
	existNode := g.GetNode(node.NodeId)
	if nil != existNode {
		for attrKey, attribute := range node.GetAttributes() {
			existAttribute := existNode.GetAttribute(attrKey)
			if nil == existAttribute {
				existNode.SetAttribute(attrKey, attribute)
			} else {
				existAttribute.SetValue(attribute.GetValue())
			}
		}
	} else {
		g.SetNode(node.NodeId, node)
	}
}

func (g *GraphImpl) Merge(graph Graph) {
	for _, node := range graph.GetNodes() {
		g.UpsertNodeToo(node)
	}

	for _, edge := range graph.GetEdges() {
		g.UpsertEdgeToo(edge)
	}
}

func (g *GraphImpl) GetNode(id NodeId) *Node {
	return g.nodes[id]
}

func (g *GraphImpl) GetNodeByTypeByKey(nodeType string, nodeId NodeId) *Node {
	typedNodeById := g.nodesByType[nodeType]
	return typedNodeById[nodeId]
}

func (g *GraphImpl) GetNodesByType(nodeType string) map[NodeId]*Node {
	return g.nodesByType[nodeType]
}

func (g *GraphImpl) GetNodes() map[NodeId]*Node {
	return g.nodes
}

func (g *GraphImpl) SetNode(id NodeId, node *Node) {
	g.nodes[id] = node
	nodeMap := g.nodesByType[node._type]
	if nil == nodeMap {
		nodeMap = make(map[NodeId]*Node)
		g.nodesByType[node._type] = nodeMap
	}
	nodeMap[id] = node
}

func (g *GraphImpl) UpsertEdge(edgeType string, edgeKey []interface{}, fromNode *Node, toNode *Node) *Edge {
	edge := NewEdge(edgeType, edgeKey, fromNode, toNode)
	if nil != g.edges[edge.EdgeId] {
		return g.edges[edge.EdgeId]
	}
	g.SetEdge(edge.EdgeId, edge)
	return edge
}

func (g *GraphImpl) UpsertEdgeToo(edge *Edge) {
	existEdge := g.GetEdge(edge.EdgeId)
	if nil != existEdge {
		for attrKey, attribute := range edge.GetAttributes() {
			existAttribute := existEdge.GetAttribute(attrKey)
			if nil == existAttribute {
				existEdge.SetAttribute(attrKey, attribute)
			} else {
				existAttribute.SetValue(attribute.GetValue())
			}
		}
	} else {
		existFromNode := g.GetNode(*edge.GetFromId())
		if nil != existFromNode {
			edge.UpdateFrom(*existFromNode)
		}
		existToNode := g.GetNode(*edge.GetToId())
		if nil != existToNode {
			edge.UpdateTo(*existToNode)
		}
		g.SetEdge(edge.EdgeId, edge)
	}
}

func (g *GraphImpl) GetEdge(id EdgeId) *Edge {
	return g.edges[id]
}

func (g *GraphImpl) SetEdge(id EdgeId, edge *Edge) {
	g.edges[id] = edge
	edgeMap := g.edgesByType[edge._type]
	if nil == edgeMap {
		edgeMap = make(map[EdgeId]*Edge)
		g.edgesByType[edge._type] = edgeMap
	}
	edgeMap[id] = edge
}

func (g *GraphImpl) GetEdges() map[EdgeId]*Edge {
	return g.edges
}

func (g *GraphImpl) GetEntityKeyNamesForNode(nodeType string) []string {
	nodeKeyMap := g.model["nodes"].(map[string]interface{})["keyMap"].(map[string][]string)
	return nodeKeyMap[nodeType]
}

func (g *GraphImpl) GetEntityKeyNamesForEdge(edgeType string) []string {
	edgeKeyMap := g.model["edges"].(map[string]interface{})["keyMap"].(map[string][]string)
	return edgeKeyMap[edgeType]
}

func (g *GraphImpl) Clear() {
	g.edges = make(map[EdgeId]*Edge)
	g.nodes = make(map[NodeId]*Node)
	g.edgesByType = make(map[string](map[EdgeId]*Edge))
	g.nodesByType = make(map[string](map[NodeId]*Node))
}

func NewGraphImpl(modelId string, id string) *GraphImpl {
	var g GraphImpl
	g.id = id
	g.modelId = modelId
	g.edges = make(map[EdgeId]*Edge)
	g.nodes = make(map[NodeId]*Node)
	g.edgesByType = make(map[string](map[EdgeId]*Edge))
	g.nodesByType = make(map[string](map[NodeId]*Node))
	return &g
}
