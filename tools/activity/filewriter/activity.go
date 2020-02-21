/*
 * Copyright © 2019. TIBCO Software Inc.
 * This file is subject to the license terms contained
 * in the license file that is distributed with this file.
 */
package filewriter

import (
	"encoding/json"
	"os"
	"strings"
	"sync"

	"github.com/TIBCOSoftware/flogo-lib/core/activity"
	"github.com/TIBCOSoftware/flogo-lib/logger"
	"github.com/TIBCOSoftware/labs-graphbuilder-lib/util"
)

const (
	OUTPUT_FILE = "outputFile"
)

var log = logger.GetLogger("tibco-activity-filewriter")

type FileWriterActivity struct {
	metadata    *activity.Metadata
	outputFiles map[string]string
	mux         sync.Mutex
}

func NewActivity(metadata *activity.Metadata) activity.Activity {
	return &FileWriterActivity{
		metadata:    metadata,
		outputFiles: make(map[string]string),
	}
}

func (a *FileWriterActivity) Metadata() *activity.Metadata {
	return a.metadata
}

func (a *FileWriterActivity) Eval(context activity.Context) (done bool, err error) {
	log.Info("(Eval) Graph to file entering ......... ")

	data := context.GetInput("Data")

	outputFile, err := a.getOutputFile(context)
	if nil != err {
		return false, err
	}

	a.mux.Lock()
	defer a.mux.Unlock()

	log.Info("(Eval) File name : ", outputFile)
	f, err := os.OpenFile(outputFile, os.O_APPEND|os.O_WRONLY, 0600)
	if err != nil {
		panic(err)
	}

	defer f.Close()

	jsonString, _ := json.Marshal(data)
	log.Debug("(Eval) Data -> " + string(jsonString))

	f.WriteString(string(jsonString) + "\r\n")

	log.Info("(Eval) write object to file exit ......... ")
	return true, nil
}

func (a *FileWriterActivity) getOutputFile(context activity.Context) (string, error) {

	myId := util.ActivityId(context)
	outputFile := a.outputFiles[myId]
	log.Info("%%%%%%%%%%%%%%%%%%%%%% getOutputFile : myId = ", myId, ", outputFile = ", outputFile)

	if "" == outputFile {
		a.mux.Lock()
		defer a.mux.Unlock()
		outputFile = a.outputFiles[myId]
		if "" == outputFile {
			log.Info("Initializing output folder start ...")

			outputFileSetting, _ := context.GetSetting(OUTPUT_FILE)
			outputFile = outputFileSetting.(string)
			var outputFolder string
			var pos int
			if strings.LastIndex(outputFile, "/") > strings.LastIndex(outputFile, "\\") {
				pos = strings.LastIndex(outputFile, "/")
			} else {
				pos = strings.LastIndex(outputFile, "\\")
			}

			if pos > 0 {
				outputFolder = outputFile[:pos-1]
			}

			//			log.Info("Output file : ", outputFile)
			//			log.Info("Output folder : ", outputFolder)

			err := os.MkdirAll(outputFolder, os.ModePerm)
			if nil != err {
				log.Error("Unable to create folder : ", err)
				return "", err
			}

			fileExist := true
			_, err = os.Stat(outputFile)
			if nil != err {
				if os.IsNotExist(err) {
					fileExist = false
				}
			}

			if fileExist {
				err = os.Remove(outputFile)
				if err != nil {
					log.Error("Unable to remove file : ", err)
				}
			}

			_, err = os.Create(outputFile)
			if nil != err {
				log.Error("Unable to create file : ", err)
				return "", err
			}

			log.Info("Initializing FileWriter Service end ...")
			a.outputFiles[myId] = outputFile
		}
	}

	return outputFile, nil
}
