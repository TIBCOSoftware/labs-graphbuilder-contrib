/*
 * Copyright © 2020. TIBCO Software Inc.
 * This file is subject to the license terms contained
 * in the license file that is distributed with this file.
 */
package graphtofile

import (
	"encoding/json"
	"fmt"
	"os"
	"sync"

	"github.com/TIBCOSoftware/flogo-lib/core/activity"
	"github.com/TIBCOSoftware/flogo-lib/logger"
	"github.com/TIBCOSoftware/labs-graphbuilder-lib/util"
)

const (
	OUTPUT_FILE = "outputFile"
)

var log = logger.GetLogger("tibco-activity-graphtofile")

type TGDBUpsertActivity struct {
	metadata    *activity.Metadata
	outputFiles map[string]string
	mux         sync.Mutex
}

func NewActivity(metadata *activity.Metadata) activity.Activity {
	return &TGDBUpsertActivity{
		metadata:    metadata,
		outputFiles: make(map[string]string),
	}
}

func (a *TGDBUpsertActivity) Metadata() *activity.Metadata {
	return a.metadata
}

func (a *TGDBUpsertActivity) Eval(context activity.Context) (done bool, err error) {
	log.Info("Graph to file entering ......... ")

	graph, exists := context.GetInput("Graph").(map[string]interface{})["graph"].(map[string]interface{})
	if !exists {
		return false, fmt.Errorf("Unable to get data from input!!")
	}

	outputFile, err := a.getOutputFile(graph["id"].(string), context)
	if nil != err {
		return false, err
	}

	a.mux.Lock()
	defer a.mux.Unlock()

	log.Debug("File name : ", outputFile)
	f, err := os.OpenFile(outputFile, os.O_APPEND|os.O_WRONLY, 0600)
	if err != nil {
		panic(err)
	}

	defer f.Close()

	nodes := graph["nodes"].(map[string]interface{})
	for _, node := range nodes {
		nodeDetail := node.(map[string]interface{})
		jsonString, _ := json.Marshal(nodeDetail)
		log.Debug("node -> " + string(jsonString))
		f.WriteString(string(jsonString) + "\r\n")
	}

	edges := graph["edges"].(map[string]interface{})
	for _, edge := range edges {
		edgeDetail := edge.(map[string]interface{})
		jsonString, _ := json.Marshal(edgeDetail)
		log.Debug("edge -> " + string(jsonString))
		f.WriteString(string(jsonString) + "\r\n")
	}

	//Set Message ID in the output
	//context.SetOutput(ovMessageId, *response.MessageId)
	log.Info("Update graph to file exit ......... ")
	return true, nil
}

func (a *TGDBUpsertActivity) getOutputFile(graphId string, context activity.Context) (string, error) {

	myId := util.ActivityId(context)
	outputFile := a.outputFiles[myId]

	if "" == outputFile {
		a.mux.Lock()
		defer a.mux.Unlock()
		outputFile = a.outputFiles[myId]
		if "" == outputFile {
			log.Info("Initializing output folder start ...")

			outputFileSetup, _ := context.GetSetting(OUTPUT_FILE)

			outputFile = util.ReplaceParameter(outputFileSetup.(string), util.GRAPH_ID, graphId)
			outputFolderPath, _ := util.SplitFilename(outputFile)

			err := os.MkdirAll(outputFolderPath, os.ModePerm)
			if nil != err {
				log.Error("Unable to create folder ...")
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
					log.Error(err)
				}
			}

			_, err = os.Create(outputFile)
			if nil != err {
				log.Error("Unable to create file ...")
				return "", err
			}

			log.Info("Initializing TGDB Service end ...")
			a.outputFiles[myId] = outputFile
		}
	}

	return outputFile, nil
}
