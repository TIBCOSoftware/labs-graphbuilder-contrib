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
export class Neo4jContributionHandler extends WiServiceHandlerContribution {
	url: string;
	user: string; 
	passsword: string; 
		
    constructor(@Inject(Injector) injector, private http: Http) {
        super(injector, http);
    }

    value = (fieldName: string, context: IActivityContribution): Observable<any> | any => {
    	
		console.log('value >>>>>>>> fieldName = ' + fieldName);
		
        if (fieldName === "neo4jConnection") {
            //Connector Type must match with the category defined in connector.json
            return Observable.create(observer => {
                let connectionRefs = [];
                
                WiContributionUtils.getConnections(this.http, "GraphBuilder_neo4j").subscribe((data: IConnectorContribution[]) => {
                    data.forEach(connection => {
                        for (let i = 0; i < connection.settings.length; i++) {
                        	 if(connection.settings[i].name === "name") {
                            	connectionRefs.push({
                                	"unique_id": WiContributionUtils.getUniqueId(connection),
                                	"name": connection.settings[i].value
                            	});
							 } else if (connection.settings[i].name === "url") {
							 } else if (connection.settings[i].name === "user") {
							 } else if (connection.settings[i].name === "password") {
                            }
                        }
                    });
                    observer.next(connectionRefs);
                });
            });
        }
        
        return null;
    }
 
    validate = (fieldName: string, context: IActivityContribution): Observable<IValidationResult> | IValidationResult => {
		
		console.log('validate >>>>>>>> fieldName = ' + fieldName);

        if (fieldName === "neo4jConnection") {
            let connection: IFieldDefinition = context.getField("neo4jConnection")
        	if (connection.value === null) {
            	return ValidationResult.newValidationResult().setError("Neo4j-MSG-1000", "Neo4j must be configured");
        	}
        }
		return null; 
    }
}