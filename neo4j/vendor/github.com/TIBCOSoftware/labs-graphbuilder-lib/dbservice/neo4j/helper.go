/*
 * Copyright © 2020. TIBCO Software Inc.
 * This file is subject to the license terms contained
 * in the license file that is distributed with this file.
 */
package neo4j

import (
	"github.com/TIBCOSoftware/labs-graphbuilder-lib/util"
)

func GetCanonicalAttributeName(
	attributeName string,
	targetRegex string,
	replacement string,
	doReplacement bool) string {

	return util.ReplaceCharacter(attributeName, targetRegex, replacement, doReplacement)
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
