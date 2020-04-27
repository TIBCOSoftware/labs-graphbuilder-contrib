---
title: "Lab2 - Query"
linkTitle: "Lab2 - Query"
weight: 2
description: >
  Build an app to query against TIBCO Graph Database
---

Create a new flogo application called "TGDB_RESTful_Service"

![Build RESTful](RESTful01.png)

Click create to build from scratch

![Build RESTful](RESTful02.png)

Create first flow to query metadata

![Build RESTful](RESTful03.png)

Define the data schema for flow input by pasting sample data (queryType in string data type). 
- queryType : metadata, edgetypes or nodetypes (for metadata flow) 

![Build RESTful](RESTful04.png)

Save sample data so schema builder can generate schema definition from it

![Build RESTful](RESTful05.png)

Define flow output schema by sample output data
- Content : contains the data of query result
- Success : true means query go through without error
- Code : error code
- Message : error message

Click save

![Build RESTful](RESTful06.png)

Save trigger schema definition generation

![Build RESTful](RESTful07.png)

Adding trigger to receive HTTP request by clicking "+" -> "ReceiveHTTPMessage"

![Build RESTful](RESTful07-5.png)

Select GET, setup resource path "/tgdb/{queryType}" then click continue

![Build RESTful](RESTful09.png)

Click "Just add the trigger" button

![Build RESTful](RESTful09-5.png)

We have a trigger with HTTP GET methods and listen on port 9999)

![Build RESTful](RESTful08.png)

Click trigger to map incoming query data to flow input data 

![Build RESTful](RESTful10.png)

In "reply Settings" set reply schema make it same as flow output data schema

![Build RESTful](RESTful11.png)

In "Map from flow outputs" mapping data.queryResult to  $flow.queryResult

![Build RESTful](RESTful12.png)

Add query activity by select GraphBuilder_TGDB -> TGDBQuery activity

![Build RESTful](RESTful13.png)

Select "TGDB" coinnection we created in lab1 so the TGDBQuery activity is going to query against the server which we've upserted Northwind data to 

![Build RESTful](RESTful14.png)

Map input data for TGDBQuery activity
- QueryType : $flow.queryType

![Build RESTful](RESTful15.png)

Adding return activity to link the query result back to HTTP trigger

![Build RESTful](RESTful16.png)

Map outputs for Return activity 
- queryResult : $activity[TGDBQuery].queryResult (map entire object)

![Build RESTful](RESTful17.png)

You've finished creating metadata query flow

![Build RESTful](RESTful18.png)

Click "Create" button to create another flow for querying data content

![Build RESTful](RESTful19.png)

Create name and description for the new flow

![Build RESTful](RESTful20.png)

Define the flow inputs data schema by sample data (schema detail see TGDB documentation)
- queryType : search (for content flow) 
- language : TGQL (TIBCO graph query language) or Gremlin
- queryString : for TGQL and Gremlin
- traversalCondition : TGQL only
- traversalDepth : TGQL only

![Build RESTful](RESTful21.png)

Click save to generate data schema definition

![Build RESTful](RESTful22.png)

Flow output data schema same as metadata flow.

![Build RESTful](RESTful22-5.png)

Add another trigger for receiving content query

![Build RESTful](RESTful23.png)

POST method for content query 

![Build RESTful](RESTful24.png)

Adding sample query for the output (to the flow) setting. To be noticed that the schema is very similar to flow input schema but grouped under "query" keyword. 

![Build RESTful](RESTful25.png)

Map to flow input
- queryType : $trigger.pathParams.queryType
- language : $trigger.body.query.language
- queryString : $trigger.body.query.queryString
- traversalCondition : $trigger.body.query.traversalCondition
- endCondition : $trigger.body.query.endCondition
- traversalDepth : $trigger.body.query.traversalDepth

Click save

![Build RESTful](RESTful26.png)

Setting the reply data (same as metadata flow)
![Build RESTful](RESTful26-5.png)

![Build RESTful](RESTful26-6.png)

Add query activity by select GraphBuilder_TGDB -> TGDBQuery activity

![Build RESTful](RESTful13.png)

Select "TGDB" coinnection we created in lab1 so the TGDBQuery activity is going to query against the server which we've upserted Northwind data to 

![Build RESTful](RESTful14.png)

Map input data for TGDBQuery activity
- QueryType : $flow.queryType
- params.language : $flow.language
- params.queryString : $flow.queryString
- params.traversalCondition : $flow.traversalCondition
- params.endCondition : $flow.endCondition
- params.traversalDepth : $flow.traversalDepth

![Build RESTful](RESTful27.png)

Adding return activity to link the query result back to HTTP trigger

![Build RESTful](RESTful16.png)

Map outputs for Return activity 
- queryResult : $activity[TGDBQuery].queryResult (map entire object)

![Build RESTful](RESTful17.png)

The TGDB_RESTful_Service is ready for query Nothwind graph

![Build RESTful](RESTful28.png)