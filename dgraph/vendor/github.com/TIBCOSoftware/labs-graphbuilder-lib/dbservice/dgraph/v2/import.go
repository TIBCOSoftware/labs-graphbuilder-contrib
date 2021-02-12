package v2

import (
	"bytes"
	"fmt"

	"github.com/TIBCOSoftware/labs-graphbuilder-lib/dbservice"
	"github.com/TIBCOSoftware/labs-graphbuilder-lib/model"
	"github.com/TIBCOSoftware/labs-graphbuilder-lib/util"
)

/*
	 	<_:firstnameMervyn> 	<value> "Mervyn" .
		<_:user0051a000001JbXj> <attr> <_:firstnameMervyn> .
		uuid:string @index(exact,term) .
		type:string @index(exact) .
		manager: uid @reverse .
		attr: uid @reverse @count .
		label: string @index(exact) .
		value: string @index(term) .
		region: uid @reverse @count .
		owns: uid @reverse @count .
		product: uid @reverse @count .
*/

const (
	timestampFormat = "2006-01-02T15:04:05.500"

	NODE_FORMAT   = "<_:%s> <%s> \"%s\" .\n"
	EDGE_FORMAT   = "<_:%s> <%s> <_:%s> %s .\n"
	SCHEMA_FORMAT = "%s:%s %s %s .\n"
)

func NewDgraphImportRDF(properties map[string]interface{}) (dbservice.ImportService, error) {
	dgraphService := &DgraphImportRDF{}

	if nil != properties["explicitType"] {
		dgraphService.explicitType = properties["explicitType"].(bool)
	} else {
		dgraphService.explicitType = false
	}

	if nil != properties["typeName"] {
		dgraphService.typeName = util.CastString(properties["typeName"])
	} else {
		dgraphService.typeName = ""
	}

	if nil != properties["addPrefixToAttr"] {
		dgraphService.addPrefixToAttr = properties["addPrefixToAttr"].(bool)
	} else {
		dgraphService.addPrefixToAttr = false
	}
	dgraphService.targetRegex = "[^A-Za-z0-9]"
	dgraphService.replacement = "_"

	if nil != properties["graphModel"] {
		err := (*dgraphService).WriteSchema(properties["schema"], properties["graphModel"].(map[string]interface{}))

		if nil != err {
			return nil, err
		}
	}

	return nil, nil
}

type DgraphImportRDF struct {
	initialized bool

	folderName               string
	fileId                   string
	compress                 bool
	currentBufferedLineCount int

	explicitType              bool
	typeName                  string
	addPrefixToAttr           bool
	targetRegex               string
	replacement               string
	doReplaceChar             bool
	nodeAttributeWrittenOrder map[string][]string
	edgeAttributeWrittenOrder map[string][]string
}

func (this *DgraphImportRDF) writeGraph(graph model.Graph) {
	model := graph.GetModel()
	//	if !this.initialized {
	//		this.writeSchema(nil, model)
	//		this.initialized = true
	//	}

	nodeModel := model["nodes"].(map[string]interface{})
	for _, nodeType := range nodeModel["types"].([]string) {
		for _, node := range graph.GetNodesByType(nodeType) {
			this.writeNode(
				node.GetType(),
				node,
				nodeModel["keyMap"].(map[string][]string),
			)
		}
	}

	edgeModel := model["edges"].(map[string]interface{})
	for _, edge := range graph.GetEdges() {
		this.writeEdge(
			edge.GetType(),
			graph.GetNodes()[*edge.GetFromId()].GetType(),
			graph.GetNodes()[*edge.GetToId()].GetType(),
			edge,
			edgeModel["directionMap"].(map[string]int),
		)
	}
}

func (this *DgraphImportRDF) WriteSchema(userSchema interface{}, graphModel map[string]interface{}) error {

	nodeModel := graphModel["nodes"].(map[string]interface{})
	edgeModel := graphModel["edges"].(map[string]interface{})
	nodeAttrTypeMap := nodeModel["attrTypeMap"].(map[string](map[string]string))
	edgeAttrTypeMap := edgeModel["attrTypeMap"].(map[string](map[string]string))

	for _, nodeType := range nodeModel["types"].([]string) {
		nodeAttributeNames := make([]string, len(nodeAttrTypeMap))
		this.nodeAttributeWrittenOrder[nodeType] = nodeAttributeNames

		index := 0
		for attrName, _ := range nodeAttrTypeMap {
			nodeAttributeNames[index] = attrName
		}
	}

	for _, edgeType := range edgeModel["types"].([]string) {
		edgeAttributeName := make([]string, len(edgeAttrTypeMap))
		this.edgeAttributeWrittenOrder[edgeType] = edgeAttributeName

		index := 0
		for attrName, _ := range edgeAttrTypeMap {
			edgeAttributeName = append(edgeAttributeName, attrName)
			edgeAttributeName[index] = attrName
		}
	}

	schema := buildSchema(
		this.explicitType,
		this.typeName,
		this.targetRegex,
		this.replacement,
		this.addPrefixToAttr,
		userSchema,
		graphModel,
	)

	fmt.Println("Schema string : ", schema)

	return nil
}

func (this *DgraphImportRDF) writeNode(nodeType string, node *model.Node, keyMap map[string][]string) {
	var nodeString bytes.Buffer
	attrMap := node.GetAttributes()
	if "" != this.typeName && nil == attrMap[this.typeName] {
		nodeString.WriteString(fmt.Sprintf(
			NODE_FORMAT,
			replaceCharacter(node.NodeId.GetKeyHash()),
			this.typeName,
			nodeType))
	}

	for _, attrName := range this.nodeAttributeWrittenOrder[nodeType] {
		nodeString.WriteString(",")
		attr := attrMap[attrName]
		if nil != attr {
			attrValue := convertToString(attr.GetType().String(), attrName, attr.GetValue())
			nodeString.WriteString(fmt.Sprintf(
				NODE_FORMAT,
				replaceCharacter(node.NodeId.GetKeyHash()),
				attrName,
				attrValue))
		}
	}

	fmt.Println("A Node : ", nodeString.String())
}

func (this *DgraphImportRDF) writeEdge(
	edgeType string,
	fromType string,
	toType string,
	edge *model.Edge,
	edgeDirectionMap map[string]int) {

	fromNodeId := edge.GetFromId().ToString()
	toNodeId := edge.GetToId().ToString()
	attrMap := edge.GetAttributes()

	// Have to remove this special "relation" attribute
	var relation string
	if nil == attrMap["relation"] {
		relation, _ = attrMap["relation"].GetValue().(string)
	} else {
		relation = edgeType
	}
	//	boolean isBidirectional = edge.isBidirectinal();

	var facetsString bytes.Buffer
	var numOfFacets int
	for _, attrName := range this.edgeAttributeWrittenOrder[edgeType] {
		if "relation" != attrName {
			if 0 != numOfFacets {
				facetsString.WriteString(", ")
			}
			attr := attrMap[attrName]
			if nil != attr {
				facetsString.WriteString(attrName)
				facetsString.WriteString("=\"")
				facetsString.WriteString(convertToString(attr.GetType().String(), attrName, attr.GetValue()))
				facetsString.WriteString("\"")
			}
			numOfFacets++
		}
	}

	var facets string
	if 0 < numOfFacets {
		facets = fmt.Sprintf("(%s)", facetsString.String())
	} else {
		facets = facetsString.String()
	}

	edgeString := fmt.Sprintf(EDGE_FORMAT, fromNodeId, relation, toNodeId, facets)

	fmt.Println("An Edge : ", edgeString)
}

func convertToString(attrDataType string, attrName string, attrValue interface{}) string {
	strValue, _ := util.ConvertToString(attrValue, timestampFormat)
	return strValue.(string)
}

func replaceCharacter(data string) string {
	return data
}
