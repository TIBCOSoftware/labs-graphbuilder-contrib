/*
 * Copyright © 2020. TIBCO Software Inc.
 * This file is subject to the license terms contained
 * in the license file that is distributed with this file.
 */
package factory

import (
	"sync"

	"github.com/TIBCOSoftware/flogo-lib/logger"
	"github.com/TIBCOSoftware/labs-graphbuilder-lib/dbservice"
	"github.com/TIBCOSoftware/labs-graphbuilder-lib/dbservice/dgraph/v2"
)

var log = logger.GetLogger("dgraph-service")

type DgraphServiceFactory struct {
	dbservice.BaseDBServiceFactory
	mux sync.Mutex
}

func (this *DgraphServiceFactory) CreateImportService(serviceId string, properties map[string]interface{}) (dbservice.ImportService, error) {
	this.mux.Lock()
	defer this.mux.Unlock()

	version := properties["version"].(string)
	log.Info("(DgraphServiceFactory.CreateImportService) API Version : ", version)

	dgraphService := this.DBServices[serviceId]
	switch version {
	case "v1":
	default:
		if nil == dgraphService {
			dgraphService, _ = v2.NewDgraphImportRDF(properties)
			this.DBServices[serviceId] = dgraphService.(*v2.DgraphImportRDF)
		}
	}
	return dgraphService.(dbservice.ImportService), nil
}

func (this *DgraphServiceFactory) CreateUpsertService(serviceId string, properties map[string]interface{}) (dbservice.UpsertService, error) {
	this.mux.Lock()
	defer this.mux.Unlock()

	version := properties["version"].(string)
	log.Info("(DgraphServiceFactory.CreateUpsertService) API Version : ", version)

	dgraphService := this.DBServices[serviceId]
	var err error
	switch version {
	case "v1":
	default:
		if nil == dgraphService {
			dgraphService, err = v2.NewDgraphService(properties)
			log.Info("(DgraphServiceFactory.CreateUpsertService) dgraphService : ", dgraphService)
			if nil != err {
				log.Info("(DgraphServiceFactory.CreateUpsertService) err : ", err)
				return nil, err
			}
			this.DBServices[serviceId] = dgraphService.(*v2.DgraphService)
		}
	}
	return dgraphService.(dbservice.UpsertService), nil
}
