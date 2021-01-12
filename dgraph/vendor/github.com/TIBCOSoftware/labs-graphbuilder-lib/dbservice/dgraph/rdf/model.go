/*
 * Copyright © 2020. TIBCO Software Inc.
 * This file is subject to the license terms contained
 * in the license file that is distributed with this file.
 */
package rdf

import (
	"bytes"
	"errors"
	"fmt"
	"strings"

	"github.com/TIBCOSoftware/flogo-lib/logger"
	"github.com/TIBCOSoftware/labs-graphbuilder-lib/dbservice/dgraph"
	"github.com/TIBCOSoftware/labs-graphbuilder-lib/model"
	"github.com/TIBCOSoftware/labs-graphbuilder-lib/util"
)

var log = logger.GetLogger("dgraph-service")

const (
	_targetRegex       = "[^A-Za-z0-9]"
	_replacement       = "_"
	IMPORT_NODE_FORMAT = "%s <%s> \"%s\" . \n"
	IMPORT_EDGE_FORMAT = "%s <%s> %s %s . \n"
)

type DNode struct {
	uid  string
	node *model.Node

	explicitType    bool
	typeName        string
	addPrefixToAttr bool
}

func NewDNode(node *model.Node, explicitType bool, typeName string, addPrefixToAttr bool) *DNode {
	if "" == typeName {
		typeName = "type"
	}
	dnode := DNode{
		node:            node,
		explicitType:    explicitType,
		typeName:        typeName,
		addPrefixToAttr: addPrefixToAttr,
	}

	return &dnode
}

func (this *DNode) Update(node *model.Node) (bool, error) {
	if node.NodeId != this.node.NodeId {
		return false, errors.New(fmt.Sprint("Update fail : old id = ", this.node.NodeId, ", new id = ", node.NodeId))
	}
	/* Assume fixed a are  */
	changed := false
	for attrName, attr := range node.GetAttributes() {
		newValue := attr.GetValue()
		oldAttr := this.node.GetAttribute(attrName)
		if nil == oldAttr && nil != attr {
			this.node.SetAttribute(attrName, attr)
			changed = true
		} else if nil != newValue && newValue != oldAttr.GetValue() {
			oldAttr.SetValue(newValue)
			changed = true
		}
	}

	return changed, nil
}

func (this *DNode) GetTypeName() string {
	return this.typeName
}

func (this *DNode) GetType() string {

	if "" != this.typeName && nil != this.node.GetAttribute(this.typeName) {
		if model.TypeString == this.node.GetAttribute(this.typeName).GetType() {
			return this.GetAttributeAsString(this.typeName)
		}
	}

	return this.node.GetType()
}

func (this *DNode) GetAttributes() map[string]*model.Attribute {
	return this.node.GetAttributes()
}

func (this *DNode) GetNode() *model.Node {
	return this.node
}

func (this *DNode) GetCanonicalAttributeName(attributeName string) string {
	return dgraph.GetCanonicalAttributeName(this.node.GetType(), attributeName, _targetRegex, _replacement, this.addPrefixToAttr, true)
}

func (this *DNode) GetAttribute(attributeName string) *model.Attribute {
	return this.node.GetAttribute(attributeName)
}

func (this *DNode) GetAttributeAsString(attributeName string) string {
	return model.ToString(this.GetAttribute(attributeName))
}

func (this *DNode) GetPrimaryKeyString() string {
	return this.node.GetKeyHash()
}

func (this *DNode) GetId() string {
	return this.node.NodeId.ToString()
}

func (this *DNode) GetEid(readable bool) string {
	if readable {
		return fmt.Sprintf("_:%s", util.ReplaceCharacter(dgraph.ReadableExternalId(this.GetType(), this.node.GetKey()), _targetRegex, _replacement, true))
	} else {
		return fmt.Sprintf("_:%s", util.ReplaceCharacter(this.node.NodeId.ToString(), _targetRegex, _replacement, true))
	}
}

func (this *DNode) GetFormatedUid() string {
	return fmt.Sprintf("<%s>", this.uid)
}

func (this *DNode) GetUid() string {
	return this.uid
}

func (this *DNode) SetUid(uid string) {
	this.uid = uid
}

func (this *DNode) Exists() bool {
	return "" != this.uid
}

func (this *DNode) ToString() string {
	return fmt.Sprintf("Node(%s)_%s", this.GetType(), this.node.GetKey())
}

func (this *DNode) ToRDF(graph model.Graph, dateTimeSample string, readable bool) []string {
	keyNames := graph.GetEntityKeyNamesForNode(this.GetType())
	nodeQuads := make([]string, 0)
	var id string
	if this.Exists() {
		id = this.GetFormatedUid()
	} else {
		id = this.GetEid(readable)
		/* model id only for insert */
		nodeQuads = append(
			nodeQuads,
			fmt.Sprintf(
				IMPORT_NODE_FORMAT,
				id,
				dgraph.GRAPH_MODEL_ID,
				util.ReplaceCharacter(graph.GetModelId(), _targetRegex, _replacement, true),
			),
		)

		var typeKey string
		var typeVal interface{}
		if !this.explicitType {
			typeKey = util.ReplaceCharacter(this.GetType(), _targetRegex, _replacement, true)
			typeVal = ""
		} else {
			typeKey = this.typeName
			typeVal = util.ReplaceCharacter(this.GetType(), _targetRegex, _replacement, true)
		}
		/* type only for insert */
		nodeQuads = append(nodeQuads, fmt.Sprintf(IMPORT_NODE_FORMAT, id, typeKey, typeVal))
	}

	for attrname, attribute := range this.node.GetAttributes() {
		if util.SliceContains(keyNames, attrname) {
			if !this.Exists() {
				attrname = this.GetCanonicalAttributeName(attrname)
			} else {
				/* key filed only for insert */
				continue
			}
		}

		if this.typeName == attrname || nil == attribute.GetValue() {
			continue
		}

		attrStringValue, err := util.ToString(attribute.GetValue(), attribute.GetType().String(), dateTimeSample)

		if nil != err {
			log.Debug("(ToRDF) Formatting error", err)
			continue
		}

		nodeQuads = append(nodeQuads,
			fmt.Sprintf(IMPORT_NODE_FORMAT, id, attrname, strings.Replace(attrStringValue, "\\", " ", -1)))
	}

	return nodeQuads
}

type DEdge struct {
	edge            *model.Edge
	explicitType    bool
	typeName        string
	addPrefixToAttr bool
	from            *DNode
	to              *DNode
}

func NewDEdge(edge *model.Edge, explicitType bool, typeName string, addPrefixToAttr bool, from *DNode, to *DNode) *DEdge {
	if "" == typeName {
		typeName = "relation"
	}
	dedge := DEdge{
		edge:            edge,
		explicitType:    explicitType,
		typeName:        typeName,
		addPrefixToAttr: addPrefixToAttr,
		from:            from,
		to:              to,
	}

	return &dedge
}

func (this *DEdge) Update(edge *model.Edge) (bool, error) {
	/* Assume fixed a are  */
	if edge.EdgeId != this.edge.EdgeId {
		return false, errors.New(fmt.Sprint("Update fail : old id = ", this.edge.EdgeId, ", new id = ", edge.EdgeId))
	}
	changed := false
	for attrName, attr := range edge.GetAttributes() {
		newValue := attr.GetValue()
		oldAttr := this.edge.GetAttribute(attrName)
		if nil != newValue && newValue != oldAttr.GetValue() {
			oldAttr.SetValue(newValue)
			changed = true
		}
	}

	return changed, nil
}

func (this *DEdge) GetType() string {
	if nil != this.edge.GetAttribute(this.typeName) {
		if model.TypeString == this.edge.GetAttribute(this.typeName).GetType() {
			return this.GetAttributeAsString(this.typeName)
		}
	}

	return this.edge.GetType()
}

func (this *DEdge) GetEdge() *model.Edge {
	return this.edge
}

func (this *DEdge) GetAttributes() map[string]*model.Attribute {
	return this.edge.GetAttributes()
}

func (this *DEdge) GetCanonicalAttributeName(attributeName string) string {
	return dgraph.GetCanonicalAttributeName(this.edge.GetType(), attributeName, _targetRegex, _replacement, this.addPrefixToAttr, true)
}

func (this *DEdge) GetAttribute(attributeName string) *model.Attribute {
	return this.edge.GetAttribute(attributeName)
}

func (this *DEdge) GetAttributeAsString(attributeName string) string {
	return model.ToString(this.GetAttribute(attributeName))
}

func (this *DEdge) GetId() string {
	return this.edge.EdgeId.ToString()
}

func (this *DEdge) Exists() bool {
	return this.GetFrom().Exists() && this.GetTo().Exists()
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

func (this *DEdge) ToRDF(graph model.Graph, dateTimeSample string, readable bool) []string {
	edgeQuads := make([]string, 0)

	var fromNodeID string
	if this.from.Exists() {
		fromNodeID = this.from.GetFormatedUid()
	} else {
		fromNodeID = this.from.GetEid(readable)
	}

	var toNodeID string
	if this.to.Exists() {
		toNodeID = this.to.GetFormatedUid()
	} else {
		toNodeID = this.to.GetEid(readable)
	}

	relation := util.ReplaceCharacter(this.GetType(), _targetRegex, _replacement, true)

	/* Building facets */
	/* Should key field allow update? */
	var buffer bytes.Buffer
	numOfFacets := 0
	attributes := this.edge.GetAttributes()
	if nil != attributes {
		for attrname, attribute := range attributes {
			if this.typeName != attrname {

				if nil == attribute.GetValue() {
					continue
				}

				attrStringValue, err := util.ToString(attribute.GetValue(), attribute.GetType().String(), dateTimeSample)
				if nil != err {
					log.Debug("(ToRDF) Formatting error : ", err, ", origData : ", attribute.GetValue())
					continue
				}

				if 0 != numOfFacets {
					buffer.WriteString(", ")
				}
				buffer.WriteString(this.GetCanonicalAttributeName(attrname))

				buffer.WriteString("=\"")
				buffer.WriteString(attrStringValue)
				buffer.WriteString("\"")

				numOfFacets++
			}
		}
	}

	var facets string
	if numOfFacets > 0 {
		facets = fmt.Sprintf("(%s)", buffer.String())
	} else {
		facets = buffer.String()
	}

	edgeQuads = append(edgeQuads, fmt.Sprintf(IMPORT_EDGE_FORMAT, fromNodeID, relation, toNodeID, facets))

	return edgeQuads
}
