/*
 * Copyright © 2020. TIBCO Software Inc.
 * This file is subject to the license terms contained
 * in the license file that is distributed with this file.
 */
package model

import (
	"encoding/json"
	"fmt"
	"strings"
)

//-================================-//
//     Define Attributefinition
//-================================-//

type Attributefinition struct {
	name     string
	dataType DataType
}

func (this *Attributefinition) GetDataType() DataType {
	return this.dataType
}

func NewAttributeModel(name string, dataTypeStr string) *Attributefinition {
	dataType, _ := ToTypeEnum(dataTypeStr)
	return &Attributefinition{name, dataType}
}

//-================================-//
//     Define EntityDefinition
//-===============================-//

type EntityDefinition struct {
	_keyDefinition []string
	_type          string
	_attributes    map[string]*Attributefinition
}

func (this *EntityDefinition) GetAttributeDefinitions() map[string]*Attributefinition {
	return this._attributes
}

func NewEntityModel(entityInfo map[string]interface{}) *EntityDefinition {
	var key []string
	if nil != entityInfo["key"] {
		keyInfo := entityInfo["key"].([]interface{})
		key = make([]string, len(keyInfo))
		for i := 0; i < len(key); i++ {
			key[i] = keyInfo[i].(string)
		}
	}

	var entityType string
	if nil != entityInfo["name"] {
		entityType = entityInfo["name"].(string)
	} else {
		fmt.Println("Type name is not defined for edge!")
	}

	attributesModel := make(map[string]*Attributefinition)

	if nil != entityInfo["attributes"] {
		attributesInfo := entityInfo["attributes"].([]interface{})
		for _, attributeInfo := range attributesInfo {
			attribute := attributeInfo.(map[string]interface{})
			attrName := attribute["name"].(string)

			var attrType string
			if nil != attribute["type"] {
				attrType = attribute["type"].(string)
			} else {
				attrType = "String"
			}
			attributesModel[attrName] = NewAttributeModel(attrName, attrType)
		}
	}

	return &EntityDefinition{key, entityType, attributesModel}
}

//-============================-//
//     Define NodeDefinition
//-============================-//

type NodeDefinition struct {
	*EntityDefinition
}

func NewNodeModel(nodeInfo map[string]interface{}) *NodeDefinition {
	var nodeModel NodeDefinition
	nodeModel.EntityDefinition = NewEntityModel(nodeInfo)
	return &nodeModel
}

//-============================-//
//     Define EdgeDefinition
//-============================-//

type EdgeDirection int

const (
	Nondirectional EdgeDirection = 0
	Directional    EdgeDirection = 1
	Bidirectional  EdgeDirection = 2
)

func (this EdgeDirection) int() int {
	index := [...]int{
		0,
		1,
		2}
	return index[this]
}

func (this EdgeDirection) String() string {
	strName := [...]string{
		"nondirectional",
		"directional",
		"bidirectional"}
	return strName[this]
}

func ToEdgeDirection(direction interface{}) (EdgeDirection, error) {
	switch direction.(type) {
	case string:
		{
			strDirection := direction.(string)
			switch strings.ToLower(strDirection) {
			case "nondirectional", "0":
				return Nondirectional, nil
			case "directional", "1":
				return Directional, nil
			case "bidirectional", "2":
				return Bidirectional, nil
			default:
				return Directional, fmt.Errorf("Undefined direction (%s), use (Directional)", direction)
			}
		}
	case float64, float32, int:
		{
			var iDirection int
			switch direction.(type) {
			case float64:
				iDirection = int(direction.(float64))
				break
			case float32:
				iDirection = int(direction.(float32))
				break
			case int:
				iDirection = direction.(int)
			}

			switch iDirection {
			case 0:
				return Nondirectional, nil
			case 1:
				return Directional, nil
			case 2:
				return Bidirectional, nil
			default:
				return Directional, fmt.Errorf("Undefined direction code (%s), use (Directional)", direction)
			}
		}
	}
	return Directional, fmt.Errorf("Direction not set (%s), use (Directional)", direction)
}

type EdgeDefinition struct {
	_direction    EdgeDirection
	_fromNodeType string
	_toNodeType   string
	*EntityDefinition
}

func NewEdgeModel(edgeInfo map[string]interface{}) *EdgeDefinition {
	var edgeModel EdgeDefinition

	var err error
	edgeModel._direction, err = ToEdgeDirection(edgeInfo["direction"])
	if nil != err {
		/* print some information? */
		fmt.Println(edgeInfo["name"], " - ", err)
	}

	edgeModel._fromNodeType = edgeInfo["from"].(string)
	edgeModel._toNodeType = edgeInfo["to"].(string)
	edgeModel.EntityDefinition = NewEntityModel(edgeInfo)

	return &edgeModel
}

//-============================-//
//     Define GraphDefinition
//-============================-//

type GraphDefinition struct {
	_id              string
	_nodeDefinitions map[string]*NodeDefinition
	_edgeDefinitions map[string]*EdgeDefinition
}

func (gd *GraphDefinition) GetId() string {
	return gd._id
}

func (gd *GraphDefinition) GetNodeDefinition(nodeType string) *NodeDefinition {
	return gd._nodeDefinitions[nodeType]
}

func (gd *GraphDefinition) GetNodeDefinitions() map[string]*NodeDefinition {
	return gd._nodeDefinitions
}

func (gd *GraphDefinition) GetEdgeDefinition(edgeType string) *EdgeDefinition {
	return gd._edgeDefinitions[edgeType]
}

func (gd *GraphDefinition) GetEdgeDefinitions() map[string]*EdgeDefinition {
	return gd._edgeDefinitions
}

func (gd *GraphDefinition) Export() map[string]interface{} {
	nodeTypes := make([]string, 0)
	nodeKeyMap := make(map[string][]string)
	attrTypeMap := make(map[string](map[string]string))
	for nodeType, definition := range gd._nodeDefinitions {
		nodeTypes = append(nodeTypes, nodeType)
		nodeKeyMap[nodeType] = definition._keyDefinition
		nodeAttrTypeMap := make(map[string]string)
		attrTypeMap[nodeType] = nodeAttrTypeMap
		for attrName, attrDef := range definition._attributes {
			nodeAttrTypeMap[attrName] = attrDef.dataType.String()
		}
	}
	nodeModels := make(map[string]interface{})
	nodeModels["types"] = nodeTypes
	nodeModels["keyMap"] = nodeKeyMap
	nodeModels["attrTypeMap"] = attrTypeMap

	edgeTypes := make([]string, 0)
	edgeDirectionMap := make(map[string]int)
	edgeKeyMap := make(map[string][]string)
	edgeVertexes := make(map[string][]string)
	attrTypeMap = make(map[string](map[string]string))
	for edgeType, definition := range gd._edgeDefinitions {
		edgeTypes = append(edgeTypes, edgeType)
		edgeDirectionMap[edgeType] = definition._direction.int()
		edgeKeyMap[edgeType] = definition._keyDefinition
		edgeVertexes[edgeType] = []string{definition._fromNodeType, definition._toNodeType}
		edgeAttrTypeMap := make(map[string]string)
		attrTypeMap[edgeType] = edgeAttrTypeMap
		for attrName, attrDef := range definition._attributes {
			edgeAttrTypeMap[attrName] = attrDef.dataType.String()
		}
	}
	edgeModels := make(map[string]interface{})
	edgeModels["types"] = edgeTypes
	edgeModels["directionMap"] = edgeDirectionMap
	edgeModels["keyMap"] = edgeKeyMap
	edgeModels["vertexes"] = edgeVertexes
	edgeModels["attrTypeMap"] = attrTypeMap

	graphModel := make(map[string]interface{})
	graphModel["nodes"] = nodeModels
	graphModel["edges"] = edgeModels

	return graphModel
}

func NewGraphModel(id string, graphmodel string) (*GraphDefinition, error) {
	var rootObject interface{}

	err := json.Unmarshal([]byte(graphmodel), &rootObject)
	if nil != err {
		return nil, err
	}

	return parseTGBModel(id, rootObject), nil
}

func parseTGBModel(id string, rootObject interface{}) *GraphDefinition {
	//fmt.Println("id = ", id, ", root obj = ", rootObject)
	dataMap := rootObject.(map[string]interface{})

	nodeModels := make(map[string]*NodeDefinition)
	nodes := dataMap["nodes"].([]interface{})
	for _, node := range nodes {
		nodeInfo := node.(map[string]interface{})
		nodeType := nodeInfo["name"].(string)
		nodeModels[nodeType] = NewNodeModel(nodeInfo)
		//fmt.Println("nodeType = ", nodeType, ", nodeInfo = ", nodeInfo)
	}

	edgeModels := make(map[string]*EdgeDefinition)
	edges := dataMap["edges"].([]interface{})
	for _, edge := range edges {
		edgeInfo := edge.(map[string]interface{})
		//fmt.Println("edgeType = ", edgeInfo["name"], ", edgeInfo = ", edgeInfo)
		edgeModels[edgeInfo["name"].(string)] = NewEdgeModel(edgeInfo)
	}
	//fmt.Println("nodeModels = ", nodeModels, ", edgeModels = ", edgeModels)

	graphModel := &GraphDefinition{id, nodeModels, edgeModels}
	//fmt.Println("graphModel = ", graphModel)
	return graphModel
}
