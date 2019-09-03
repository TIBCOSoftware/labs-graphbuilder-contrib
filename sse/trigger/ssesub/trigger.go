/*
 * Copyright © 2019. TIBCO Software Inc.
 * This file is subject to the license terms contained
 * in the license file that is distributed with this file.
 */
package ssesub

import (
	"context"
	"sync"

	"github.com/TIBCOSoftware/flogo-lib/core/activity"
	"github.com/TIBCOSoftware/flogo-lib/core/data"
	"github.com/TIBCOSoftware/flogo-lib/core/trigger"
	"github.com/TIBCOSoftware/flogo-lib/logger"
	"github.com/TIBCOSoftware/labs-graphbuilder-lib/internet/sseservice"
)

var log = logger.GetLogger("trigger-sse")

const (
	Connection = "sseConnection"
)

//-============================================-//
//   Entry point create a new Trigger factory
//-============================================-//

func NewFactory(md *trigger.Metadata) trigger.Factory {
	return &SSESubscriberFactory{metadata: md}
}

//-===============================-//
//     Define Trigger Factory
//-===============================-//

type SSESubscriberFactory struct {
	metadata *trigger.Metadata
}

func (t *SSESubscriberFactory) New(config *trigger.Config) trigger.Trigger {
	return &SSESubscriber{metadata: t.metadata, config: config}
}

//-=========================-//
//      Define Trigger
//-=========================-//

type SSESubscriber struct {
	metadata   *trigger.Metadata
	config     *trigger.Config
	sseService sse.SSEService
	mux        sync.Mutex

	handlers []*trigger.Handler
}

// implements trigger.Trigger.Metadata (trigger.go)
func (this *SSESubscriber) Metadata() *trigger.Metadata {
	return this.metadata
}

// implements trigger.Initializable.Initialize
func (this *SSESubscriber) Initialize(ctx trigger.InitContext) error {

	this.handlers = ctx.GetHandlers()
	return nil
}

// implements ext.Trigger.Start
func (this *SSESubscriber) Start() error {

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
				if setting["name"] == "url" {
					properties["url"], _ = data.CoerceToString(setting["value"])
				} else if setting["name"] == "resource" {
					properties["resource"], _ = data.CoerceToString(setting["value"])
				} else if setting["name"] == "accessToken" {
					properties["accessToken"], _ = data.CoerceToString(setting["value"])
				}
			}
		}
		log.Info(properties)

		this.sseService = sse.NewSSEServiceFactory().GetService(properties)
		this.sseService.SetEventListener(this)
		go this.sseService.Start()

	}

	return nil
}

// implements ext.Trigger.Stop
func (this *SSESubscriber) Stop() error {
	this.sseService.Stop()
	return nil
}

func (this *SSESubscriber) ProcessEvent(event string) error {
	//	log.Info("Even : ", event)
	output := make(map[string]interface{})
	output["Event"] = event
	log.Debug(output)
	_, err := this.handlers[0].Handle(context.Background(), output)
	if nil != err {
		log.Info("Error -> ", err)
	}

	return err
}
