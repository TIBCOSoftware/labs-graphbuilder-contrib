/*
 * Copyright © 2020. TIBCO Software Inc.
 * This file is subject to the license terms contained
 * in the license file that is distributed with this file.
 */
package v2

import (
	"bytes"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/TIBCOSoftware/labs-graphbuilder-lib/dbservice/dgraph"
	"github.com/TIBCOSoftware/labs-graphbuilder-lib/util"
)

func buildSchema(
	explicitType bool,
	typeName string,
	targetRegex string,
	replacement string,
	addPrefixToAttr bool,
	userSchema interface{},
	model map[string]interface{},
) string {
	var schema map[string]interface{}
	var schemaObject interface{}

	if nil != userSchema {
		err := json.Unmarshal([]byte(userSchema.(string)), &schemaObject)
		if nil != err {
			fmt.Println(err)
		}
	}

	if nil != schemaObject {
		schema = schemaObject.(map[string]interface{})
	} else {
		schema = make(map[string]interface{})
	}
	fmt.Println("^^^^^^^^^^^^ User Schema ^^^^^^^^^^^^^^^^")
	fmt.Println(schema)
	fmt.Println("^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^")

	nodeModel := model["nodes"].(map[string]interface{})
	for _, nodeType := range nodeModel["types"].([]string) {
		attrTypeMap := nodeModel["attrTypeMap"].(map[string](map[string]string))[nodeType]
		for _, pkey := range nodeModel["keyMap"].(map[string][]string)[nodeType] {
			attributeName := GetCanonicalAttributeName(nodeType, pkey, targetRegex, replacement, addPrefixToAttr, true)
			dgraphDataType := GetDgraphType(attrTypeMap[pkey])

			var properties map[string]interface{}
			if nil != schema[attributeName] {
				properties = schema[attributeName].(map[string]interface{})
			} else {
				properties = make(map[string]interface{})
				schema[attributeName] = properties
			}

			if nil == properties["type"] {
				properties["type"] = dgraphDataType
			}

			if nil == properties["index"] {
				indexAttrs := make([]interface{}, 1)
				if "string" == dgraphDataType {
					indexAttrs[0] = "exact"
				} else {
					indexAttrs[0] = dgraphDataType

				}
				properties["index"] = indexAttrs
			}
		}
	}

	edgeModel := model["edges"].(map[string]interface{})
	edgeDirectionMap := edgeModel["directionMap"].(map[string]int)
	for _, edgeType := range edgeModel["types"].([]string) {
		attrTypeMap := edgeModel["attrTypeMap"].(map[string](map[string]string))[edgeType]

		allowRreverse := false
		if 2 == edgeDirectionMap[edgeType] {
			allowRreverse = true
		}

		canonicalEdgeType := util.ReplaceCharacter(edgeType, targetRegex, replacement, true)

		var eProperties map[string]interface{}
		if nil != schema[canonicalEdgeType] {
			eProperties = schema[canonicalEdgeType].(map[string]interface{})
		} else {
			eProperties = make(map[string]interface{})
			schema[canonicalEdgeType] = eProperties
		}

		if nil == eProperties["type"] {
			eProperties["type"] = "[uid]"
		}

		if nil == eProperties["extraFlags"] {
			if allowRreverse {
				extraFlags := make([]interface{}, 1)
				extraFlags[0] = "reverse"
				eProperties["extraFlags"] = extraFlags
			}
		}

		for _, pkey := range edgeModel["keyMap"].(map[string][]string)[edgeType] {
			attributeName := GetCanonicalAttributeName(edgeType, pkey, targetRegex, replacement, addPrefixToAttr, true)
			dgraphDataType := GetDgraphType(attrTypeMap[pkey])

			var properties map[string]interface{}
			if nil != schema[attributeName] {
				properties = schema[attributeName].(map[string]interface{})
			} else {
				properties = make(map[string]interface{})
				schema[attributeName] = properties
			}

			if nil == properties["type"] {
				properties["type"] = dgraphDataType
			}

			if nil == properties["index"] {
				indexAttrs := make([]interface{}, 1)
				if "string" == dgraphDataType {
					indexAttrs[0] = "exact"
				} else {
					indexAttrs[0] = dgraphDataType

				}
				properties["index"] = indexAttrs
			}
		}
	}

	fmt.Println("^^^^^^^^^^^^ final Schema ^^^^^^^^^^^^^^^")
	fmt.Println(schema)
	fmt.Println("^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^")

	var aSchema string
	processedSchemas := make([]string, 0)
	aSchema = fmt.Sprintf("%s: string @index(exact) . \n", dgraph.GRAPH_MODEL_ID)
	processedSchemas = append(processedSchemas, aSchema)
	if explicitType {
		/* It means there will be an universal tag (for all entities)
		   in DB to indicate type(label) of entity */
		if "" == typeName {
			typeName = "type"
		}
		aSchema = fmt.Sprintf("%s: string @index(exact) . \n", typeName)
		processedSchemas = append(processedSchemas, aSchema)
	}

	for name, value := range schema {
		properties := value.(map[string]interface{})
		if nil != properties["index"] {
			index := properties["index"].([]interface{})
			var indexAttrs bytes.Buffer
			for i, indexAttr := range index {
				if 0 != i {
					indexAttrs.WriteString(", ")
				}
				indexAttrs.WriteString(indexAttr.(string))
			}

			aSchema = fmt.Sprintf(
				"%s: %s @index(%s) . \n",
				name,
				properties["type"],
				indexAttrs.String(),
			)
		} else {
			if nil != properties["extraFlags"] {
				extraFlags := properties["extraFlags"].([]interface{})
				var extraFlagsB bytes.Buffer
				for i, extraFlag := range extraFlags {
					if 0 != i {
						extraFlagsB.WriteString(" ")
					}
					extraFlagsB.WriteString("@")
					extraFlagsB.WriteString(extraFlag.(string))
				}

				aSchema = fmt.Sprintf(
					"%s: %s %s . \n",
					name,
					properties["type"],
					extraFlagsB.String(),
				)
			} else {
				aSchema = fmt.Sprintf(
					"%s: %s . \n",
					name,
					properties["type"],
				)
			}
		}
		if !util.SliceContains(processedSchemas, aSchema) {
			processedSchemas = append(processedSchemas, aSchema)
		}
	}

	var schemas bytes.Buffer
	for _, schema := range processedSchemas {
		schemas.WriteString(schema)
	}
	fmt.Println("***************** schema query ********************")
	fmt.Println(schemas.String())
	fmt.Println("***************************************************")

	return schemas.String()
}

func ReadableExternalId(nodeType string, key []interface{}) string {
	var eid bytes.Buffer
	eid.WriteString(nodeType)
	for _, element := range key {
		eid.WriteString("_")
		element, _ = util.ConvertToString(element, "")
		eid.WriteString(element.(string))
	}

	return eid.String()
}

func GetCanonicalAttributeName(
	entityType string,
	attributeName string,
	targetRegex string,
	replacement string,
	addPrefixToAttr bool,
	doReplacement bool) string {

	if addPrefixToAttr {
		return util.ReplaceCharacter(fmt.Sprintf("%s_%s", entityType, attributeName), targetRegex, replacement, doReplacement)
	} else {
		return util.ReplaceCharacter(attributeName, targetRegex, replacement, doReplacement)
	}
}

func TrimWhiteSpace(data string) string {
	return strings.NewReplacer(" ", "").Replace(data)
}

func GetDgraphType(dataType string) string {
	switch dataType {
	case "String":
		return "string"
	case "Integer":
		return "int"
	case "Long":
		return "int"
	case "Boolean":
		return "bool"
	case "Double":
		return "float"
	case "Date": /* eg: 2006-01-02T15:04:05.999999999+10:00 or 2006-01-02T15:04:05.999999999 */
		return "dateTime"
	}

	// how about dgraph geo type?
	// geometries stored using go-geom

	return "string"
}
