# Meet Up Event

## Overview

This example is created in TIBCO Flogo® Enterprise 2.8.0 studio. 

## Create Graph Model

![create_connection](create_connection.png)

## Create Application

![create_application](create_application.png)

### Create Flow for MeetUp Event 

![create_application2](create_application2.png)

#### Configure flow inputs and outputs

- **input sample** 
```
{
    "EventString" : ""
}
```
#### Add activities

- **Activity 1 :**
GraphBuilder_Tools -> JSONDeserializer

- **Activity 2 :**
GraphBuilder_Builder -> BuildGraph

- **Activity 3-1 :**
GraphBuilder_SSE -> SSEEndPoint

- **Activity 3-2 :**
GraphBuilder_TGDB -> TGDBUpsert

#### Add a trigger (GraphBuilder_SSE -> SSESubscriber)

- **output**

$trigger.Event map to flow input

### Create Flow for Serving Streaming Graph Data 

![create_application3](create_application3.png)

#### Configure flow inputs and outputs

Data Flow Comes from SSEEndPoint (MeetUp Event Flow) 

#### Add a trigger (Receive HTTP Message)

- **Incoming Query**

HTTP GET with resource path /sse/{streamId}

- **reply**

$flow.queryResult

sample : 
```
{
 "graph":{
  "edges":{},
  "id":"GeographyInfo",
  "model":{
   "edges":{
    "attrTypeMap":{"in_Continent":{}},
    "directionMap":{"in_Continent":1},
    "keyMap":{"in_Continent":null},
    "types":["in_Continent"],
    "vertexes":{"in_Continent":["City","Continent"]}
   },
   "nodes":{
    "attrTypeMap":{"Continent":{"Name":"String"},"Country":{"Country_Code":"String"}},
    "keyMap":{"Continent":["Name"],"Country":["Country_Code"]},"types":["Country","Continent"]
   }
  },
  "modelId":"GeographyInfo",
  "nodes":{
   "Continent_0ecff3229a1a13980689def44b2c66e1":{
    "attributes":{"Name":{"name":"Name","type":"String","value":"North_America"}},
    "key":["North_America"],
    "keyAttributeName":["Name"],
    "type":"Continent"
   },
   "Country_5181a8acdef7be40dfbf3ec66bee2b20":{
    "attributes":{"Country_Code":{"name":"Country_Code","type":"String","value":"us"}},
    "key":["us"],
    "keyAttributeName":["Country_Code"],
    "type":"Country"
   }
  }
 }
}
```
