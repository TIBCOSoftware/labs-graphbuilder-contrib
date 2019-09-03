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
export class SQLSubscriberContributionHandler extends WiServiceHandlerContribution {
	url: string;
	user: string; 
	passsword: string; 
		
    constructor(@Inject(Injector) injector, private http: Http) {
        super(injector, http);
    }

    value = (fieldName: string, context: IActivityContribution): Observable<any> | any => {
    	
		console.log('value >>>>>>>> fieldName = ' + fieldName);
		
        if (fieldName === "sqlConnection") {
            //Connector Type must match with the category defined in connector.json
            return Observable.create(observer => {
                let connectionRefs = [];
                
                WiContributionUtils.getConnections(this.http, "GraphBuilder_SQL").subscribe((data: IConnectorContribution[]) => {
                    data.forEach(connection => {
                        for (let i = 0; i < connection.settings.length; i++) {
                        		if(connection.settings[i].name === "name") {
                            		connectionRefs.push({
                                		"unique_id": WiContributionUtils.getUniqueId(connection),
                                		"name": connection.settings[i].value
                            		});
                            }
                        }
                    });
                    observer.next(connectionRefs);
                });
            });
        } else if (fieldName === "DataRow") {
            var attrJsonSchema = {};
            let attrNames: IFieldDefinition = context.getField("outputFieldMap");
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
		
		console.log('validate >>>>>>>> fieldName = ' + fieldName);

        if (fieldName === "sqlConnection") {
            let connection: IFieldDefinition = context.getField("sqlConnection")
        		if (connection.value === null) {
        			return ValidationResult.newValidationResult().setError("SQL-SERVICE-MSG-1000", "Connector must be configured");
        		}
        }
		return null; 
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