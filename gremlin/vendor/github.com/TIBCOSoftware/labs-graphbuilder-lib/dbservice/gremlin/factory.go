/*
 * Copyright © 2020. TIBCO Software Inc.
 * This file is subject to the license terms contained
 * in the license file that is distributed with this file.
 */
package gremlin

import (
	"sync"

	"github.com/TIBCOSoftware/flogo-lib/logger"
	"github.com/TIBCOSoftware/labs-graphbuilder-lib/dbservice"
)

var log = logger.GetLogger("gremlin")

type GremlinServiceFactory struct {
	dbservice.BaseDBServiceFactory
	mux sync.Mutex
}

func (this *GremlinServiceFactory) GetImportService(serviceId string) dbservice.ImportService {
	dbService := this.DBServices[serviceId]
	if nil != dbService {
		return dbService.(dbservice.ImportService)
	}
	return nil
}

func (this *GremlinServiceFactory) CreateImportService(serviceId string, properties map[string]interface{}) (dbservice.ImportService, error) {
	return nil, nil
}

func (this *GremlinServiceFactory) GetUpsertService(serviceId string) dbservice.UpsertService {
	dbService := this.DBServices[serviceId]
	if nil != dbService {
		return dbService.(dbservice.UpsertService)
	}
	return nil
}

func (this *GremlinServiceFactory) CreateUpsertService(serviceId string, properties map[string]interface{}) (dbservice.UpsertService, error) {
	this.mux.Lock()
	defer this.mux.Unlock()

	gremlinService := &GremlinService{}
	if nil != properties["url"] {
		gremlinService._url = properties["url"].(string)
	}
	if nil != properties["user"] {
		gremlinService._user = properties["user"].(string)
	}
	if nil != properties["password"] {
		gremlinService._password = properties["password"].(string)
	}

	//gremlinService._typeName = util.CastString(properties["typeName"])
	//gremlinService._addPrefixToAttr = properties["addPrefixToAttr"].(bool)
	gremlinService._targetRegex = "[^A-Za-z0-9]"
	gremlinService._replacement = "_"

	return gremlinService, nil
}

func NewGremlinServiceFactory() *GremlinServiceFactory {
	return &GremlinServiceFactory{}
}
