/*
 * Copyright © 2020. TIBCO Software Inc.
 * This file is subject to the license terms contained
 * in the license file that is distributed with this file.
 */
package tgdb

import (
	//"fmt"
	"reflect"
	"strings"

	"tgdb"
	"tgdb/impl"
)

func BuildMetadata(metadata impl.GraphMetadata) map[string]interface{} {

	data := make(map[string]interface{})

	nodeId2NameMap := make(map[int]string)
	nodeTypes, _ := metadata.GetNodeTypes()
	nodeTypeInfos := make([]map[string]interface{}, len(nodeTypes))
	data["nodeTypes"] = nodeTypeInfos
	for index, nodeType := range nodeTypes {
		nodeTypeInfo := make(map[string]interface{})
		id := nodeType.GetEntityTypeId()
		name := nodeType.GetName()
		nodeTypeInfo["id"] = id
		nodeTypeInfo["name"] = name
		nodeTypeInfo["systemType"] = nodeType.GetSystemType()
		nodeId2NameMap[id] = name

		pkeyAttributeDescriptors := nodeType.GetPKeyAttributeDescriptors()
		pkeyAttributeDescriptorInfos := make([]map[string]interface{}, len(pkeyAttributeDescriptors))
		nodeTypeInfo["pkeyAttributeDescriptors"] = pkeyAttributeDescriptorInfos
		for index, pkeyAttributeDescriptor := range pkeyAttributeDescriptors {
			pkeyAttributeDescriptorInfo := make(map[string]interface{})
			pkeyAttributeDescriptorInfo["name"] = pkeyAttributeDescriptor.GetName()
			pkeyAttributeDescriptorInfo["type"] = pkeyAttributeDescriptor.GetAttrType()
			pkeyAttributeDescriptorInfo["attributeId"] = pkeyAttributeDescriptor.GetAttributeId()
			pkeyAttributeDescriptorInfo["scale"] = pkeyAttributeDescriptor.GetScale()
			pkeyAttributeDescriptorInfo["precision"] = pkeyAttributeDescriptor.GetPrecision()
			pkeyAttributeDescriptorInfo["systemType"] = pkeyAttributeDescriptor.GetSystemType()
			pkeyAttributeDescriptorInfos[index] = pkeyAttributeDescriptorInfo
		}

		attributeDescriptors := nodeType.GetAttributeDescriptors()
		attributeDescriptorInfos := make([]map[string]interface{}, len(attributeDescriptors))
		nodeTypeInfo["attributeDescriptors"] = attributeDescriptorInfos
		for index, attributeDescriptor := range attributeDescriptors {
			attributeDescriptorInfo := make(map[string]interface{})
			attributeDescriptorInfo["name"] = attributeDescriptor.GetName()
			attributeDescriptorInfo["type"] = attributeDescriptor.GetAttrType()
			attributeDescriptorInfo["attributeId"] = attributeDescriptor.GetAttributeId()
			attributeDescriptorInfo["scale"] = attributeDescriptor.GetScale()
			attributeDescriptorInfo["precision"] = attributeDescriptor.GetPrecision()
			attributeDescriptorInfo["systemType"] = attributeDescriptor.GetSystemType()
			attributeDescriptorInfos[index] = attributeDescriptorInfo
		}

		nodeTypeInfos[index] = nodeTypeInfo
	}

	edgeTypes, _ := metadata.GetEdgeTypes()
	edgeTypeInfos := make([]map[string]interface{}, len(edgeTypes))
	data["edgeTypes"] = edgeTypeInfos
	for index, edgeType := range edgeTypes {
		edgeTypeInfo := make(map[string]interface{})
		edgeTypeInfo["id"] = edgeType.GetEntityTypeId()
		edgeTypeInfo["name"] = edgeType.GetName()
		edgeTypeInfo["systemType"] = edgeType.GetSystemType()

		fromNode := edgeType.GetFromNodeType()
		//		fmt.Println("============================")
		//		fmt.Println(edgeType.GetEntityTypeId())
		//		fmt.Println(edgeType.GetFromTypeId())
		//		fmt.Println(edgeType.GetFromNodeType())
		//		fmt.Println("============================")
		fromNodeInfo := make(map[string]interface{})
		edgeTypeInfo["fromNodeType"] = fromNodeInfo
		if nil != fromNode {
			fromNodeInfo["id"] = fromNode.GetEntityTypeId()
			fromNodeInfo["name"] = fromNode.GetName()
		} else {
			fromNodeId := edgeType.GetFromTypeId()
			fromNodeInfo["id"] = fromNodeId
			fromNodeInfo["name"] = nodeId2NameMap[fromNodeId]
		}

		toNode := edgeType.GetToNodeType()
		toNodeInfo := make(map[string]interface{})
		edgeTypeInfo["toNodeType"] = toNodeInfo
		if nil != toNode {
			toNodeInfo["id"] = toNode.GetEntityTypeId()
			toNodeInfo["name"] = toNode.GetName()
		} else {
			toNodeId := edgeType.GetToTypeId()
			toNodeInfo["id"] = toNodeId
			toNodeInfo["name"] = nodeId2NameMap[toNodeId]
		}

		attributeDescriptors := edgeType.GetAttributeDescriptors()
		attributeDescriptorInfos := make([]map[string]interface{}, len(attributeDescriptors))
		edgeTypeInfo["attributeDescriptors"] = attributeDescriptorInfos
		for index, attributeDescriptor := range attributeDescriptors {
			attributeDescriptorInfo := make(map[string]interface{})
			attributeDescriptorInfo["name"] = attributeDescriptor.GetName()
			attributeDescriptorInfo["type"] = attributeDescriptor.GetAttrType()
			attributeDescriptorInfo["attributeId"] = attributeDescriptor.GetAttributeId()
			attributeDescriptorInfo["scale"] = attributeDescriptor.GetScale()
			attributeDescriptorInfo["precision"] = attributeDescriptor.GetPrecision()
			attributeDescriptorInfo["systemType"] = attributeDescriptor.GetSystemType()
			attributeDescriptorInfos[index] = attributeDescriptorInfo
		}

		edgeTypeInfos[index] = edgeTypeInfo
	}

	/*
		attributeDescriptors, _ := metadata.GetAttributeDescriptors()
		attributeDescriptorInfos := make([]map[string]interface{}, len(attributeDescriptors))
		data["attributeDescriptors"] = attributeDescriptorInfos
		for index, attributeDescriptor := range attributeDescriptors {
			attributeDescriptorInfo := make(map[string]interface{})
			attributeDescriptorInfo["name"] = attributeDescriptor.GetName()
			attributeDescriptorInfo["type"] = attributeDescriptor.GetAttrType()
			attributeDescriptorInfo["attributeId"] = attributeDescriptor.GetAttributeId()
			attributeDescriptorInfo["scale"] = attributeDescriptor.GetScale()
			attributeDescriptorInfo["precision"] = attributeDescriptor.GetPrecision()
			attributeDescriptorInfo["systemType"] = attributeDescriptor.GetSystemType()
			attributeDescriptorInfos[index] = attributeDescriptorInfo
		}
	*/
	return data
}

/*
func BuildGraph(tgdb *TGDBService, entity types.TGEntity, result map[string]interface{}) {
	tgResult := make(map[string][]types.TGEntity)
	tgResult["nodes"] = make([]types.TGEntity, 0)
	tgResult["edges"] = make([]types.TGEntity, 0)
	//	fmt.Println("------------------------------> Node : ", node)
	traverse(tgResult, entity, 0)
	//	fmt.Println("------------------------------> tgResult : ", tgResult)
	buildResult(tgdb, tgResult, result)
}
*/

func traverse(
	tgResult map[string]([]tgdb.TGEntity),
	entity tgdb.TGEntity,
	currDepth int) {
	//	fmt.Println("0000000000000 Entity Type : ", reflect.TypeOf(entity).String())
	if "*model.Node" == reflect.TypeOf(entity).String() {
		//fmt.Println("Add Node : ", entity.GetVirtualId())
		node := entity.(tgdb.TGNode)
		tgResult["nodes"] = append(tgResult["nodes"], node)
		for _, edge := range node.GetEdges() {
			if nil != edge {
				if !contains(tgResult["edges"], edge) {
					currDepth += 1
					traverse(tgResult, edge, currDepth)
				}
			}
		}
	} else if "*model.Edge" == reflect.TypeOf(entity).String() {
		//fmt.Println("Add Edge : ", entity.GetVirtualId())
		edge := entity.(tgdb.TGEdge)
		tgResult["edges"] = append(tgResult["edges"], edge)
		if nil != edge.GetVertices() {
			for _, node := range edge.GetVertices() {
				if nil != node {
					if !contains(tgResult["nodes"], node) {
						traverse(tgResult, node, currDepth)
					}
				}
			}
		}
	}
}

func BuildNode(tgdb *TGDBService, node tgdb.TGNode) map[string]interface{} {
	aNode := make(map[string]interface{})
	aNode["_type"] = node.GetEntityType().GetName()
	id := ExtractNodeKeyAttrValue(tgdb, node)
	if nil != id {
		aNode["id"] = id
	}
	attributes, _ := node.GetAttributes()
	for _, attr := range attributes {
		aNode[attr.GetAttributeDescriptor().GetName()] = attr.GetValue()
	}
	return aNode
}

func BuildResult(tgdbs *TGDBService, tgResult map[string]interface{}) map[string]interface{} {
	result := make(map[string]interface{})
	if nil != tgResult {
		if nil != tgResult["result"] {
			result["result"] = tgResult["result"]
		}

		if nil != tgResult["nodes"] {
			nodeInfos := make([]map[string]interface{}, 0)
			//fmt.Print("(1)")
			for _, entity := range tgResult["nodes"].(map[int64]tgdb.TGEntity) {
				//fmt.Println("Transform Node : ", entity.GetVirtualId())
				node := entity.(tgdb.TGNode)
				//fmt.Print("(2)")
				if nil != node && nil != node.GetEntityType() {
					aNode := make(map[string]interface{})
					//fmt.Print("(3)")
					aNode["_type"] = node.GetEntityType().GetName()
					//fmt.Print("(4)")
					id := ExtractNodeKeyAttrValue(tgdbs, node)
					//fmt.Print("(5)")
					if nil != id {
						aNode["id"] = id
					}
					attributes, _ := node.GetAttributes()
					//fmt.Print("(6)")
					for _, attr := range attributes {
						aNode[attr.GetAttributeDescriptor().GetName()] = attr.GetValue()
					}
					//fmt.Print("(7)")
					nodeInfos = append(nodeInfos, aNode)
					result["nodes"] = nodeInfos
				}
			}
		}

		if nil != tgResult["edges"] {
			edgeInfos := make([]map[string]interface{}, 0)
			//fmt.Print("(8)")
			for _, entity := range tgResult["edges"].(map[int64]tgdb.TGEntity) {
				//fmt.Println("Transform Edge : ", entity.GetVirtualId())
				edge := entity.(tgdb.TGEdge)
				//fmt.Print("(9)")
				if nil != edge {
					if nil != edge.GetEntityType() &&
						nil != edge.GetVertices() &&
						nil != edge.GetVertices()[0] &&
						nil != edge.GetVertices()[1] {
						//fmt.Print("(9.1) = ", edge.GetVertices())
						fromNode := edge.GetVertices()[0]
						fromNodeId := ExtractNodeKeyAttrValue(tgdbs, fromNode)
						//fmt.Print("(9.2)")
						toNode := edge.GetVertices()[1]
						toNodeId := ExtractNodeKeyAttrValue(tgdbs, toNode)
						//fmt.Print("(10)")
						if nil == fromNodeId || nil == toNodeId {
							continue
						}

						anEdge := make(map[string]interface{})
						anEdge["_type"] = edge.(tgdb.TGEdge).GetEntityType().GetName()
						anEdge["fromNode"] = fromNodeId
						anEdge["fromNodeType"] = fromNode.GetEntityType().GetName()
						anEdge["toNode"] = toNodeId
						anEdge["toNodeType"] = toNode.GetEntityType().GetName()
						//fmt.Print("(11)")
						attributes, _ := edge.GetAttributes()
						for _, attr := range attributes {
							//fmt.Print("(12)")
							anEdge[attr.GetAttributeDescriptor().GetName()] = attr.GetValue()
						}
						//fmt.Print("(13)")
						edgeInfos = append(edgeInfos, anEdge)
						result["edges"] = edgeInfos
					}
				}
			}
		}
	}

	//fmt.Print("result = ", result)
	return result
}

func contains(arrays []tgdb.TGEntity, target tgdb.TGEntity) bool {
	for _, element := range arrays {
		//fmt.Println("element.GetVirtualId() = ", element.GetVirtualId(), ", target.GetVirtualId() = ", target.GetVirtualId())
		if element.GetVirtualId() == target.GetVirtualId() {
			//fmt.Println("contains = true ")
			return true
		}
	}
	//fmt.Println("contains = false ")
	return false
}

func ExtractNodeKeyAttrValue(tgdb *TGDBService, node tgdb.TGNode) []interface{} {
	//fmt.Println("node.GetEntityType() = ", node.GetEntityType())
	//fmt.Println("node.GetEntityType().GetName() = ", node.GetEntityType().GetName())
	keyFields := tgdb.GetNodeKeyfields(node.GetEntityType().GetName())
	//fmt.Println("keyFields = ", keyFields)
	if nil == keyFields {
		return nil
	}
	key := make([]interface{}, 0)
	for _, keyField := range keyFields {
		//fmt.Println("keyField = ", keyField)
		keyObj := node.GetAttribute(keyField)
		if nil != keyObj {
			key = append(key, keyObj.GetValue())
		}
	}

	return key
}

func ExtractEdgeKey(tgdb *TGDBService, edge tgdb.TGEdge) []interface{} {

	return nil
}

var replacer = strings.NewReplacer("'", "\\'")

func EscapeIllegalChar(data string) string {
	return replacer.Replace(data)
}

func TrimWhiteSpace(data string) string {
	return strings.NewReplacer(" ", "").Replace(data)
}

func GetTGDBType(dataType string) string {
	switch dataType {
	case "String":
		return "string"
	case "Integer":
		return "int"
	case "Long":
		return "int"
	case "Boolean":
		return "boolean"
	case "Double":
		return "double"
	case "Date": /* eg: 2006-01-02T15:04:05.999999999+10:00 or 2006-01-02T15:04:05.999999999 */
		return "date"
	}

	return "string"
}
