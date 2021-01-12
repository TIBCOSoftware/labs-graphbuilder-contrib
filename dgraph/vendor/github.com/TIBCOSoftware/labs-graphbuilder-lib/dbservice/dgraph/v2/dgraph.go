/*
 * Copyright © 2020. TIBCO Software Inc.
 * This file is subject to the license terms contained
 * in the license file that is distributed with this file.
 */
package v2

import (
	"bytes"
	"context"

	"encoding/json"
	"fmt"

	"sync"

	"github.com/TIBCOSoftware/labs-graphbuilder-lib/dbservice"
	"github.com/TIBCOSoftware/labs-graphbuilder-lib/dbservice/dgraph"
	"github.com/TIBCOSoftware/labs-graphbuilder-lib/dbservice/dgraph/cache"
	"github.com/TIBCOSoftware/labs-graphbuilder-lib/dbservice/dgraph/rdf"

	"github.com/TIBCOSoftware/labs-graphbuilder-lib/model"
	"github.com/TIBCOSoftware/labs-graphbuilder-lib/util"

	"github.com/dgraph-io/dgo/v2"
	"github.com/dgraph-io/dgo/v2/protos/api"
	"github.com/dgraph-io/dgraph/x"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"

	"github.com/TIBCOSoftware/flogo-lib/logger"
)

var log = logger.GetLogger("dgraph-service")

func NewDgraphService(properties map[string]interface{}) (dbservice.UpsertService, error) {
	dgraphService := &DgraphService{}
	if nil != properties["url"] {
		dgraphService._url = properties["url"].(string)
	}
	if nil != properties["user"] {
		dgraphService._user = properties["user"].(string)
	}
	if nil != properties["password"] {
		dgraphService._password = properties["password"].(string)
	}
	if nil != properties["tlsEnabled"] {
		dgraphService._tlsEnabled = properties["tlsEnabled"].(bool)
	}
	if nil != properties["tls"] {
		var tlsObject interface{}
		err := json.Unmarshal([]byte(properties["tls"].(string)), &tlsObject)
		if nil != err {
			fmt.Println(err)
		}
		if nil != tlsObject {
			dgraphService._tlsUserCfg = tlsObject.(map[string]interface{})
		}
	}

	if nil != properties["explicitType"] {
		dgraphService._explicitType = properties["explicitType"].(bool)
	} else {
		dgraphService._explicitType = false
	}

	if nil != properties["cacheSize"] {
		cacheSize := properties["cacheSize"].(int)
		dgraphService._cachedNodes = cache.NewCache(cacheSize)
		dgraphService._cachedEdges = cache.NewCache(cacheSize)
	} else {
		dgraphService._cachedNodes = cache.NewCache(-1)
		dgraphService._cachedEdges = cache.NewCache(-1)
	}

	if nil != properties["readableExternalId"] {
		dgraphService._readableExternalId = properties["readableExternalId"].(bool)
	} else {
		dgraphService._readableExternalId = true
	}

	if nil != properties["typeName"] {
		dgraphService._typeName = util.CastString(properties["typeName"])
	} else {
		dgraphService._typeName = ""
	}

	if nil != properties["addPrefixToAttr"] {
		dgraphService._addPrefixToAttr = properties["addPrefixToAttr"].(bool)
	} else {
		dgraphService._addPrefixToAttr = false
	}
	dgraphService._targetRegex = "[^A-Za-z0-9]"
	dgraphService._replacement = "_"

	var err error
	if "no" != properties["schemaGen"] && nil != properties["graphModel"] {
		err = dgraphService.BuildSchema(properties["schema"], properties["graphModel"].(map[string]interface{}))
	} else {
		err = dgraphService.ensureConnection()
	}

	if nil != err {
		dgraphService = nil
	}

	return dgraphService, err
}

type DgraphService struct {
	_url                string
	_user               string
	_password           string
	_explicitType       bool
	_typeName           string
	_addPrefixToAttr    bool
	_doReplaceChar      bool
	_targetRegex        string
	_replacement        string
	_readableExternalId bool

	_tlsEnabled   bool
	_tlsUserCfg   map[string]interface{}
	_dgConnection *grpc.ClientConn
	_dgraphClient *dgo.Dgraph

	_cachedNodes *cache.Cache
	_cachedEdges *cache.Cache

	_mux sync.Mutex
}

func (this *DgraphService) ensureConnection() error {
	var connErr error

	if nil == this._dgConnection {
		this._mux.Lock()
		defer this._mux.Unlock()
		if nil == this._dgConnection {
			if nil != this._dgConnection {
				this._dgConnection.Close()
			}

			fmt.Println("[DgraphService::ensureConnection] Will try to connect ..........")
			fmt.Println("[DgraphService::ensureConnection] url = " + this._url)
			fmt.Println("[DgraphService::ensureConnection] user = " + this._user)
			fmt.Println("[DgraphService::ensureConnection] password = " + this._password)
			fmt.Println("[DgraphService::ensureConnection] tlsUserCfg = ", this._tlsUserCfg)

			if !this._tlsEnabled {
				this._dgConnection, connErr = grpc.Dial(this._url, grpc.WithInsecure())
			} else {
				helperConfig := &x.TLSHelperConfig{}
				if nil != this._tlsUserCfg["tlsCertDir"] {
					helperConfig.CertDir = this._tlsUserCfg["tlsCertDir"].(string)
				}
				if nil != this._tlsUserCfg["tlsCertRequired"] {
					helperConfig.CertRequired = this._tlsUserCfg["tlsCertRequired"].(bool)
				}
				if nil != this._tlsUserCfg["tlsCert"] {
					helperConfig.Cert = this._tlsUserCfg["tlsCert"].(string)
				}
				if nil != this._tlsUserCfg["tlsKey"] {
					helperConfig.Key = this._tlsUserCfg["tlsKey"].(string)
				}

				if nil != this._tlsUserCfg["tlsServerName"] {
					helperConfig.ServerName = this._tlsUserCfg["tlsServerName"].(string)
				}
				if nil != this._tlsUserCfg["tlsRootCACert"] {
					helperConfig.RootCACert = this._tlsUserCfg["tlsRootCACert"].(string)
				}
				if nil != this._tlsUserCfg["tlsClientAuth"] {
					helperConfig.ClientAuth = this._tlsUserCfg["tlsClientAuth"].(string)
				}
				if nil != this._tlsUserCfg["tlsUseSystemCACerts"] {
					helperConfig.UseSystemCACerts = this._tlsUserCfg["tlsUseSystemCACerts"].(bool)
				}

				tlsCfg, err := x.GenerateClientTLSConfig(helperConfig)
				if nil != err {
					log.Error("[DgraphService::ensureConnection] Unable to configure TLS connection !!! Will not connect ......")
					this._dgConnection = nil
					return connErr
				}

				this._dgConnection, connErr = grpc.Dial(this._url, grpc.WithTransportCredentials(credentials.NewTLS(tlsCfg)))
			}

			if nil != connErr {
				log.Error("[DgraphService::ensureConnection] Unable to create connection !!! Will not connect ......")
				this._dgConnection = nil
				return connErr
			}

			this._dgraphClient = dgo.NewDgraphClient(api.NewDgraphClient(this._dgConnection))
		}
	}

	return nil
}

func (this *DgraphService) Destroy() {
	this._dgConnection.Close()
}

func (this *DgraphService) BuildSchema(userSchema interface{}, model map[string]interface{}) error {
	schema := buildSchema(
		this._explicitType,
		this._typeName,
		this._targetRegex,
		this._replacement,
		this._addPrefixToAttr,
		userSchema,
		model,
	)
	fmt.Println("***************** schema query ********************")
	fmt.Println(schema)
	fmt.Println("***************************************************")

	this.ensureConnection()

	err := this._dgraphClient.Alter(
		context.Background(),
		&api.Operation{Schema: schema},
	)

	return err
}

func (this *DgraphService) DropGraph() {
	this._dgraphClient.Alter(context.Background(), &api.Operation{DropAll: true})
}

func (this *DgraphService) DeleteGraph(filter int, graphToo map[string]interface{}) error {
	return nil
}

//-====================-//
//    Query Graph
//-====================-//

func (this *DgraphService) Query(query string) (interface{}, error) {
	this.ensureConnection()

	res, err := this._dgraphClient.NewTxn().QueryWithVars(context.Background(), query, make(map[string]string))
	if nil != err {
		return "{}", err
	}

	var rootObject interface{}

	err = json.Unmarshal(res.GetJson(), &rootObject)
	if nil != err {
		return nil, err
	}

	return rootObject, nil
}

//-====================-//
//    Upsert Graph
//-====================-//

func (this *DgraphService) UpsertGraph(graph model.Graph, graphToo map[string]interface{}) error {

	log.Debug("(UpsertGraph) begin - graph = ", graph)

	log.Debug("graph.GetNodes() = ", graph.GetNodes())

	pendingNodes := make(map[string]*rdf.DNode)
	for id, node := range graph.GetNodes() {
		log.Debug("node id = ", id, ", node = ", node)
		pendingNodes[id.ToString()] = rdf.NewDNode(node, this._explicitType, this._typeName, this._addPrefixToAttr)
	}

	log.Debug("graph.GetEdges() = ", graph.GetEdges())
	pendingEdges := make(map[string]*rdf.DEdge)
	for id, edge := range graph.GetEdges() {
		log.Debug("edge id = ", id, ", edge = ", edge)
		from := pendingNodes[edge.GetFromId().ToString()]
		to := pendingNodes[edge.GetToId().ToString()]
		pendingEdges[id.ToString()] = rdf.NewDEdge(edge, this._explicitType, this._typeName, this._addPrefixToAttr, from, to)
	}

	err := this.Commit(pendingNodes, pendingEdges, graph)

	log.Debug("(UpsertGraph) Done ! ")

	return err
}

/* Need to revisit the transaction scope */
func (this *DgraphService) Commit(
	pendingNodes map[string]*rdf.DNode,
	pendingEdges map[string]*rdf.DEdge,
	graph model.Graph) error {

	this._mux.Lock()
	defer this._mux.Unlock()

	log.Debug("(Commit) begin - ", graph)

	ctx := context.Background()
	txn := this._dgraphClient.NewTxn()
	defer txn.Discard(ctx)

	for id, dEdge := range pendingEdges {
		log.Debug("XXXXXXXXXXXXXXX edge id : ", dEdge.GetId())
		cachedEdge := this._cachedEdges.Get(id)

		log.Debug("newEdge = ", dEdge, ", cached edge = ", cachedEdge)

		if nil == cachedEdge {
			log.Debug("(commit) edge NOT in cache, eid = ", id)

			/* edge not found in cache */
			/* check exists from remote server */
			err := this.checkEdge(txn, graph, dEdge)
			if nil != err {
				log.Error("Error from checkEdge : ", err)
				return err
			}
			/* Update cache */
			this._cachedEdges.Add(id, dEdge)

			/* populate dgraph uid for from node */
			from := dEdge.GetFrom()
			fromInPending := pendingNodes[from.GetId()]
			if nil != fromInPending {
				if from.Exists() {
					fromInPending.SetUid(from.GetUid())
					/* overwrite from node in cache */
					this._cachedNodes.Add(from.GetId(), fromInPending)
				}
			}

			/* populate dgraph uid for to node */
			to := dEdge.GetTo()
			toInPending := pendingNodes[to.GetId()]
			if nil != toInPending {
				if to.Exists() {
					toInPending.SetUid(to.GetUid())
					/* overwrite to node in cache */
					this._cachedNodes.Add(to.GetId(), toInPending)
				}
			}
		} else {
			log.Debug("(commit) edge FOUND in cache, eid = ", id)
			edgeChanged, err := cachedEdge.(*rdf.DEdge).Update(dEdge.GetEdge())
			if nil != err {
				log.Error("Error : ", err)
				return err
			}
			if !edgeChanged {
				/* edge not changed */
				/* no need for upsert */
				delete(pendingEdges, id)
			}
		}
	}

	var Nquads bytes.Buffer
	for id, dNode := range pendingNodes {
		if !dNode.Exists() {
			log.Info("(commit) node has NO UID ! nid = ", id)
			cachedNode := this._cachedNodes.Get(id)

			log.Info("newNode = ", dNode, ", cached node = ", cachedNode)

			if nil == cachedNode || "" == cachedNode.(*rdf.DNode).GetUid() {
				log.Debug("(commit) node NOT in cache or cached node has no UID !")
				err := this.checkNode(txn, graph, dNode)
				if nil != err {
					log.Error("Error from checkNode : ", err)
					return err
				}
				/* Add to cache */
				this._cachedNodes.Add(id, dNode)
			} else {
				log.Debug("(commit) node FOUND in cache !")
				nodeChanged, err := cachedNode.(*rdf.DNode).Update(dNode.GetNode())
				if nil != err {
					log.Error("Error : ", err)
					return err
				}
				dNode.SetUid(cachedNode.(*rdf.DNode).GetUid())
				if !nodeChanged {
					/* node has no change */
					/* don't need to upsert */
					continue
				}
			}
		} else {
			log.Debug("(commit) node has UID ! nid = ", id)
		}

		for _, data := range dNode.ToRDF(graph, dgraph.DATE_TIME_SAMPLE, this._readableExternalId) {
			Nquads.WriteString(data)
		}
	}

	for _, edge := range pendingEdges {
		for _, data := range edge.ToRDF(graph, dgraph.DATE_TIME_SAMPLE, this._readableExternalId) {
			Nquads.WriteString(data)
		}
	}

	log.Info("========================= Nquads.String() ===========================")
	log.Info(Nquads.String())
	log.Info("=====================================================================")

	if 0 == Nquads.Len() {
		log.Debug("(DgraphService::Commit) No data commited !!")
		return nil
	}

	_, err := txn.Mutate(
		ctx,
		&api.Mutation{
			SetNquads: []byte(Nquads.String()),
		},
	)

	if nil != err {
		return err
	}

	err = txn.Commit(ctx)

	/* Since it's local data structure, why delete? */
	if nil != err {
		/* keep all pending entities */
		log.Error("(Commit) Fail with error : ", err)
		return err
	}

	log.Debug("(Commit) Done ! ")

	return nil
}

func (this *DgraphService) checkNode(txn *dgo.Txn, graph model.Graph, node *rdf.DNode) error {

	var query bytes.Buffer
	query.WriteString("query {\n")
	query.WriteString("node(func: eq(")
	query.WriteString(node.GetTypeName())
	query.WriteString(", \"")
	query.WriteString(util.ReplaceCharacter(node.GetType(), this._targetRegex, this._replacement, true))
	query.WriteString("\"")
	query.WriteString("))")

	count := 0
	for _, keyAttributeName := range graph.GetEntityKeyNamesForNode(node.GetType()) {
		if 0 == count {
			query.WriteString(" @filter(")
		} else {
			query.WriteString(" and")
		}
		query.WriteString(" eq(")
		query.WriteString(GetCanonicalAttributeName(node.GetType(), keyAttributeName, this._targetRegex, this._replacement, this._addPrefixToAttr, true))
		query.WriteString(", \"")
		query.WriteString(node.GetAttributeAsString(keyAttributeName))
		query.WriteString("\")")
		count++
	}

	query.WriteString(") {\n")
	query.WriteString("  uid\n  ")
	query.WriteString("  }\n")
	query.WriteString("}\n")

	log.Info("************** query node ****************")
	log.Info(query.String())
	log.Info("******************************************")

	res, err := txn.QueryWithVars(context.Background(), query.String(), make(map[string]string))
	if nil != err {
		return err
	}

	targetNodeData := make(map[string]interface{})
	err = json.Unmarshal(res.GetJson(), &targetNodeData)
	if nil != err {
		return err
	}

	log.Debug("Trying to find node : ", node.ToString())
	nodeArray := targetNodeData["node"].([]interface{})
	if nil != nodeArray && 0 < len(nodeArray) {

		if 1 < len(nodeArray) {
			log.Warnf("%s is found duplicated ! \n", node.ToString())
		}

		node.SetUid(nodeArray[0].(map[string]interface{})["uid"].(string))
	}

	if !node.Exists() {
		log.Debugf("Node is not found : %s \n", node.ToString())
	}

	return nil
}

func (this *DgraphService) checkEdge(txn *dgo.Txn, graph model.Graph, edge *rdf.DEdge) error {
	var query bytes.Buffer
	query.WriteString("query {\n")
	query.WriteString("edge(func: has(")
	query.WriteString(util.ReplaceCharacter(edge.GetType(), this._targetRegex, this._replacement, true))
	query.WriteString(")) @filter(eq(")

	query.WriteString(edge.GetFrom().GetTypeName())
	query.WriteString(", \"")
	query.WriteString(util.ReplaceCharacter(edge.GetFrom().GetType(), this._targetRegex, this._replacement, true))
	query.WriteString("\")")

	for _, keyAttributeName := range graph.GetEntityKeyNamesForNode(edge.GetFrom().GetType()) {
		query.WriteString(" and eq(")
		query.WriteString(GetCanonicalAttributeName(edge.GetFrom().GetType(), keyAttributeName, this._targetRegex, this._replacement, this._addPrefixToAttr, true))
		query.WriteString(", \"")
		query.WriteString(edge.GetFrom().GetAttributeAsString(keyAttributeName))
		query.WriteString("\")")
	}

	query.WriteString(") {\n")
	query.WriteString("  ")
	query.WriteString(edge.GetFrom().GetTypeName())
	query.WriteString("\n  ")
	query.WriteString("  uid\n  ")
	query.WriteString("  ")
	/* edge */
	query.WriteString(util.ReplaceCharacter(edge.GetType(), this._targetRegex, this._replacement, true))

	count := 0
	for _, keyAttributeName := range graph.GetEntityKeyNamesForEdge(edge.GetType()) {
		if 0 == count {
			query.WriteString(" @facets(eq(")
			query.WriteString(keyAttributeName)
			query.WriteString(", \"")
			query.WriteString(edge.GetAttributeAsString(keyAttributeName))
			query.WriteString("\")")
		} else {
			query.WriteString(" and eq(")
			query.WriteString(keyAttributeName)
			query.WriteString(", \"")
			query.WriteString(edge.GetAttributeAsString(keyAttributeName))
			query.WriteString("\")")
		}
		count++
	}

	count = 0
	for _, keyAttributeName := range graph.GetEntityKeyNamesForNode(edge.GetTo().GetType()) {
		if 0 == count {
			query.WriteString(" @filter(eq(")
			query.WriteString(GetCanonicalAttributeName(edge.GetTo().GetType(), keyAttributeName, this._targetRegex, this._replacement, this._addPrefixToAttr, true))
			query.WriteString(", \"")
			query.WriteString(edge.GetTo().GetAttributeAsString(keyAttributeName))
			query.WriteString("\")")
		} else {
			query.WriteString(" and eq(")
			query.WriteString(GetCanonicalAttributeName(edge.GetTo().GetType(), keyAttributeName, this._targetRegex, this._replacement, this._addPrefixToAttr, true))
			query.WriteString(", \"")
			query.WriteString(edge.GetTo().GetAttributeAsString(keyAttributeName))
			query.WriteString("\")")
		}
		count++
	}

	if 0 != count {
		query.WriteString(")")
	}

	query.WriteString(" {\n")
	query.WriteString("        type\n  ")
	query.WriteString("        uid\n  ")
	query.WriteString("      }\n  ")
	query.WriteString("    }\n  ")
	query.WriteString("  }\n  ")

	log.Info("************** query edge ****************")
	log.Info(query.String())
	log.Info("******************************************")

	res, err := txn.QueryWithVars(context.Background(), query.String(), make(map[string]string))
	if nil != err {
		return err
	}

	targetEdgeData := make(map[string]interface{})
	err = json.Unmarshal(res.GetJson(), &targetEdgeData)
	if nil != err {
		return err
	}

	log.Debug("Trying to find edge : ", edge.ToString())

	edgeArray := targetEdgeData["edge"].([]interface{})
	if nil != edgeArray && 0 < len(edgeArray) {

		if 1 < len(edgeArray) {
			log.Warnf("%s is found duplicated !\n", edge.ToString())
		}

		fromObj := edgeArray[0].(map[string]interface{})
		if nil != fromObj && edge.GetFrom().GetType() == fromObj["type"] {
			/************************************/
			/* We get internal id for from node */
			/************************************/
			edge.GetFrom().SetUid(fromObj["uid"].(string))
		}

		if nil != fromObj[edge.GetType()] {
			innerEdgeArray := fromObj[edge.GetType()].([]interface{})
			if nil != innerEdgeArray && 0 < len(innerEdgeArray) {
				for i := 0; i < len(innerEdgeArray); i++ {
					toObj := innerEdgeArray[i].(map[string]interface{})
					if nil != toObj && edge.GetTo().GetType() == toObj["type"] {
						/**********************************/
						/* We get internal id for to node */
						/**********************************/
						edge.GetTo().SetUid(toObj["uid"].(string))
						break
					}
				}
			}
		}

		if !edge.GetFrom().Exists() {
			log.Warn("From node is not found !")
		}

		if !edge.GetTo().Exists() {
			log.Warn("To node is not found !")
		}

	} else {
		log.Debug("Edge not found !")
	}

	return nil
}
