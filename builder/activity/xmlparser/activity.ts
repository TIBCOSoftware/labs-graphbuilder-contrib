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
export class XMLParserContributionHandler extends WiServiceHandlerContribution {
	filename: string;
	content: string; 

    constructor(@Inject(Injector) injector, private http: Http) {
        super(injector, http);
    }

    value = (fieldName: string, context: IActivityContribution): Observable<any> | any => {
    	
		console.log('value ?????????? fieldName = ' + fieldName);
		
		if (fieldName === "Data") {
            var attrJsonSchema = {};
            let attrNames: IFieldDefinition = context.getField("OutputFieldnames");
            if (attrNames.value) {
                let data = JSON.parse(attrNames.value);
                for (var i = 0; i < data.length; i++) {
                    if (data[i].Type === "String") {
                        attrJsonSchema[data[i].AttributeName] = "abc";
                    } else if (data[i].Type === "Integer") {
                        attrJsonSchema[data[i].AttributeName] = 0;
                    }
                }
            }
            return JSON.stringify(attrJsonSchema);
        }
        
        return null;
    }
 
    validate = (fieldName: string, context: IActivityContribution): Observable<IValidationResult> | IValidationResult => {
		
		console.log('validate >>>>>>>> fieldName = ' + fieldName);

		return null; 
    }
}