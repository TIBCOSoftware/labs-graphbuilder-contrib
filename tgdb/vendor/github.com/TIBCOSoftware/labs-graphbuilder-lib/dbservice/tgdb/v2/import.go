package tgdb

import (
	"bytes"
	"fmt"

	"github.com/TIBCOSoftware/labs-graphbuilder-lib/model"
	"github.com/TIBCOSoftware/labs-graphbuilder-lib/util"
)

const (
	EXT_PROCESS_NAME   = "TGDB"
	timestampFormat    = "2006-01-02T15:04:05.500"
	undirectedEdgeCode = "1024"
	directedEdgeCode   = "1025"
	bidirectedEdgeCode = "1026"
)

func NewTGDBImportCSV() *TGDBImportCSV {
	return &TGDBImportCSV{
		callTGDBForImport: true,
		initialized:       false,
		/*
			# DEFAULT EDGETYPES:
			# undirected: 1024
			# directed: 1025
			# bidirected: 1026  */
		entityIDGen:         0,
		edgeTypeImportIdGen: 1040,

		fromNodes:                 make(map[string]bool),
		nodeAttributeWrittenOrder: make(map[string][]string),
		edgeAttributeWrittenOrder: make(map[string][]string),
		edgeTypeImportIds:         make(map[string]string),
	}
}

type TGDBImportCSV struct {
	callTGDBForImport bool
	initialized       bool

	importFileFolder string
	tgdbHome         string
	importCommand    string

	entityIDGen         int
	edgeTypeImportIdGen int

	fromNodes                 map[string]bool
	nodeAttributeWrittenOrder map[string][]string
	edgeAttributeWrittenOrder map[string][]string
	edgeTypeImportIds         map[string]string
}

func (this *TGDBImportCSV) WriteGraph(graph model.Graph) error {
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
	return nil
}

func (this *TGDBImportCSV) writeConfiguration(graphModel map[string]interface{}) {

	nodeModel := graphModel["nodes"].(map[string]interface{})
	edgeModel := graphModel["edges"].(map[string]interface{})
	nodeKeyMap := nodeModel["keyMap"].(map[string][]string)

	nodeAttrTypeMap := nodeModel["attrTypeMap"].(map[string](map[string]string))
	masterAttrTypeMap := make(map[string]string)
	for _, attrTypeMap := range nodeAttrTypeMap {
		for attrName, attrType := range attrTypeMap {
			if "" != masterAttrTypeMap[attrName] {
				if masterAttrTypeMap[attrName] != attrType {
					fmt.Println("Duplicate attribute : ", attrName)
				} else {
					continue
				}
			}
			masterAttrTypeMap[attrName] = attrType
		}
	}

	edgeAttrTypeMap := nodeModel["attrTypeMap"].(map[string](map[string]string))
	for _, attrTypeMap := range edgeAttrTypeMap {
		for attrName, attrType := range attrTypeMap {
			if "" != masterAttrTypeMap[attrName] {
				if masterAttrTypeMap[attrName] != attrType {
					fmt.Println("Duplicate attribute : ", attrName)
				} else {
					continue
				}
			}
			masterAttrTypeMap[attrName] = attrType
		}
	}

	var confString bytes.Buffer

	/* [attrtypes] */
	confString.WriteString("[attrtypes]\n")
	for attrName, attrType := range masterAttrTypeMap {
		confString.WriteString(fmt.Sprintf("%s = @type:%s\n", TrimWhiteSpace(attrName), attrType))
	}
	confString.WriteString("\n")

	/* [nodetypes] */
	confString.WriteString("[nodetypes]\n")
	for _, nodeType := range nodeModel["types"].([]string) {
		confString.WriteString(TrimWhiteSpace(nodeType))
		confString.WriteString(" = ")

		nodeAttributeNames := make([]string, len(nodeAttrTypeMap))
		this.nodeAttributeWrittenOrder[nodeType] = nodeAttributeNames

		index := 0
		for attrName, _ := range nodeAttrTypeMap {
			nodeAttributeNames[index] = attrName
			if index != 0 {
				confString.WriteString(",")
			} else {
				confString.WriteString("@attrs:")
			}
			confString.WriteString(TrimWhiteSpace(attrName))
		}
		if nil != nodeKeyMap[nodeType] {
			confString.WriteString(" @pkey:")
			for index, pkey := range nodeKeyMap[nodeType] {
				if index != 0 {
					confString.WriteString(",")
				}
				confString.WriteString(TrimWhiteSpace(pkey))
			}
		}
		confString.WriteString("\n")
	}
	confString.WriteString("\n")

	/* [edgetypes] */
	confString.WriteString("[edgetypes]\n")
	edgeVertexesMap := edgeModel["vertexes"].(map[string][]string)
	for _, edgeType := range edgeModel["types"].([]string) {

		/* @direction:undirected */
		edgeDirectionMap := edgeModel["directionMap"].(map[string]string)
		var edgeDirection string
		strDirection, _ := model.ToEdgeDirection(edgeDirectionMap[edgeType])
		switch strDirection {
		case model.Nondirectional:
			edgeDirection = "undirected"
		case model.Directional:
			edgeDirection = "directed"
		case model.Bidirectional:
			edgeDirection = "bidirected"
		default:
			edgeDirection = "directed"
		}
		confString.WriteString(TrimWhiteSpace(edgeType))
		confString.WriteString(" = @direction:")
		confString.WriteString(edgeDirection)
		confString.WriteString(" ")

		/* @fromnode:houseMemberType */
		edgeVertexes := edgeVertexesMap[edgeType]
		fromNodeType := edgeVertexes[0]
		toNodeType := edgeVertexes[1]
		if "*" != fromNodeType {
			confString.WriteString("@fromnode:")
			confString.WriteString(TrimWhiteSpace(fromNodeType))
			confString.WriteString(" ")
		}

		/* @tonode:houseMemberType */
		if "*" != toNodeType {
			confString.WriteString("@tonode:")
			confString.WriteString(TrimWhiteSpace(toNodeType))
			confString.WriteString(" ")
		}

		edgeAttributeName := make([]string, 0)
		this.edgeAttributeWrittenOrder[fromNodeType] = edgeAttributeName

		index := 0
		for attrName, _ := range edgeAttrTypeMap {
			edgeAttributeName = append(edgeAttributeName, attrName)
			edgeAttributeName[index] = attrName
			if index != 0 {
				confString.WriteString(",")
			} else {
				confString.WriteString("@attrs:")
			}
			confString.WriteString(TrimWhiteSpace(attrName))
		}

		confString.WriteString(" ")

		/* @importid:1042 */
		confString.WriteString("@importid:")
		edgeTypeImportId, _ := util.ConvertToString(this.edgeTypeImportIdGen, "")
		this.edgeTypeImportIds[edgeType] = edgeTypeImportId.(string)
		confString.WriteString(this.edgeTypeImportIds[edgeType])
		this.edgeTypeImportIdGen++

		confString.WriteString("\n")
	}
	confString.WriteString("\n")

	/* [import] */
	confString.WriteString("[import]\n")
	confString.WriteString("loadopts = upsert\n")
	for _, nodeType := range nodeModel["types"].([]string) {
		confString.WriteString(nodeType)
		confString.WriteString(" = ")

		for index, attrName := range this.nodeAttributeWrittenOrder[nodeType] {
			if index != 0 {
				confString.WriteString(",")
			} else {
				confString.WriteString("@attrs:")
			}
			confString.WriteString(TrimWhiteSpace(attrName))
		}
		confString.WriteString(" ")

		for index, attrName := range this.edgeAttributeWrittenOrder[nodeType] {
			if index != 0 {
				confString.WriteString(",")
			} else {
				confString.WriteString("@edgeattrs:")
			}
			confString.WriteString(TrimWhiteSpace(attrName))
		}
		confString.WriteString(" ")

		confString.WriteString("@files:")
		confString.WriteString(nodeType + ".csv")
		if this.fromNodes[nodeType] {
			confString.WriteString(",")
			confString.WriteString(nodeType + "$edges.csv")
		}
		confString.WriteString("\n")
	}

	fmt.Println("Configuration : ", confString.String())
}

func (this *TGDBImportCSV) writeNode(nodeType string, node *model.Node, keyMap map[string][]string) {
	//	filename := node.GetType() + ".csv"
	var nodeString bytes.Buffer

	this.entityIDGen += 1
	strID, _ := util.ConvertToString(this.entityIDGen, timestampFormat)

	nodeString.WriteString(strID.(string))

	attrMap := node.GetAttributes()
	for _, attrName := range this.nodeAttributeWrittenOrder[nodeType] {
		nodeString.WriteString(",")
		attr := attrMap[attrName]
		if nil != attr {
			nodeString.WriteString(convertToString(attr.GetType().String(), attrName, attr.GetValue()))
		}
	}
	nodeString.WriteString("\n")

	fmt.Println("A Node : ", nodeString.String())

	/*
		        try {
		        	String filename = nodeTypeName+".csv";
		        	fileWriter = new FileWriter(importFileFolder+"/"+filename);
		        	boolean headerPrinted = false;
		    		ArrayList<String> globalAttributeNames = nodeAttributeWrittenOrder.get(nodeTypeName);
		        	for(Node node : graph.getNodesByType(nodeTypeName)) {
		        		if(!headerPrinted) {
		        			List<Object> header = new ArrayList<Object>();
		        			header.add("entityID");
		        			header.addAll(globalAttributeNames);
		        			headerPrinted = true;
		        		}
		        		int entityID = ++entityIDGen;
		        		print(fileWriter, entityID);
		        		for(String attrName : globalAttributeNames) {
		        			Object attrValue = node.getAttribute(attrName);
		        			fileWriter.write(',');
		        			print(fileWriter, convertToString(node.getType(), attrName, attrValue));
		        		}
		        		fileWriter.write('\n');
		        		node.setExtraInfo(entityID);
					}
		//    		System.out.println("[TGDBImportCSV::writeNodesForSameType] CSV file was created successfully !!!");
		    		super.reportStatus("DataSink", "TGDBImportCSV", String.format("%s.csv was created successfully !!!", nodeTypeName));
		    		return filename;
		        } catch (Exception e) {
		        	e.printStackTrace();
			        throw e;
		        } finally {
		            try {
		                fileWriter.flush();
		                fileWriter.close();
		                //csvFilePrinter.close();
		            } catch (IOException e) {
		                System.out.println("[TGDBImportCSV::writeNodesForSameType] Error while flushing/closing fileWriter/csvPrinter !!!");
		//                e.printStackTrace();
		            }
		        }*/
}

func (this *TGDBImportCSV) writeEdge(
	edgeType string,
	fromType string,
	toType string,
	edge *model.Edge,
	edgeDirectionMap map[string]int) {

	var edgeString bytes.Buffer

	this.entityIDGen += 1
	strID, _ := util.ConvertToString(this.entityIDGen, timestampFormat)
	edgeString.WriteString(strID.(string))
	edgeString.WriteString(",")

	var edgeTypeImportId string
	strDirection, _ := model.ToEdgeDirection(edgeDirectionMap[edgeType])
	switch strDirection {
	case model.Nondirectional:
		edgeTypeImportId = undirectedEdgeCode
	case model.Directional:
		edgeTypeImportId = directedEdgeCode
	case model.Bidirectional:
		edgeTypeImportId = bidirectedEdgeCode
	default:
		edgeTypeImportId = directedEdgeCode
	}

	edgeString.WriteString(edgeTypeImportId)
	edgeString.WriteString(",")
	edgeString.WriteString(edge.GetFromId().ToString())
	edgeString.WriteString(",")
	edgeString.WriteString(edge.GetToId().ToString())

	attrMap := edge.GetAttributes()
	for _, attrName := range this.edgeAttributeWrittenOrder[edgeType] {
		edgeString.WriteString(",")
		attr := attrMap[attrName]
		if nil != attr {
			edgeString.WriteString(convertToString(attr.GetType().String(), attrName, attr.GetValue()))
		}
	}
	edgeString.WriteString("\n")
	fmt.Println("An Edge : ", edgeString.String())

	/*
	       				String edgeType = edge.getType();
	       				List<Object> processedEdgeTypeKey = new ArrayList<Object>();
	       				processedEdgeTypeKey.add(edgeType);
	       				processedEdgeTypeKey.add(edge.getFromNode().getType());
	       				processedEdgeTypeKey.add(edge.getToNode().getType());
	       				processedEdgeTypeKey.add(edge.getDefinition().getDirection());
	   					EdgeSchema edgeSchema = null;
	   					if(!processedEdges.contains(processedEdgeTypeKey)) {
	       					if(null==(edgeSchema=edgeDefinitions.get(edgeType))) {
	       						edgeSchema = new EdgeSchema(edgeType);
	           					edgeDefinitions.put(edgeType, edgeSchema);
	           					edgeTypeImportIds.put(edgeType, edgeTypeImportIdGen++);
	       					}
	       					edgeSchema.addEdgeDefinition(edge.getDefinition());
	       					processedEdges.add(processedEdgeTypeKey);
	       				}

	               		Object fromNodeID = edge.getFromNode().getExtraInfo();
	               		Object toNodeID = edge.getToNode().getExtraInfo();
	               		Integer edgeTypeImportId = edgeTypeImportIds.get(edgeType);
	               		if(null==edgeTypeImportId) {
	                   		int direction = edge.getDefinition().getDirection();
	                   		switch(direction) {
	                   			case 0 : edgeTypeImportId = undirectedEdgeCode ; break;
	                   			case 1 : edgeTypeImportId = directedEdgeCode ; break;
	                   			case 2 : edgeTypeImportId = directedEdgeCode ; break;
	                   			case 3 : edgeTypeImportId = bidirectedEdgeCode ; break;
	                   			case 4 : edgeTypeImportId = bidirectedEdgeCode ; break;
	                   		}
	               		}
	               		int entityID = ++entityIDGen;
	               		print(fileWriter, entityID);
	           			fileWriter.write(',');
	               		print(fileWriter, edgeTypeImportId);
	           			fileWriter.write(',');
	               		print(fileWriter, fromNodeID);
	           			fileWriter.write(',');
	               		print(fileWriter, toNodeID);
	               		for(String attrName : globalAttributeNames) {
	               			Object attrValue = edge.getAttribute(attrName);
	               			fileWriter.write(',');
	               			print(fileWriter, convertToString(edge.getType(), attrName, attrValue));
	               		}
	               		fileWriter.write('\n');	 */
}

func convertToString(attrDataType string, attrName string, attrValue interface{}) string {
	strValue, _ := util.ConvertToString(attrValue, timestampFormat)
	return strValue.(string)
}
