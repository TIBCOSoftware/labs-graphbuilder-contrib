/*
 * Copyright © 2020. TIBCO Software Inc.
 * This file is subject to the license terms contained
 * in the license file that is distributed with this file.
 */
package query

import (
	"reflect"

	"github.com/TIBCOSoftware/labs-graphbuilder-lib/model"
	"github.com/TIBCOSoftware/labs-graphbuilder-lib/util"
)

func printTraversal(tNode *model.TraversalNode) {
	tgResult := make(map[string][]interface{})
	tgResult["nodes"] = make([]interface{}, 0)
	tgResult["edges"] = make([]interface{}, 0)
	traverse(tgResult, tNode, 0)
}

func traverse(
	tgResult map[string]([]interface{}),
	entity interface{},
	currDepth int) {
	//	fmt.Println("Entity Type : ", reflect.TypeOf(entity).String())
	if "*model.TraversalNode" == reflect.TypeOf(entity).String() {
		node := entity.(*model.TraversalNode)
		//		fmt.Println("A Node --> ", node.NodeId.ToString())
		tgResult["nodes"] = append(tgResult["nodes"], node)
		for _, edge := range node.GetAllEdges() {
			if nil != edge {
				if !util.Contains(tgResult["edges"], edge) {
					currDepth += 1
					traverse(tgResult, edge, currDepth)
				}
			}
		}
	} else if "*model.TraversalEdge" == reflect.TypeOf(entity).String() {
		edge := entity.(*model.TraversalEdge)
		//		fmt.Println("An Edge ----> ", edge.EdgeId.ToString())
		tgResult["edges"] = append(tgResult["edges"], edge)
		if nil != edge.GetAllNodes() {
			for _, node := range edge.GetAllNodes() {
				if nil != node {
					if !util.Contains(tgResult["nodes"], node) {
						traverse(tgResult, node, currDepth)
					}
				}
			}
		}
	}
}
