/*
 * Copyright © 2019. TIBCO Software Inc.
 * This file is subject to the license terms contained
 * in the license file that is distributed with this file.
 */
package filesreader

import (
	"context"
	"fmt"
	"reflect"

	"github.com/TIBCOSoftware/flogo-lib/core/trigger"
	"github.com/TIBCOSoftware/flogo-lib/logger"
	"github.com/TIBCOSoftware/graph-builder-flogo/lib/file"
)

var log = logger.GetLogger("trigger-files")

const (
	filename            = "Filename"
	emitPerLine         = "EmitPerLine"
	asynchronous        = "Asynchronous"
	output_FileContent  = "FileContent"
	output_ModifiedTime = "ModifiedTime"
	output_LineNumber   = "LineNumber"
)

//-============================================-//
//   Entry point create a new Trigger factory
//-============================================-//

func NewFactory(md *trigger.Metadata) trigger.Factory {
	return &FilesReaderFactory{metadata: md}
}

//-===============================-//
//     Define Trigger Factory
//-===============================-//

type FilesReaderFactory struct {
	metadata *trigger.Metadata
}

func (t *FilesReaderFactory) New(config *trigger.Config) trigger.Trigger {
	return &FilesReader{metadata: t.metadata, config: config}
}

//-=========================-//
//      Define Trigger
//-=========================-//

type FilesReader struct {
	metadata *trigger.Metadata
	config   *trigger.Config

	contentHandlers []*FileContentHandler
}

// implements trigger.Trigger.Metadata (trigger.go)
func (t *FilesReader) Metadata() *trigger.Metadata {
	return t.metadata
}

// implements trigger.Initializable.Initialize
func (t *FilesReader) Initialize(ctx trigger.InitContext) error {
	triggerHandlers := ctx.GetHandlers()
	t.contentHandlers = make([]*FileContentHandler, len(triggerHandlers))
	for index, handler := range triggerHandlers {
		t.contentHandlers[index] = &FileContentHandler{handler: handler}
	}
	return nil
}

// implements ext.Trigger.Start
func (t *FilesReader) Start() error {

	log.Info("FilesReader Start")

	for _, contentHandler := range t.contentHandlers {
		filename := contentHandler.handler.GetStringSetting(filename)
		log.Info("filename = ", filename)

		emitPerLine, _ := contentHandler.handler.GetSetting(emitPerLine)

		if "" == filename {
			return fmt.Errorf("Filename not set yet!")
		}

		reader, err := file.NewFileWatcher(filename, emitPerLine.(bool))
		if err != nil {
			log.Error("File reading error", err)
			return err
		}

		log.Info("reader = ", reader)

		asynchronous, _ := contentHandler.handler.GetSetting(asynchronous)

		if asynchronous.(bool) {
			go reader.Start(contentHandler)
		} else {
			err = reader.Start(contentHandler)
		}

		if nil != err {
			return err
		}
	}
	return nil
}

// implements ext.Trigger.Stop
func (t *FilesReader) Stop() error {

	log.Debug("Stopping endpoints")

	return nil
}

type FileContentHandler struct {
	handler *trigger.Handler
}

func (this *FileContentHandler) HandleContent(content string, modifiedTime int64, lineNumber int) {
	outputData := make(map[string]interface{})
	outputData[output_FileContent] = content
	outputData[output_ModifiedTime] = modifiedTime
	outputData[output_LineNumber] = lineNumber

	log.Info("(HandleContent) - Trigger start for content lineNumber : ", lineNumber)
	log.Info("(HandleContent) - modifiedTime : ", modifiedTime, ", type : ", reflect.TypeOf(modifiedTime).String())
	log.Debug("(HandleContent) - content : ", content)

	_, err := this.handler.Handle(context.Background(), outputData)

	if nil != err {
		log.Errorf("Run action for handler [%s] failed for reason [%s] message lost", this.handler, err)
	}
	log.Infof("(HandleContent) - Trigger done for content lineNumber : %d \n\n", lineNumber)
}
