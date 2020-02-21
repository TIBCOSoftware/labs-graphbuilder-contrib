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
		console.log('------------- SSEConnectorContribution validate --------------');
		console.log(context);
		let outbound: IFieldDefinition = context.getField("outbound")
		let tlsEnabled: IFieldDefinition = context.getField("tlsEnabled")
		let uploadCRT: IFieldDefinition = context.getField("uploadCRT")

		if(name === "Connect") {
			return ValidationResult.newValidationResult().setReadOnly(false);
		} else if( name === "url"||name === "resource"||name === "accessToken") {
        		if (outbound.value === true) {
            		return ValidationResult.newValidationResult().setVisible(true);
        		} else {
				return ValidationResult.newValidationResult().setVisible(false);
			}
		} else if( name === "port"||name === "path"||name === "tlsEnabled") {
        		if (outbound.value === false) {
            		return ValidationResult.newValidationResult().setVisible(true);
        		} else {
				return ValidationResult.newValidationResult().setVisible(false);
			}
		} else if( name === "uploadCRT") {
        		if (outbound.value === false&&tlsEnabled.value === true) {
            		return ValidationResult.newValidationResult().setVisible(true);
        		} else {
				return ValidationResult.newValidationResult().setVisible(false);
			}
		} else if( name === "tlsCRT"||name === "tlsKey") {
        		if (outbound.value === false&&tlsEnabled.value === true&&uploadCRT.value === true) {
            		return ValidationResult.newValidationResult().setVisible(true);
        		} else {
				return ValidationResult.newValidationResult().setVisible(false);
			}
		} else if( name === "tlsCRTPath"||name === "tlsKeyPath") {
        		if (outbound.value === false&&tlsEnabled.value === true) {
            		return ValidationResult.newValidationResult().setVisible(true);
        		} else {
				return ValidationResult.newValidationResult().setVisible(false);
			}
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