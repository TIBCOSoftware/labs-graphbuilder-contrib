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
export class TableUpsertContributionHandler extends WiServiceHandlerContribution {
	selectedConnector: string;

    constructor(@Inject(Injector) injector, private http: Http) {
        super(injector, http);
    }

    value = (fieldName: string, context: IActivityContribution): Observable<any> | any => {
    	
		console.log('[TableContributionHandler::value] fieldName = ' + fieldName);
		
        if (fieldName === "Table") {
            let allowedConnectors = context.getField("Table").allowed;	
			let selectedConnectorId = context.getField("Table").value;
			for(let allowedConnector of allowedConnectors) {
				if(allowedConnector["unique_id"] === selectedConnectorId) {
					this.selectedConnector = allowedConnector["name"]
				}
			}
            
            return Observable.create(observer => {
            		//Connector Type must match with the category defined in connector.json
                WiContributionUtils.getConnections(this.http, "GraphBuilder_Tools").subscribe((data: IConnectorContribution[]) => {
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
        } else if (fieldName === "Mapping" || fieldName === "Data") {	
        		return buildSchema(this.http, this.selectedConnector, (schema : string) => {
            		var attrJsonSchema = {};
            		if (schema) {
                		let data = JSON.parse(schema);
                		for (var i = 0; i < data.length; i++) {
                			attrJsonSchema[data[i].Name] = populateAttribute(data[i].Type);
                		}
            		}
            		return attrJsonSchema;
			});
        }
        
        return null;
    }
 
    validate = (fieldName: string, context: IActivityContribution): Observable<IValidationResult> | IValidationResult => {
		
		console.log('[TableContributionHandler::validate] fieldName = ' + fieldName);

		return null; 
    }
}

function buildSchema(http, selectedConnector, builder) : Observable<any> {
	return Observable.create(observer => {
		WiContributionUtils.getConnections(http, "GraphBuilder_Tools").subscribe((data: IConnectorContribution[]) => {
			var schema;
        		data.forEach(connection => {
				var currentConnector;
            		for (let setting of connection.settings) {
					if(setting.name === "name") {
						currentConnector = setting.value
					}else if (setting.name === "schema"&&
						selectedConnector === currentConnector) {
						schema = setting.value;
                		}
            		}
        		});
			observer.next(JSON.stringify(builder(schema)));
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