/*
 * Copyright © 2019. TIBCO Software Inc.
 * This file is subject to the license terms contained
 * in the license file that is distributed with this file.
 */
package builder

import (
	b64 "encoding/base64"
	"strings"
	"sync"

	"github.com/TIBCOSoftware/flogo-lib/core/activity"
	"github.com/TIBCOSoftware/flogo-lib/core/data"
	"github.com/TIBCOSoftware/flogo-lib/logger"
	"github.com/TIBCOSoftware/graph-builder-flogo/lib/model"
	"github.com/TIBCOSoftware/graph-builder-flogo/lib/util"
)

var log = logger.GetLogger("tibco-activity-graphbuilder")

var initialized bool = false

const (
	GraphModel = "GraphModel"
	Nodes      = "Nodes"
	Edges      = "Edges"
)

type BuilderActivity struct {
	metadata        *activity.Metadata
	inMemoryGraph   bool
	graphBuilder    *model.GraphBuilder
	activityToModel map[string]string
	models          map[string]*model.GraphDefinition
	mux             sync.Mutex
}

func NewActivity(metadata *activity.Metadata) activity.Activity {
	aBuilderActivity := &BuilderActivity{
		metadata:        metadata,
		inMemoryGraph:   false,
		graphBuilder:    model.NewGraphBuilder(),
		activityToModel: make(map[string]string),
		models:          make(map[string]*model.GraphDefinition),
	}

	return aBuilderActivity
}

func (a *BuilderActivity) Metadata() *activity.Metadata {
	return a.metadata
}

func (a *BuilderActivity) Eval(context activity.Context) (done bool, err error) {

	log.Info("[BuilderActivity:Eval] entering ........ ")

	graphModel, err := a.getGraphModel(context)

	if nil != err {
		return false, err
	}

	graph := a.graphBuilder.CreateGraph(graphModel.GetId(), graphModel)

	a.graphBuilder.BuildGraph(
		graph,
		graphModel,
		context.GetInput(Nodes).(*data.ComplexObject).Value,
		context.GetInput(Edges).(*data.ComplexObject).Value,
	)

	data := make(map[string]interface{})
	data["graph"] = a.graphBuilder.Export(graph, graphModel)

	//if a.inMemoryGraph {
	//	model.GetGraphManager().GetGraph(graphModel.GetId(), graphModel.GetId()).UpsertGraph(data)
	//}

	context.SetOutput("Graph", data)

	log.Info("[BuilderActivity:Eval] Exit ........ ")

	return true, nil
}

func (a *BuilderActivity) getGraphModel(context activity.Context) (*model.GraphDefinition, error) {
	var graphModel *model.GraphDefinition

	myId := util.ActivityId(context)
	graphModel = a.models[a.activityToModel[myId]]

	if nil == graphModel {
		a.mux.Lock()
		defer a.mux.Unlock()
		graphModel = a.models[a.activityToModel[myId]]
		if nil == graphModel {
			graphmodel, exist := context.GetSetting(GraphModel)
			if !exist {
				return nil, activity.NewError("GraphModel is not configured", "GRAPHBUILDER-4002", nil)
			}

			//Read graph model details
			connectionInfo, _ := data.CoerceToObject(graphmodel)
			if connectionInfo == nil {
				return nil, activity.NewError("Unable extract model", "GRAPHBUILDER-4001", nil)
			}

			var jsonmodel []byte
			var modelName string
			connectionSettings, _ := connectionInfo["settings"].([]interface{})
			if connectionSettings != nil {
				for _, v := range connectionSettings {
					setting, _ := data.CoerceToObject(v)

					if nil != setting {
						if setting["name"] == "model" {
							modelcontent, _ := data.CoerceToObject(setting["value"])
							jsonmodel, _ = b64.StdEncoding.DecodeString(strings.Split(modelcontent["content"].(string), ",")[1])
						} else if setting["name"] == "inMemory" {
							a.inMemoryGraph = setting["value"].(bool)
						} else if setting["name"] == "name" {
							modelName = setting["value"].(string)
						}
					}
				}
			}

			if "" == modelName {
				return nil, activity.NewError("Unable to get builder name", "GRAPHBUILDER-4003", nil)
			}

			if nil == jsonmodel {
				return nil, activity.NewError("Unable to get model string", "GRAPHBUILDER-4004", nil)
			}

			graphModel = model.NewGraphModel(modelName, string(jsonmodel))
			a.models[modelName] = graphModel
			a.activityToModel[myId] = modelName
		}
	}

	return graphModel, nil
}
