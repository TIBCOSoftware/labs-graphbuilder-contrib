/*
 * Copyright © 2020. TIBCO Software Inc.
 * This file is subject to the license terms contained
 * in the license file that is distributed with this file.
 */
package gremlin

import (
	"bytes"
	"encoding/json"
	"fmt"

	"strconv"

	"github.com/TIBCOSoftware/labs-graphbuilder-lib/model"
	"github.com/TIBCOSoftware/labs-graphbuilder-lib/util"
	gmodel "github.com/northwesternmutual/grammes/model"
)

const (
	_targetRegex string = "[^A-Za-z0-9]"
	_replacement string = "_"
	typeName     string = "type"
)

type EntityKind int

const (
	Node EntityKind = 0
	Edge EntityKind = 1
)

func (this EntityKind) int() int {
	index := [...]int{0, 1}
	return index[this]
}

type GremlinEntity interface {
	Kind() EntityKind
	SetGremlinId(gremlinId interface{})
	GetGremlinId() interface{}
	GetGremlinIdString() string
	ToGremlin(graph model.Graph, dateTimeSample string) string
}

type GremlinEntityImpl struct {
	gremlinId       interface{}
	addPrefixToAttr bool
}

func (this *GremlinEntityImpl) SetGremlinId(gremlinId interface{}) {
	this.gremlinId = gremlinId
	log.Debug("(DNode::GetGremlinId) gremlinId = ", this.gremlinId)
}

func (this *GremlinEntityImpl) GetGremlinId() interface{} {
	log.Debug("(DNode::GetGremlinId) gremlinId = ", this.gremlinId)
	return this.gremlinId
}

func (this *GremlinEntityImpl) GetGremlinIdString() string {

	log.Info("(DNode::GetGremlinId) gremlinId = ", this.gremlinId)

	strData, ok := this.gremlinId.(string)
	if !ok {
		strData = strconv.FormatInt(int64(this.gremlinId.(map[string]interface{})["@value"].(float64)), 10)
	} else {
		strData = fmt.Sprintf("'%s'", strData)
	}
	return strData
}

type DNode struct {
	GremlinEntityImpl
	node *model.Node
}

func NewDNode(node *model.Node, typeName string, addPrefixToAttr bool) *DNode {
	dnode := DNode{
		node: node,
	}
	dnode.addPrefixToAttr = addPrefixToAttr

	return &dnode
}

func (this *DNode) Kind() EntityKind {
	return Node
}

func (this *DNode) GetType() string {

	if "" != typeName && nil != this.node.GetAttribute(typeName) {
		if model.TypeString == this.node.GetAttribute(typeName).GetType() {
			return this.GetAttribute(typeName).GetValue().(string)
		}
	}

	return this.node.GetType()
}

func (this *DNode) GetAttributes() map[string]*model.Attribute {
	return this.node.GetAttributes()
}

func (this *DNode) GetCanonicalAttributeName(attributeName string) string {
	return "" //GetCanonicalAttributeName(this.node.GetType(), attributeName, _targetRegex, _replacement, this.addPrefixToAttr, true)
}

func (this *DNode) GetAttribute(attributeName string) *model.Attribute {
	return this.node.GetAttribute(attributeName)
}

func (this *DNode) GetAttributeValueString(attributeName string, DateTimeSample string) (string, error) {
	attribute := this.node.GetAttribute(attributeName)
	if nil == attribute {
		return "", nil
	}
	dataType := attribute.GetType().String()
	strValue, err := util.ToString(attribute.GetValue(), dataType, DateTimeSample)

	if nil != err {
		log.Error(err)
	}

	if "String" == dataType || "Date" == dataType {
		strValue = util.ReplaceCharacter(strValue, "'", "\"", true)
		strValue = fmt.Sprintf("'%s'", strValue)
	}

	return strValue, nil
}

func (this *DNode) GetPrimaryKeyString() string {
	return this.node.GetKeyHash()
}

func (this *DNode) GetId() string {
	return this.node.NodeId.ToString()
}

func (this *DNode) GetEid() string {
	return fmt.Sprintf("_:%s", util.ReplaceCharacter(this.node.NodeId.ToString(), _targetRegex, _replacement, true))
}

func (this *DNode) Exists() bool {
	return nil != this.gremlinId
}

func (this *DNode) ToString() string {
	return fmt.Sprintf("Node(%s)_%s", this.GetType(), this.node.GetKey())
}

func (this *DNode) ToGremlin(graph model.Graph, dateTimeSample string) string {
	log.Debug("(node.ToGremlin) begin - ", this)

	var query bytes.Buffer
	keyNames := graph.GetEntityKeyNamesForNode(this.GetType())
	if this.Exists() {
		query.WriteString("g.addV('")
		query.WriteString(util.ReplaceCharacter(this.GetType(), _targetRegex, _replacement, true))
		query.WriteString("')")

		log.Debug("(node.ToGremlin) ", this.node.GetAttributes())

		if 0 < len(this.node.GetAttributes()) {
			var count int
			for attrname, _ := range this.node.GetAttributes() {
				if util.SliceContains(keyNames, attrname) {
					if !this.Exists() {
						attrname = this.GetCanonicalAttributeName(attrname)
					} else {
						// key filed only for insert
						continue
					}
				}

				attrval, err := this.GetAttributeValueString(attrname, dateTimeSample)
				if "" == attrval || nil != err {
					continue
				}

				if 0 == count {
					query.WriteString(".property('")
				} else {
					query.WriteString("property('")
				}
				query.WriteString("'")
				query.WriteString(attrname)
				query.WriteString("', ")
				query.WriteString(attrval)
				count += 1
				query.WriteString(")")
			}
		}
	} else {
		query.WriteString("g.addV('")
		query.WriteString(util.ReplaceCharacter(this.GetType(), _targetRegex, _replacement, true))
		query.WriteString("')")
		var count int
		for attrname, _ := range this.node.GetAttributes() {
			attrval, err := this.GetAttributeValueString(attrname, dateTimeSample)
			if "" == attrval || nil != err {
				continue
			}

			if 0 == count {
				query.WriteString(".property('")
			} else {
				query.WriteString("property('")
			}
			query.WriteString(attrname)
			query.WriteString("', ")
			query.WriteString(attrval)
			count += 1
			query.WriteString(")")
		}
	}

	return query.String()
}

type DEdge struct {
	GremlinEntityImpl
	edge     *model.Edge
	typeName string
	from     *DNode
	to       *DNode
}

func NewDEdge(edge *model.Edge, typeName string, addPrefixToAttr bool, from *DNode, to *DNode) *DEdge {
	dedge := DEdge{
		edge:     edge,
		typeName: typeName,
		from:     from,
		to:       to,
	}
	dedge.addPrefixToAttr = addPrefixToAttr

	return &dedge
}

func (this *DEdge) Kind() EntityKind {
	return Edge
}

func (this *DEdge) GetType() string {
	var edgeType string
	if nil != this.GetAttribute("relation") {
		edgeType = this.GetAttribute("relation").GetValue().(string)
	}

	if "" == edgeType {
		return this.edge.GetType()
	}

	return edgeType
}

func (this *DEdge) GetAttributes() map[string]*model.Attribute {
	return this.edge.GetAttributes()
}

func (this *DEdge) GetCanonicalAttributeName(attributeName string) string {
	return "" //GetCanonicalAttributeName(this.edge.GetType(), attributeName, _targetRegex, _replacement, this.addPrefixToAttr, true)
}

func (this *DEdge) GetAttribute(attributeName string) *model.Attribute {
	return this.edge.GetAttribute(attributeName)
}

func (this *DEdge) GetAttributeValueString(attributeName string, DateTimeSample string) (string, error) {
	attribute := this.edge.GetAttribute(attributeName)
	if nil == attribute {
		return "", nil
	}
	dataType := attribute.GetType().String()
	strValue, err := util.ToString(attribute.GetValue(), dataType, DateTimeSample)

	if nil != err {
		log.Error(err)
	}

	if "String" == dataType || "Date" == dataType {
		strValue = util.ReplaceCharacter(strValue, "'", "\"", true)
		strValue = fmt.Sprintf("'%s'", strValue)
	}

	return strValue, nil
}
func (this *DEdge) GetId() string {
	return this.edge.EdgeId.ToString()
}

func (this *DEdge) Exists() bool {
	return nil != this.gremlinId
}

func (this *DEdge) GetFrom() *DNode {
	return this.from
}

func (this *DEdge) GetTo() *DNode {
	return this.to
}

func (this *DEdge) ToString() string {
	return fmt.Sprintf("Edge(%s):from(%s):to(%s)", this.GetType(), this.GetFrom().ToString(), this.GetTo().ToString())
}

func (this *DEdge) ToGremlin(graph model.Graph, dateTimeSample string) string {
	log.Debug("(edge.ToGremlin) begin - ", this)

	var query bytes.Buffer
	if this.Exists() {
		log.Info("edge exists")
		query.WriteString("g.V().hasId(")
		query.WriteString(this.from.GetGremlinIdString())
		query.WriteString(").outE('")
		query.WriteString(util.ReplaceCharacter(this.GetType(), _targetRegex, _replacement, true))
		query.WriteString("').inV().hasId(")
		query.WriteString(this.to.GetGremlinIdString())
		query.WriteString(").inE('")
		query.WriteString(util.ReplaceCharacter(this.GetType(), _targetRegex, _replacement, true))
		query.WriteString("')")

		if 0 < len(this.GetAttributes()) {
			var count int
			for attrname, _ := range this.GetAttributes() {
				if "relation" == attrname {
					continue
				}

				attrval, err := this.GetAttributeValueString(attrname, dateTimeSample)
				if "" == attrval || nil != err {
					continue
				}

				if 0 == count {
					query.WriteString(".property('")
				} else {
					query.WriteString("property('")
				}
				query.WriteString(attrname)
				query.WriteString("', ")
				query.WriteString(attrval)
				query.WriteString(")")
				count += 1
			}
		}
	} else {
		//g.V().has('recipe','name','Beef Bourguignon').as('a').
		//V().has('ingredient','name','beef').
		//addE('includes').from('a').property('amount','2 lbs')

		log.Debug("edge not exists, from : ", this.from, ", to : ", this.to, ", type : ", this.GetType())
		query.WriteString("g.V().hasId(")
		query.WriteString(this.from.GetGremlinIdString())
		query.WriteString(").as('fromNode').V().hasId(")
		query.WriteString(this.to.GetGremlinIdString())
		query.WriteString(").addE('")
		query.WriteString(util.ReplaceCharacter(this.GetType(), _targetRegex, _replacement, true))
		query.WriteString("').from('fromNode')")

		if 0 < len(this.GetAttributes()) {
			var count int
			for attrname, _ := range this.GetAttributes() {
				attrval, err := this.GetAttributeValueString(attrname, dateTimeSample)
				if "" == attrval || nil != err {
					continue
				}

				if 0 == count {
					query.WriteString(".property('")
				} else {
					query.WriteString("property('")
				}
				query.WriteString(attrname)
				query.WriteString("', ")
				query.WriteString(attrval)
				query.WriteString(")")
				count += 1
			}
		}
	}

	return query.String()
}

type Path struct {
	listOfEdges gmodel.List
	Edges       []gmodel.Edge
}

func (l *Path) UnmarshalJSON(data []byte) error {
	if err := json.Unmarshal(data, &l.listOfEdges); err == nil {
		if data, err = json.Marshal(l.listOfEdges.Value); err != nil {
			return err
		}
	}

	return json.Unmarshal(data, &l.Edges)
}
