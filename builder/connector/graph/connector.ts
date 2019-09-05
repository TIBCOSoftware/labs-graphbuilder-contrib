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
export class TibcoGraphContribution extends WiServiceHandlerContribution {
    constructor() {
        super();
    }

    value = (fieldName: string, context: IConnectorContribution): Observable<any> | any => {
        return null;
    }
 
	validate = (name: string, context: IConnectorContribution): Observable<IValidationResult> | IValidationResult => {
		console.log('------------- validate --------------');
		console.log(context);
		console.log('>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>> name = ' + name);
		let modelSource: IFieldDefinition = context.getField("modelSource")
		console.log(modelSource);
 		if (name === "url") {
        		if (modelSource.value === "TGDB") {
            		return ValidationResult.newValidationResult().setVisible(true);
        		} else {
				return ValidationResult.newValidationResult().setVisible(false);
			}
		} else if (name === "model") {			
        		if (modelSource.value === "Local File") {
            		return ValidationResult.newValidationResult().setVisible(true);
        		} else {
				return ValidationResult.newValidationResult().setVisible(false);
			}
		} else if(name === "Connect") {
			let model: IFieldDefinition;
			let filename: string;
			let content: string; 

			for (let configuration of context.settings) {
				if( modelSource.value === "Local File" && configuration.name === "model") {
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
					}
					
					if(!model.value) {
						return ValidationResult.newValidationResult().setReadOnly(true);
					}
				}
			}
			return ValidationResult.newValidationResult().setReadOnly(false);
		}
 		return null;
	}

	action = (actionName: string, context: IConnectorContribution): Observable<IActionResult> | IActionResult => {
		
		if (actionName == "Connect") {
			console.log('------------- action click connect --------------');
            return Observable.create(observer => {
			console.log('------------- action callback --------------');
                /**
                 * These are the two fields that need to be checked whether they're filled in or not
                 */
                let model: IFieldDefinition;
 				let modelSource: IFieldDefinition = context.getField("modelSource")
				if (modelSource.value === "Local File") {
			console.log('----------------------------------------------->' + modelSource.value);
					model = context.getField("model");
					console.log('content = ' + model.value);
				} else if (modelSource.value === "TGDB") {
			console.log('----------------------------------------------->' + modelSource.value);
					model = context.getField("model");
					model.value.filename = "ServerModel";
					model.value.content = "{}";
					console.log('filename = ' + model.value.filename);
					console.log('content = ' + model.value.content);
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