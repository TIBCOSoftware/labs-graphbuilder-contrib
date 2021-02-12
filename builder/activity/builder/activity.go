/*
 * Copyright © 2020. TIBCO Software Inc.
 * This file is subject to the license terms contained
 * in the license file that is distributed with this file.
 */
package builder

import (
	//b64 "encoding/base64"
	"fmt"
	"sync"

	"github.com/TIBCOSoftware/flogo-lib/core/activity"
	"github.com/TIBCOSoftware/flogo-lib/core/data"
	"github.com/TIBCOSoftware/flogo-lib/logger"
	"github.com/TIBCOSoftware/labs-graphbuilder-lib/model"
	"github.com/TIBCOSoftware/labs-graphbuilder-lib/util"
)

var log = logger.GetLogger("tibco-activity-graphbuilder")

var initialized bool = false

const (
	GraphModel         = "GraphModel"
	AllowNullKey       = "AllowNullKey"
	Nodes              = "Nodes"
	Edges              = "Edges"
	PassThroughData    = "PassThroughData"
	BatchEnd           = "BatchEnd"
	PassThroughDataOut = "PassThroughDataOut"
)

type BuilderActivity struct {
	metadata           *activity.Metadata
	inMemoryGraph      bool
	graphBuilder       *model.GraphBuilder
	activityToModel    map[string]string
	models             map[string]*model.GraphDefinition
	activeGraphs       map[string]*model.Graph
	passThroughDataDef map[string]map[string]*Field
	mux                sync.Mutex
}

func NewActivity(metadata *activity.Metadata) activity.Activity {
	aBuilderActivity := &BuilderActivity{
		metadata:           metadata,
		inMemoryGraph:      false,
		graphBuilder:       model.NewGraphBuilder(),
		activityToModel:    make(map[string]string),
		models:             make(map[string]*model.GraphDefinition),
		activeGraphs:       make(map[string]*model.Graph),
		passThroughDataDef: make(map[string]map[string]*Field),
	}

	return aBuilderActivity
}

func (a *BuilderActivity) Metadata() *activity.Metadata {
	return a.metadata
}

func (a *BuilderActivity) Eval(context activity.Context) (done bool, err error) {

	log.Info("[BuilderActivity:Eval] entering ........ ")

	tempGraph, graphModel, passThroughDataDef, err := a.getGraphModel(context)

	if nil != err {
		return false, err
	}

	allowNullKey, exists := context.GetSetting(AllowNullKey)
	if !exists {
		allowNullKey = false
	}

	log.Debug("[BuilderActivity:Eval] BatchEnd : ", context.GetInput(BatchEnd))
	log.Info("[BuilderActivity:Eval] Nodes : ", context.GetInput(Nodes).(*data.ComplexObject).Value)
	log.Info("[BuilderActivity:Eval] Edges : ", context.GetInput(Edges).(*data.ComplexObject).Value)

	graphId := graphModel.GetId()
	deltaGraph := a.graphBuilder.CreateGraph(graphId, graphModel)
	err = a.graphBuilder.BuildGraph(
		&deltaGraph,
		graphModel,
		context.GetInput(Nodes).(*data.ComplexObject).Value,
		context.GetInput(Edges).(*data.ComplexObject).Value,
		allowNullKey.(bool),
	)

	if nil != err {
		return false, err
	}

	if a.inMemoryGraph {
		theGraph := model.GetGraphManager().GetGraph(model.GRAPH, graphId, graphId).(*model.Graph)
		(*theGraph).Merge(deltaGraph)
	}

	(*tempGraph).Merge(deltaGraph)
	if nil == context.GetInput(BatchEnd) || context.GetInput(BatchEnd).(bool) {
		graphData := make(map[string]interface{})
		graphData["graph"] = a.graphBuilder.Export(tempGraph, graphModel)

		log.Debug("[BuilderActivity:Eval] Graph : ", graphData)

		context.SetOutput("Graph", graphData)

		if 0 != len(passThroughDataDef) {
			log.Info("[BuilderActivity:Eval] PassThroughData : ", context.GetInput(PassThroughData).(*data.ComplexObject).Value)
			passThroughData := context.GetInput(PassThroughData).(*data.ComplexObject).Value.(map[string]interface{})
			passThroughDataOut := make(map[string]interface{})
			for name, attrDef := range passThroughDataDef {
				value := passThroughData[name]
				defaultDV := attrDef.GetDValue()
				log.Debug("[BuilderActivity:Eval] name : ", name, ", value : ", value, ", default : ", defaultDV, ", optional : ", attrDef.IsOptional())
				if nil == value && !attrDef.IsOptional() {
					if nil != defaultDV {
						value = defaultDV
					} else {
						return false, fmt.Errorf("Data (%s)  should not be nil!", name)
					}

				}
				passThroughDataOut[name] = value
			}
			context.SetOutput(PassThroughDataOut, passThroughDataOut)
		}
		/* clear graph data */
		(*tempGraph).Clear()
	}

	log.Info("[BuilderActivity:Eval] Exit ........ ")

	return true, nil
}

func (a *BuilderActivity) getGraphModel(context activity.Context) (*model.Graph, *model.GraphDefinition, map[string]*Field, error) {
	var graphModel *model.GraphDefinition
	var passThroughData map[string]*Field

	myId := util.ActivityId(context)
	graphModel = a.models[a.activityToModel[myId]]

	if nil == graphModel {
		a.mux.Lock()
		defer a.mux.Unlock()
		graphModel = a.models[a.activityToModel[myId]]
		if nil == graphModel {
			graphmodel, exist := context.GetSetting(GraphModel)
			if !exist {
				return nil, nil, nil, activity.NewError("GraphModel is not configured", "GRAPHBUILDER-4002", nil)
			}

			//Read graph model details
			connectionInfo, _ := data.CoerceToObject(graphmodel)
			if connectionInfo == nil {
				return nil, nil, nil, activity.NewError("Unable extract model", "GRAPHBUILDER-4001", nil)
			}

			var jsonmodel string
			var modelName string
			connectionSettings, _ := connectionInfo["settings"].([]interface{})
			if connectionSettings != nil {
				for _, v := range connectionSettings {
					setting, _ := data.CoerceToObject(v)
					if nil != setting {
						if setting["name"] == "model" {
							//modelcontent, _ := data.CoerceToObject(setting["value"])
							//jsonmodel, _ = b64.StdEncoding.DecodeString(strings.Split(modelcontent["content"].(string), ",")[1])
						} else if setting["name"] == "metadata" {
							jsonmodel = setting["value"].(string)
						} else if setting["name"] == "inMemory" {
							a.inMemoryGraph = setting["value"].(bool)
						} else if setting["name"] == "name" {
							modelName = setting["value"].(string)
						}
					}
				}
			}

			if "" == modelName {
				return nil, nil, nil, activity.NewError("Unable to get builder name", "GRAPHBUILDER-4003", nil)
			}

			if "" == jsonmodel {
				return nil, nil, nil, activity.NewError("Unable to get model string", "GRAPHBUILDER-4004", nil)
			}

			log.Debug("Model = ", jsonmodel)

			var err error
			graphModel, err = model.NewGraphModel(modelName, jsonmodel)
			if nil != err {
				return nil, nil, nil, err
			}

			a.models[modelName] = graphModel
			a.activityToModel[myId] = modelName

			/* create once */
			if nil == a.activeGraphs[myId] {
				graph := a.graphBuilder.CreateGraph(graphModel.GetId(), graphModel)
				a.activeGraphs[myId] = &graph
			}

			passThroughData = a.buildPassThroughData(myId, context)
			a.passThroughDataDef[myId] = passThroughData
		}
	}

	return a.activeGraphs[myId], graphModel, passThroughData, nil
}

func (a *BuilderActivity) buildPassThroughData(myId string, context activity.Context) map[string]*Field {
	passThroughData := make(map[string]*Field)
	passThroughFieldnames, _ := context.GetSetting("PassThrough")
	log.Info("Processing handlers : PassThroughData = ", passThroughFieldnames)

	for _, passThroughFieldname := range passThroughFieldnames.([]interface{}) {
		passThroughFieldnameInfo := passThroughFieldname.(map[string]interface{})
		attribute := &Field{}
		attribute.SetName(passThroughFieldnameInfo["FieldName"].(string))
		attribute.SetType(passThroughFieldnameInfo["Type"].(string))
		attribute.SetOptional(nil != passThroughFieldnameInfo["Optional"] && "no" == passThroughFieldnameInfo["Optional"].(string))
		if nil != passThroughFieldnameInfo["Default"] && "" != passThroughFieldnameInfo["Default"].(string) {
			//attribute.SetDValue()
		}
		passThroughData[attribute.GetName()] = attribute
	}
	return passThroughData
}

type Field struct {
	name     string
	dValue   interface{}
	dataType string
	optional bool
}

func (this *Field) SetName(name string) {
	this.name = name
}

func (this *Field) GetName() string {
	return this.name
}

func (this *Field) SetDValue(dValue string) {
	this.dValue = dValue
}

func (this *Field) GetDValue() interface{} {
	return this.dValue
}

func (this *Field) SetType(dataType string) {
	this.dataType = dataType
}

func (this *Field) GetType() string {
	return this.dataType
}

func (this *Field) SetOptional(optional bool) {
	this.optional = optional
}

func (this *Field) IsOptional() bool {
	return this.optional
}
