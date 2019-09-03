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
export class CloudConversationContributionHandler extends WiServiceHandlerContribution {
	url: string;
	user: string; 
	passsword: string; 
	databaseType: string
		
    constructor(@Inject(Injector) injector, private http: Http) {
        super(injector, http);
    }

    value = (fieldName: string, context: IActivityContribution): Observable<any> | any => {
    	
		console.log('value >>>>>>>> fieldName = ' + fieldName);
		
		if (fieldName === "databaseConnection") {
            //Connector Type must match with the category defined in connector.json
            return Observable.create(observer => {
                let connectionRefs = [];
                let databaseType = context.getField("databaseType").value;;
                let connectionName: string = "tibco-tgdb";
                if (databaseType === "Dgraph") {
                		connectionName = "tibco-dgraph";
                } else if(databaseType === "Neo4j") {
                		connectionName = "tibco-neo4j";
                }
                 
                WiContributionUtils.getConnections(this.http).subscribe((data: IConnectorContribution[]) => {
                    data.forEach(connection => {
                    		console.log(connection.name);
						if (connectionName===connection.name) {
	                        for (let i = 0; i < connection.settings.length; i++) {
								if(connection.settings[i].name === "name") {
									connectionRefs.push({
										"unique_id": WiContributionUtils.getUniqueId(connection),
										"name": connection.settings[i].value
									});
								}
                        		}						
						}
                    });
                    observer.next(connectionRefs);
                });
            });
        } else if (fieldName === "queryResult") {
			var queryResult = {
    				"data" : {},
    				"success": true,
    				"error": {
        				"code" : 101,
        				"message" : "string"
    				}
			};
			return JSON.stringify(queryResult);
        } else if (fieldName === "pathParams") {
        		var pathParams = {
				"queryType" : "queryType",
				"entityType" :"entityType"	
			};
			return JSON.stringify(pathParams);
        } else if (fieldName === "queryParams") {
        		var queryParams = {
			};
			return JSON.stringify(queryParams);
        }
        
        return null;
    }
 
    validate = (fieldName: string, context: IActivityContribution): Observable<IValidationResult> | IValidationResult => {
		
		console.log('validate >>>>>>>> fieldName = ' + fieldName);

        if (fieldName === "tgdbConnection") {
            let connection: IFieldDefinition = context.getField("tgdbConnection")
        	if (connection.value === null) {
            	return ValidationResult.newValidationResult().setError("TGDBUpsert-MSG-1000", "Graph model must be configured");
        	}
        }
		return null; 
    }
}