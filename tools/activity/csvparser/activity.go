/*
 * Copyright © 2019. TIBCO Software Inc.
 * This file is subject to the license terms contained
 * in the license file that is distributed with this file.
 */
package csvparser

import (
	"encoding/csv"
	"io"
	"sync"

	"strconv"
	"strings"

	"github.com/TIBCOSoftware/flogo-lib/core/activity"
	"github.com/TIBCOSoftware/flogo-lib/core/data"
	"github.com/TIBCOSoftware/flogo-lib/logger"
	"github.com/TIBCOSoftware/labs-graphbuilder-lib/util"
)

// activityLogger is the default logger for the Filter Activity
var log = logger.GetLogger("activity-csvparser")

const (
	setting_OutputFieldnames = "OutputFieldnames"
	setting_DateFormat       = "DateFormat"
	setting_FirstRowIsHeader = "FirstRowIsHeader"
	input_CSVString          = "CSVString"
	input_SequenceNumber     = "SequenceNumber"
	output                   = "Data"
)

// CSVParserActivity is an Activity that is used to Filter a message to the console
type CSVParserActivity struct {
	metadata     *activity.Metadata
	mux          sync.Mutex
	workingDatas map[string]*CSVParserWorkingData
}

// NewActivity creates a new AppActivity
func NewActivity(metadata *activity.Metadata) activity.Activity {
	aCSVParserActivity := &CSVParserActivity{
		metadata:     metadata,
		workingDatas: make(map[string]*CSVParserWorkingData),
	}
	return aCSVParserActivity
}

// Metadata returns the activity's metadata
func (a *CSVParserActivity) Metadata() *activity.Metadata {
	return a.metadata
}

// Eval implements api.Activity.Eval - Filters the Message
func (a *CSVParserActivity) Eval(ctx activity.Context) (done bool, err error) {

	CSVString := ctx.GetInput(input_CSVString).(string)
	workingData, err := a.getWorkingData(ctx)

	if workingData.firstRowIsHeader {
		sequenceNumber := ctx.GetInput(input_SequenceNumber).(int)
		if 1 == sequenceNumber {
			log.Info("Skip header !!")
			return false, nil
		}
	}

	reader := csv.NewReader(strings.NewReader(CSVString))
	csvdataArray := make([]map[string]interface{}, 0)
	for {
		line, err := reader.Read()
		if err == io.EOF {
			break
		} else if err != nil {
			log.Error(err)
		}

		csvdata := make(map[string]interface{})
		for i := 0; i < len(workingData.indexToFieldname); i++ {
			log.Debug("Processing handlers : indexToFieldname[i] = ", workingData.indexToFieldname[i], ", line[i] = ", line[i], "")
			csvdata[workingData.indexToFieldname[i]], err = util.StringToTypes(line[i], workingData.indexToFieldtype[i], workingData.dateFromatString)
			if nil != err {
				log.Error(err)
				//return false, err
			}
		}

		log.Info("(CSVParserActivity.Eval) - csvdata : ", csvdata)
		csvdataArray = append(csvdataArray, csvdata)
	}

	complexcsvdata := &data.ComplexObject{Metadata: "Data", Value: csvdataArray}
	ctx.SetOutput("Data", complexcsvdata)
	return true, nil
}

type CSVParserWorkingData struct {
	dateFromatString string
	firstRowIsHeader bool
	indexToFieldname map[int]string
	indexToFieldtype map[int]string
	initialized      bool
}

func (a *CSVParserActivity) getWorkingData(ctx activity.Context) (*CSVParserWorkingData, error) {
	myId := util.ActivityId(ctx)
	workingData := a.workingDatas[myId]
	log.Info("workingDatas : ", a.workingDatas, ", myId : ", myId)
	if nil == workingData {
		a.mux.Lock()
		defer a.mux.Unlock()
		workingData = a.workingDatas[myId]
		if nil == workingData {
			workingData := &CSVParserWorkingData{}

			dateFormat, _ := ctx.GetSetting(setting_DateFormat)
			if nil != dateFormat {
				workingData.dateFromatString = dateFormat.(string)
			}
			outputFieldnames, _ := ctx.GetSetting(setting_OutputFieldnames)
			firstRowIsHeader, _ := ctx.GetSetting(setting_FirstRowIsHeader)
			workingData.firstRowIsHeader = firstRowIsHeader.(bool)
			workingData.indexToFieldname = make(map[int]string)
			workingData.indexToFieldtype = make(map[int]string)
			for i, outputFieldname := range outputFieldnames.([]interface{}) {
				outputFieldnameInfo := outputFieldname.(map[string]interface{})
				log.Info(">>>>>>>>>>>> outputFieldnameInfo : ", outputFieldnameInfo)
				//i, _ := strconv.Atoi(outputFieldnameInfo["CSVFieldName"].(string))
				workingData.indexToFieldname[i] = outputFieldnameInfo["AttributeName"].(string)
				workingData.indexToFieldtype[i] = outputFieldnameInfo["Type"].(string)
				log.Info("Processing handlers : parameter.Column-1 = ", strconv.Itoa(i), ", parameter.Name = ", outputFieldnameInfo["AttributeName"])
			}
			a.workingDatas[myId] = workingData

			return workingData, nil
		}
	}

	return workingData, nil
}
