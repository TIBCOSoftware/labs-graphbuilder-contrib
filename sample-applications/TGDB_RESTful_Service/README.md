# TGDB RESTful Service

## Overview

This example is created in TIBCO Flogo® Enterprise 2.8.0 studio. 

## Create TGDB Connection



## Create Application


### Create Flow for querying Metadata 

#### Configure flow inputs and outputs

1) input sample
```
{
    "queryType" : ""
}
```
2) output sample
```
{
    "queryResult": {
        "content": {},
        "success": true,
        "error": {
            "code": 101,
            "message": "Not found"
        }
    }
}
```
#### Add activities

1) GraphBuilder_TGDB -> TGDBQuery

2) Default -> Return

#### Add a trigger (Receive HTTP Message)

1) output

$trigger.pathParams.queryType

2) reply

$flow.queryResult

sample : 
```
{
    "queryResult": {
        "content": {},
        "success": true,
        "error": {
            "code": 101,
            "message": "Not found"
        }
    }
}
```
### Create Flow for Querying Data 

#### Configure flow inputs and outputs

1) input sample
```
{
    "queryType" : "",
    "language": "",
    "queryString": "",
    "traversalCondition": "",
    "endCondition": "",
    "traversalDepth": 1
}
```
2) output sample
```
{
    "queryResult": {
        "content": {},
        "success": true,
        "error": {
            "code": 101,
            "message": "Not found"
        }
    }
}
```
#### Add activities

1) GraphBuilder_TGDB -> TGDBQuery

2) Default -> Return

#### Add a trigger (Receive HTTP Message)

1) output

$trigger.pathParams.queryType
 and 
$trigger.body

sample :
```
{
  "query": {
  	"language" : "tgql",
    "queryString" : "@nodetype = 'houseMemberType' and memberName = 'Napoleon Bonaparte';",
    "traversalCondition" : "@edgetype = 'relation' and relation = 'spouse' and @isfromedge = 1 and @degree = 1;",
    "endCondition" : "",
	"traversalDepth" : 1
  }
}
```
2) reply

$flow.queryResult

sample : 
```
{
    "queryResult": {
        "content": {},
        "success": true,
        "error": {
            "code": 101,
            "message": "Not found"
        }
    }
}
```
