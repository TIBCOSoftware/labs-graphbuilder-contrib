/*
 * Copyright © 2019. TIBCO Software Inc.
 * This file is subject to the license terms contained
 * in the license file that is distributed with this file.
 */
import {Observable} from "rxjs/Observable";
import {Injectable, Injector, Inject} from "@angular/core";
import {Http} from "@angular/http";
import {
    WiContrib,
    WiServiceHandlerContribution,
    IValidationResult,
    ValidationResult,
    IFieldDefinition,
    IActivityContribution,
    IConnectorContribution,
    WiContributionUtils
} from "wi-studio/app/contrib/wi-contrib";

@WiContrib({})
@Injectable()
export class GraphBuilderActivityContributionHandler extends WiServiceHandlerContribution {
	selectedConnector: string;
		
    constructor(@Inject(Injector) injector, private http: Http) {
        super(injector, http);
    }

    value = (fieldName: string, context: IActivityContribution): Observable<any> | any => {

		console.log('[GraphBuilder::value] Build field : ', fieldName);
		console.log(context);
		
        if (fieldName === "GraphModel") {
            let allowedConnectors = context.getField("GraphModel").allowed;	
			let selectedConnectorId = context.getField("GraphModel").value;
			for(let allowedConnector of allowedConnectors) {
				if(allowedConnector["unique_id"] === selectedConnectorId) {
					this.selectedConnector = allowedConnector["name"]
				}
			}
            
            return Observable.create(observer => {
            		//Connector Type must match with the category defined in connector.json
                WiContributionUtils.getConnections(this.http, "GraphBuilder").subscribe((data: IConnectorContribution[]) => {
                		let connectionRefs = [];
                    data.forEach(connection => {
                        for (let setting of connection.settings) {
							if(setting.name === "name") {
								connectionRefs.push({
									"unique_id": WiContributionUtils.getUniqueId(connection),
									"name": setting.value
								});
							}
                        }
                    });
                    observer.next(connectionRefs);
                		observer.complete();
                });
            });
        } else if (fieldName === "Nodes") {
        	    return buildEntity(this.http, this.selectedConnector, (content : string) => {
				let multiinstancesDef: IFieldDefinition = context.getField("Multiinstances");
				let nodesConfiguration: IFieldDefinition = context.getField("Nodes");
				return buildUIForNode(content , multiinstancesDef, nodesConfiguration);
			});
        } else if (fieldName === "Edges") {
        		return buildEntity(this.http, this.selectedConnector, (content : string) => {
				let multiinstancesDef: IFieldDefinition = context.getField("Multiinstances");
				let edgesConfiguration: IFieldDefinition = context.getField("Edges")
				return buildUIForEdge(content, multiinstancesDef, edgesConfiguration);
			});
        } else if (fieldName === "PassThroughData") {
            var attrJsonSchema = {};
            let attrNames: IFieldDefinition = context.getField("PassThrough");
            if (attrNames.value) {
                let data = JSON.parse(attrNames.value);
                for (var i = 0; i < data.length; i++) {
                	    attrJsonSchema[data[i].FieldName] = populateAttribute(data[i].Type);
                }
            }
            return JSON.stringify(attrJsonSchema);
        } else if (fieldName === "PassThroughDataOut") {
            var attrJsonSchema = {};
            let attrNames: IFieldDefinition = context.getField("PassThrough");
            if (attrNames.value) {
                let data = JSON.parse(attrNames.value);
                for (var i = 0; i < data.length; i++) {
                	    attrJsonSchema[data[i].FieldName] = populateAttribute(data[i].Type);
                }
            }
            return JSON.stringify(attrJsonSchema);
        }

        return null;
    }
 
    validate = (fieldName: string, context: IActivityContribution): Observable<IValidationResult> | IValidationResult => {

		console.log('[GraphBuilder::value] Validate field : ', fieldName);

        if (fieldName === "GraphModel") {
            let connection: IFieldDefinition = context.getField("GraphModel")
        		if (connection.value === null) {
            		return ValidationResult.newValidationResult().setError("GraphBuilder-MSG-1000", "Graph model must be configured");
        		}
        } else if (fieldName === "BatchEnd") {			
            let batchMode: IFieldDefinition = context.getField("BatchMode")
		console.log('[GraphBuilder::value] batchMode : ', batchMode);
        		if (batchMode.value === true) {
            		return ValidationResult.newValidationResult().setVisible(true);
        		} else {
				return ValidationResult.newValidationResult().setVisible(false);

			}
        }
		return null; 
    }
}

function buildEntity(http, selectedConnector, builder) : Observable<any> {
	return Observable.create(observer => {
		WiContributionUtils.getConnections(http, "GraphBuilder").subscribe((data: IConnectorContribution[]) => {
			let content : string;
        		data.forEach(connection => {
				var currentConnector;
            		for (let setting of connection.settings) {
					if(setting.name === "name") {
						currentConnector = setting.value
					}else if (setting.name === "metadata"&&
						selectedConnector === currentConnector) {
						content = setting.value;
						console.log("~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~")
						console.log(content)
						console.log("~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~")
                		}
            		}
        		});
			observer.next(JSON.stringify(builder(content)));
			observer.complete();
		});
	});			
}

function populateAttribute(attrType) : any {
	switch(attrType) {
		case "Double" :
    			return 2.0;
		case "Integer":
			return 2;
		case "Long":
			return 2;
		case "Boolean":
			return true;
		case "Date":
			return 2;
		default:
    			return "2";
	}
}

function buildUIForEdges(content : string, multiinstancesDef: IFieldDefinition, edgesConfiguration: IFieldDefinition): any {
	var edges = [{}];
	if(content) {
		var instanceSizeMap = {};
		if (multiinstancesDef.value) {
			let items = JSON.parse(multiinstancesDef.value);
			for (var i = 0; i < items.length; i++) {
				if("Edge"===items[i].EntityType) {
					instanceSizeMap[items[i].Name] = items[i].NumberOfInstances;
				}
			}
		}
					
		console.log(instanceSizeMap);
					
		if(edgesConfiguration.value) {
			var graphModel = JSON.parse(content);
			var nodeKeys = {};
			for(let node of graphModel["nodes"]) {
				var nodeName = node["name"];
				var keyInfo = [];
				for(let keyElement of node["key"]) {
					for(let attr of node["attributes"]) {
						if(keyElement===attr["name"]) {
							keyInfo.push(attr);
						}
					}
				}
								
				nodeKeys[nodeName] = keyInfo;
			}
							
			console.log("\n\nnodeKeys : " + nodeKeys)
							
			var relations = graphModel["edges"];
			for(let relation of relations) {	
				var relationName;
				var instanceSize = 1;
				if(instanceSizeMap[relation["name"]]) {
					instanceSize =  parseInt(instanceSizeMap[relation["name"]], 10);
				}
				for (var i=0; i<instanceSize; i++) {
					if(1<instanceSize) {
						relationName = relation["name"] + "_" + i;
					} else {
						relationName = relation["name"];
					}

					edges[0][relationName]=[{}];
					edges[0][relationName][0]["_skipCondition"] = false;
					var fromNode = relation["from"];
					var toNode = relation["to"];
					console.log("from : " + fromNode + ", to : " + toNode)
					edges[0][relationName][0]["vertices"] = [{}];
					edges[0][relationName][0]["vertices"][0]["from"] = "string"
					edges[0][relationName][0]["vertices"][0]["to"] = "string"
								
					if(relation["attributes"]) {
						for(let attr of relation["attributes"]) {
							var attrName = attr["name"];
							edges[0][relationName][0][attrName] = populateAttribute(attr["type"]);
						}
					}
				}
			}
		}
	}
            
    return edges;
}

function buildUIForEdge(content : string, multiinstancesDef: IFieldDefinition, edgesConfiguration: IFieldDefinition): any {
	var edges = {};
	if(content) {
		var instanceSizeMap = {};
		if (multiinstancesDef.value) {
			let items = JSON.parse(multiinstancesDef.value);
			for (var i = 0; i < items.length; i++) {
				if("Edge"===items[i].EntityType) {
					instanceSizeMap[items[i].Name] = items[i].NumberOfInstances;
				}
			}
		}
					
		console.log(instanceSizeMap);
					
		if(edgesConfiguration.value) {
			var graphModel = JSON.parse(content);
			var nodeKeys = {};
			for(let node of graphModel["nodes"]) {
				var nodeName = node["name"];
				var keyInfo = [];
				for(let keyElement of node["key"]) {
					for(let attr of node["attributes"]) {
						if(keyElement===attr["name"]) {
							keyInfo.push(attr);
						}
					}
				}

				nodeKeys[nodeName] = keyInfo;
			}
							
			console.log("\n\nnodeKeys : " + nodeKeys)
							
			var relations = graphModel["edges"];
			for(let relation of relations) {	
				var relationName;
				var instanceSize = 1;
				if(instanceSizeMap[relation["name"]]) {
					instanceSize =  parseInt(instanceSizeMap[relation["name"]], 10);
				}
				for (var i=0; i<instanceSize; i++) {
					if(1<instanceSize) {
						relationName = relation["name"] + "_" + i;
					} else {
						relationName = relation["name"];
					}

					edges[relationName]={};
					edges[relationName]["_skipCondition"] = false;
					var fromNode = relation["from"];
					var toNode = relation["to"];
					console.log("from : " + fromNode + ", to : " + toNode)
					edges[relationName]["vertices"] = {};
					edges[relationName]["vertices"]["from"] = "string"
					edges[relationName]["vertices"]["to"] = "string"
								
					if(relation["attributes"]) {
						for(let attr of relation["attributes"]) {
							var attrName = attr["name"];
							edges[relationName][attrName] = populateAttribute(attr["type"]);
						}
					}
				}
			}
		}
	}
            
    return edges;
}

function buildUIForNodes(content : string, multiinstancesDef: IFieldDefinition, nodesConfiguration: IFieldDefinition): any {
	let nodes = [{}];
    if(content && nodesConfiguration.value) {
		var instanceSizeMap = {};
		if (multiinstancesDef.value) {
			let items = JSON.parse(multiinstancesDef.value);
			for (var i = 0; i < items.length; i++) {
				if("Node"===items[i].EntityType) {
					instanceSizeMap[items[i].Name] = items[i].NumberOfInstances;
				}
			}
		}
					
		console.log(instanceSizeMap);
		let graphModel = JSON.parse(content);
					
		var entities = graphModel["nodes"];
							
		console.log(instanceSizeMap);

		for(let entity of entities) {
			var entityName;
			var instanceSize = 1;
			if(instanceSizeMap[entity["name"]]) {
				instanceSize =  parseInt(instanceSizeMap[entity["name"]], 10);
			}
			for (var i=0; i<instanceSize; i++) {
				if(1<instanceSize) {
					entityName = entity["name"] + "_" + i;
				} else {
					entityName = entity["name"];
				}
				nodes[0][entityName]=[{}];
				nodes[0][entityName][0]["_skipCondition"] = false;
				for(let attr of entity["attributes"]) {
					var attrName = attr["name"];
					nodes[0][entityName][0][attrName] = populateAttribute(attr["type"]);
				}
			}
		}
	}
    return nodes;
}

function buildUIForNode(content : string, multiinstancesDef: IFieldDefinition, nodesConfiguration: IFieldDefinition): any {
	let nodes = {};
    if(content && nodesConfiguration.value) {
		var instanceSizeMap = {};
		if (multiinstancesDef.value) {
			let items = JSON.parse(multiinstancesDef.value);
			for (var i = 0; i < items.length; i++) {
				if("Node"===items[i].EntityType) {
					instanceSizeMap[items[i].Name] = items[i].NumberOfInstances;
				}
			}
		}
					
		console.log(instanceSizeMap);
		let graphModel = JSON.parse(content);
					
		var entities = graphModel["nodes"];
							
		console.log(instanceSizeMap);

		for(let entity of entities) {
			var entityName;
			var instanceSize = 1;
			if(instanceSizeMap[entity["name"]]) {
				instanceSize =  parseInt(instanceSizeMap[entity["name"]], 10);
			}
			for (var i=0; i<instanceSize; i++) {
				if(1<instanceSize) {
					entityName = entity["name"] + "_" + i;
				} else {
					entityName = entity["name"];
				}
				nodes[entityName]={};
				nodes[entityName]["_skipCondition"] = false;
				for(let attr of entity["attributes"]) {
					var attrName = attr["name"];
					nodes[entityName][attrName] = populateAttribute(attr["type"]);
				}
			}
		}
		
		console.log(nodes);

	}
    return nodes;
}