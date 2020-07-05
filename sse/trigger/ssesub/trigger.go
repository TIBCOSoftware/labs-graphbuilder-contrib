/*
 * Copyright © 2020. TIBCO Software Inc.
 * This file is subject to the license terms contained
 * in the license file that is distributed with this file.
 */
package ssesub

import (
	"context"
	"sync"

	"github.com/TIBCOSoftware/flogo-lib/core/activity"
	"github.com/TIBCOSoftware/flogo-lib/core/data"
	"github.com/TIBCOSoftware/labs-graphbuilder-lib/internet/sseservice"
	"github.com/project-flogo/core/data/metadata"
	"github.com/project-flogo/core/support/log"
	"github.com/project-flogo/core/trigger"
)

const (
	Connection = "sseConnection"
)

//-============================================-//
//   Entry point register Trigger & factory
//-============================================-//

var triggerMd = trigger.NewMetadata(&Settings{}, &HandlerSettings{}, &Output{})

func init() {
	_ = trigger.Register(&SSESubscriber{}, &Factory{})
}

//-===============================-//
//     Define Trigger Factory
//-===============================-//

type Factory struct {
}

// Metadata implements trigger.Factory.Metadata
func (*Factory) Metadata() *trigger.Metadata {
	return triggerMd
}

// New implements trigger.Factory.New
func (*Factory) New(config *trigger.Config) (trigger.Trigger, error) {
	settings := &Settings{}
	err := metadata.MapToStruct(config.Settings, settings, true)
	if err != nil {
		return nil, err
	}

	return &SSESubscriber{settings: settings}, nil
}

//-=========================-//
//      Define Trigger
//-=========================-//

var logger log.Logger

type SSESubscriber struct {
	metadata   *trigger.Metadata
	config     *trigger.Config
	sseService sse.SSEService
	mux        sync.Mutex

	settings *Settings
	handlers []trigger.Handler
}

// implements trigger.Initializable.Initialize
func (this *SSESubscriber) Initialize(ctx trigger.InitContext) error {

	this.handlers = ctx.GetHandlers()
	logger = ctx.Logger()

	return nil
}

// implements ext.Trigger.Start
func (this *SSESubscriber) Start() error {

	logger.Debug("Start")
	handlers := this.handlers

	logger.Debug("Processing handlers")

	connection, exist := handlers[0].Settings()[Connection]
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
		logger.Info(properties)

		this.sseService = sse.NewSSEServiceFactory().GetService(properties)
		this.sseService.SetEventListener(this)
		go this.sseService.Start()

	}

	return nil
}

// implements ext.Trigger.Stop
func (this *SSESubscriber) Stop() error {
	logger.Debug("Stopping endpoints")
	this.sseService.Stop()
	return nil
}

func (this *SSESubscriber) ProcessEvent(event string) error {
	this.mux.Lock()
	defer this.mux.Unlock()
	logger.Debug("Got SSE Even : ", event)
	outputData := &Output{}
	outputData.Event = event
	logger.Debug("Send SSE Even out : ", outputData)

	_, err := this.handlers[0].Handle(context.Background(), outputData)
	if nil != err {
		logger.Info("Error -> ", err)
	}

	return err
}
