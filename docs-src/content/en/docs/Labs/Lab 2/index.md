---
title: "Lab2 - Query"
linkTitle: "Lab2 - Query"
weight: 2
description: >
  Build an application to query against TIBCO® Graph Database
---

Create a new Flogo application called “TGDB_RESTful_Service”

![Build RESTful](RESTful01.png)

Click “+ Create” button to build from scratch

![Build RESTful](RESTful02.png)

Create the first flow for querying metadata

![Build RESTful](RESTful03.png)

Define the data schema for the input of current flow sample data (queryType in string data type).
- queryType: the value could be “metadata”, “edgetypes” or “nodetypes” (metadata querying flow)

![Build RESTful](RESTful04.png)

Saving sample data will evoke schema builder to generate the schema definition from it

![Build RESTful](RESTful05.png)

Define the output schema for the current flow by pasting sample output data
- Content: contains the data of query result
- Success: true means query go through without error
- Code: error code
- Message: error message

Click “Save” button

![Build RESTful](RESTful06.png)

Clicking “Save” button triggers schema definition generation

![Build RESTful](RESTful07.png)

Add a trigger to receive HTTP request by clicking “+” -> “ReceiveHTTPMessage”

![Build RESTful](RESTful07-5.png)

Select GET, enter resource path “/tgdb/{queryType}” then click “Finish”

![Build RESTful](RESTful09.png)

Now we have a trigger with HTTP GET methods and listen on port 9999)

![Build RESTful](RESTful08.png)

Click the icon of trigger to map incoming query data to flow input data

![Build RESTful](RESTful10.png)

In “Reply Settings” set reply schema make it same as flow output data schema

![Build RESTful](RESTful11.png)

In “Map from flow outputs” mapping data.queryResult to $flow.queryResult

![Build RESTful](RESTful12.png)

Add query activity by selecting GraphBuilder_TGDB -> TGDBQuery activity

![Build RESTful](RESTful13.png)

Select the “TGDB” connection that was created in Lab1 so the TGDBQuery activity executes against the server where the Northwind data was inserted

![Build RESTful](RESTful14.png)

Map input data for TGDBQuery activity
- QueryType : $flow.queryType

![Build RESTful](RESTful15.png)

Add return activity to link the query result back to HTTP trigger

![Build RESTful](RESTful16.png)

Map outputs for Return activity
- queryResult : $activity[TGDBQuery].queryResult (map entire object)

![Build RESTful](RESTful17.png)

You've finished creating metadata query flow

![Build RESTful](RESTful18.png)

Click “Create” button to create another flow for querying the content of Northwind graph

![Build RESTful](RESTful19.png)

Add a name and description for the new flow

![Build RESTful](RESTful20.png)

Define the flow inputs data schema by pasting sample data (for schema detail; see TGDB documentation)
- queryType: search (for content flow)
- language: TGQL (TIBCO graph query language) or Gremlin
- queryString: for TGQL and Gremlin
- traversalCondition: TGQL only
- traversalDepth: TGQL only

![Build RESTful](RESTful21.png)

Click save to generate data schema definition

![Build RESTful](RESTful22.png)

Flow output data schema same as metadata flow

![Build RESTful](RESTful22-5.png)

Add another trigger for receiving content query

![Build RESTful](RESTful23.png)

POST method for content query

![Build RESTful](RESTful24.png)

Adding sample query for the output (to the flow) setting. To be noticed that the schema is very similar to flow input schema but grouped under “query” keyword.

![Build RESTful](RESTful25.png)

Map to flow input
- queryType: $trigger.pathParams.queryType
- language: $trigger.body.query.language
- queryString: $trigger.body.query.queryString
- traversalCondition: $trigger.body.query.traversalCondition
- endCondition: $trigger.body.query.endCondition
- traversalDepth: $trigger.body.query.traversalDepth

Click save

![Build RESTful](RESTful26.png)

Set the reply data (same as metadata flow)

![Build RESTful](RESTful26-5.png)

![Build RESTful](RESTful26-6.png)

Add query activity by select GraphBuilder_TGDB -> TGDBQuery activity

![Build RESTful](RESTful13.png)

Select “TGDB” connection that was created in Lab1 so the TGDBQuery activity queries the same server updated Northwind data was inserted into

![Build RESTful](RESTful14.png)

Map input data for TGDBQuery activity
- QueryType: $flow.queryType
- params.language: $flow.language
- params.queryString: $flow.queryString
- params.traversalCondition: $flow.traversalCondition
- params.endCondition: $flow.endCondition
- params.traversalDepth: $flow.traversalDepth

![Build RESTful](RESTful27.png)

Add “Return” activity to link the query result back to HTTP trigger

![Build RESTful](RESTful16.png)

Map outputs for Return activity
- queryResult: $activity[TGDBQuery].queryResult (map entire object)

![Build RESTful](RESTful17.png)

The TGDB_RESTful_Service is configured and it's ready for query Nothwind graph

![Build RESTful](RESTful28.png)

Test TGDB_RESTful_Service so you can see Nothwind data after querying against TGDB server

For building Flogo application
1. In project click “Build” button
2. Select the build target OS (in my case Darwin/amd64) then click to build

![Build RESTful](BuildRESTful01.png) 

Once finished you can get your executable in your browser's download folder

![Build RESTful](BuildRESTful02.png)

Find your executable and change its permission to executable then run it

![Build RESTful](Launch_RESTfulService.png)

Switch to local labs -> utilities -> lite folder
- Launch UI tool by type “npm start”
- It is required to have npm and lite-server installed before using this tool

![Build RESTful](Launch_Lite_Server.png)

Upon launching the server, the default browser will pop up and show Project GraphBuilder UI utility. For querying data against TGDB server, click “TGDB Data” tab

![Build RESTful](Launch_UI_01.png)

A query to TGDB using TGQL expression can be made as shown in screenshot below

![Build RESTful](Launch_UI_02.png)

Now, Northwind data from TGDB server can be seen