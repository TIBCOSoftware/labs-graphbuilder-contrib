/*
 * Copyright © 2020. TIBCO Software Inc.
 * This file is subject to the license terms contained
 * in the license file that is distributed with this file.
 */
package tgdb

import (
	"errors"
	"sync"

	"tgdb"
	"tgdb/impl"

	log "github.com/TIBCOSoftware/flogo-lib/logger"
	"github.com/TIBCOSoftware/labs-graphbuilder-lib/dbservice"
)

const (
	Query                    = "query"
	Query_Language           = "language"
	Query_QueryString        = "queryString"
	Query_TraversalCondition = "traversalCondition"
	Query_EdgeFilter         = "edgeFilter"
	Query_EndCondition       = "endCondition"
	Query_OPT_PrefetchSize   = "prefetchSize"
	Query_OPT_EdgeLimit      = "edgeLimit"
	Query_OPT_TraversalDepth = "traversalDepth"
)

var logger = log.GetLogger("tibco-tgdb-service")

//var logger = DefaultTGLogManager().GetLogger()

type TGDBServiceFactory struct {
	dbservice.BaseDBServiceFactory
	mux sync.Mutex
}

func (this *TGDBServiceFactory) Initialize() {
	impl.DefaultTGLogManager().GetLogger().SetLogLevel(tgdb.ErrorLog)
}

func (this *TGDBServiceFactory) GetImportService(serviceId string) dbservice.ImportService {
	return this.DBServices[serviceId].(dbservice.ImportService)
}

func (this *TGDBServiceFactory) CreateImportService(serviceId string, properties map[string]interface{}) (dbservice.ImportService, error) {

	this.mux.Lock()
	defer this.mux.Unlock()
	tgdbService := this.DBServices[serviceId].(*TGDBImportCSV)

	if nil == tgdbService {
		tgdbService = NewTGDBImportCSV()
		if nil != properties["outputFolder"] {
			tgdbService.importFileFolder = properties["outputFolder"].(string)
		} else {
			return nil, errors.New("outputFolder not set!")
		}

		this.DBServices[serviceId] = tgdbService
	}

	return tgdbService, nil
}

func (this *TGDBServiceFactory) GetUpsertService(serviceId string) dbservice.UpsertService {
	if nil != this.DBServices[serviceId] {
		return this.DBServices[serviceId].(dbservice.UpsertService)
	} else {
		return nil
	}
}

func (this *TGDBServiceFactory) CreateUpsertService(serviceId string, properties map[string]interface{}) (dbservice.UpsertService, error) {

	this.mux.Lock()
	defer this.mux.Unlock()

	if nil != this.DBServices[serviceId] {
		return this.DBServices[serviceId].(*TGDBService), nil
	}

	var tgdbService *TGDBService
	if nil == tgdbService {
		tgdbService = &TGDBService{}
		if nil != properties["url"] {
			tgdbService._url = properties["url"].(string)
		}
		if nil != properties["user"] {
			tgdbService._user = properties["user"].(string)
		}
		if nil != properties["password"] {
			tgdbService._password = properties["password"].(string)
		}
		if nil != properties["connectionProps"] {
			tgdbService._connectionProps = properties["connectionProps"].(map[string]string)
		} else {
			tgdbService._connectionProps = make(map[string]string)
		}
		if nil != properties["keepAlive"] {
			tgdbService._keepAlive = properties["keepAlive"].(bool)
		} else {
			tgdbService._keepAlive = false
		}
		if nil != properties["allowEmptyStringKey"] {
			tgdbService._allowEmptyKey = properties["allowEmptyStringKey"].(bool)
		} else {
			tgdbService._allowEmptyKey = false
		}
		tgdbService._keyMap = make(map[string][]string)
		tgdbService._nodeAttrMap = make(map[string]map[string]tgdb.TGAttributeDescriptor)
		tgdbService._edgeAttrMap = make(map[string]map[string]tgdb.TGAttributeDescriptor)

		err := tgdbService.ensureConnection()
		if nil != err {
			return nil, err
		}

		this.DBServices[serviceId] = tgdbService
	}

	return tgdbService, nil
}
