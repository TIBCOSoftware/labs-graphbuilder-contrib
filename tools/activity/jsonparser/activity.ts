/*
 * Copyright © 2020. TIBCO Software Inc.
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
export class JSONParserContributionHandler extends WiServiceHandlerContribution {

    constructor(@Inject(Injector) injector, private http: Http) {
        super(injector, http);
    }

    value = (fieldName: string, context: IActivityContribution): Observable<any> | any => {
    	
		console.log('value ?????????? fieldName = ' + fieldName);
		let serveGraphData: IFieldDefinition = context.getField("ServeGraphData")
		let attrNames: IFieldDefinition = context.getField("OutputFieldnames");
		let graphModel: IFieldDefinition = context.getField("GraphModel");
    		let data = [];
		if(attrNames.value) {
			data = JSON.parse(attrNames.value)
		}

		if (serveGraphData.value) {
			let previousConnector = context.getField("PreviousConnector").value
			let selectedConnector : string;
            	let allowedConnectors = context.getField("GraphModel").allowed;	
			let selectedConnectorId = context.getField("GraphModel").value;
			for(let allowedConnector of allowedConnectors) {
				if(allowedConnector["unique_id"] === selectedConnectorId) {
					selectedConnector = allowedConnector["name"]
				}
			}
			if (fieldName === "GraphModel") {
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
				//[{"parameterName":"","type":"string","repeating":"false","required":"false","isEditable":true,"AttributeName":"aa","JSONPath":"aa","Default":"a","Type":"String","Optional":"yes"}]
        		    return buildData(this.http, selectedConnector, (content : string) => {
					console.log('in value >>>>>>>> selectedConnector = ' + selectedConnector + ', previousConnector = ' + previousConnector);
        		        	if(content && (0===data.length||selectedConnector!=previousConnector)) {
						data = [];
						let graphModel = JSON.parse(content);
						buildAttributeDataForEntity("node", data, graphModel["nodes"])
						buildAttributeDataForEntity("edge", data, graphModel["edges"])
					}
					return data;
				});
        		} else if (fieldName === "PreviousConnector"){
				return selectedConnector;
			}
		}
		
		if (fieldName === "Data") {
            var attrJsonSchema = [{}];
            let attrNames: IFieldDefinition = context.getField("OutputFieldnames");
            if (attrNames.value) {
                let data = JSON.parse(attrNames.value);
                for (var i = 0; i < data.length; i++) {
                	    attrJsonSchema[0][data[i].AttributeName] = populateAttribute(data[i].Type);
                }
                attrJsonSchema[0]["LastElement"] = populateAttribute("Boolean");
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
		} else if (fieldName === "PreviousConnector") {
			return ValidationResult.newValidationResult().setVisible(false);
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
                		}
            		}
        		});
			console.log("OOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOO")
			console.log(content)
			console.log("OOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOO")
			console.log(builder(content))
			console.log("OOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOO")
			observer.next(JSON.stringify(builder(content)));
			observer.complete();
		});
	});			
}

function buildAttributeDataForEntity(entityType, data, entities) {
	for(let entity of entities) {
		let entityName = entity["name"];
		if(!entity["attributes"]) {
			continue;
		}
		for(let attr of entity["attributes"]) {
			let attrName = attr["name"];
			data.push({
				"parameterName":"",
				"type":"string",
				"repeating":"false",
				"required":"false",
				"isEditable":true,
				"AttributeName":entityType+"_"+entityName+"_"+attrName,
				"JSONPath":"",
				"Default":"",
				"Type":attr["type"],
				"Optional":"yes"
			});
		}
	}
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