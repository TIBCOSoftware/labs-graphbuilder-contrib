/*
 * Copyright © 2020. TIBCO Software Inc.
 * This file is subject to the license terms contained
 * in the license file that is distributed with this file.
 */
package dgraph

import (
	"bytes"
	"fmt"
	"strings"

	"github.com/TIBCOSoftware/labs-graphbuilder-lib/util"
)

const (
	GRAPH_MODEL_ID   = "graph_builder_model_id"
	DATE_TIME_SAMPLE = "2006-01-02"
)

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
