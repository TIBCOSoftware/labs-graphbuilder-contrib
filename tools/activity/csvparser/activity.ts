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
export class CSVParserContributionHandler extends WiServiceHandlerContribution {
	selectedConnector: string;

    constructor(@Inject(Injector) injector, private http: Http) {
        super(injector, http);
    }

    value = (fieldName: string, context: IActivityContribution): Observable<any> | any => {
    	
		console.log('value ?????????? fieldName = ' + fieldName);
		let serveGraphData: IFieldDefinition = context.getField("ServeGraphData")
		
		if (fieldName === "GraphModel" && serveGraphData.value) {
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
        } else if (fieldName === "OutputFieldnames") {
        	    let attrNames: IFieldDefinition = context.getField("OutputFieldnames");
			console.log(attrNames);

			//[{"parameterName":"","type":"string","repeating":"false","required":"false","isEditable":true,"AttributeName":"aa","CSVFieldName":"aa","Default":"a","Type":"String","Optional":"yes"}]
        	    return buildData(this.http, this.selectedConnector, (content : string) => {
                	if(content) {
					let data = [];
					let graphModel = JSON.parse(content);
					let nodes = graphModel["nodes"];
					for(let node of nodes) {
						if(!node["attributes"]) {
							continue;
						}
						let nodeName = node["name"];
						for(let attr of node["attributes"]) {
							let attrName = attr["name"];
							data.push({
								"parameterName":"",
								"type":"string",
								"repeating":"false",
								"required":"false",
								"isEditable":true,
								"AttributeName":"node_"+nodeName+"_"+attrName,
								"CSVFieldName":"",
								"Default":"",
								"Type":attr["type"],
								"Optional":"yes"
							});
						}
					}
					let edges = graphModel["edges"];
					for(let edge of edges) {
						if(!edge["attributes"]) {
							continue;
						}
						let edgeName = edge["name"];
						for(let attr of edge["attributes"]) {
							let attrName = attr["name"];
							data.push({
								"parameterName":"",
								"type":"string",
								"repeating":"false",
								"required":"false",
								"isEditable":true,
								"AttributeName":"edge_"+edgeName+"_"+attrName,
								"CSVFieldName":"",
								"Default":"",
								"Type":"String",
								"Optional":attr["type"]
							});
						}
					}
					return data;
				}
				return attrNames.value
			});
			
        } else if (fieldName === "Data") {
            var attrJsonSchema = [{}];
            let attrNames: IFieldDefinition = context.getField("OutputFieldnames");
            if (attrNames.value) {
                let data = JSON.parse(attrNames.value);
                for (var i = 0; i < data.length; i++) {
                	    attrJsonSchema[0][data[i].AttributeName] = populateAttribute(data[i].Type);
                }
            }
            return JSON.stringify(attrJsonSchema);
        }
        
        return null;
    }
 
    validate = (fieldName: string, context: IActivityContribution): Observable<IValidationResult> | IValidationResult => {
		
		console.log('validate >>>>>>>> fieldName = ' + fieldName);

 		if (fieldName === "GraphModel") {
			let serveGraphData: IFieldDefinition = context.getField("ServeGraphData")
			console.log(serveGraphData)
        		if (serveGraphData.value) {
            		return ValidationResult.newValidationResult().setVisible(true);
        		} else {
				return ValidationResult.newValidationResult().setVisible(false);
			}
		}

		return null; 
    }
}

function buildData(http, selectedConnector, builder) : Observable<any> {
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