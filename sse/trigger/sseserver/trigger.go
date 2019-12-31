/*
 * Copyright © 2019. TIBCO Software Inc.
 * This file is subject to the license terms contained
 * in the license file that is distributed with this file.
 */
package sseserver

import (
	"context"
	b64 "encoding/base64"
	"strings"
	"sync"

	"github.com/TIBCOSoftware/flogo-lib/core/activity"
	"github.com/TIBCOSoftware/flogo-lib/core/data"
	"github.com/TIBCOSoftware/labs-graphbuilder-lib/internet/sseserver"
	"github.com/project-flogo/core/data/metadata"
	"github.com/project-flogo/core/support/log"
	"github.com/project-flogo/core/trigger"
)

const (
	cConnection     = "sseConnection"
	cConnectionName = "name"
)

//-============================================-//
//   Entry point register Trigger & factory
//-============================================-//

var triggerMd = trigger.NewMetadata(&Settings{}, &HandlerSettings{}, &Output{})

func init() {
	_ = trigger.Register(&SSEServer{}, &Factory{})
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

	return &SSEServer{settings: settings}, nil
}

//-=========================-//
//      Define Trigger
//-=========================-//

var logger log.Logger

type SSEServer struct {
	metadata *trigger.Metadata
	config   *trigger.Config
	server   *sseserver.Server
	mux      sync.Mutex

	settings *Settings
	handlers []trigger.Handler
}

// implements trigger.Initializable.Initialize
func (this *SSEServer) Initialize(ctx trigger.InitContext) error {

	this.handlers = ctx.GetHandlers()
	logger = ctx.Logger()

	return nil
}

// implements ext.Trigger.Start
func (this *SSEServer) Start() error {

	logger.Info("Start")
	handlers := this.handlers

	logger.Info("Processing handlers")

	connection, exist := handlers[0].Settings()[cConnection]
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
			setting, err := data.CoerceToObject(v)
			if nil != err {
				continue
			}

			if nil != setting {
				var err error
				if setting["name"] == sseserver.ServerPort {
					properties[sseserver.ServerPort], err = data.CoerceToString(setting["value"])
					if nil != err {
						return err
					}
				} else if setting["name"] == cConnectionName {
					serverId = setting["value"].(string)
				} else if setting["name"] == sseserver.ConnectionPath {
					properties[sseserver.ConnectionPath], err = data.CoerceToString(setting["value"])
					if nil != err {
						return err
					}
				} else if setting["name"] == sseserver.ConnectionTlsEnabled {
					properties[sseserver.ConnectionTlsEnabled], err = data.CoerceToBoolean(setting["value"])
					if nil != err {
						return err
					}
				} else if setting["name"] == sseserver.ConnectionUploadCRT {
					properties[sseserver.ConnectionUploadCRT], err = data.CoerceToBoolean(setting["value"])
					if nil != err {
						return err
					}
				} else if setting["name"] == sseserver.ConnectionTlsCRT {
					tlsCRT, err := data.CoerceToObject(setting["value"])
					if nil != err {
						properties[sseserver.ConnectionTlsCRT], _ = b64.StdEncoding.DecodeString(strings.Split(tlsCRT["content"].(string), ",")[1])
						properties[sseserver.ConnectionTlsCRTPath], _ = tlsCRT["filename"].(string)
					}
				} else if setting["name"] == sseserver.ConnectionTlsKey {
					tlsKey, err := data.CoerceToObject(setting["value"])
					if nil != err {
						properties[sseserver.ConnectionTlsKey], _ = b64.StdEncoding.DecodeString(strings.Split(tlsKey["content"].(string), ",")[1])
						properties[sseserver.ConnectionTlsKeyPath], _ = tlsKey["filename"].(string)
					}
				} else if setting["name"] == sseserver.ConnectionTlsCRTPath {
					if nil != err {
						properties[sseserver.ConnectionTlsCRTPath], err = data.CoerceToString(setting["value"])
					}
				} else if setting["name"] == sseserver.ConnectionTlsKeyPath {
					if nil != err {
						properties[sseserver.ConnectionTlsKeyPath], err = data.CoerceToString(setting["value"])
					}
				}
			}

		}
		logger.Debug(properties)

		var err error
		this.server, err = sseserver.GetFactory().CreateServer(serverId, properties, this)
		if nil != err {
			return err
		}
		logger.Info("Server = ", *this.server)
		go this.server.Start()
	}

	return nil
}

// implements ext.Trigger.Stop
func (this *SSEServer) Stop() error {
	this.server.Stop()
	return nil
}

func (this *SSEServer) ProcessRequest(request string) error {
	this.mux.Lock()
	defer this.mux.Unlock()
	logger.Debug("Got SSE Request : ", request)
	outputData := &Output{}
	outputData.Request = request
	logger.Debug("Send SSE Request out : ", outputData)

	_, err := this.handlers[0].Handle(context.Background(), outputData)
	if nil != err {
		logger.Info("Error -> ", err)
	}

	return err
}
