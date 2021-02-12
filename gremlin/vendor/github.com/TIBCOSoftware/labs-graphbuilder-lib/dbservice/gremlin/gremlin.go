/*
 * Copyright © 2020. TIBCO Software Inc.
 * This file is subject to the license terms contained
 * in the license file that is distributed with this file.
 */
package gremlin

import (
	"bytes"
	"encoding/json"
	"fmt"

	"sync"

	gremlin "github.com/northwesternmutual/grammes"
	gmodel "github.com/northwesternmutual/grammes/model"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"

	"github.com/TIBCOSoftware/labs-graphbuilder-lib/model"
	"github.com/TIBCOSoftware/labs-graphbuilder-lib/util"
)

const (
	GRAPH_MODEL_ID   = "graph_builder_model_id"
	DATE_TIME_SAMPLE = "2006-01-02"
)

type GremlinService struct {
	_url             string
	_user            string
	_password        string
	_typeName        string
	_addPrefixToAttr bool
	_doReplaceChar   bool
	_targetRegex     string
	_replacement     string

	_client *gremlin.Client

	_pendingNodes map[string]*DNode
	_pendingEdges map[string]*DEdge

	_mux sync.Mutex
}

func (this *GremlinService) ensureConnection() error {

	if nil == this._client {
		this._mux.Lock()
		defer this._mux.Unlock()
		if nil == this._client {

			log.Info("(ensureConnection) Will try to connect ..........")
			log.Info("(ensureConnection) url = ", this._url)
			log.Info("(ensureConnection) user = ", this._user)
			log.Info("(ensureConnection) password = ", this._password)
			client, err := gremlin.DialWithWebSocket(this._url)
			if err != nil {
				log.Error("(ensureConnection) Error while creating client: ", err.Error())
			}
			this._client = client
			if err != nil {
				panic(err)
			}
			if nil != err {
				log.Error("(ensureConnection) Unable to create connection !!! Will not connect ......")
				this._client = nil
				return err
			}
		}
	}

	return nil
}

func (this *GremlinService) Destroy() {
}

//-====================-//
//       Query
//-====================-//

func (this *GremlinService) Query(query string) (*zapcore.Field, error) {
	log.Info("(execute) query = ", query)

	responses, err := this._client.ExecuteStringQuery(query)
	if err != nil {
		log.Error("(execute) Error querying server", zap.Error(err))
		return nil, err
	}

	retult := zap.ByteStrings("result", responses)

	return &retult, nil
}

//-====================-//
//    Delete Graph
//-====================-//

func (this *GremlinService) DeleteGraph(filter int, graphToo map[string]interface{}) error {
	return nil
}

//-====================-//
//    Upsert Graph
//-====================-//

func (this *GremlinService) UpsertGraph(graph model.Graph, graphToo map[string]interface{}) error {

	this.ensureConnection()

	this._pendingNodes = make(map[string]*DNode)
	this._pendingEdges = make(map[string]*DEdge)

	this._mux.Lock()
	defer this._mux.Unlock()

	log.Debug("(UpsertGraph) begin - graph = ", graph)

	log.Debug("(UpsertGraph) _pendingNodes = ", this._pendingNodes)
	log.Debug("(UpsertGraph) graph.GetNodes() = ", graph.GetNodes())

	this._pendingNodes = make(map[string]*DNode)
	for id, node := range graph.GetNodes() {
		log.Info("node id = ", id, ", node = ", node)
		this._pendingNodes[id.ToString()] = NewDNode(node, this._typeName, this._addPrefixToAttr)
	}

	log.Debug("(UpsertGraph) _pendingEdges = ", this._pendingEdges)
	log.Debug("(UpsertGraph) graph.GetEdges() = ", graph.GetEdges())
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

func (this *GremlinService) Commit(graph model.Graph) error {

	log.Info("(Commit) begin: edges = ", this._pendingEdges, ", nodes = ", this._pendingNodes)

	for _, edge := range this._pendingEdges {
		this.checkEdge(graph, edge)
		log.Info(
			"(Commit) after check edge : edgeId = ", edge.GetGremlinId(),
			", fromId = ", edge.GetFrom().GetGremlinId(),
			", toId = ", edge.GetTo().GetGremlinId())

		from := edge.GetFrom()
		fromInPending := this._pendingNodes[from.GetId()]
		if nil != fromInPending {
			if from.Exists() {
				fromInPending.SetGremlinId(from.GetGremlinId())
			}
		}

		to := edge.GetTo()
		toInPending := this._pendingNodes[to.GetId()]
		if nil != toInPending {
			if to.Exists() {
				toInPending.SetGremlinId(to.GetGremlinId())
			}
		}
	}

	for key, node := range this._pendingNodes {
		if !node.Exists() {
			this.checkNode(graph, node)
		}

		if !node.Exists() {
			err := this.execute(graph, node)
			if nil != err {
				log.Info(err)
			}
		}
		delete(this._pendingNodes, key)
	}

	for key, edge := range this._pendingEdges {
		err := this.execute(graph, edge)
		if nil != err {
			log.Info(err)
		}
		delete(this._pendingEdges, key)
	}

	return nil
}

func (this *GremlinService) execute(graph model.Graph, entity GremlinEntity) error {
	query := entity.ToGremlin(graph, DATE_TIME_SAMPLE)
	log.Info("(execute) query = ", query)

	responses, err := this._client.ExecuteStringQuery(query)
	if err != nil {
		log.Error("(execute) Error querying server", zap.Error(err))
	}

	switch entity.Kind() {
	case Node:
		//{"result": "{\"@type\":\"g:List\",\"@value\":[{\"@type\":\"g:Vertex\",\"@value\":{\"id\":{\"@type\":\"g:Int64\",\"@value\":83034328},\"label\":\"houseMemberType\",\"properties\":{\"memberName\":[{\"@type\":\"g:VertexProperty\",\"@value\":{\"id\":{\"@type\":\"janusgraph:RelationIdentifier\",\"@value\":{\"relationId\":\"1d5na3-1dfpp4-1l1\"}},\"label\":\"memberName\",\"value\":\"Marie Louise of Austria\"}}]}}}]}"}
		var result []gmodel.Vertex
		result, err = buildNode(responses[0])
		if nil == err {
			entity.SetGremlinId(result[0].ID())
			log.Info("(execute) node result : ", result[0].ID())
		}
	case Edge:
		//{"result": "{\"@type\":\"g:List\",\"@value\":[{\"@type\":\"g:Edge\",\"@value\":{\"id\":{\"@type\":\"janusgraph:RelationIdentifier\",\"@value\":{\"relationId\":\"ovj0q-ozryo-2dx-1dfpp4\"}},\"inV\":{\"@type\":\"g:Int64\",\"@value\":83034328},\"inVLabel\":\"houseMemberType\",\"label\":\"relation\",\"outV\":{\"@type\":\"g:Int64\",\"@value\":41979984},\"outVLabel\":\"houseMemberType\",\"properties\":{\"relType\":{\"@type\":\"g:Property\",\"@value\":{\"key\":\"relType\",\"value\":\"spouse\"}}}}}]}"}
		var result []gmodel.Edge
		result, err = buildEdge(responses[0])
		if nil == err {
			if nil == entity.GetGremlinId() {
				entity.SetGremlinId(result[0].ID())
				log.Info("(execute) edge created : ", result[0].ID())
			} else {
				log.Info("(execute) edge updated : ", entity.GetGremlinId())
			}
		}
	default:
		err = fmt.Errorf("Execution fail, no entity id return!")
	}

	return err
}

func (this *GremlinService) checkNode(graph model.Graph, node *DNode) error {
	log.Info("(checkNode) begin - ", node)
	var query bytes.Buffer

	query.WriteString("g.V()")
	this.nodeHas(&query, graph, node)

	log.Info("(execute) check node query - ", query.String())

	vertices, err := this._client.VerticesByString(query.String())
	if err != nil {
		log.Error("Couldn't gather vertices", zap.Error(err))
	}

	log.Debug("(execute) check node result - ", vertices)

	switch len(vertices) {
	case 1:
		node.SetGremlinId(vertices[0].ID())
		log.Info("(execute) check node id - ", vertices[0].ID())
	default:
		log.Info("(execute) check node id - node not found!!")
	}

	return nil
}

func (this *GremlinService) checkEdge(graph model.Graph, edge *DEdge) error {
	log.Info("(checkEdge) begin - ", edge)
	log.Info("(edge.GetFrom().GetType()) - ", edge.GetFrom().GetType())

	var query bytes.Buffer
	query.WriteString("g.V()")
	this.nodeHas(&query, graph, edge.GetFrom())
	query.WriteString(".outE('")
	query.WriteString(util.ReplaceCharacter(edge.GetType(), this._targetRegex, this._replacement, true))
	log.Info("(execute) edgeType - ", util.ReplaceCharacter(edge.GetType(), this._targetRegex, this._replacement, true))

	for index, keyAttributeName := range graph.GetEntityKeyNamesForEdge(edge.GetType()) {
		if 0 == index {
			query.WriteString("').has('")
		}
		query.WriteString(keyAttributeName)
		query.WriteString(", '")
		keyAttributeValue, _ := edge.GetAttributeValueString(keyAttributeName, DATE_TIME_SAMPLE)
		query.WriteString(keyAttributeValue)
	}
	query.WriteString("').inV()")
	this.nodeHas(&query, graph, edge.GetTo())
	query.WriteString(".inE('")
	query.WriteString(util.ReplaceCharacter(edge.GetType(), this._targetRegex, this._replacement, true))
	for index, keyAttributeName := range graph.GetEntityKeyNamesForEdge(edge.GetType()) {
		if 0 == index {
			query.WriteString("').has('")
		}
		query.WriteString(keyAttributeName)
		query.WriteString(", '")
		keyAttributeValue, _ := edge.GetAttributeValueString(keyAttributeName, DATE_TIME_SAMPLE)
		query.WriteString(keyAttributeValue)
	}
	query.WriteString("')")

	log.Info("(execute) check edge query - ", query.String())

	responses, err := this._client.ExecuteStringQuery(query.String())
	if err != nil {
		log.Error("Error querying server", zap.Error(err))
	}

	switch len(responses) {
	case 1:
		if 4 == len(responses[0]) {
			log.Info("(execute) checkEdge result :", string(responses[0]))
		} else {
			gedge, err := buildEdge(responses[0])
			if nil == err {
				edge.SetGremlinId(gedge[0].ID())
				edge.from.SetGremlinId(gedge[0].InVertexID())
				edge.to.SetGremlinId(gedge[0].OutVertexID())
				log.Info("checkEdge result :", gedge[0].ID(), " - ", gedge[0].InVertexID(), " - ", gedge[0].OutVertexID())
			}
		}
	default:
		log.Info("(execute) checkEdge result : unknown !!!")
	}

	return nil
}

func (this *GremlinService) nodeHas(query *bytes.Buffer, graph model.Graph, node *DNode) {
	for index, keyAttributeName := range graph.GetEntityKeyNamesForNode(node.GetType()) {
		query.WriteString(".has('")
		if 0 == index {
			query.WriteString(node.GetType())
			query.WriteString("', '")
		}
		query.WriteString(GetCanonicalAttributeName(keyAttributeName, this._targetRegex, this._replacement, true))
		query.WriteString("', ")
		keyAttributeValue, _ := node.GetAttributeValueString(keyAttributeName, DATE_TIME_SAMPLE)
		query.WriteString(keyAttributeValue)
		query.WriteString(")")
	}
}

func buildResult(entity []byte) map[string]interface{} {
	res := make(map[string]interface{})
	json.Unmarshal(entity, &res)
	return res
}

func buildNode(node []byte) ([]gmodel.Vertex, error) {
	var vertices gmodel.VertexList
	// Unmarshal the response into the structs.
	err := json.Unmarshal(node, &vertices)
	if err != nil {
		return nil, err
	}
	log.Info("(execute) buildNode : ", vertices)
	return vertices.Vertices, err
}

func buildEdge(edge []byte) ([]gmodel.Edge, error) {
	var edges gmodel.EdgeList
	// Unmarshal the response into the structs.
	err := json.Unmarshal(edge, &edges)
	if err != nil {
		return nil, err
	}
	log.Info("(execute) buildEdge : ", edges)
	return edges.Edges, err
}
