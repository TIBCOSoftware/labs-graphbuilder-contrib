/*
 * Copyright © 2020. TIBCO Software Inc.
 * This file is subject to the license terms contained
 * in the license file that is distributed with this file.
 */
package xml

import (
	"bytes"
	"encoding/xml"
	"fmt"
	"strings"
)

type XMLParser struct {
	arrtibuteMap map[string]string
	inScope      bool
}

func NewXMLParser(arrtibuteMap map[string]string) *XMLParser {
	jsonParser := &XMLParser{}
	jsonParser.arrtibuteMap = arrtibuteMap

	return jsonParser
}

func (this *XMLParser) Parse(xmlString string) map[string]interface{} {
	handler := GraphDataCollector{
		arrtibuteMap: this.arrtibuteMap,
		graphData:    make(map[string]interface{})}
	jsonWalker := NewXMLWalker(handler)

	jsonWalker.Start(xmlString)
	return handler.graphData
}

type XMLDataHandler interface {
	HandleAttributes(attributeId AttributeId, attribute interface{}, dataType interface{})
}

type GraphDataCollector struct {
	arrtibuteMap map[string]string
	graphData    map[string]interface{}
}

func (this GraphDataCollector) HandleAttributes(attributeId AttributeId, attribute interface{}, dataType interface{}) {
	//fmt.Println("Handle : id = ", attributeId.GetId(), ", attribute = ", attribute, ", type = ", dataType)
	attributeName := this.arrtibuteMap[attributeId.GetId()]
	if "" != attributeName {
		this.graphData[attributeName] = attribute
	}
}

type AttributeId struct {
	namespace []string
}

func (this *AttributeId) GetId() string {
	var buffer bytes.Buffer
	for index := range this.namespace {
		buffer.WriteString(this.namespace[index])
		if index < len(this.namespace)-1 {
			buffer.WriteString(".")
		}
	}
	return buffer.String()
}

func (this *AttributeId) enterScope(scopename string) {
	this.namespace = append(this.namespace, scopename)
}

func (this *AttributeId) leaveScope(scopename string) {
	this.namespace = this.namespace[:len(this.namespace)-1]
}

type XMLWalker struct {
	AttributeId
	currentLevel int
	jsonHandler  XMLDataHandler
}

func NewXMLWalker(jsonHandler XMLDataHandler) XMLWalker {
	jsonWalker := XMLWalker{
		currentLevel: 0,
		jsonHandler:  jsonHandler}
	jsonWalker.AttributeId = AttributeId{
		namespace: make([]string, 0)}

	return jsonWalker
}

func (this *XMLWalker) Start(xmlData string) {
	this.walk(xmlData)
}

func (this *XMLWalker) walk(data string) {

	decoder := xml.NewDecoder(strings.NewReader(data))
	var inElement string
	for {
		t, _ := decoder.Token()
		if t == nil {
			break
		}
		switch se := t.(type) {
		case xml.StartElement:
			{
				inElement = se.Name.Local
				this.startElement(inElement)

				/*


									                         if Element.Name.Local == "loc" {
					                                 fmt.Println("Element name is : ", Element.Name.Local)

					                                 err := decoder.DecodeElement(&l, &Element)
					                                 if err != nil {
					                                         fmt.Println(err)
					                                 }

					                                 fmt.Println("Element value is : ", l.Loc)
					                         }




				*/

			}
		case xml.EndElement:
			{
				this.endElement(se.Name.Local)
			}
		case xml.CharData:
			{
				this.jsonHandler.HandleAttributes(this.AttributeId, string([]byte(se)), "string")
			}
		default:
			{
				fmt.Println("4", t.([]byte))
			}
		}

	}
}

func (this *XMLWalker) startElement(name string) {
	this.AttributeId.enterScope(name)
	//fmt.Println("Start Object -> ", name, ", ", this.namespace)
}

func (this *XMLWalker) endElement(name string) {
	//fmt.Println("End Object -> ", name)
	this.AttributeId.leaveScope(name)
}
