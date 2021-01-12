/*
 * Copyright © 2020. TIBCO Software Inc.
 * This file is subject to the license terms contained
 * in the license file that is distributed with this file.
 */
package neo4j

import (
	"bytes"
	"fmt"

	"strconv"

	"github.com/TIBCOSoftware/labs-graphbuilder-lib/model"
	"github.com/TIBCOSoftware/labs-graphbuilder-lib/util"
)

const (
	_targetRegex string = "[^A-Za-z0-9]"
	_replacement string = "_"
	NODE_FORMAT  string = "%s <%s> \"%s\" . \n"
	EDGE_FORMAT  string = "%s <%s> %s %s . \n"
	typeName     string = "type"
)

type DNode struct {
	neoId           interface{}
	node            *model.Node
	addPrefixToAttr bool
}

func NewDNode(node *model.Node, typeName string, addPrefixToAttr bool) *DNode {
	dnode := DNode{
		node:            node,
		addPrefixToAttr: addPrefixToAttr,
	}

	return &dnode
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
		return strValue, err
	}

	if "String" == dataType || "Date" == dataType {
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

func (this *DNode) SetNeoId(neoId interface{}) {
	this.neoId = neoId
	log.Debug("(DNode::GetNeoId) neoId = ", this.neoId)
}

func (this *DNode) GetNeoId() interface{} {
	log.Debug("(DNode::GetNeoId) neoId = ", this.neoId)
	return this.neoId
}

func (this *DNode) GetNeoIdString() string {

	log.Debug("(DNode::GetNeoIdString) neoId = ", this.neoId)

	strData := strconv.FormatInt(this.neoId.(int64), 10)
	return strData
}

func (this *DNode) Exists() bool {
	return nil != this.neoId
}

func (this *DNode) ToString() string {
	return fmt.Sprintf("Node(%s)_%s", this.GetType(), this.node.GetKey())
}

func (this *DNode) ToCypher(graph model.Graph, dateTimeSample string) string {
	log.Debug("(node.ToCypher) begin - ", this)

	var query bytes.Buffer
	keyNames := graph.GetEntityKeyNamesForNode(this.GetType())
	if this.Exists() {
		//MATCH (node:nodeType) Where id(node)=100 SET n.surname = 'Taylor'
		query.WriteString("MATCH (node:")
		query.WriteString(util.ReplaceCharacter(this.GetType(), _targetRegex, _replacement, true))
		query.WriteString(")\n")
		query.WriteString("WHERE id(node) = ")
		query.WriteString(this.GetNeoIdString())
		query.WriteString("\n")

		log.Debug(">>>>>>>>>>>>", this.node.GetAttributes())

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
					query.WriteString("SET ")
				} else {
					query.WriteString(", ")
				}
				query.WriteString("node.")
				query.WriteString(attrname)
				query.WriteString(" = ")
				query.WriteString(attrval)
				count += 1
				query.WriteString("\n")
			}
		}
		query.WriteString("RETURN id(node) \n")
	} else {
		//CREATE (node:nodeType { name: 'Andy', title: 'Developer' })
		query.WriteString("CREATE (node:")
		query.WriteString(util.ReplaceCharacter(this.GetType(), _targetRegex, _replacement, true))
		query.WriteString("{ ")
		var count int
		for attrname, _ := range this.node.GetAttributes() {
			attrval, err := this.GetAttributeValueString(attrname, dateTimeSample)
			if "" == attrval || nil != err {
				continue
			}

			if 0 != count {
				query.WriteString(", ")
			}
			query.WriteString(attrname)
			query.WriteString(" : ")
			query.WriteString(attrval)
			count += 1
		}
		query.WriteString(" }) RETURN id(node) \n")
	}

	return query.String()
}

type DEdge struct {
	neoId           interface{}
	edge            *model.Edge
	typeName        string
	addPrefixToAttr bool
	from            *DNode
	to              *DNode
}

func NewDEdge(edge *model.Edge, typeName string, addPrefixToAttr bool, from *DNode, to *DNode) *DEdge {
	dedge := DEdge{
		edge:            edge,
		typeName:        typeName,
		addPrefixToAttr: addPrefixToAttr,
		from:            from,
		to:              to,
	}

	return &dedge
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
		return strValue, err
	}

	if "String" == dataType || "Date" == dataType {
		strValue = fmt.Sprintf("'%s'", strValue)
	}
	return strValue, nil
}
func (this *DEdge) GetId() string {
	return this.edge.EdgeId.ToString()
}

func (this *DEdge) SetNeoId(neoId interface{}) {
	this.neoId = neoId
	log.Debug("(DEdge::GetNeoId) neoId = ", this.neoId)
}

func (this *DEdge) GetNeoId() interface{} {
	log.Debug("(DEdge::GetNeoId) neoId = ", this.neoId)
	return this.neoId
}

func (this *DEdge) GetNeoIdString() string {
	strData := strconv.FormatInt(this.neoId.(int64), 10)
	return strData
}

func (this *DEdge) Exists() bool {
	return nil != this.neoId
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

func (this *DEdge) ToCypher(graph model.Graph, dateTimeSample string) string {
	log.Debug("(edge.ToCypher) begin - ", this)

	var query bytes.Buffer
	if this.Exists() {
		//MATCH (from:nodeType)-[edge:edgeType]->(to:nodeType)
		//WHERE id(edge) = 100
		//Set edge.attrname = "attrval"
		log.Debug("edge esists")
		query.WriteString(fmt.Sprintf(
			"MATCH (from:%s)-[edge:%s]->(to:%s)\n",
			util.ReplaceCharacter(this.from.GetType(), _targetRegex, _replacement, true),
			util.ReplaceCharacter(this.GetType(), _targetRegex, _replacement, true),
			util.ReplaceCharacter(this.to.GetType(), _targetRegex, _replacement, true),
		))
		query.WriteString("WHERE id(edge) = ")
		query.WriteString(this.GetNeoIdString())
		query.WriteString("\n")

		log.Debug("========>>>>>>>>>>>>", this.GetAttributes())

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
					query.WriteString("SET ")
				} else {
					query.WriteString(", ")
				}
				query.WriteString("edge.")
				query.WriteString(attrname)
				query.WriteString(" = ")
				query.WriteString(attrval)
				count += 1
			}
		}
		query.WriteString("RETURN id(edge)\n")
	} else {
		//MATCH (from:nodeType),(to:nodeType)
		//WHERE id(from) = 100 AND id(to) = 101
		//CREATE (from)-[edge:edgeType{ attrname:attrval }]->(to)
		//RETURN type(r)

		//MATCH (from:nodeType),(to:nodeType)
		//WHERE id(from) = 150 AND id(to) = 109
		//CREATE (from)-[edge:Sold By]->(to)
		//RETURN id(edge)

		log.Debug("edge not exists, from : ", this.from, ", to : ", this.to, ", type : ", this.GetType())
		query.WriteString(fmt.Sprintf(
			"MATCH (from:%s),(to:%s)\n",
			util.ReplaceCharacter(this.from.GetType(), _targetRegex, _replacement, true),
			util.ReplaceCharacter(this.to.GetType(), _targetRegex, _replacement, true),
		))
		query.WriteString("WHERE id(from) = ")
		query.WriteString(this.from.GetNeoIdString())
		query.WriteString(" AND id(to) = ")
		query.WriteString(this.to.GetNeoIdString())
		query.WriteString("\nCREATE (from)-[edge:")
		query.WriteString(
			util.ReplaceCharacter(this.GetType(), _targetRegex, _replacement, true),
		)

		if 0 < len(this.GetAttributes()) {
			query.WriteString("{ ")
			var count int
			for attrname, _ := range this.GetAttributes() {
				attrval, err := this.GetAttributeValueString(attrname, dateTimeSample)
				if "" == attrval || nil != err {
					continue
				}

				if 0 != count {
					query.WriteString(", ")
				}
				query.WriteString(attrname)
				query.WriteString(" : ")
				query.WriteString(attrval)
				count += 1
			}
			query.WriteString(" }")
		}
		query.WriteString("]->(to)\nRETURN id(edge)\n")
	}

	return query.String()
}
