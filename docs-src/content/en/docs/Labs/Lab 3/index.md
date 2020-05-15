---
title: "Lab3 - real-time"
linkTitle: "Lab3 - real-time"
weight: 3
description: >
  Build an application to insert/update real-time data to TIBCO® Graph Database
---

Now let's create an application which receives real-time order events from Kafka topic then build graph entities (nodes, edges and their attributes), insert/update entities to TGDB then serve real-time graph entities as a streaming server.

First of all create an internal "Server Sent Event (SSE)" connection to link between order event flow and SSE server flow (for serving streaming graph entities to external client).

In "Connections" tab select GraphBuilder_SSE -> Server-sent Events Connection

![Realtime](realtime02.png)

Connection settings (following settings match the client tool which is provided for browsing real-time graph entity update)
- Connection Name : Set name to "EventServer"
- Outbound : Set false as it's a server
- Server port : 8888
- Path : It's URI path "/sse/"
- TLS enabled : false

Click "Connect"

![Realtime](realtime01.png)

Back to Northwind application to create a new flow called "Order Event Server"

![Realtime](realtime07.png)

Select a "SSE Server" trigger to serve graph entities (come from order event flow) for streaming client

![Realtime](realtime03.png)

Settings
- Connection Name : Select the "EventServer" connection which we just created

Click "Save"

![Realtime](realtime04.png)

This simple flow will be serving streaming graph entities

![Realtime](realtime04-5.png)

Now Add the last flow for Northwind application. It is "Order Data Flow" which listen to Kafla topic to consume order events as input data of the flow. 

Before we create it, we need to create a "Kafka Connection". In connection tab select "Appach Kafka Client Configuration"  

![Realtime](realtime06.png)

Configure Appach Kafka Client as following screenshot then save it

![Realtime](realtime05.png)

Back to application to create new flow called "Order Event"

![Realtime](realtime08.png)

Click "Flow Inputs & Outputs" (vertical blue bar) to define schema between flow and trigger. Set following data sample then click save

![Realtime](realtime12.png)

After clicking save button schema generator would convert sample data to schema definition like before

![Realtime](realtime13.png)

Click "+" to add trigger (Kafka Consumer) 

![Realtime](realtime09.png)

Select "Northwind Orders" configuration we just created then click continue.

![Realtime](realtime10.png)

Select "Just Add Trigger" button to add trigger

![Realtime](realtime11.png)

Finish the trigger setting as it shown below in screenshot

![Realtime](realtime14.png)

Map OrderString to $trigger.stringValue

![Realtime](realtime15.png)

Add CSVParser to convert incoming CVS string to system object

![Realtime](realtime16.png)

Follow the instruction in lab1 define the mapping between CSV fields and attribute of system object. Use the column field name as attribute name.

Make sure "First Row Is Header" set to false

![Realtime](realtime17.png)

Configure the input
- CSVString : $flow.OrderString
- Leave SequenceNumber not mapped

![Realtime](realtime18.png)

After the data has bean transform to the object which could be recognized by the system. The next step is to convert data to graph entities (nodes, edges and their attributes). We use the core activity "Build Graph" to perform this tranformation.
Let's select GraphBuilder -> Bruild Graph and configue it.

![Realtime](realtime18-5.png)

Follow lab1's instruction to turn on the "iterator" for iterating through upstream output data (at runtime) then map with input data of BuildGraph activity. Here is the mapping

Product node
- ProductID -> $iteration[value].ProductID

Employee node
- EmployeeID -> $iteration[value].EmployeeID

Customer node
- CustomerID -> $iteration[value].CustomerID

Order node
- OrderID -> $iteration[value].OrderID
- CustomerID -> $iteration[value].CustomerID
- EmployeeID- > $iteration[value].EmployeeID
- OrderDate -> $iteration[value].OrderDate
- RequiredDate -> $iteration[value].RequiredDate
- ShippedDate -> $iteration[value].ShippedDate
- ShipVia -> $iteration[value].ShipVia
- Freight -> $iteration[value].Freight
- ShipName -> $iteration[value].ShipName
- ShipAddress -> $iteration[value].ShipAddress
- ShipCity -> $iteration[value].ShipCity
- ShipRegion -> $iteration[value].ShipRegion
- ShipPostalCode - > $iteration[value].ShipPostalCode
- ShipCountry -> $iteration[value].ShipCountry

Suborder node
- OrderID -> $iteration[value].OrderID
- ProductID -> $iteration[value].ProductID
- UnitPrice -> $iteration[value].UnitPrice
- Quantity -> $iteration[value].Quantity
- Discount -> $iteration[value].Discount

Region node
- RegionName -> $iteration[value].RegionName
- Country -> $iteration[value].Country

Since one order can be splited to multiple order events (with different product sold). We create two types of order nodes 1. Odrer node (main order) with OrderID as its primary key and 2. Suborder node with OrderID, ProductID as primary key. BuildGraph activity would link (via edge) All Suborder nodes to Order node by matching their the OrderID. (See following screenshot)

Order : 

![Realtime](realtime19.png)

Suborder : 

![Realtime](realtime20.png)

Follow lab1's instruction to add TGDBUpsert activity

![Realtime](realtime20-2.png)

Select Connetion

![Realtime](realtime20-3.png)

Map input data

![Realtime](realtime20-4.png)

Now adding a new type of activity called SSEEndPoint which sends graph entities to SSEServer for serving streaming client.

Select SSEEndPoint activity from GraphBuilder_SSE.

![Realtime](realtime21.png)

Select "SSEConnection" we created and used in SSEServer earlier then the new SSEEndPoint is connected to SSEServer now.

![Realtime](realtime22.png)

Setup SessionId to "order" so the complete URI to access to this event flow would be /sse/order

![Realtime](realtime23.png)

Map input data to Graph object from BuildGraph activity

![Realtime](realtime24.png)

We can add log and GraphtoFile activities like previous configured flows.

![Realtime](realtime25.png)

Now we have finish the last flow for our Northwind application.

![Realtime](realtime26.png)

This is the final version of flogo Northwind application 

![Realtime](realtime27.png)

Let's rebuild application for further testing

![Realtime](BuildNorthwind_02.png)

We are going to install Kafka message bus for providing order event. <a href="https://kafka.apache.org/quickstart" target="_blank">Here</a> is the instalation instructions.

After downloading and installing Kafka we can start Kafka
- Start zoo keeper

![Realtime](StartKafka01.png)
![Realtime](StartKafka02.png)

- Start server

![Realtime](StartKafka03.png)

- Create "test" topic

![Realtime](StartKafka04.png)

Let's restart Northwind appliction executable. 
- Switch to the folder which contains Northwind appliction executable (Northwind-darwin_amd64). 
- Change Northwind-darwin_amd64's permission to executable 
- Run Northwind-darwin_amd64

This time you'll see two extra information while Northwind application starting
- Kafka consumer (the trigger of order event flow) is up and listening
- SSEServer is up and waiting for client (UI utility) to connect 

![Realtime](StartKafka05.png)

Here it the our test (see screenshot)
- Make sure TGDB, TGDB_RESTful_Service, Kafka (server, zoo keeper, producer) and UI utility are running 
- On the upper/middle left of screenshot open oerders.csv file 
- On the lower left of screenshot start Kafka producer and keep it opened
- On the right follow the instruction to 1. Click "Realtime Data" 2. Click "Connect" to connect to SSE server in Northwind application 3. Copy & paste order to Kafka producer then press enter 4 ~ 6. Each time you send one order you will see the corresponding graph entities showing on the UI.

Send more order as your wish. 

![Realtime](FinalTest.png)

After the streaming testing we also want to see the order in TGDB. Follow the instrctions in lab2
- Click "TGDB Data" button
- Use the default query setup but make traversalDepth = 5
- Click "Make Query" button

You'll see the oder with OrderID = 10248 and its associated graph entities on the UI

![Realtime](QueryGraph.png)

The last test is about traversal query. We are going to find all companies which supply prodcuts within the order from the company 'Vins et alcools Chevalier'. We are going to use Postman and TGDB_RESTful_Service to query against TGDB server
- Open a postman and setup a POST query
- The gremlin query is "g.V().has('Company', 'CompanyID', 'Vins et alcools Chevalier').in('Customer_Company').in('SoldTo').out('Includes').out('Contains').in('Supplies').out('Supplier_Company');"

You should get "Formaggi Fortini s.r.l.", "Leka Trading" and "Cooperativa de Quesos 'Las Cabras'" in your result

![Realtime](QueryPostMan.png)

Observe the traversal request on the UI utility and verify the correctness of the query

![Realtime](QueryPostManOnGraph.png)

Congradulations! Now you've finished all three labs