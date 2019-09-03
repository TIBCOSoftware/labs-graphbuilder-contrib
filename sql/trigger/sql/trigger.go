/*
 * Copyright © 2019. TIBCO Software Inc.
 * This file is subject to the license terms contained
 * in the license file that is distributed with this file.
 */
package sql

import (
	"context"
	"encoding/json"
	"sync"

	"github.com/TIBCOSoftware/flogo-lib/core/activity"
	"github.com/TIBCOSoftware/flogo-lib/core/data"
	"github.com/TIBCOSoftware/flogo-lib/core/trigger"
	"github.com/TIBCOSoftware/flogo-lib/logger"
	"github.com/TIBCOSoftware/labs-graphbuilder-lib/sql"
)

var log = logger.GetLogger("trigger-sql")

type Field struct {
	DBColumn  string `json:"DBColumn"`
	FieldName string `json:"FieldName"`
	Type      string `json:"Type"`
}

const (
	SQLString      = "queryString"
	OutputFieldMap = "outputFieldMap"
	Connection     = "sqlConnection"
)

//-============================================-//
//   Entry point create a new Trigger factory
//-============================================-//

func NewFactory(md *trigger.Metadata) trigger.Factory {
	return &SQLSubscriberFactory{metadata: md}
}

//-===============================-//
//     Define Trigger Factory
//-===============================-//

type SQLSubscriberFactory struct {
	metadata *trigger.Metadata
}

func (t *SQLSubscriberFactory) New(config *trigger.Config) trigger.Trigger {
	return &SQLSubscriber{metadata: t.metadata, config: config}
}

//-=========================-//
//      Define Trigger
//-=========================-//

type SQLSubscriber struct {
	metadata   *trigger.Metadata
	config     *trigger.Config
	sqlService sql.SQLService
	mux        sync.Mutex

	handlers []*trigger.Handler
}

// implements trigger.Trigger.Metadata (trigger.go)
func (this *SQLSubscriber) Metadata() *trigger.Metadata {
	return this.metadata
}

// implements trigger.Initializable.Initialize
func (this *SQLSubscriber) Initialize(ctx trigger.InitContext) error {

	this.handlers = ctx.GetHandlers()
	return nil
}

// implements ext.Trigger.Start
func (this *SQLSubscriber) Start() error {

	log.Debug("Start")
	handlers := this.handlers

	log.Debug("Processing handlers")

	connection, exist := handlers[0].GetSetting(Connection)
	if !exist {
		return activity.NewError("SSE connection is not configured", "TGDB-SSE-4001", nil)
	}

	connectionInfo, _ := data.CoerceToObject(connection)
	if connectionInfo == nil {
		return activity.NewError("SSE connection not able to be parsed", "TGDB-SSE-4002", nil)
	}

	properties := make(map[string]interface{})
	connectionSettings, _ := connectionInfo["settings"].([]interface{})
	if connectionSettings != nil {
		for _, v := range connectionSettings {
			setting, _ := data.CoerceToObject(v)
			if setting != nil {
				if setting["name"] == "vendor" {
					properties["vendor"], _ = data.CoerceToString(setting["value"])
				} else if setting["name"] == "url" {
					properties["url"], _ = data.CoerceToString(setting["value"])
				} else if setting["name"] == "user" {
					properties["user"], _ = data.CoerceToString(setting["value"])
				} else if setting["name"] == "password" {
					properties["password"], _ = data.CoerceToString(setting["value"])
				} else if setting["name"] == "resource" {
					properties["resource"], _ = data.CoerceToString(setting["value"])
				}
			}
		}
		log.Info(properties)

		sqlString, exist := handlers[0].GetSetting(SQLString)
		if !exist {
			return activity.NewError("Query string is not configured", "SQL-SERVICE-4003", nil)
		}

		outputFieldMap := handlers[0].GetStringSetting(OutputFieldMap)
		var outputFields []Field
		json.Unmarshal([]byte(outputFieldMap), &outputFields)

		log.Info("outputFields = ", outputFields)

		fieldMap := make(map[string](map[string]string))
		for _, outputField := range outputFields {
			field := make(map[string]string)
			field["FieldName"] = outputField.FieldName
			field["Type"] = outputField.Type
			fieldMap[outputField.DBColumn] = field
		}

		log.Info("(Start) - fieldMap = ", fieldMap)

		this.sqlService = sql.NewSQLServiceFactory().GetService(properties)
		this.sqlService.SetQueryResultListener(this)
		go this.sqlService.Start(sqlString.(string), fieldMap)

	}

	return nil
}

// implements ext.Trigger.Stop
func (this *SQLSubscriber) Stop() error {
	this.sqlService.Stop()
	return nil
}

func (this *SQLSubscriber) ProcessRow(dataRow map[string]interface{}) error {
	//og.Info("Data row : ", dataRow)
	output := make(map[string]interface{})
	//	output["DataRow"] = dataRow

	complexcsvdata := &data.ComplexObject{Metadata: "DataRow", Value: dataRow}
	complexDataRow, _ := data.NewAttribute(complexcsvdata.Metadata, data.TypeComplexObject, complexcsvdata)
	output["dataRow"] = complexDataRow
	log.Info("Output : ", output)

	_, err := this.handlers[0].Handle(context.Background(), output)
	if nil != err {
		log.Info("Error -> ", err)
	}

	return err
}
