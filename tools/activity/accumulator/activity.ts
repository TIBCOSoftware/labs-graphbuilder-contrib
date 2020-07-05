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
export class AccumulatorContributionHandler extends WiServiceHandlerContribution {
	filename: string;
	content: string; 

    constructor(@Inject(Injector) injector, private http: Http) {
        super(injector, http);
    }

    value = (fieldName: string, context: IActivityContribution): Observable<any> | any => {
    	
		console.log('[AccumulatorContributionHandler::value] fieldName = ' + fieldName);
		var attrJsonSchema;
		if (fieldName === "Input") {	
			let arrayMode: IFieldDefinition = context.getField("ArrayMode")
			if (arrayMode.value === true) {
            		attrJsonSchema = [{}];
            		let attrNames: IFieldDefinition = context.getField("DataFields");
            		if (attrNames.value) {
            		    let data = JSON.parse(attrNames.value);
            		    for (var i = 0; i < data.length; i++) {
            		    		attrJsonSchema[0][data[i].Name] = populateAttribute(data[i].Type);
            		    }
            		}
        		} else {
            		attrJsonSchema = {};
            		let attrNames: IFieldDefinition = context.getField("DataFields");
            		if (attrNames.value) {
            		    let data = JSON.parse(attrNames.value);
            		    for (var i = 0; i < data.length; i++) {
            		    		attrJsonSchema[data[i].Name] = populateAttribute(data[i].Type);
            		    }
            		}
			} 		
            return JSON.stringify(attrJsonSchema);
        } else if (fieldName === "Output") {			
            attrJsonSchema = [{}];
            let attrNames: IFieldDefinition = context.getField("DataFields");
            if (attrNames.value) {
                let data = JSON.parse(attrNames.value);
                for (var i = 0; i < data.length; i++) {
                		attrJsonSchema[0][data[i].Name] = populateAttribute(data[i].Type);
                }
                attrJsonSchema[0]["LastElement"] = populateAttribute("Boolean");
            }
            return JSON.stringify(attrJsonSchema);
        }
        
        return null;
    }
 
    validate = (fieldName: string, context: IActivityContribution): Observable<IValidationResult> | IValidationResult => {
		
		console.log('[AccumulatorContributionHandler::validate] fieldName = ' + fieldName);
		
		if (fieldName === "WindowSize") {			
            let arrayMode: IFieldDefinition = context.getField("ArrayMode")
            
			console.log('[AccumulatorContributionHandler::WindowSize] arrayMode : ', arrayMode);
			
        		if (arrayMode.value === true) {
            		return ValidationResult.newValidationResult().setVisible(false);
        		} else {
				return ValidationResult.newValidationResult().setVisible(true);

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