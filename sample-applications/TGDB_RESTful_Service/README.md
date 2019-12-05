# TGDB RESTful Service

## Overview

## Create TGDB Connection

## Create Application

As discussed above this implementation is written with Golang within the Flogo ecosystem.  As such CatalystML can be used with the flogo command line interface (with a flogo.json) or the Golang Flogo API (library).  Two examples for each of the CLI or the API are discussed below.

### Create Flow for querying Metadata 

1) Configure flow inputs and outputs

input sample

{
    "queryType" : ""
}

output sample

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

2) Add activities

GraphBuilder_TGDB -> TGDBQuery

Default -> Return

3) Add a trigger (Receive HTTP Message)

output : 

$trigger.pathParams.queryType

reply : 

$flow.queryResult

reply sample : 

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

### Create Flow for Querying Data 

1) Configure flow inputs and outputs

input sample

{
    "queryType" : "",
    "language": "",
    "queryString": "",
    "traversalCondition": "",
    "endCondition": "",
    "traversalDepth": 1
}

output sample

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

2) Add activities

GraphBuilder_TGDB -> TGDBQuery

Default -> Return

3) Add a trigger (Receive HTTP Message)

output : 

$trigger.pathParams.queryType

and 

{
  "query": {
  	"language" : "tgql",
    "queryString" : "@nodetype = 'houseMemberType' and memberName = 'Napoleon Bonaparte';",
    "traversalCondition" : "@edgetype = 'relation' and relation = 'spouse' and @isfromedge = 1 and @degree = 1;",
    "endCondition" : "",
	"traversalDepth" : 1
  }
}

reply : 

$flow.queryResult

reply sample : 

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