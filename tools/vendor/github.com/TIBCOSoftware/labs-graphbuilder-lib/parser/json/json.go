/*
 * Copyright © 2020. TIBCO Software Inc.
 * This file is subject to the license terms contained
 * in the license file that is distributed with this file.
 */
package json

import (
	"bytes"
	"encoding/json"
	"fmt"
	"strconv"
	"strings"

	"github.com/TIBCOSoftware/flogo-lib/logger"
	"github.com/TIBCOSoftware/labs-graphbuilder-lib/util"
)

var log = logger.GetLogger("graphbuilder-json")

type Attribute struct {
	name          string
	dValue        interface{}
	dataType      string
	optional      bool
	multiInstance bool
}

func (this *Attribute) SetName(name string) {
	this.name = name
}

func (this *Attribute) GetName() string {
	return this.name
}

func (this *Attribute) SetDValue(dValue string) {
	this.dValue = dValue
}

func (this *Attribute) GetDValue() interface{} {
	return this.dValue
}

func (this *Attribute) SetType(dataType string) {
	this.dataType = dataType
}

func (this *Attribute) GetType() string {
	return this.dataType
}

func (this *Attribute) SetOptional(optional bool) {
	this.optional = optional
}

func (this *Attribute) IsOptional() bool {
	return this.optional
}

func (this *Attribute) SetMultiInstance(multiInstance bool) {
	this.multiInstance = multiInstance
}

func (this *Attribute) IsMultiInstance() bool {
	return this.multiInstance
}

type JSONParser struct {
	arrtibuteMap            map[string]*Attribute
	mandatoryAttrs          map[string]bool
	multiInstancePathPrefix string
	dateFromatString        string
	inScope                 bool
}

func NewJSONParser(arrtibuteMap map[string]*Attribute, mandatoryAttrs map[string]bool, multiInstancePathPrefix string) *JSONParser {
	return &JSONParser{
		arrtibuteMap:            arrtibuteMap,
		mandatoryAttrs:          mandatoryAttrs,
		multiInstancePathPrefix: multiInstancePathPrefix,
	}
}

func (this *JSONParser) SetDateFromatString(dateFromatString string) {
	this.dateFromatString = dateFromatString
}

func (this *JSONParser) Parse(jsonString []byte) []map[string]interface{} {
	var jsonData interface{}
	json.Unmarshal(jsonString, &jsonData)

	//log.Info("(JSONParser::Parse) jsonString : ", jsonData)

	handler := GraphDataCollector{
		dateFromatString:        this.dateFromatString,
		arrtibuteMap:            this.arrtibuteMap,
		multiInstancePathPrefix: this.multiInstancePathPrefix,
		commonFields:            make(map[string]interface{}),
		multiInstanceFields:     make(map[interface{}]map[string]interface{}),
	}

	jsonWalker := NewJSONWalker(handler)
	jsonWalker.Start(jsonData)
	tupleArray := handler.GetData()

	if 0 == len(tupleArray) || 0 == len(tupleArray[0]) {
		return tupleArray
	}

	for index, tuple := range tupleArray {
		for _, attrDef := range this.arrtibuteMap {
			attrName := attrDef.GetName()
			if nil == tuple[attrName] {
				if !attrDef.IsOptional() {
					log.Info("Not valid data : missing -> ", attrName, ", Tuple : ", tuple)
					return nil
				} else if nil != attrDef.GetDValue() {
					tuple[attrName] = attrDef.GetDValue()
				}
			}
		}
		if index >= len(tupleArray)-1 {
			tuple["LastElement"] = true
		} else {
			tuple["LastElement"] = false
		}
	}

	log.Info("(JSONParser::Parse) tuples - ", tupleArray)

	return tupleArray
}

type JSONDataHandler interface {
	HandleAttributes(namespace AttributeId, attribute interface{}, dataType interface{})
	GetData() []map[string]interface{}
}

type GraphDataCollector struct {
	dateFromatString        string
	multiInstancePathPrefix string
	arrtibuteMap            map[string]*Attribute
	commonFields            map[string]interface{}
	multiInstanceFields     map[interface{}]map[string]interface{}
}

func (this GraphDataCollector) GetData() []map[string]interface{} {
	log.Debug("(GraphDataCollector::GetData) this.multiInstanceFields - ", this.multiInstanceFields, ", this.commonFields - ", this.commonFields)
	tupleArray := make([]map[string]interface{}, 0)
	if 0 != len(this.multiInstanceFields) {
		for _, tuple := range this.multiInstanceFields {
			for key, value := range this.commonFields {
				tuple[key] = value
			}
			tupleArray = append(tupleArray, tuple)
		}
	} else {
		if 0 != len(this.commonFields) {
			tupleArray = append(tupleArray, this.commonFields)
		}
	}
	return tupleArray
}

func (this GraphDataCollector) HandleAttributes(namespace AttributeId, attribute interface{}, dataType interface{}) {
	attributeIds := namespace.GetId()
	//fmt.Println("\n\n\nHandle : id = ", attributeIds, ", attribute = ", attribute, ", type = ", dataType)

	for _, attributeId := range attributeIds {
		var err error
		attributeDef := this.arrtibuteMap[attributeId]
		if nil != attributeDef {
			attributeName := attributeDef.GetName()
			attributeDataType := attributeDef.GetType()
			this.commonFields[attributeName], err = util.TypeConversion(attribute, attributeDataType, this.dateFromatString)
			//fmt.Println(
			//	"Name : ", attributeName,
			//	", Defined Type : ", attributeDataType,
			//	", Original Type : ", reflect.TypeOf(attribute).String(),
			//	", Value (before) : ", attribute,
			//	", converted Type : ", reflect.TypeOf(this.graphData[attributeName]).String(),
			//	", Value (after) : ", this.graphData[attributeName],
			//)

			if nil != err {
				fmt.Println("Data conversion error : ", err)
			}
			break
		} else {
			pos := strings.Index(attributeId, this.multiInstancePathPrefix)
			if 0 <= pos {
				attributeIdSuffix := string(attributeId[pos+len(this.multiInstancePathPrefix):])
				pos2 := strings.Index(attributeIdSuffix, "]")
				if 0 <= pos2 {
					newAttributeId := this.multiInstancePathPrefix + string(attributeIdSuffix[pos2:])
					index, _ := strconv.Atoi(string(attributeIdSuffix[0:pos2]))
					attributeDef = this.arrtibuteMap[newAttributeId]

					if nil != attributeDef {
						attributeName := attributeDef.GetName()
						attributeDataType := attributeDef.GetType()
						currentInstance := this.multiInstanceFields[index]
						if nil == currentInstance {
							currentInstance = make(map[string]interface{})
							this.multiInstanceFields[index] = currentInstance
						}
						currentInstance[attributeName], err = util.TypeConversion(attribute, attributeDataType, this.dateFromatString)
						break
					}
				}
			}
		}
	}
}

type Scope struct {
	index    int
	maxIndex int
	name     string
	array    bool
}

type AttributeId struct {
	namespace []Scope
	name      string
}

func (this *AttributeId) GetIndex() int {
	return this.namespace[len(this.namespace)-1].index
}

func (this *AttributeId) GetId() []string {
	ids := make([]string, 1)

	var buffer bytes.Buffer
	arrayElement := false
	for i := range this.namespace {
		if !arrayElement {
			if 0 != i {
				buffer.WriteString(".")
			}
			buffer.WriteString(this.namespace[i].name)
		} else {
			arrayElement = false
		}

		if this.namespace[i].array {
			buffer.WriteString("[")
			buffer.WriteString(strconv.Itoa(this.namespace[i].index))
			buffer.WriteString("]")
			arrayElement = true
		}
	}
	buffer.WriteString(".")
	buffer.WriteString(this.name)
	ids[0] = buffer.String()
	return ids
}

func (this *AttributeId) SetName(name string) {
	this.name = name
}

func (this *AttributeId) updateIndex(index int, maxIndex int) {
	//	fmt.Println("   Before updateIndex : ", this.namespace, ", index : ", index)
	this.namespace[len(this.namespace)-1].index = index
	this.namespace[len(this.namespace)-1].maxIndex = maxIndex
	//	fmt.Println("   After updateIndex : ", this.namespace, ", index : ", index)
}

func (this *AttributeId) enterScope(scopename string, isArray bool) {
	//	fmt.Println("Before enterScope : ", this.namespace) //, ", index : ", this.namespace[len(this.namespace)-1].index)
	this.namespace = append(this.namespace, Scope{name: scopename, array: isArray, index: -1, maxIndex: -1})
}

func (this *AttributeId) leaveScope(scopename string, isArray bool) {
	//	fmt.Println("Before leaveScope : ", this.namespace) //, ", index : ", this.namespace[len(this.namespace)-1].index)
	this.namespace = this.namespace[:len(this.namespace)-1]
}

type JSONWalker struct {
	AttributeId
	currentLevel int
	jsonHandler  JSONDataHandler
}

func NewJSONWalker(jsonHandler JSONDataHandler) JSONWalker {
	jsonWalker := JSONWalker{
		currentLevel: 0,
		jsonHandler:  jsonHandler}
	jsonWalker.AttributeId = AttributeId{
		namespace: make([]Scope, 0),
	}

	return jsonWalker
}

func (this *JSONWalker) Start(jsonData interface{}) {
	this.walk("root", jsonData)
}

func (this *JSONWalker) walk(name string, data interface{}) {

	switch dataType := data.(type) {
	case []interface{}:
		{
			this.startArray(name)
			dataArray := data.([]interface{})
			maxIndex := len(dataArray) - 1
			for index, subdata := range dataArray {
				this.updateIndex(index, maxIndex)
				this.walk(name, subdata)
			}
			this.updateIndex(-1, -1)
			this.endArray(name)
			break
		}
	case map[string]interface{}:
		{
			this.startObject(name)
			for subname, subdata := range data.(map[string]interface{}) {
				this.walk(subname, subdata)
			}
			this.endObject(name)
			break
		}
	default:
		{
			this.AttributeId.SetName(name)
			this.jsonHandler.HandleAttributes(this.AttributeId, data, dataType)
		}
	}
}

func (this *JSONWalker) startArray(name string) {
	this.AttributeId.enterScope(name, true)
	//	fmt.Println("Start Array -> ", name, ", ", this.namespace)
}

func (this *JSONWalker) endArray(name string) {
	//	fmt.Println("End Array -> ", name)
	this.AttributeId.leaveScope(name, true)
}

func (this *JSONWalker) startObject(name string) {
	this.AttributeId.enterScope(name, false)
	//	fmt.Println("Start Object -> ", name, ", ", this.namespace)
}

func (this *JSONWalker) endObject(name string) {
	//	fmt.Println("End Object -> ", name)
	this.AttributeId.leaveScope(name, false)
}
