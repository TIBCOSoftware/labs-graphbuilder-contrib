/*
 * Copyright © 2020. TIBCO Software Inc.
 * This file is subject to the license terms contained
 * in the license file that is distributed with this file.
 */
import {Injectable, Injector, Inject} from "@angular/core";
import {Http} from "@angular/http";

import {
	WiContrib, 
	WiServiceHandlerContribution, 
	AUTHENTICATION_TYPE,
    IActivityContribution,
    WiContributionUtils
} from "wi-studio/app/contrib/wi-contrib";

import {IConnectorContribution, IFieldDefinition, IActionResult, ActionResult} from "wi-studio/common/models/contrib";
import {Observable} from "rxjs/Observable";
import {IValidationResult, ValidationResult, ValidationError} from "wi-studio/common/models/validation";

@WiContrib({})
@Injectable()
export class TibcoGraphContribution extends WiServiceHandlerContribution {
	
    /*constructor() {
        super();
    }*/
    constructor(@Inject(Injector) injector, private http: Http) {
        super(injector, http);
    }

    value = (fieldName: string, context: IConnectorContribution): Observable<any> | any => {
		console.log('[GraphBuilder:Connector:value] Build field : ', fieldName);
		console.log(context);
		let modelSource: IFieldDefinition = context.getField("modelSource")
		let model: IFieldDefinition = context.getField("model")
        if (fieldName === "metadata") {
			if(modelSource.value === "TGDB") {
				return Observable.create(observer => {
					let url = context.getField("url");
					this.http.get(url.value).subscribe((res)=>{
						observer.next(JSON.stringify(metadataToModel(res.text())));
                			observer.complete();
        				});
				});
			} else {
				if(model.value) {
					let content = model.value.content;
					if(content) {
						content = content.substr(content.indexOf(',')+1);
						content = atob(content);
					}
					console.log('content = ' + content);
					return content
				}
			}
        }
        return null;
    }

	validate = (name: string, context: IConnectorContribution): Observable<IValidationResult> | IValidationResult => {
		console.log('------------- validate --------------');
		console.log(context);
		console.log('>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>> name = ' + name);
		let modelSource: IFieldDefinition = context.getField("modelSource")
		let metadata: IFieldDefinition = context.getField("metadata")
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
		} else if (name === "metadata") {			
			return ValidationResult.newValidationResult().setVisible(false);
		} else if(name === "Connect") {
			//setTimeout(()=>{
			//	if(!metadata.value) {
			//		return ValidationResult.newValidationResult().setReadOnly(true);
			//	}
				return ValidationResult.newValidationResult().setReadOnly(false);
			//}, 10000)
		}
 		return null;
	}

	action = (actionName: string, context: IConnectorContribution): Observable<IActionResult> | IActionResult => {
		
		if (actionName == "Connect") {
			console.log('------------- action click connect --------------');
            return Observable.create(observer => {
                /**
                 * These are the two fields that need to be checked whether they're filled in or not
                 */
                let model: IFieldDefinition = context.getField("model");
                let metadata: IFieldDefinition = context.getField("metadata");
 				let modelSource: IFieldDefinition = context.getField("modelSource")
				if (modelSource.value === "Local File") {
					console.log('content = ' + model.value);
				} else if (modelSource.value === "TGDB") {
					console.log('metadata = ' + metadata.value);
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

function metadataToModel(metadata) : any {
	let metadataObj = JSON.parse(metadata);
	console.log(metadataObj.data.data)
	let modelObj = {"nodes":[], "edges":[]}
	for(let edgeType of metadataObj.data.data.edgeTypes) {
		let edgeName = edgeType.name;
		if('$'===edgeName[0]) {
			continue;
		}
		
		modelObj.edges.push(buildEdge(edgeType))
	}
	
	for(let nodeType of metadataObj.data.data.nodeTypes) {
		modelObj.nodes.push(buildNode(nodeType))
	}

	console.log(modelObj)

	return modelObj;
}

function buildNode(nodeType) : any {
	let nodeTypeObj = {
		"name" : nodeType.name,
		"key" : [],
		"attributes" : []
	}
	
	for(let pkeyType of nodeType.pkeyAttributeDescriptors) {
		nodeTypeObj.key.push(buildPKey(pkeyType))
	}
	
	for(let attrType of nodeType.attributeDescriptors) {
		nodeTypeObj.attributes.push(buildAttribute(attrType))
	}
	return nodeTypeObj;
}

function buildEdge(edgeType) : any {
	let edgeTypeObj = {
		"name" : edgeType.name,
		"from" : edgeType.fromNodeType.name,
		"to" : edgeType.toNodeType.name,
		"attributes" : []
	}

	for(let attrType of edgeType.attributeDescriptors) {
		edgeTypeObj.attributes.push(buildAttribute(attrType))
	}
	return edgeTypeObj;
}

function buildPKey(keyAttributeType) : any {
	return keyAttributeType.name
}

function buildAttribute(attributeType) : any {
	return { "name" : attributeType.name, "type" : populateAttributeType(attributeType.type) }
}

function populateAttributeType(attrType) : string {
	switch(attrType) {
		case 0: //AttributeTypeInvalid
			return "Invalid"
		case 1: //AttributeTypeBoolean
			return "Boolean";
		case 2: //AttributeTypeByte
			return "Byte";
		case 3: //AttributeTypeChar
			return "Char";
		case 4: //AttributeTypeShort
			return "Short";
		case 5: //AttributeTypeInteger
    			return "Integer";
		case 6: //AttributeTypeLong
			return "Long";
		case 7: //AttributeTypeFloat
			return "Float";
		case 8: //AttributeTypeDouble
			return "Double";
		case 9: //AttributeTypeNumber
			return "Number";
		case 10: //AttributeTypeString
    			return "String";
		case 11: //AttributeTypeDate
			return "Date";
		case 12: //AttributeTypeTime
			return "Time";
		case 13: //AttributeTypeTimeStamp
			return "TimeStamp";
		case 14: //AttributeTypeBlob
			return "Blob";
		case 15: //AttributeTypeClob
			return "Clob";
		default:
    			return "String";
	}
}
