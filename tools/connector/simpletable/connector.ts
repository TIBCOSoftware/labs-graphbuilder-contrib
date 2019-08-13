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
export class TibcoGraphBuilderContribution extends WiServiceHandlerContribution {
    constructor() {
        super();
    }

    value = (fieldName: string, context: IConnectorContribution): Observable<any> | any => {
		console.log('------------- value --------------');
		console.log(context);
		console.log('%%%%%%%%%% fieldName = ' + fieldName);
		if (fieldName === "key") {
			let allowedKeys = [];
			let columnNames: IFieldDefinition = context.getField("schema");
			if (columnNames.value) {
                let data = JSON.parse(columnNames.value);
                for (var i = 0; i < data.length; i++) {
					allowedKeys.push({
						"unique_id": data[i].Name,
						"name": data[i].Name
					});
                }
            }
			return allowedKeys;
        }
        return null;
    }
 
	validate = (name: string, context: IConnectorContribution): Observable<IValidationResult> | IValidationResult => {
		console.log('------------- validate --------------');
		console.log(context);
		console.log('%%%%%%%%%% name = ' + name);
		
		if(name === "create") {
			let key: IFieldDefinition;

			for (let configuration of context.settings) {
				if( configuration.name === "key") {
					key = configuration
					if(configuration.value) {
						console.log('key = ', configuration.value);
					}			
				}
			}

			if( key.value) {
				return ValidationResult.newValidationResult().setReadOnly(false);
			} else {
				return ValidationResult.newValidationResult().setReadOnly(true);

			}
		}
 		return null;
	}

	action = (actionName: string, context: IConnectorContribution): Observable<IActionResult> | IActionResult => {
		
		if (actionName == "create") {
			console.log('------------- action click connect --------------');
            return Observable.create(observer => {
			console.log('------------- action callback --------------');
                /**
                 * These are the two fields that need to be checked whether they're filled in or not
                 */
                let key: IFieldDefinition;
 
                for (let configuration of context.settings) {
                    if (configuration.name === "key") {
                        key = configuration
						console.log(key.value);
                    }
                }

                /**
                 * Set the action result to save the configuration data
                 */
                let actionResult = {
                    context: context,
                    authType: AUTHENTICATION_TYPE.BASIC,
                    authData: {}
                }

                /**
                 * Call the observer and tell it the validation was sucessful and the data should be saved
                 */
                observer.next(ActionResult.newActionResult().setSuccess(true).setResult(actionResult));
            });
        }
		
		return null;
    }
}