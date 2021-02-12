/*
 * Copyright © 2020. TIBCO Software Inc.
 * This file is subject to the license terms contained
 * in the license file that is distributed with this file.
 */
package neo4j

import (
	"sync"

	"github.com/TIBCOSoftware/flogo-lib/logger"
	"github.com/TIBCOSoftware/labs-graphbuilder-lib/dbservice"
)

var log = logger.GetLogger("neo4jupsert")

type Neo4jServiceFactory struct {
	dbservice.BaseDBServiceFactory
	mux sync.Mutex
}

func (this *Neo4jServiceFactory) GetImportService(serviceId string) dbservice.ImportService {
	dbService := this.DBServices[serviceId]
	if nil != dbService {
		return dbService.(dbservice.ImportService)
	}
	return nil
}

func (this *Neo4jServiceFactory) CreateImportService(serviceId string, properties map[string]interface{}) (dbservice.ImportService, error) {
	this.mux.Lock()
	defer this.mux.Unlock()

	neo4jService := &Neo4jImportCSV{}
	neo4jService.initialized = false
	neo4jService.nodeFileNames = make(map[string][]string)
	neo4jService.edgeFileNames = make(map[string][]string)

	neo4jService.label = make(map[string]string)
	neo4jService.ignoredNodeAttributes = make(map[string]map[string]bool)
	neo4jService.ignoredEdgeAttributes = make(map[string]map[string]bool)
	neo4jService.nodeAttributeWrittenOrder = make(map[string][]string)
	neo4jService.edgeAttributeWrittenOrder = make(map[string][]string)

	if nil != properties["database"] {
		neo4jService.database = properties["database"].(string)
	}

	if nil != properties["typeName"] {
		neo4jService.typeName = properties["typeName"].(string)
	}

	if nil != properties["outputFolder"] {
		neo4jService.outputFolder = properties["outputFolder"].(string)
	}

	return nil, nil
}

func (this *Neo4jServiceFactory) GetUpsertService(serviceId string) dbservice.UpsertService {
	if nil != this.DBServices[serviceId] {
		return this.DBServices[serviceId].(dbservice.UpsertService)
	} else {
		return nil
	}
}

func (this *Neo4jServiceFactory) CreateUpsertService(serviceId string, properties map[string]interface{}) (dbservice.UpsertService, error) {
	this.mux.Lock()
	defer this.mux.Unlock()

	neo4jService := &Neo4jService{}
	if nil != properties["url"] {
		neo4jService._url = properties["url"].(string)
	}
	if nil != properties["user"] {
		neo4jService._user = properties["user"].(string)
	}
	if nil != properties["password"] {
		neo4jService._password = properties["password"].(string)
	}

	//neo4jService._typeName = util.CastString(properties["typeName"])
	//neo4jService._addPrefixToAttr = properties["addPrefixToAttr"].(bool)
	neo4jService._targetRegex = "[^A-Za-z0-9]"
	neo4jService._replacement = "_"

	this.DBServices[serviceId] = neo4jService

	return neo4jService, nil
}

func NewNeo4jServiceFactory() *Neo4jServiceFactory {
	return &Neo4jServiceFactory{}
}
