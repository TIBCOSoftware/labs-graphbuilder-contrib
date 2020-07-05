/*
 * Copyright © 2020. TIBCO Software Inc.
 * This file is subject to the license terms contained
 * in the license file that is distributed with this file.
 */
package csvfilewriter

import (
	"encoding/csv"
	"os"
	"strings"
	"sync"

	"github.com/TIBCOSoftware/flogo-lib/core/activity"
	"github.com/TIBCOSoftware/flogo-lib/core/data"
	"github.com/TIBCOSoftware/flogo-lib/logger"
	"github.com/TIBCOSoftware/labs-graphbuilder-lib/util"
)

// activityLogger is the default logger for the Filter Activity
var log = logger.GetLogger("activity-csvfilewriter")

const (
	input  = "CSVFields"
	output = "Data"
)

// CSVFileWriter is an Activity that is used to Filter a message to the console
type CSVFileWriter struct {
	metadata     *activity.Metadata
	mux          sync.Mutex
	workingDatas map[string]*CSVFileWriterWorkingData
}

// NewActivity creates a new AppActivity
func NewActivity(metadata *activity.Metadata) activity.Activity {
	aCSVParserActivity := &CSVFileWriter{
		metadata:     metadata,
		workingDatas: make(map[string]*CSVFileWriterWorkingData),
	}
	return aCSVParserActivity
}

// Metadata returns the activity's metadata
func (a *CSVFileWriter) Metadata() *activity.Metadata {
	return a.metadata
}

// Eval implements api.Activity.Eval - Filters the Message
func (a *CSVFileWriter) Eval(ctx activity.Context) (done bool, err error) {

	workingData, err := a.getWorkingData(ctx)

	if nil != err {
		return false, err
	}

	dataTuple := ctx.GetInput(input).(map[string]interface{})
	log.Info("input data = ", dataTuple)

	columns := make([]string, workingData.columnSize)
	index := 0
	for _, dataValue := range dataTuple {
		stringValue, err := data.CoerceToString(dataValue)
		if nil != err {
			log.Warn("Unable to coerce to string for ", dataValue)
		}
		columns[index] = stringValue
		index++
	}
	workingData.writer.Write(columns)

	return true, nil
}

type CSVFileWriterWorkingData struct {
	outputFile *os.File
	writer     *csv.Writer
	columnSize int
}

func (a *CSVFileWriter) getWorkingData(context activity.Context) (*CSVFileWriterWorkingData, error) {
	myId := util.ActivityId(context)
	workingData := a.workingDatas[myId]
	if nil == workingData {
		a.mux.Lock()
		defer a.mux.Unlock()
		workingData = a.workingDatas[myId]
		if nil == workingData {
			workingData := &CSVFileWriterWorkingData{}

			outputFilename, _ := context.GetSetting("CSVFilename")
			outputFilenameString := outputFilename.(string)
			result := strings.LastIndex(outputFilenameString, "/")
			outputFolderPathString := outputFilenameString[0:result]
			err := os.MkdirAll(outputFolderPathString, os.ModePerm)
			if nil != err {
				log.Info("Unable to create folder ...")
				return nil, err
			}

			fileExist := true
			_, err = os.Stat(outputFilenameString)
			if nil != err {
				if os.IsNotExist(err) {
					fileExist = false
				}
			}

			if fileExist {
				err = os.Remove(outputFilenameString)
				if err != nil {
					log.Error(err)
				}
			}

			workingData.outputFile, err = os.Create(outputFilenameString)
			if nil != err {
				log.Info("Unable to create file ...")
				return nil, err
			}

			outputFieldnames, _ := context.GetSetting("CSVFieldnames")
			outputFieldnameList := outputFieldnames.([]interface{})
			workingData.columnSize = len(outputFieldnameList)

			workingData.writer = csv.NewWriter(workingData.outputFile)
			printHeader, _ := context.GetSetting("WithHeader")
			if printHeader.(bool) {
				columnNames := make([]string, workingData.columnSize)
				for index, outputFieldname := range outputFieldnameList {
					outputFieldnameInfo := outputFieldname.(map[string]interface{})
					columnNames[index] = outputFieldnameInfo["Name"].(string)
				}
				workingData.writer.Write(columnNames)
			}
			a.workingDatas[myId] = workingData
		}
	}

	return workingData, nil
}
