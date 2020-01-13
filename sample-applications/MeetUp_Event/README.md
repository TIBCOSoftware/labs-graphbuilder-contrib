# Meet Up Event

## Overview

- This example is created in TIBCO Flogo® Enterprise 2.8.0 studio. 

- This example uses Meetup open event through Meetup API see https://www.meetup.com/meetup_api/

## Create Graph Model

![create_connection](create_connection.png)
### Setting
- **Graph Name:** -> Meetup
- **Model Source:** -> Select Local File
- **Graph Model:** -> Select sample-applications/Meetup_Event/Model_Meetup.json

## Create Connection to subscribe MeetUp event 

![create_connection2](create_connection2.png)
### Setting
- **Connection Name:** -> Meetup_Event
- **Outbound:** -> Sellect "true" for connecting to Meetup service
- **Server URL:** -> http://stream.meetup.com/
- **Resource Name:** -> 2/open_events
- **Access Token:** -> not required for accessing open event

## Create Connection to serve streaming graph data

![create_connection3](create_connection3.png)
### Setting
- **Connection Name:** -> EventServer
- **Outbound:** -> select "false" since it's a server
- **Server port:** -> any available port (8888 for this example)
- **Path:** -> /sse/ (client connect http://[host]:[port]/sse/meetup to subscribe "meetup" graph stream)

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
### Add Activity 1
Select GraphBuilder_Tools -> JSONDeserializer
- **JSON Data Sample:** -> Select sample-applications/Meetup_Event/.json
- **Default Values:** -> Set "na" as default for venue.address_1, category.name

### Add Activity 2
Select GraphBuilder_Builder -> BuildGraph
- **Graph Model:** -> Select "Meetup" (the connection we created previously)
- **Configure Model:** -> Map attributes to input data fields (for nodes and edges) 

### Add Activity 3-1
Select GraphBuilder_SSE -> SSEEndPoint
- **SSE Connection:** -> Select "EventServer" for serving streaming data(the connection we created previously)
- **Avtivity Input 1:** set StreamId to "meetup" (the resource name for client to subscribe)
- **Avtivity Input 2:** map required Data object to $activity[BuildGraph].Graph (output of BuildGraph activity)

### Add Activity 3-2
Select GraphBuilder_TGDB -> TGDBUpsert
- **TGDB Connection:** -> Select "TGDB" for upserting streaming data to TGDB(the connection we created in TGDB_RESTful_Service sample application)
- **Avtivity Input 1:** set required Graph object to $activity[BuildGraph].Graph

#### Add a trigger 
Select GraphBuilder_SSE -> SSESubscriber
- **SSE Connection(outbound request):** -> Select "Meetup_Event" for consuming open event from Meetup web site
- **Flow Input:** -> Map EventString to $trigger.Event (This is the output of SSESubscriber)

$trigger.Event map to flow input

### Create Flow for Serving Streaming Graph Data 

![create_application3](create_application3.png)

#### Configure flow inputs and outputs

No configuration is required here since the data flow comes from SSEEndPoint of Meetup Event Flow directly

#### Add a trigger 
 GraphBuilder_SSE -> SSESubscriber
- **SSE Connection(inbound requests):** -> Select "EventServer" for serving streaming data(so now SSEEndPoint connected)
- **Flow Input:** -> Map EventString to $trigger.Event (This is the output of SSESubscriber)

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
