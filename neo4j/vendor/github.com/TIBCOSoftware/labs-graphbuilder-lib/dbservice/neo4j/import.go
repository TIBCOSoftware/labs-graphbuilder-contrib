package neo4j

import (
	"bytes"
	"fmt"
	"strings"

	"github.com/TIBCOSoftware/labs-graphbuilder-lib/model"
	"github.com/TIBCOSoftware/labs-graphbuilder-lib/util"
)

const (
	timestampFormat      = "2006-01-02T15:04:05.500"
	EXT_PROCESS_NAME     = "neo4j"
	CONF_KEYWORD_TYPE    = "ENTITY_TYPE"
	CONF_KEYWORD_DEFAULT = "*"
)

type Neo4jImportCSV struct {
	callNeo4jForImport bool
	initialized        bool
	typeName           string

	database string

	executionFolder string
	outputFolder    string

	importCommand  string
	recycleCommand string
	fileWriter     string

	nodeFileNames map[string][]string
	edgeFileNames map[string][]string

	label                     map[string]string
	ignoredNodeAttributes     map[string]map[string]bool
	ignoredEdgeAttributes     map[string]map[string]bool
	nodeAttributeWrittenOrder map[string][]string
	edgeAttributeWrittenOrder map[string][]string
}

func (this *Neo4jImportCSV) writeGraph(graph model.Graph) {
	model := graph.GetModel()
	if !this.initialized {
		this.writeConfiguration(model)
		this.initialized = true
	}

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

func (this *Neo4jImportCSV) writeConfiguration(graphModel map[string]interface{}) {

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

}

func (this *Neo4jImportCSV) writeNode(nodeType string, node *model.Node, keyMap map[string][]string) {
	//	filename := node.GetType() + ".csv"
	var nodeString bytes.Buffer

	/* Write header */
	nodeString.WriteString(fmt.Sprintf(
		"%sId:ID(%s)",
		strings.ToLower(this.typeName),
		this.typeName,
	))

	for _, attrName := range this.nodeAttributeWrittenOrder[nodeType] {
		nodeString.WriteString(",")
		nodeString.WriteString(attrName)
	}

	labelDefinition := this.label[this.typeName]
	if "" == labelDefinition {
		labelDefinition = this.label[CONF_KEYWORD_DEFAULT]
	}

	if "" != labelDefinition {
		nodeString.WriteString(",")
		nodeString.WriteString(":LABEL")
	}

	nodeString.WriteString("\n")

	/* Write data */
	nodeString.WriteString(node.NodeId.ToString())
	attrMap := node.GetAttributes()
	for _, attrName := range this.nodeAttributeWrittenOrder[nodeType] {
		/*Label either defined as attribute with key=tgTypeAttributeName
		  or node type. */
		nodeString.WriteString(",")
		attr := attrMap[attrName]
		if nil != attr {
			nodeString.WriteString(convertToString(attr.GetType().String(), attrName, attr.GetValue()))
		}
	}

	if "" != labelDefinition {
		nodeString.WriteString(",")
		nodeString.WriteString(getLabel(this.typeName, labelDefinition, node))
	}

	nodeString.WriteString("\n")

	fmt.Println("A Node : ", nodeString.String())

	/*

		        	//Write header
		        	print(fileWriter, String.format("%sId:ID(%s)", nodeTypeName.toLowerCase(), nodeTypeName));
		    		for(String attrName : nodeAttributeNames) {
		    			fileWriter.write(',');
		            	print(fileWriter, attrName);
		    		}

		    		if(null!=labelDefinition) {
		    			fileWriter.write(',');
		    			print(fileWriter, ":LABEL");
		    		}

		    		fileWriter.write('\n');

		    		// Write data
		        	for(Node node : graph.getNodesByType(nodeTypeName)) {

		        		// Label either defined as attribute with key=tgTypeAttributeName
		        		// or node type.

		        		Object entityID = extractId(node.getPKey());
		        		print(fileWriter, entityID);
		        		for(String attrName : nodeAttributeNames) {
		        			fileWriter.write(',');
		        			print(fileWriter, convertToString(node.getAttribute(attrName)));
		        		}


		        		if(null!=labelDefinition) {
		        			fileWriter.write(',');
		        			print(fileWriter, getLabel(labelDefinition, node));
		        		}

		        		fileWriter.write('\n');
		        		node.setExtraInfo(entityID);
					}

		        	filenames.add(filename);
		    		System.out.println("[Neo4jImport::writeNodesForSameType] CSV file was created successfully !!!");
	*/
}

func (this *Neo4jImportCSV) writeEdge(
	edgeType string,
	fromType string,
	toType string,
	edge *model.Edge,
	edgeDirectionMap map[string]int) {

	var edgeString bytes.Buffer

	/* Write header */
	edgeString.WriteString(fmt.Sprintf("%s:START_ID(%s)", "", fromType))

	attrMap := edge.GetAttributes()
	for _, attrName := range this.edgeAttributeWrittenOrder[edgeType] {
		edgeString.WriteString(",")
		edgeString.WriteString(attrName)
	}

	edgeString.WriteString(",")
	edgeString.WriteString(fmt.Sprintf("%s:TYPE", ""))

	edgeString.WriteString(",")
	edgeString.WriteString(fmt.Sprintf("%s:END_ID(%s)", "", toType))

	edgeString.WriteString("\n")

	/* Write data */

	// START_ID
	edgeString.WriteString(edge.GetFromId().ToString())

	// attributes
	for _, attrName := range this.edgeAttributeWrittenOrder[edgeType] {
		edgeString.WriteString(",")
		attr := attrMap[attrName]
		if nil != attr {
			edgeString.WriteString(convertToString(attr.GetType().String(), attrName, attr.GetValue()))
		}
	}

	// TYPE
	edgeString.WriteString(",")
	if "" != edgeType {
		edgeString.WriteString(edgeType)
	}

	// END_ID
	edgeString.WriteString(",")
	edgeString.WriteString(edge.GetToId().ToString())
	edgeString.WriteString("\n")

	fmt.Println("An Edge : ", edgeString.String())

	/*
		    	HashMap<String, ArrayList<String>> edgeAttributeNameMap = new HashMap<String, ArrayList<String>>();
		    	ArrayList<String> filenames = new ArrayList<String>();
		    	HashMap<String, FileWriter> filewriters = new HashMap<String, FileWriter>();
		        try {
					ArrayList<String> edgeAttributeNames = null;
		        	for(Node node : graph.getNodesByType(nodeTypeName)) {
		        		for(Edge edge : node.allEdges()) {
		    				String edgeType = edge.getType();
		        			if(null==(edgeAttributeNames=edgeAttributeNameMap.get(edgeType))) {
		        				edgeAttributeNames = new ArrayList<String>();
		        				edgeAttributeNameMap.put(edgeType, edgeAttributeNames);
		        				for(String name : graph.edgeAttributeNames(edge.getType())) {
		                    		HashSet<String> ignoredAttributes = ignoredEdgeAttributes.get(edgeType);
		                			if(null==ignoredAttributes||!ignoredAttributes.contains(name)) {
		            					edgeAttributeNames.add(name);
		                			}
		        				}
		        			}

		            		Object fromNodeID = null;
		            		Object toNodeID = null;
		                	int direction = edge.getDefinition().getDirection();
		                	switch(direction) {
		                		case 0 : break;
		                		case 1 :
		                		case 3 : fromNodeID = edge.getFromNode().getExtraInfo();
		                				 toNodeID = edge.getToNode().getExtraInfo();
		                				 break;
		                		case 2 :
		                		case 4 : toNodeID = edge.getFromNode().getExtraInfo();
		       					 		 fromNodeID = edge.getToNode().getExtraInfo();
		       					 		 break;
		                	}

		        			if(edge.getFromNode().equals(node)) {
		        				String toNodeType = edge.getToNode().getType();
		        				String writerKey = String.format("%s_%s", toNodeType, String.valueOf(edgeType.hashCode()));
		        				if(null==(fileWriter=filewriters.get(writerKey))) {
		        					String filename = String.format("%s_%s.csv", nodeTypeName, writerKey);
		        					filenames.add(filename);
		        					fileWriter = new FileWriter(String.format("%s/%s", importFileFolder, filename));
		        					filewriters.put(writerKey, fileWriter);

		        		        	// Write header
		        		        	print(fileWriter, String.format("%s:START_ID(%s)", "", (2==direction||4==direction)?toNodeType:nodeTypeName));
		        		    		for(String attrName : edgeAttributeNames) {
		        		    			fileWriter.write(',');
		        		            	print(fileWriter, attrName);
		        		    		}

		        		    		fileWriter.write(',');
		        		            print(fileWriter, String.format("%s:TYPE", ""));

		        		    		fileWriter.write(',');
		        		        	print(fileWriter, String.format("%s:END_ID(%s)", "", (2==direction||4==direction)?nodeTypeName:toNodeType));
		        		    		fileWriter.write('\n');
		        				}

		        				// Write data

		                		// START_ID
		                		print(fileWriter, fromNodeID);

		                		// attributes
		                		for(String attrName : edgeAttributeNames) {
		            				fileWriter.write(',');
		            				print(fileWriter, convertToString(edge.getAttribute(attrName)));
		                		}

		                		// TYPE
		            			fileWriter.write(',');
		                		if(null!=edgeType) {
		                        	print(fileWriter, edgeType);
		                		}

		                		// END_ID
		            			fileWriter.write(',');
		                		print(fileWriter, toNodeID);
		                		fileWriter.write('\n');
		        			}
		        		}
					}
	*/
}

func convertToString(attrDataType string, attrName string, attrValue interface{}) string {
	strValue, _ := util.ConvertToString(attrValue, timestampFormat)
	return strValue.(string)
}

func getLabel(typeName string, labelDef string, node *model.Node) string {
	var neoLabel string
	attributeNames := node.GetAttributes()
	typeAttr := attributeNames[typeName]
	if CONF_KEYWORD_TYPE == labelDef {
		if "" != typeName && nil != typeAttr {
			neoLabel = typeAttr.GetValue().(string)
		} else {
			neoLabel = node.GetType()
		}
	} else {
		neoLabel = typeAttr.GetValue().(string)
	}
	return neoLabel
}
