/*
 * Copyright © 2019. TIBCO Software Inc.
 * This file is subject to the license terms contained
 * in the license file that is distributed with this file.
 */
package sseserver

import (
	//	"context"
	"sync"

	//	"encoding/json"
	//	"fmt"
	//	"time"

	"github.com/TIBCOSoftware/flogo-lib/core/activity"
	"github.com/TIBCOSoftware/flogo-lib/core/data"
	"github.com/TIBCOSoftware/flogo-lib/core/trigger"
	"github.com/TIBCOSoftware/flogo-lib/logger"
	"github.com/TIBCOSoftware/labs-graphbuilder-lib/internet/sseserver"
)

var log = logger.GetLogger("trigger-sse-server")

const (
	cServerPort     = "port"
	cConnection     = "sseConnection"
	cConnectionName = "name"
)

//-============================================-//
//   Entry point create a new Trigger factory
//-============================================-//

func NewFactory(md *trigger.Metadata) trigger.Factory {
	return &SSEServerFactory{metadata: md}
}

//-===============================-//
//     Define Trigger Factory
//-===============================-//

type SSEServerFactory struct {
	metadata *trigger.Metadata
}

func (t *SSEServerFactory) New(config *trigger.Config) trigger.Trigger {
	return &SSEServer{metadata: t.metadata, config: config}
}

//-=========================-//
//      Define Trigger
//-=========================-//

type SSEServer struct {
	metadata *trigger.Metadata
	config   *trigger.Config
	server   *sseserver.Server
	mux      sync.Mutex

	handlers []*trigger.Handler
}

// implements trigger.Trigger.Metadata (trigger.go)
func (this *SSEServer) Metadata() *trigger.Metadata {
	return this.metadata
}

// implements trigger.Initializable.Initialize
func (this *SSEServer) Initialize(ctx trigger.InitContext) error {

	this.handlers = ctx.GetHandlers()
	return nil
}

// implements ext.Trigger.Start
func (this *SSEServer) Start() error {

	log.Debug("Start")
	handlers := this.handlers

	log.Debug("Processing handlers")

	connection, exist := handlers[0].GetSetting(cConnection)
	if !exist {
		return activity.NewError("SSE connection is not configured", "TGDB-SSE-4001", nil)
	}

	connectionInfo, _ := data.CoerceToObject(connection)
	if connectionInfo == nil {
		return activity.NewError("SSE connection not able to be parsed", "TGDB-SSE-4002", nil)
	}

	var serverId string
	properties := make(map[string]interface{})
	connectionSettings, _ := connectionInfo["settings"].([]interface{})
	if connectionSettings != nil {
		for _, v := range connectionSettings {
			setting, _ := data.CoerceToObject(v)
			if setting != nil {
				if setting["name"] == cServerPort {
					properties[cServerPort], _ = data.CoerceToString(setting["value"])
				} else if setting["name"] == cConnectionName {
					serverId = setting["value"].(string)
				}
			}
		}
		log.Info(properties)

		this.server, _ = sseserver.GetFactory().CreateServer(serverId, properties)
		go this.server.Start()
	}

	return nil
}

// implements ext.Trigger.Stop
func (this *SSEServer) Stop() error {
	this.server.Stop()
	return nil
}
