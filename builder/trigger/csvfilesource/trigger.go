/*
 * Copyright © 2019. TIBCO Software Inc.
 * This file is subject to the license terms contained
 * in the license file that is distributed with this file.
 */
package csvfilesource

import (
	"bufio"
	"context"
	"encoding/csv"
	"encoding/json"
	"io"
	"os"
	"strconv"
	"strings"

	"github.com/TIBCOSoftware/flogo-lib/core/data"
	"github.com/TIBCOSoftware/flogo-lib/core/trigger"
	"github.com/TIBCOSoftware/flogo-lib/logger"
	"github.com/carlescere/scheduler"
)

var log = logger.GetLogger("trigger-file")

type Parameter struct {
	ParameterName string `json:"parameterName"`
	Type          string `json:"type"`
	Repeating     string `json:"repeating"`
	Required      string `json:"required"`
	IsEditable    string `json:"isEditable"`
	Column        string `json:"Column"`
	Name          string `json:"Name"`
	CapType       string `json:"Type"`
}

//-============================================-//
//   Entry point create a new Trigger factory
//-============================================-//

func NewFactory(md *trigger.Metadata) trigger.Factory {
	return &CSVFileSourceFactory{metadata: md}
}

//-===============================-//
//     Define Trigger Factory
//-===============================-//

type CSVFileSourceFactory struct {
	metadata *trigger.Metadata
}

func (t *CSVFileSourceFactory) New(config *trigger.Config) trigger.Trigger {
	return &CSVFileSource{metadata: t.metadata, config: config, skipFirstLine: false}
}

//-=========================-//
//      Define Trigger
//-=========================-//

type CSVFileSource struct {
	metadata      *trigger.Metadata
	config        *trigger.Config
	timers        []*scheduler.Job
	skipFirstLine bool
	handlers      []*trigger.Handler
}

// implements trigger.Trigger.Metadata (trigger.go)
func (t *CSVFileSource) Metadata() *trigger.Metadata {
	return t.metadata
}

// implements trigger.Initializable.Initialize
func (t *CSVFileSource) Initialize(ctx trigger.InitContext) error {

	t.handlers = ctx.GetHandlers()
	return nil
}

// implements ext.Trigger.Start
func (t *CSVFileSource) Start() error {

	log.Debug("Start")
	handlers := t.handlers

	log.Debug("Processing handlers")

	for _, handler := range handlers {
		csvfilename := handler.GetStringSetting("CSVFilename")
		//		log.Info(strings.Join([]string{"Processing handlers : CSVFilename = ", csvfilename}, ""))

		outputFieldnames := handler.GetStringSetting("OutputFieldnames")
		//		log.Info(strings.Join([]string{"Processing handlers : outputFieldnames = ", outputFieldnames}, ""))

		log.Info("WithHeader : ", handler.GetStringSetting("WithHeader"))
		skipFirstLine, _ := handler.GetSetting("WithHeader")
		t.skipFirstLine = skipFirstLine.(bool)

		var parameters []Parameter
		json.Unmarshal([]byte(outputFieldnames), &parameters)

		indexToFieldname := make(map[int]string)
		for _, parameter := range parameters {
			i, _ := strconv.Atoi(parameter.Column)
			indexToFieldname[i-1] = parameter.Name
			log.Info(strings.Join([]string{"Processing handlers : parameter.Column-1 = ", strconv.Itoa(i - 1), ", parameter.Name = ", parameter.Name}, ""))
		}

		csvFile, _ := os.Open(csvfilename)
		reader := csv.NewReader(bufio.NewReader(csvFile))

		count := 0
		for {
			line, error := reader.Read()
			if error == io.EOF {
				break
			} else if error != nil {
				log.Error(error)
			}

			if t.skipFirstLine && 0 == count {
				count++
				continue
			} else {
				count++
			}

			output := make(map[string]interface{})
			output["Sourcename"] = handler.GetStringSetting("Sourcename")
			output["Queryname"] = handler.GetStringSetting("Queryname")

			csvdata := make(map[string]interface{})
			for i := 0; i < len(line); i++ {
				log.Debug(strings.Join([]string{"Processing handlers : indexToFieldname[i] = ", indexToFieldname[i], ", line[i] = ", line[i]}, ""))
				csvdata[indexToFieldname[i]] = line[i]
			}

			complexcsvdata := &data.ComplexObject{Metadata: "Data", Value: csvdata}
			attr, _ := data.NewAttribute(complexcsvdata.Metadata, data.TypeComplexObject, complexcsvdata)
			output["Data"] = attr
			log.Debug(output)

			_, err := handler.Handle(context.Background(), output)

			if err != nil {
				log.Errorf("Run action for handler [%s] failed for reason [%s] message lost", handler, err)
			}
		}
	}

	return nil
}

// implements ext.Trigger.Stop
func (t *CSVFileSource) Stop() error {

	log.Debug("Stopping endpoints")
	for _, timer := range t.timers {

		if timer.IsRunning() {
			//log.Debug("Stopping timer for : ", k)
			timer.Quit <- true
		}
	}

	t.timers = nil

	return nil
}
