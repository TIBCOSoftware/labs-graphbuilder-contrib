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
export class JSONDeserializerContributionHandler extends WiServiceHandlerContribution {
	filename: string;
	content: string; 

    constructor(@Inject(Injector) injector, private http: Http) {
        super(injector, http);
    }

    value = (fieldName: string, context: IActivityContribution): Observable<any> | any => {
    	
		console.log('fieldName = ' + fieldName);
		console.log(context.settings);
		
		if(fieldName === "Data") {
			let model: IFieldDefinition;
			let filename: string;
			let content: string; 
			
			for (let configuration of context.settings) {
				if( configuration.name === "sample") {
					let schemaObj: any
					model = configuration
					if(configuration.value) {
						filename = configuration.value.filename;
						content = configuration.value.content;
						if(content) {
							content = content.substr(content.indexOf(',')+1);
							content = atob(content);
						}
						console.log('filename = ' + filename);
						console.log('content = ' + content);
						return JSON.stringify(JSON.parse(content));
					}
				}
			}
			
			return JSON.stringify({})
		}
        
        return null;
    }
 
    validate = (fieldName: string, context: IActivityContribution): Observable<IValidationResult> | IValidationResult => {
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

/*
function buildJSON(jsonelement) : any {
	
	switch(typeof jsonelement) {
		case object :
    			return 2.0;
		case array :
			return 2;
		default:
    			return "2";
	}
}
*/