/*
 * Copyright © 2019. TIBCO Software Inc.
 * This file is subject to the license terms contained
 * in the license file that is distributed with this file.
 */
import {Injectable} from "@angular/core";
import {WiContrib, WiServiceHandlerContribution, AUTHENTICATION_TYPE} from "wi-studio/app/contrib/wi-contrib";
import {IConnectorContribution, IFieldDefinition, IActionResult, ActionResult} from "wi-studio/common/models/contrib";
import {Observable} from "rxjs/Observable";
import {IValidationResult, ValidationResult, ValidationError} from "wi-studio/common/models/validation";

@WiContrib({})
@Injectable()
export class SSEConnectorContribution extends WiServiceHandlerContribution {
    constructor() {
        super();
    }

    value = (fieldName: string, context: IConnectorContribution): Observable<any> | any => {
        return null;
    }

    validate = (name: string, context: IConnectorContribution): Observable<IValidationResult> | IValidationResult => {
		console.log('------------- validate --------------');
		console.log(context);
		console.log('%%%%%%%%%% name = ' + name);
		
		if(name === "Connect") {
			let url: string;
			let port: string;
			let outbound: string; 
			let resource: string; 
			let accessKey: string; 

			for (let configuration of context.settings) {
				if( configuration.name === "url") {
					url = configuration.value	
				} else if( configuration.name === "port") {
					port = configuration.value
				} else if( configuration.name === "outbound") {
					outbound = configuration.value
				} else if( configuration.name === "resource") {
					resource = configuration.value
				} else if( configuration.name === "accessKey") {
					accessKey = configuration.value
				}
			}

			//if( url && resource ) {
				return ValidationResult.newValidationResult().setReadOnly(false);
			//} else {
			//	return ValidationResult.newValidationResult().setReadOnly(true);
			//}
		}
 		return null;
    }

    action = (actionName: string, context: IConnectorContribution): Observable<IActionResult> | IActionResult => {
		if (actionName == "Connect") {
            return Observable.create(observer => {
                let actionResult = {
                    context: context,
                    authType: AUTHENTICATION_TYPE.BASIC,
                    authData: {}
                }
                observer.next(ActionResult.newActionResult().setSuccess(true).setResult(actionResult));
            });
        }
		
		return null;
    }
}