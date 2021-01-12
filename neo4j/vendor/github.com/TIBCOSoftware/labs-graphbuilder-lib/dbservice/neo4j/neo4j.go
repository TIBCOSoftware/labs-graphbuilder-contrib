/*
 * Copyright © 2019. TIBCO Software Inc.
 * This file is subject to the license terms contained
 * in the license file that is distributed with this file.
 */
package neo4j

import (
	"bytes"
	"fmt"

	"sync"

	"github.com/TIBCOSoftware/labs-graphbuilder-lib/model"
	"github.com/TIBCOSoftware/labs-graphbuilder-lib/util"
	bolt "github.com/johnnadratowski/golang-neo4j-bolt-driver"
)

const (
	GRAPH_MODEL_ID   = "graph_builder_model_id"
	DATE_TIME_SAMPLE = "2006-01-02"
)

type Neo4jService struct {
	_url             string
	_user            string
	_password        string
	_typeName        string
	_addPrefixToAttr bool
	_doReplaceChar   bool
	_targetRegex     string
	_replacement     string

	_neoConnection *bolt.Conn

	_pendingNodes map[string]*DNode
	_pendingEdges map[string]*DEdge

	_mux sync.Mutex
}

func (this *Neo4jService) ensureConnection() error {

	if nil == this._neoConnection {
		this._mux.Lock()
		defer this._mux.Unlock()
		if nil == this._neoConnection {
			if nil != this._neoConnection {
				(*this._neoConnection).Close()
			}

			fmt.Println("[Neo4jService::ensureConnection] Will try to connect ..........")
			fmt.Println("[Neo4jService::ensureConnection] url = " + this._url)
			fmt.Println("[Neo4jService::ensureConnection] user = " + this._user)
			fmt.Println("[Neo4jService::ensureConnection] password = " + this._password)
			driver := bolt.NewDriver()
			conn, err := driver.OpenNeo(fmt.Sprintf("bolt://%s:%s@%s", this._user, this._password, this._url))
			this._neoConnection = &conn
			if err != nil {
				panic(err)
			}
			if nil != err {
				fmt.Println("[Neo4jService::ensureConnection] Unable to create connection !!! Will not connect ......")
				this._neoConnection = nil
				return err
			}
		}
	}

	return nil
}

func (this *Neo4jService) Destroy() {
	(*this._neoConnection).Close()
}

//-====================-//
//    Delete Graph
//-====================-//

func (this *Neo4jService) DeleteGraph(filter int, graphToo map[string]interface{}) error {
	return nil
}

//-====================-//
//    Upsert Graph
//-====================-//

func (this *Neo4jService) UpsertGraph(graph model.Graph, graphToo map[string]interface{}) error {

	this.ensureConnection()

	this._mux.Lock()
	defer this._mux.Unlock()

	fmt.Println("\n\n\n\n\n\n*************************************************")

	log.Info("(UpsertGraph) begin - graph = ", graph)

	log.Info("(UpsertGraph) _pendingNodes = ", this._pendingNodes)
	log.Info("(UpsertGraph) graph.GetNodes() = ", graph.GetNodes())

	this._pendingNodes = make(map[string]*DNode)
	for id, node := range graph.GetNodes() {
		log.Info("node id = ", id, ", node = ", node)
		this._pendingNodes[id.ToString()] = NewDNode(node, this._typeName, this._addPrefixToAttr)
	}

	log.Info("(UpsertGraph) _pendingEdges = ", this._pendingEdges)
	log.Info("(UpsertGraph) graph.GetEdges() = ", graph.GetEdges())
	this._pendingEdges = make(map[string]*DEdge)
	for id, edge := range graph.GetEdges() {
		from := this._pendingNodes[edge.GetFromId().ToString()]
		to := this._pendingNodes[edge.GetToId().ToString()]
		this._pendingEdges[id.ToString()] = NewDEdge(edge, this._typeName, this._addPrefixToAttr, from, to)
	}

	err := this.Commit(graph)

	log.Info("(UpsertGraph) Done ! ")

	return err
}

func (this *Neo4jService) Commit(graph model.Graph) error {

	log.Info("(Commit) begin - ", graph)

	for _, edge := range this._pendingEdges {
		this.checkEdge(graph, edge)

		from := edge.GetFrom()
		fromInPending := this._pendingNodes[from.GetId()]
		if nil != fromInPending {
			if from.Exists() {
				fromInPending.SetNeoId(from.GetNeoId())
			}
		}

		to := edge.GetTo()
		toInPending := this._pendingNodes[to.GetId()]
		if nil != toInPending {
			if to.Exists() {
				toInPending.SetNeoId(to.GetNeoId())
			}
		}
	}

	for key, node := range this._pendingNodes {
		if !node.Exists() {
			this.checkNode(graph, node)
		}

		query := node.ToCypher(graph, DATE_TIME_SAMPLE)
		log.Info("(Commit) query for upsert node : ", query)
		neoId, err := this.execute(query)
		if nil != err {
			log.Info(err)
		}
		node.SetNeoId(neoId)
		delete(this._pendingNodes, key)
	}

	for key, edge := range this._pendingEdges {
		query := edge.ToCypher(graph, DATE_TIME_SAMPLE)
		log.Info("(Commit) query for upsert edge : ", query)
		neoId, err := this.execute(query)
		if nil != err {
			log.Info(err)
		}
		edge.SetNeoId(neoId)
		delete(this._pendingEdges, key)
	}

	return nil
}

func (this *Neo4jService) execute(query string) (interface{}, error) {
	log.Info("(execute) query = ", query)
	stmt, err := (*this._neoConnection).PrepareNeo(query)
	defer stmt.Close()
	if err != nil {
		panic(err)
	}

	result, err := stmt.QueryNeo(nil)
	if err != nil {
		panic(err)
	}

	dataArray, _, _ := result.All()
	log.Info("check execute result - ", dataArray)
	if 1 <= len(dataArray) {
		return dataArray[0][0], nil
	}

	return nil, fmt.Errorf("Execution fail, no entity id return!")
}

func (this *Neo4jService) checkNode(graph model.Graph, node *DNode) error {
	log.Info("(checkNode) begin - ", node)

	//MATCH (node:NodeType { pkKey: pkValue }) RETURN id(node)
	var query bytes.Buffer
	query.WriteString("MATCH (node:")
	query.WriteString(util.ReplaceCharacter(node.GetType(), this._targetRegex, this._replacement, true))
	query.WriteString("{")
	for index, keyAttributeName := range graph.GetEntityKeyNamesForNode(node.GetType()) {
		if 0 != index {
			query.WriteString(", ")
		}
		query.WriteString(GetCanonicalAttributeName(keyAttributeName, this._targetRegex, this._replacement, true))
		query.WriteString(": ")
		keyAttributeValue, _ := node.GetAttributeValueString(keyAttributeName, DATE_TIME_SAMPLE)
		query.WriteString(keyAttributeValue)
	}
	query.WriteString("}) RETURN id(node)")
	log.Info("check node query - ", query.String())

	stmt, err := (*this._neoConnection).PrepareNeo(query.String())
	defer stmt.Close()
	if err != nil {
		panic(err)
	}

	result, err := stmt.QueryNeo(nil)
	if err != nil {
		panic(err)
	}

	dataArray, _, _ := result.All()
	log.Info("check node result - ", dataArray)
	if 1 <= len(dataArray) {
		node.SetNeoId(dataArray[0][0])
	}

	return nil
}

func (this *Neo4jService) checkEdge(graph model.Graph, edge *DEdge) error {
	log.Info("(checkEdge) begin - ", edge)

	//MATCH (from:houseMemberType { memberName: 'Carlo Bonaparte' })-[edge:relation]->(to:houseMemberType { memberName: 'Letizia Ramolino' }) RETURN id(from), id(edge), id(to)
	var query bytes.Buffer
	query.WriteString("MATCH (from:")
	log.Info("(edge.GetFrom().GetType()) - ", edge.GetFrom().GetType())
	query.WriteString(edge.GetFrom().GetType())
	query.WriteString("{")
	for index, keyAttributeName := range graph.GetEntityKeyNamesForNode(edge.GetFrom().GetType()) {
		log.Info("keyAttributeName - ", keyAttributeName)
		if 0 != index {
			query.WriteString(", ")
		}
		query.WriteString(GetCanonicalAttributeName(keyAttributeName, this._targetRegex, this._replacement, true))
		query.WriteString(": ")
		keyAttributeValue, _ := edge.GetFrom().GetAttributeValueString(keyAttributeName, DATE_TIME_SAMPLE)
		query.WriteString(keyAttributeValue)
	}
	query.WriteString("})-[edge:")
	query.WriteString(util.ReplaceCharacter(edge.GetType(), this._targetRegex, this._replacement, true))
	log.Info("edgeType - ", util.ReplaceCharacter(edge.GetType(), this._targetRegex, this._replacement, true))

	for index, keyAttributeName := range graph.GetEntityKeyNamesForEdge(edge.GetType()) {
		if 0 == index {
			query.WriteString(", ")
		}
		query.WriteString(keyAttributeName)
		query.WriteString(": ")
		keyAttributeValue, _ := edge.GetAttributeValueString(keyAttributeName, DATE_TIME_SAMPLE)
		query.WriteString(keyAttributeValue)
	}
	query.WriteString("]->(to:")
	query.WriteString(edge.GetTo().GetType())
	query.WriteString("{")
	for index, keyAttributeName := range graph.GetEntityKeyNamesForNode(edge.GetTo().GetType()) {
		if 0 != index {
			query.WriteString(", ")
		}
		query.WriteString(GetCanonicalAttributeName(keyAttributeName, this._targetRegex, this._replacement, true))
		query.WriteString(": ")
		keyAttributeValue, _ := edge.GetTo().GetAttributeValueString(keyAttributeName, DATE_TIME_SAMPLE)
		query.WriteString(keyAttributeValue)
	}
	query.WriteString("}) RETURN id(from), id(edge), id(to)")
	log.Info("check edge query - ", query.String())

	stmt, err := (*this._neoConnection).PrepareNeo(query.String())
	defer stmt.Close()
	if err != nil {
		panic(err)
	}

	result, err := stmt.QueryNeo(nil)
	if err != nil {
		panic(err)
	}

	resultArray, _, _ := result.All()
	log.Info("check edge result - ", resultArray)
	if 1 <= len(resultArray) {
		edge.GetFrom().SetNeoId(resultArray[0][0])
		edge.SetNeoId(resultArray[0][1])
		edge.GetTo().SetNeoId(resultArray[0][2])
	}

	return nil
}
