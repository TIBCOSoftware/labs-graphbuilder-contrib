---
title: "Lab1 - CSV"
linkTitle: "Lab1 - CSV"
weight: 1
description: >
  Build an application to populate TIBCO® Graph Database with csv data from files
---

Let's start with building a graph model for Northwind dataset. In connection tab select Graph to host graph model for your flogo application.

![Import Extension](createModel01.png)

In the diaog box 
- Set model name
- Select "Local File"
- select and upload northwind_model.json (Northwind model descriptor) from your download folder 
- Click connect

The Northwind model descriptor file has attached to your graph model 

![Import Extension](createModel02.png)

Now let's select "Apps" tab and click "Create" button to start building your first allication

![Import Extension](createApp01.png)

Name the application "Northwind" then create it

![Import Extension](createApp02.png)

Select create to build it from scratch

![Import Extension](createApp03.png)

Flogo® Enterprise studio brings you to the dialog for creating the first flow. According to Northwind dataset we have, we are going to create five flows to process data from customers.csv, suppliers.csv, employees.csv, categories.csv and products.csv respectively. We'll start from Building the customer data flow. 

![Import Extension](createApp04.png)

In the empty flow panel click "Flow Inputs & Outputs" verticle bar to generate data schema for current flow.

![Import Extension](createApp05.png)

The flow starts from processing CSV data rows from a file one line each time (we'll set it up in FileReader Trigger later). Just set a sample of useful data fields from incoming data. In the sample "FileContent" field represents a row of CSV data and "LineNumber" filelds is current "sequence number" of that row.

![Import Extension](createApp06.png)

After clik "Save" buttom the schema generator of studio converts data sample to its schema definition.

![Import Extension](createApp07.png)

Let's add a trigger (data source of the flow) by clicking "+" buttom on the left and select GraphBuilder_Tools -> FileReader trigger.

![Import Extension](createApp08.png)

Filling the "Trigger Settings"
- Filename : point to the customers.csv in your download folder
- Asynchronous : set it true so all triggers for different data could run simutanously
- Emit per Line : set it true to make sure each time only one row of data get sent to flow
- Max Number of Line : set to negative means no limit.

Click "Save" after you finish it

![Import Extension](createApp09.png)

Now switch to "Map to Flow Inputs" and make following mapping
- FileContent (defined in schema) -> $trigger.FileContent
- LineNumber (defined in schema) -> $trigger.LineNumber

Click "Save" button

![Import Extension](createApp10.png)

Back to flow and adding first activity to the flow. Select GraphBuilder_Tools -> CSVParser for converting CSV text to system object.

![Import Extension](createApp11.png)

Filling Settings
- Date Format Sample : 2006-01-02 (Data format setup for underlining GOLang code) 
- Serve Graph Data : set to false since we are not using it
- Output Field Names : One line of setting for each data column. AttributeName is the attribute name in generated system object and CSVFieldName is the column name in CSV data row. Set optional to "false" for all key element fields. Click "Save" after finish configuring each line.
- First Row is Heade : Set it true since the data file we use has a header

Click "Save" button

![Import Extension](createApp12.png)

![Import Extension](createApp14.png)

Switch to Inputs and map current input data to output data from upstream
- CSVString -> $flow.FileContent
- SequenceNumber -> $flow.LineNumber

Click "Save" when finishing it

![Import Extension](createApp13.png)

Now the data has been transform to the system object which could be recognized by the system. The next step is to convert plain object data to graph entities (nodes, edges and their attributes). We are going to use the core activity clled "BuildGraph" to perform this tranformation.

Let's select GraphBuilder -> Bruild Graph and configue it.

![Import Extension](createApp15.png)

Filling setting
- Graph Model : Select "Northwind" connector which we just created. The "Northwind" graph model now associated with this activity which means  BuildGraph activity take "Northwind" graph model to build the structure of its input data schema. You would see this when we setup "Inputs" data mapping later.
- Allow Null Key : set it "true" will make it generate node, edge even their primary key contains null elements.
- Batch Mode : set it "false" since we process one data each time.
- Pass Through Fields : (leave it empty)
- Modify Size of Instances : (leave it empty, will use it in Employee data setup)

Click "Save" button

![Import Extension](createApp16.png)

Before we can map the input data, let's a take look of the output schema of "CSVParser" (it's the upstream data for current "BuildGraph" activity). Since "CSVParser" has ability to handle multiple CSV rows, the output data structure is an array of object not just a single object.

![Import Extension](createApp16-9.png)

For procees the incoming array type of data, we need to turn on the iterator to iterate through upstream output data. Even there is only one element in the array for the current case. Following screenshot shows how to do it.

![Import Extension](createApp17.png)

For mapping the input data you may notice that 1. the "Northwind" graph model has been brought to this activity as input schema, 2. The mapping target is not to "CSVParser" anymore but the local interation. For the data coming from customers.csv you can populate more than one type of nodes which are deinfed in Northwind graph. Here is the nodes (Customer, Company and Region) which will be set.

Customer node
- _skipCondition -> null==$iteration[value].CustomerID
- CustomerID -> $iteration[value].CustomerID
- CustomerName -> $iteration[value].CustomerName
- ContactName -> $iteration[value].ContactName
- ContactTitle -> $iteration[value].ContactTitle
- City -> $iteration[value].City
- RegionName -> $iteration[value].RegionName
- RegionCode -> $iteration[value].RegionCode
- Country -> $iteration[value].Country
- Phone -> $iteration[value].Phone
- Fax -> $iteration[value].Fax

Company node
- _skipCondition -> null==$iteration[value].CompanyID
- CompanyID -> $iteration[value].CompanyID
- CompanyName -> $iteration[value].CompanyName

Region node
- RegionName -> $iteration[value].RegionName
- Country -> $iteration[value].Country

You don't have to setup mapping for edge if you don't have attribute need to be setup for them (we don't configue "label" attribute for edges now since TGDB doesn't need it). BuildGraph activity is going to use the edge definded in graph model to create edge between nodes automatically. 

![Import Extension](createApp18.png)

![Import Extension](createApp17-5.png)

![Import Extension](createApp19.png)

After we convert data to graph entities we can insert them to TIBCO® Graph Database. Let's create TIBCO® Graph Database connection first. In "Connections" tab select Add Connection -> TGDB Connector

![Import Extension](createModel03.png)

In the diaog box filling following information
- Connection name (for example "TGDB")
- TGDB Server URL
- Username
- Password
- Keep Connection Alive : select "true"

Click "Connect" button

![Import Extension](createModel04.png)

Now back to application's "Customer Data" flow to add TGDB activity. Slect GraphBuilder_TGDB -> TGDBUpsert.

![Import Extension](createApp20.png)

Filling Setting for
- TGDB connection : Select the "TGDB" Connection we just created
- Set Allow empty sting key to true (so a node with empty string key still get inserted)

Click "Save"

![Import Extension](createApp21.png)

Map input data
- Graph{} - $activity[BuildGraph].Graph{}

Since the Graph object is immutable, you are not allowed to access the detail of its internal structure. 

![Import Extension](createApp22.png)

You can insert a built-in "Log" activity by following steps:
Make room for "Log" activity by shifting activities one position to the right.

![Import Extension](createApp23.png)

Add "Log" activity by select Default -> Log

![Import Extension](createApp24.png)

Setup message for priniting (you can apply built-in fuctions to incoming data fields)
- message : string.concat(string.tostring($flow.LineNumber), " - ", $flow.FileContent)

![Import Extension](createApp25.png)

You can write the entities (which are generated by BuildGraph activity) to file by adding GraphBuilder -> GraphtoFile activity

![Import Extension](createApp26.png)

Specify the output folder and filename for GraphtoFile activity

![Import Extension](createApp27.png)

Input data is and only can be Graph. The input setup same as TGDBUpsert 

![Import Extension](createApp28.png)

Congradulations you have finish the first data flow for the application

![Import Extension](createApp29.png)

Now you can follow the same steps to finish all the rest of flows

![Import Extension](createApp30.png)

When you work on "Employee flow", please pay attention to following steps.  

In employee data there are two fields called EmployeeID and ReportTo each of them represents one indivisual employee. It implies that from the infomation of one employee data we can populate two employee nodes. One for empoyee himself/herself and the other one for his/her manager. We have to incresae the instance of employee node for such data mapping.
- Modify size of instances : Add one entry for "Employee" node and make the number of instances to 2

Click "Save"

![Import Extension](createApp31.png)

Switch to Inputs you will see two employee nodes appears (Employee0 and Employee1). Let's make employee0 the employee (not manager) so all data can be populate to this node.

![Import Extension](createApp32.png)

We make the emplyee1 node represent the manager of employee0 node so the only information we have for it (in the data) is "ReportTo" which will populate employee1's EmployeeID.

![Import Extension](createApp33.png)

Then we need to tell BuildGrap activity the relation between employee0 and employee1.

![Import Extension](createApp34.png)

Now we can test out Northwind application by sending data to it then we'll verify if data get inserted to TIBCO® Graph Database server

For building Northwind flogo application
1. In project click "Build" button
2. Select the build target OS (in my case Darwin/amd64) then click to build

![Build RESTful](BuildNorthwind_01.png) 

Once finished you can get your executable (Northwind-darwin_amd64) in your browser download folder

![Build RESTful](BuildNorthwind_02.png)

Then we need to setup a TIBCO® Graph Database. Currently Project GraphBuilder "only" support TIBCO® Graph Database 2.0.1 (both Enterprise Edition and Community Edition are supported). 
You can get a Community version from <a href="http://community.tibco.com/products/tibco-graph-database" target="_blank">here</a>.

![Build RESTful](TGDB_01.png)

Follow instructions in the download file to install TIBCO® Graph Database server then copy artifacts from your download folder
- Northwind/tgdb/northwind -> tgdb/2.0/examples
- Northwind/tgdb/init_northwind_with_data_definition.sh -> tgdb/2.0/bin/
- Northwind/tgdb/run_northwind.sh -> tgdb/2.0/bin/

![Build RESTful](TGDB_02.png)

In terminal switch to tgdb/2.0/bin folder then
- execute ./init_northwind_with_data_definition.sh to initialize tgdb with Northwind schema

![Build RESTful](TGDB_03.png)

![Build RESTful](TGDB_04.png)

- execute ./run_northwind.sh to run tgdb server

![Build RESTful](TGDB_05.png)

![Build RESTful](TGDB_06.png)

Open a new terminal and switch to the folder which contains Northwind appliction executable (Northwind-darwin_amd64). 
- Change Northwind-darwin_amd64's permission to executable 
- Run Northwind-darwin_amd64

![Build RESTful](LaunchNorthwind.png)

Open a new terminal and switch to TIBCO® Graph Database bin folder
- run tgdb-admin
- make query to get all categories

![Build RESTful](Query.png)

We've proved that data has been inserted to TIBCO® Graph Database server