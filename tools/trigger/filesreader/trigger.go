/*
 * Copyright © 2020. TIBCO Software Inc.
 * This file is subject to the license terms contained
 * in the license file that is distributed with this file.
 */
package filesreader

import (
	"context"
	"errors"
	"sync"

	"github.com/TIBCOSoftware/labs-graphbuilder-lib/file"
	"github.com/project-flogo/core/data/metadata"
	"github.com/project-flogo/core/support/log"
	"github.com/project-flogo/core/trigger"
)

var triggerMd = trigger.NewMetadata(&Settings{}, &HandlerSettings{}, &Output{})

func init() {
	_ = trigger.Register(&Trigger{}, &Factory{})
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

	return &Trigger{settings: settings}, nil
}

//-=========================-//
//      Define Trigger
//-=========================-//

var logger log.Logger

type Trigger struct {
	settings *Settings
	handlers []trigger.Handler
	mux      sync.Mutex
}

// Init implements trigger.Init
func (this *Trigger) Initialize(ctx trigger.InitContext) error {

	this.handlers = ctx.GetHandlers()
	logger = ctx.Logger()

	return nil
}

// Start implements ext.Trigger.Start
func (this *Trigger) Start() error {

	for handlerId, handler := range this.handlers {

		logger.Info("Start handler : name =  ", handler.Name())

		handlerSetting := &HandlerSettings{}
		err := metadata.MapToStruct(handler.Settings(), handlerSetting, true)
		if err != nil {
			return err
		}

		if "" == handlerSetting.Filename {
			return errors.New("Filename not set yet!")
		}

		reader, err := file.NewFileWatcher(handlerId, handlerSetting.Filename, handlerSetting.EmitPerLine, handlerSetting.MaxNumberOfLine)
		if err != nil {
			logger.Error("File reading error", err)
			return err
		}

		logger.Info("reader = ", reader)

		if handlerSetting.Asynchronous {
			go reader.Start(this)
		} else {
			err = reader.Start(this)
		}

		if nil != err {
			return err
		}
	}

	return nil
}

// Stop implements ext.Trigger.Stop
func (t *Trigger) Stop() error {
	logger.Debug("Stopping endpoints")
	return nil
}

func (t *Trigger) HandleContent(handlerId int, id string, content string, modifiedTime int64, lineNumber int, endOfFile bool) {
	t.mux.Lock()
	defer t.mux.Unlock()
	outputData := &Output{}
	outputData.MessageId = id
	outputData.FileContent = content
	outputData.ModifiedTime = modifiedTime
	outputData.LineNumber = lineNumber
	outputData.EndOfFile = endOfFile

	logger.Info("(FileContentHandler.HandleContent) - Trigger sends MessageId : ", id, ", handlerId : ", handlerId, ", lineNumber : ", lineNumber, ", modifiedTime : ", modifiedTime)
	logger.Debug("(FileContentHandler.HandleContent) - content : ", content)

	_, err := t.handlers[handlerId].Handle(context.Background(), outputData)

	if nil != err {
		logger.Errorf("Run action for handler [%s] failed for reason [%s] message lost", t.handlers[handlerId], err)
	}
	logger.Infof("(FileContentHandler.HandleContent) - Trigger done for content lineNumber : %d ", lineNumber)
}
