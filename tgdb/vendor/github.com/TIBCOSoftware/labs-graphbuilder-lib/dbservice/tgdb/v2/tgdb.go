/*
 * Copyright © 2020. TIBCO Software Inc.
 * This file is subject to the license terms contained
 * in the license file that is distributed with this file.
 */
package tgdb

import (
	"errors"
	"fmt"
	"strings"
	"sync"
	"time"

	"github.com/TIBCOSoftware/labs-graphbuilder-lib/model"
	"github.com/TIBCOSoftware/tgdb-client/client/goAPI/connection"
	"github.com/TIBCOSoftware/tgdb-client/client/goAPI/query"
	"github.com/TIBCOSoftware/tgdb-client/client/goAPI/types"
)

type TGDBService struct {
	_url             string
	_user            string
	_password        string
	_keepAlive       bool
	_allowEmptyKey   bool
	_connectionProps map[string]string
	_tgConnection    types.TGConnection
	_gof             types.TGGraphObjectFactory
	_gmd             types.TGGraphMetadata
	_keyMap          map[string][]string
	_nodeAttrMap     map[string]map[string]types.TGAttributeDescriptor
	_edgeAttrMap     map[string]map[string]types.TGAttributeDescriptor
	_mux             sync.Mutex
}

func (this *TGDBService) ensureConnection() types.TGError {
	var err types.TGError
	if nil == this._tgConnection {
		logger.Debug("XXXXXXXXXXXXXXXXXXXXXX Start connecting XXXXXXXXXXXXXXXXXXXXXXXXXXXXXX")
		logger.Debug("[TGDBService::ensureConnection] Will try to connect ..........")
		logger.Debug(fmt.Sprintf("[TGDBService::ensureConnection] url = %s ", this._url))
		logger.Debug(fmt.Sprintf("[TGDBService::ensureConnection] user = %s ", this._user))
		logger.Debug(fmt.Sprintf("[TGDBService::ensureConnection] password = %s ", this._password))
		logger.Debug(fmt.Sprintf("[TGDBService::ensureConnection] connectionProps = %s ", this._connectionProps))
		defer logger.Debug("XXXXXXXXXXXXXXXXXXXXXX End connecting XXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXX")

		this._tgConnection, err = connection.NewTGConnectionFactory().CreateConnection(this._url, this._user, this._password, this._connectionProps)
		if nil != err {
			this._tgConnection = nil
			return err
		}
		err = this._tgConnection.Connect()
		if nil != err {
			this._tgConnection = nil
			this._gof = nil
			return err
		}

		err := this.fetchMetadata()
		if nil != err {
			return err
		}

		this._gof, err = this._tgConnection.GetGraphObjectFactory()
		if nil != err {
			this._tgConnection = nil
			return err
		}
		this._tgConnection.SetExceptionListener(this)
	}
	return nil
}

//-====================-//
//    Metadata API
//-====================-//

func (this *TGDBService) GetFactory() (types.TGGraphObjectFactory, error) {
	if nil != this._gof {
		return this._gof, nil
	} else {
		err := this.ensureConnection()
		return this._gof, err
	}
}

func (this *TGDBService) GetMetadata() (types.TGGraphMetadata, error) {
	if nil != this._gmd {
		return this._gmd, nil
	} else {
		err := this.ensureConnection()
		return this._gmd, err
	}
}

func (this *TGDBService) GetNodeType(nodeTypeStr string) (types.TGNodeType, error) {
	gdm, err := this.GetMetadata()
	if nil != err {
		return nil, err
	}
	return gdm.GetNodeType(nodeTypeStr)
}

func (this *TGDBService) GetEdgeType(edgeTypeStr string) (types.TGEdgeType, error) {
	gdm, err := this.GetMetadata()
	if nil != err {
		return nil, err
	}

	return gdm.GetEdgeType(edgeTypeStr)
}

func (this *TGDBService) GetNodeKeyfields(nodeTypeStr string) []string {
	gdm, err := this.GetMetadata()
	if nil != err || nil == gdm {
		return nil
	}

	nodeType, err := gdm.GetNodeType(nodeTypeStr)
	if nil != err || nil == nodeType {
		return nil
	}

	pkeyAttributeDescriptors := nodeType.GetPKeyAttributeDescriptors() //.GetAttributeDescriptors()
	attrLenth := len(pkeyAttributeDescriptors)

	key := make([]string, attrLenth)
	for i := 0; i < attrLenth; i++ {
		key[i] = pkeyAttributeDescriptors[i].GetName()
	}
	return key
}

func (this *TGDBService) commit(readyEntity EntityKeeper) types.TGError {

	logger.Debug("[TGDBService::commit] Entering ........ ")
	logger.Debug("[TGDBService::commit] Exit ........ ")

	readyEntity.Populate(this._tgConnection)

	logger.Debug("=========================== Connection ===============================")
	logger.Debug(this._tgConnection)
	logger.Debug("======================================================================")

	_, err := this._tgConnection.Commit()
	if nil != err {
		logger.Error(fmt.Sprintf("[TGDBService::commit] Error while comitting : Message = %v ", err.GetErrorMsg()))
		rollbackErr := this._tgConnection.Rollback()
		if nil != rollbackErr {
			logger.Error(fmt.Sprintf("[TGDBService::commit] Error while rollback : Message = %v ", rollbackErr.GetErrorMsg()))
		}
		logger.Warn("[TGDBService::commit] Reset connection to null after error happens!!")
		this.ResetConnection()
		logger.Warn("[TGDBService::commit] Exit with error........ ")
		return err
	}

	if !this._keepAlive {
		logger.Warn("[TGDBService::commit] Not keep alive so will reset connection to null !!")
		this.ResetConnection()
	}

	return nil
}

func (this *TGDBService) ResetConnection() {

	logger.Warn("[TGDBService::ResetConnection] Entering ...")
	defer logger.Warn("[TGDBService::ResetConnection] Exit ...")

	if nil != this._tgConnection {
		logger.Info("[TGDBService::ResetConnection] Will disconnect from TGDB server !!")
		this._tgConnection.Disconnect()
		logger.Info("[TGDBService::ResetConnection] Disconnected !!")
		this._tgConnection = nil
	}
}

//-====================-//
//       Delete
//-====================-//

func (this *TGDBService) DeleteGraph(filter int, graph map[string]interface{}) error {
	this._mux.Lock()
	defer this._mux.Unlock()
	logger.Debug("[TGDBService::DeleteGraph] Entering ........ ")
	defer logger.Debug("[TGDBService::DeleteGraph] Exit ........ ")

	var readyEntities EntityKeeper
	switch filter {
	case 0:
		readyEntities = NewDeleteEntityKeeper(true, false)
	case 1:
		readyEntities = NewDeleteEntityKeeper(false, true)
	case 2:
		readyEntities = NewDeleteEntityKeeper(true, true)
	}

	nodes := graph["nodes"].(map[string]interface{})
	edges := graph["edges"].(map[string]interface{})
	for edgeId, edge := range edges {
		edgeDetail := edge.(map[string]interface{})
		fromNodeId := edgeDetail["from"].(string)
		toNodeId := edgeDetail["to"].(string)
		err := this.prepareEdge(
			readyEntities,
			edgeId,
			edgeDetail,
			nodes[fromNodeId].(map[string]interface{}),
			nodes[toNodeId].(map[string]interface{}),
			true)

		if nil != err {
			logger.Error(fmt.Sprintf("Unable to prepare edge, error = '%s", err.Error()))
			return err
		}
	}

	for _, node := range nodes {
		nodeDetail := node.(map[string]interface{})
		err := this.prepareNode(readyEntities, nodeDetail)
		if nil != err {
			logger.Error(fmt.Sprintf("Unable to prepare node, error = '%s", err.Error()))
			return err
		}
	}

	return this.commit(readyEntities)
}

//-====================-//
//       Upsert
//-====================-//

func (this *TGDBService) UpsertGraph(graphToo model.Graph, graph map[string]interface{}) error {
	this._mux.Lock()
	defer this._mux.Unlock()
	logger.Debug("[TGDBService::UpsertGraph] Entering ........ ")
	defer logger.Debug("[TGDBService::UpsertGraph] Exit ........ ")

	nodes := graph["nodes"].(map[string]interface{})
	readyEntities := NewUpsertEntityKeeper()
	for edgeId, edge := range graph["edges"].(map[string]interface{}) {
		edgeDetail := edge.(map[string]interface{})
		fromNodeId := edgeDetail["from"].(string)
		toNodeId := edgeDetail["to"].(string)
		err := this.prepareEdge(
			readyEntities,
			edgeId,
			edgeDetail,
			nodes[fromNodeId].(map[string]interface{}),
			nodes[toNodeId].(map[string]interface{}),
			true)

		if nil != err {
			logger.Error(fmt.Sprintf("Unable to prepare edge, error = '%s", err.Error()))
			continue
		}
	}

	for _, node := range nodes {
		nodeDetail := node.(map[string]interface{})
		err := this.prepareNode(readyEntities, nodeDetail)
		if nil != err {
			logger.Error(fmt.Sprintf("Unable to prepare node, error = '%s", err.Error()))
			continue
		}
	}

	return this.commit(readyEntities)
}

func (this *TGDBService) prepareNode(readyEntity EntityKeeper, nodeData map[string]interface{}) error {
	logger.Debug("[TGDBService::prepareNode] Entering ........ ", nodeData)
	defer logger.Debug("[TGDBService::prepareNode] Exit ........ ")

	nodeTypeStr := nodeData["type"].(string)

	var attributes map[string]interface{}
	if nil != nodeData["attributes"] {
		attributes = nodeData["attributes"].(map[string]interface{})
		err := this.checkNode(nodeTypeStr, attributes)
		if nil != err {
			return err
		}
	}

	nodeKey, err := this.buildNodeKey(nodeTypeStr, nodeData)

	if nil != err {
		return err
	}

	var node types.TGNode
	var isNew bool
	nodeWrapper := readyEntity.GetNode(nodeTypeStr, nodeKey)
	if nil == nodeWrapper {
		node, err = this.fetchNode(nodeTypeStr, nodeData)
		if nil != err {
			return err
		}

		if nil != node {
			isNew = false
		} else {
			isNew = true
			node, err = this.createNode(nodeTypeStr)
			if nil != err {
				return err
			}
		}
		nodeWrapper = NewTGNodeWrapper(node, isNew)
		readyEntity.AddNode(nodeTypeStr, nodeKey, nodeWrapper)
	} else {
		isNew = false
		node = *nodeWrapper.node
	}

	if nil != attributes {
		this.populateAttributes(node, attributes)
	}

	return nil
}

func (this *TGDBService) checkNode(nodeType string, attributes map[string]interface{}) error {
	logger.Debug("[TGDBService::checkNode] Entering : nodeType = ", nodeType, ", attributes = ", attributes)
	defer logger.Debug("[TGDBService::checkNode] Exit ........ ")

	nodeTypeObj, err := this.GetNodeType(nodeType)
	if nil != err {
		return err
	}
	if nil == nodeTypeObj {
		return errors.New(fmt.Sprintf("Node type %s is not defined in TGDB.", nodeType))
	}

	attrDescs := nodeTypeObj.GetPKeyAttributeDescriptors()
	var keyAttributeData map[string]interface{}
	for _, attrDesc := range attrDescs {
		name := attrDesc.GetName()
		attrType := attrDesc.GetAttrType()
		/* missing key attribute */
		if nil == attributes[name] {
			if this._allowEmptyKey && 10 == attrType /* string type */ {
				keyAttributeData = make(map[string]interface{})
				keyAttributeData["name"] = name
				keyAttributeData["type"] = "string"
				keyAttributeData["value"] = ""
				attributes[name] = keyAttributeData
			} else {
				return errors.New("(checkNode) Key attribute should not be nil!")
			}
		}
		keyAttributeData = attributes[name].(map[string]interface{})
		if nil == keyAttributeData["value"] && 10 == attrType /* string type */ {
			if this._allowEmptyKey {
				keyAttributeData["value"] = ""
			} else {
				return errors.New("(checkNode) Key attribute should not be nil!")
			}
		}
	}

	return nil
}

func (this *TGDBService) prepareEdge(
	readyEntities EntityKeeper,
	id interface{},
	edgeData map[string]interface{},
	firstNodeData map[string]interface{},
	secondNodeData map[string]interface{},
	autoCommit bool) error {

	logger.Debug("[TGDBService::prepareEdge] Entering : ", firstNodeData, " - ", edgeData, " - ", secondNodeData)
	defer logger.Debug("[TGDBService::prepareEdge] Exit ........ ")

	edgeTypeStr := edgeData["type"].(string)
	attributesData := edgeData["attributes"].(map[string]interface{})
	//	allowDuplicate := false
	//	if nil != edgeData["allowDuplicate"] {
	//		allowDuplicate == edgeData["allowDuplicate"].(bool)
	//	}

	iDirection := edgeData["direction"]

	firstNodeType := firstNodeData["type"].(string)
	if nil != firstNodeData["attributes"] {
		attributes := firstNodeData["attributes"].(map[string]interface{})
		err := this.checkNode(firstNodeType, attributes)
		if nil != err {
			return err
		}
	}

	secondNodeType := secondNodeData["type"].(string)
	if nil != secondNodeData["attributes"] {
		attributes := secondNodeData["attributes"].(map[string]interface{})
		err := this.checkNode(secondNodeType, attributes)
		if nil != err {
			return err
		}
	}

	edgeType, err := this.GetEdgeType(edgeTypeStr)
	if nil != err {
		return err
	}
	var targetEdge types.TGEdge
	edgeWrapper := readyEntities.GetEdge(edgeTypeStr, id)
	if nil == edgeWrapper {
		firstNodeKey, err := this.buildNodeKey(firstNodeType, firstNodeData)
		if nil != err {
			return err
		}
		secondNodeKey, err := this.buildNodeKey(secondNodeType, secondNodeData)
		if nil != err {
			return err
		}

		// Check edge existence in tgdb
		var iEdgeType interface{}
		if nil != edgeType {
			iEdgeType = edgeTypeStr
		}

		var keyAttributes []string
		if nil != edgeData["keyAttributeName"] {
			keyAttributes = edgeData["keyAttributeName"].([]string)
		}

		/* Got all required information so now we querying an edge */
		var fromNodeWrapper *TGNodeWrapper
		var toNodeWrapper *TGNodeWrapper

		fromNode, searchErr := this.search(firstNodeType, firstNodeData, secondNodeType, secondNodeData, iEdgeType, keyAttributes, attributesData)
		if nil != searchErr {
			return searchErr
		}

		var toNode types.TGNode
		if nil != fromNode {
			fromNodeWrapper = NewTGNodeWrapper(fromNode, false)
			edges := fromNode.GetEdges()
			for _, edge := range edges {
				toNode = edge.GetVertices()[1]
				if nil != toNode {
					targetEdge = edge
					break
				}
			}

			if nil != toNode {
				toNodeWrapper = NewTGNodeWrapper(toNode, false)
			}

			if nil != targetEdge {
				edgeWrapper = NewTGEdgeWrapper(targetEdge, false)
			}
		}

		if nil == fromNodeWrapper {
			fromNodeWrapper, err = this.findNode(readyEntities, firstNodeType, firstNodeKey, firstNodeData)
			if nil != err {
				return err
			}
		}

		if nil == toNodeWrapper {
			toNodeWrapper, err = this.findNode(readyEntities, secondNodeType, secondNodeKey, secondNodeData)
			if nil != err {
				return err
			}
		}

		if nil == edgeWrapper {
			ptFromNode := fromNodeWrapper.GetNode()
			ptToNode := toNodeWrapper.GetNode()
			if nil != edgeType {
				targetEdge, err = this._gof.CreateEdgeWithEdgeType(*ptFromNode, *ptToNode, edgeType)
			} else {
				direction := types.DirectionTypeUnDirected
				if nil != iDirection {
					switch iDirection.(int) {
					case 0:
						direction = types.DirectionTypeUnDirected
						break
					case 1:
						direction = types.DirectionTypeDirected
						break
					case 2:
						direction = types.DirectionTypeBiDirectional
						break
					}
				}
				targetEdge, err = this._gof.CreateEdgeWithDirection(*ptFromNode, *ptToNode, direction)
			}

			if nil != err {
				return err
			}
			edgeWrapper = NewTGEdgeWrapper(targetEdge, true)
		}

		if nil != fromNodeWrapper {
			readyEntities.AddNode(firstNodeType, firstNodeKey, fromNodeWrapper)
		}

		if nil != toNodeWrapper {
			readyEntities.AddNode(secondNodeType, secondNodeKey, toNodeWrapper)
		}

		if nil != edgeWrapper {
			readyEntities.AddEdge(edgeTypeStr, id, edgeWrapper)
		}
	} else {
		targetEdge = *edgeWrapper.edge
	}

	if nil != targetEdge && nil != attributesData {
		this.populateAttributes(targetEdge, attributesData)
	}

	return nil
}

//-====================-//
//       Query
//-====================-//

func (this *TGDBService) GetNode(nodeType string, parameter map[string]interface{}) (types.TGEntity, error) {
	return this.fetchNode(nodeType, parameter)
}

func (this *TGDBService) GremlinQuery(para map[string]interface{}) (types.TGResultSet, types.TGError) {
	var queryString string
	if nil == para["query"] {
		return nil, &types.TGDBError{ErrorMsg: "No query parameter defined!"}
	}
	queryObj := para["query"].(map[string]interface{})
	if nil != queryObj {
		if nil != queryObj[Query_QueryString] {
			queryString = queryObj[Query_QueryString].(string)
		}
	}

	option := this.buildQueryOption(para)

	logger.Debug("======================== query ========================")
	logger.Debug(fmt.Sprintf("* queryString           = %s ", queryString))
	logger.Debug(fmt.Sprintf("* Option.prefetchSize   = %d ", option.GetPreFetchSize()))
	logger.Debug(fmt.Sprintf("* Option.edgeLimit      = %d ", option.GetEdgeLimit()))
	logger.Debug(fmt.Sprintf("* Option.traversalDepth = %d ", option.GetTraversalDepth()))
	logger.Debug("-------------------------------------------------------")

	err := this.ensureConnection()
	if nil != err {
		return nil, err
	}

	resultSet, err := this._tgConnection.ExecuteQuery(
		"gremlin://"+queryString,
		option,
	)

	if nil != err {
		return nil, err
	}

	return resultSet, nil
}

func (this *TGDBService) TGQLQuery(para map[string]interface{}) (types.TGResultSet, types.TGError) {
	var queryString string
	var edgeFilter string
	var traversalCondition string
	var endCondition string

	if nil == para["query"] {
		return nil, &types.TGDBError{ErrorMsg: "No query parameter defined!"}
	}
	queryObj := para["query"].(map[string]interface{})
	if nil != queryObj {
		if nil != queryObj[Query_QueryString] {
			queryString = queryObj[Query_QueryString].(string)
		}

		if nil != queryObj[Query_EdgeFilter] {
			edgeFilter = queryObj[Query_EdgeFilter].(string)
		}

		if nil != queryObj[Query_TraversalCondition] {
			traversalCondition = queryObj[Query_TraversalCondition].(string)
		}

		if nil != queryObj[Query_EndCondition] {
			endCondition = queryObj[Query_EndCondition].(string)
		}
	}

	option := this.buildQueryOption(para)

	logger.Debug("======================== query ========================")
	logger.Debug(fmt.Sprintf("* queryString           = %s ", queryString))
	logger.Debug(fmt.Sprintf("* edgeFilter            = %s ", edgeFilter))
	logger.Debug(fmt.Sprintf("* traversalCondition    = %s ", traversalCondition))
	logger.Debug(fmt.Sprintf("* endCondition          = %s ", endCondition))
	logger.Debug(fmt.Sprintf("* Option.prefetchSize   = %d ", option.GetPreFetchSize()))
	logger.Debug(fmt.Sprintf("* Option.edgeLimit      = %d ", option.GetEdgeLimit()))
	logger.Debug(fmt.Sprintf("* Option.traversalDepth = %d ", option.GetTraversalDepth()))
	logger.Debug("-------------------------------------------------------")

	err := this.ensureConnection()
	if nil != err {
		return nil, err
	}

	resultSet, err := this._tgConnection.ExecuteQueryWithFilter(
		queryString,
		edgeFilter,
		traversalCondition,
		endCondition,
		option,
	)

	if nil != err {
		return nil, err
	}
	logger.Debug(fmt.Sprintf("Query Result = %v ", resultSet))

	return resultSet, err
}

//-====================-//
//    TGDB Callback
//-====================-//

func (this *TGDBService) OnException(ex types.TGError) {
	logger.Warn("***************************************************************")
	logger.Warn(fmt.Sprintf("[TGDBService::onException] Exception Happens : %v ", ex.GetErrorMsg()))
	if nil != this._tgConnection {
		logger.Warn("[TGDBService::onException] Connection not nil, will call disconnect !!! ")
		this._tgConnection.Disconnect()
	}
	logger.Warn("[TGDBService::onException] Set connection to nil !!! ")
	this._tgConnection = nil
	logger.Warn("***************************************************************")
}

//-====================-//
//       private
//-====================-//

func (this *TGDBService) buildNodeKey(nodeTypeStr string, nodeData map[string]interface{}) (string, error) {

	attribuesData := nodeData["attributes"].(map[string]interface{})
	logger.Debug(fmt.Sprintf("[TGDBService.buildNodeKey] Type = %s ,nodeData = %s, attribuesData = %s", nodeTypeStr, nodeData, attribuesData))

	nodeType, err := this.GetNodeType(nodeTypeStr)
	if nil != err {
		return "", err
	}

	pKeyAttributeDescriptors := nodeType.GetPKeyAttributeDescriptors()
	attrLength := len(pKeyAttributeDescriptors)
	key := make([]interface{}, attrLength)
	for i := 0; i < attrLength; i++ {
		name := pKeyAttributeDescriptors[i].GetName()
		if nil != attribuesData[name] {
			attributeData := attribuesData[name].(map[string]interface{})
			key[i] = attributeData["value"]
		} else {
			key[i] = nil
		}
	}
	logger.Debug(fmt.Sprintf("[TGDBService.buildNodeKey] Key = %s", key))

	return model.Hash(key), nil
}

func (this *TGDBService) fetchNode(fromType string, fromNode map[string]interface{}) (types.TGNode, error) {
	logger.Debug("[TGDBService::fetchNode] Entering ........ ")
	defer logger.Debug("[TGDBService::fetchNode] Exit ........ ")
	return this.search(fromType, fromNode, "", nil, nil, nil, nil)
}

func extractAttribute(entity map[string]interface{}) map[string]interface{} {
	attributes := entity["attributes"]
	if nil == attributes {
		attributes = make(map[string]interface{})
	}
	return attributes.(map[string]interface{})
}

func (this *TGDBService) search(
	fromType string,
	fromNode map[string]interface{},
	toType string,
	toNode map[string]interface{},
	edgeType interface{},
	edgeKeyAttributes []string,
	edgeKey map[string]interface{}) (types.TGNode, error) {

	fromAttributes := extractAttribute(fromNode)
	toAttributes := extractAttribute(toNode)

	logger.Debug(fmt.Sprintf("fromType = %s, fromAttributes = %s", fromType, fromAttributes))
	logger.Debug(fmt.Sprintf("toType = %s, toAttributes = %s", toType, toAttributes))
	logger.Debug(fmt.Sprintf("edgeType = %s, edgeKeyAttributes = %s, edgeKey = %s", edgeType, edgeKeyAttributes, edgeKey))

	nodeType, err := this.GetNodeType(fromType)
	if nil != err {
		return nil, err
	}

	var fromNodeCondition strings.Builder
	attrDescs := nodeType.GetPKeyAttributeDescriptors()
	for index := range attrDescs {
		name := attrDescs[index].GetName()
		fromNodeCondition.WriteString(" and ")
		fromNodeCondition.WriteString(name)
		fromNodeCondition.WriteString(" = '")
		if nil == fromAttributes[name] {
			return nil, errors.New("(search) Key attribute should not be nil!")
		}
		fromKeyAttributeData := fromAttributes[name].(map[string]interface{})
		if nil == fromKeyAttributeData["value"] {
			return nil, errors.New("(search) Key attribute should not be nil!")
		}
		fromNodeCondition.WriteString(EscapeIllegalChar(fromKeyAttributeData["value"].(string)))

		fromNodeCondition.WriteString("'")
	}

	var edgeCondition strings.Builder
	if nil != edgeKeyAttributes {
		for index := range edgeKeyAttributes {
			edgeCondition.WriteString(" and @edge.")
			edgeCondition.WriteString(edgeKeyAttributes[index])
			edgeCondition.WriteString(" = '")
			if nil == edgeKey[edgeKeyAttributes[index]] {
				return nil, errors.New("(search) Key attribute should not be nil!")
			}
			edgeKeyAttributeData := edgeKey[edgeKeyAttributes[index]].(map[string]interface{})
			if nil == edgeKeyAttributeData["value"] {
				return nil, errors.New("(search) Key attribute should not be nil!")
			}
			edgeCondition.WriteString(EscapeIllegalChar(edgeKeyAttributeData["value"].(string)))
			edgeCondition.WriteString("'")
		}
	}

	var toNodeCondition strings.Builder
	if "" != toType {
		toNodeType, err := this.GetNodeType(toType)
		if nil != err {
			return nil, err
		}
		attrDescs := toNodeType.GetPKeyAttributeDescriptors()
		for index := range attrDescs {
			name := attrDescs[index].GetName()
			toNodeCondition.WriteString(" and @tonode.")
			toNodeCondition.WriteString(name)
			toNodeCondition.WriteString(" = '")
			if nil == toAttributes[name] {
				return nil, errors.New("(search) Key attribute should not be nil!")
			}
			toKeyAttributeData := toAttributes[name].(map[string]interface{})
			if nil == toKeyAttributeData["value"] {
				return nil, errors.New("(search) Key attribute should not be nil!")
			}
			toNodeCondition.WriteString(EscapeIllegalChar(toKeyAttributeData["value"].(string)))
			toNodeCondition.WriteString("'")
		}
	}

	para := make(map[string]interface{})
	query := make(map[string]interface{})
	query["queryString"] = fmt.Sprintf("@nodetype = '%s' %s;", fromType, fromNodeCondition.String())
	if nil != edgeType {
		query["traversalCondition"] = fmt.Sprintf("@edgetype = '%s' %s and @tonodetype = '%s' %s and @isfromedge = 1 and @degree = 1;", edgeType, edgeCondition.String(), toType, toNodeCondition.String())
	} else if nil != edgeKey {
		/* need to search by edge index */
	}

	logger.Info("[TGDBService.search] ", query)

	para[Query] = query
	para[Query_OPT_PrefetchSize] = 500
	para[Query_OPT_TraversalDepth] = 1
	para[Query_OPT_EdgeLimit] = 500

	var startingNode types.TGNode
	resultSet, err := this.TGQLQuery(para)
	if nil == err {
		if nil != resultSet {
			result := resultSet.Next()
			if nil == result {
				return nil, nil
			}
			startingNode = result.(types.TGNode)
		}
	} else {
		return nil, err
	}

	return startingNode, err
}

func (this *TGDBService) populateAttributes(
	entity types.TGEntity,
	attributesData map[string]interface{}) types.TGEntity {
	logger.Info(fmt.Sprintf("[TGDBService.populateAttributes] entity = %s, attributesData = %s ", entity.GetEntityType().GetName(), attributesData))
	entityType := entity.GetEntityType()
	if nil != entityType {
		attrDescs := entityType.GetAttributeDescriptors()
		for _, attrDesc := range attrDescs {
			name := attrDesc.GetName()
			if nil != attributesData[name] {
				attributeData := attributesData[name].(map[string]interface{})
				if nil != attributeData["value"] {
					if "Integer" == attributeData["type"] {
						entity.SetOrCreateAttribute(name, int(attributeData["value"].(int32)))
					} else if "Date" == attributeData["type"] {
						entity.SetOrCreateAttribute(name, attributeData["value"].(time.Time).Unix())
					} else {
						entity.SetOrCreateAttribute(name, attributeData["value"])
					}
				}
			}
		}
	} else {
		for name := range attributesData {
			attributeData := attributesData[name].(map[string]interface{})
			if nil != attributeData["value"] {
				if "Integer" == attributeData["type"] {
					entity.SetOrCreateAttribute(name, int(attributeData["value"].(int32)))
				} else if "Date" == attributeData["type"] {
					entity.SetOrCreateAttribute(name, attributeData["value"].(time.Time).Unix())
				} else {
					entity.SetOrCreateAttribute(name, attributeData["value"])
				}
			}
		}
	}

	return entity
}

func (this *TGDBService) createNode(nodeTypeStr string) (types.TGNode, types.TGError) {

	nodeType, err := this._gmd.GetNodeType(nodeTypeStr)
	if nil == nodeType {
		return nil, err
	}

	node, _ := this._gof.CreateNodeInGraph(nodeType)

	return node, nil
}

func (this *TGDBService) findNode(
	readyEntity EntityKeeper,
	nodeType string,
	nodeKey string,
	nodeData map[string]interface{}) (*TGNodeWrapper, error) {

	nodeWrapper := readyEntity.GetNode(nodeType, nodeKey)
	var node types.TGNode
	var err error
	if nil == nodeWrapper {
		// Check node existence in tgdb
		node, err = this.fetchNode(nodeType, nodeData)
		if nil != err {
			return nil, err
		}
		if nil != node {
			nodeWrapper = NewTGNodeWrapper(node, false)
			logger.Debug(fmt.Sprintf("/* Node found in DB (%s) */", nodeWrapper))
		} else {
			node, err = this.createNode(nodeType)
			if nil != err {
				logger.Error("Error to build Node !!!")
				return nil, err
			}
			if nil == node {
				logger.Error("Unable to build Node !!!")
				return nil, nil
			}
			nodeWrapper = NewTGNodeWrapper(node, true)
			logger.Debug(fmt.Sprintf("/* Node not found in DB so create Node (%s) */", nodeWrapper))
		}
	} else {
		node = *nodeWrapper.node
		logger.Debug(fmt.Sprintf("/* Node exist locallly - might be created from previous edge (%s) */", nodeWrapper))
	}

	if nil != nodeData["attributes"] {
		this.populateAttributes(node, nodeData["attributes"].(map[string]interface{}))
	}

	return nodeWrapper, nil
}

func (this *TGDBService) fetchMetadata() types.TGError {
	gmd, err := this._tgConnection.GetGraphMetadata(true)

	if nil != err {
		return err
	}

	this._gmd = gmd

	for value := range this._keyMap {
		delete(this._keyMap, value)
	}

	nodeTypes, err2 := this._gmd.GetNodeTypes()

	if nil != err2 {
		return err2
	}

	for _, nodeType := range nodeTypes {
		keys := make([]string, 0)
		attrDescs := nodeType.GetPKeyAttributeDescriptors()
		for index := range attrDescs {
			keys = append(keys, attrDescs[index].GetName())
		}
		this._keyMap[nodeType.GetName()] = keys

		attrs := make(map[string]types.TGAttributeDescriptor)
		attrDescs = nodeType.GetAttributeDescriptors()
		for index := range attrDescs {
			attrs[attrDescs[index].GetName()] = attrDescs[index]
		}
		this._nodeAttrMap[nodeType.GetName()] = attrs
	}

	edgeTypes, err3 := this._gmd.GetEdgeTypes()

	if nil != err3 {
		return err3
	}

	for _, edgeType := range edgeTypes {
		attrs := make(map[string]types.TGAttributeDescriptor)
		attrDescs := edgeType.GetAttributeDescriptors()
		for index := range attrDescs {
			attrs[attrDescs[index].GetName()] = attrDescs[index]
		}
		this._edgeAttrMap[edgeType.GetName()] = attrs
	}

	return nil
}

func (this *TGDBService) buildQueryOption(para map[string]interface{}) types.TGQueryOption {
	option := query.NewQueryOption()

	logger.Debug(fmt.Sprintf("Parameter : %v ", para))

	if nil == para[Query_OPT_PrefetchSize] || 0 == para[Query_OPT_PrefetchSize].(int) {
		option.SetPreFetchSize(500)
	} else {
		option.SetPreFetchSize(para[Query_OPT_PrefetchSize].(int))
	}

	if nil == para[Query_OPT_TraversalDepth] || 0 == para[Query_OPT_TraversalDepth].(int) {
		option.SetTraversalDepth(5)
	} else {
		option.SetTraversalDepth(para[Query_OPT_TraversalDepth].(int))
	}

	if nil == para[Query_OPT_EdgeLimit] || 0 == para[Query_OPT_EdgeLimit].(int) {
		option.SetEdgeLimit(100)
	} else {
		option.SetEdgeLimit(para[Query_OPT_EdgeLimit].(int))
	}

	return option
}
