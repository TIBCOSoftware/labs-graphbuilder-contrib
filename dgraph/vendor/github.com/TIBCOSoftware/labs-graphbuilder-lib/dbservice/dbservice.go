/*
 * Copyright © 2020. TIBCO Software Inc.
 * This file is subject to the license terms contained
 * in the license file that is distributed with this file.
 */
package dbservice

import (
	"github.com/TIBCOSoftware/labs-graphbuilder-lib/model"
)

type DBType int

const (
	TGDB   DBType = 0
	Dgraph DBType = 1
	Neo4j  DBType = 2
)

func (this DBType) int() int {
	index := [...]int{0, 1, 2}
	return index[this]
}

type DBServiceFactory interface {
	Initialize()
	GetUpsertService(serviceId string) UpsertService
	CreateUpsertService(serviceId string, properties map[string]interface{}) (UpsertService, error)
	GetImportService(serviceId string) ImportService
	CreateImportService(serviceId string, properties map[string]interface{}) (ImportService, error)
}

type DBService interface {
}

type UpsertService interface {
	DBService
	UpsertGraph(graph model.Graph, graphToo map[string]interface{}) error
	DeleteGraph(filter int, graphToo map[string]interface{}) error
}

type ImportService interface {
	DBService
	WriteGraph(graph model.Graph) error
}

type BaseDBServiceFactory struct {
	DBServices map[string]DBService
}

func (this *BaseDBServiceFactory) Initialize() {
}

func (this *BaseDBServiceFactory) GetUpsertService(serviceId string) UpsertService {
	dbService := this.DBServices[serviceId]
	if nil != dbService {
		return dbService.(UpsertService)
	}
	return nil
}

func (this *BaseDBServiceFactory) GetImportService(serviceId string) ImportService {
	dbService := this.DBServices[serviceId]
	if nil != dbService {
		return dbService.(ImportService)
	}
	return nil
}
