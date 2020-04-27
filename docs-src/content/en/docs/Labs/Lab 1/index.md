---
title: "Lab1 - CSV"
linkTitle: "Lab1 - CSV"
weight: 1
description: >
  Build an app to populate TIBCO Graph Database with data in csv files
---

Let's build a graph model for Northwind data. In connection tab select Graph to host graph model for your flogo application.
![Import Extension](createModel01.png)

In the diaog box 
- Set model name
- Select "Local File"
- select and upload northwind_model.json from your download folder 
- Click connect

The Northwind model descriptor file has attached to your graph model 
![Import Extension](createModel02.png)

Now let's select Application table then create to start building an allication
![Import Extension](createApp01.png)

Name the application "Northwind" then create it
![Import Extension](createApp02.png)

Select create to make it from scratch
![Import Extension](createApp03.png)

Flogo studio brings you to the dialog for creating the first flow. According to Northwind data set we have, we are going to create five flows to process data from customers.csv, suppliers.csv, employees.csv, categories.csv and products.csv respectively. We start from Building the customer data flow. 
![Import Extension](createApp04.png)

In the empty flow panel click "Flow Inputs & Outputs" verticle bar to generate data schema for current flow.
![Import Extension](createApp05.png)

The flow processes a CSV data row from file each time. Just set a sample of incoming data. In the sample "FileContent" is a row of CSV data and "LineNumber" is current "row number" of the row.
![Import Extension](createApp06.png)

After clik "Save" buttom schema generator converts data sample to schema definition.
![Import Extension](createApp07.png)

Let's add a trigger (data source) for the flow by clicking "+" buttom on the left and select GraphBuilder_Tools -> FileReader trigger.
![Import Extension](createApp08.png)

Filling the "Trigger Settings"
- Filename : point to customers.csv in your download folder
- Asychroous : make it true so all triggers for different data would run simutanously
- Emit per Line : set true to make sure each time only one row of data sending to flow
- Max Number of Line : set to negative means no limit.

Click save after you finish it
![Import Extension](createApp09.png)

Now switch to "Map to Flow Inputs" and make following mapping
- FileContent (defined in schema) -> $trigger.FileContent
- LineNumber (defined in schema) -> $trigger.LineNumber

Then click save button
![Import Extension](createApp10.png)

Back to flow and adding first activities. Select GraphBuilder_Tools -> CSVParser for converting CSV text to system object.
![Import Extension](createApp11.png)

Filling Settings
- Date Format Sample : 2006-01-02 (GOLang data format) 
- Serve Graph Data : set false since we are not using it
- Output Field Names : One line of setting for each data column. AttributeName is the field name in system object and CSVFieldName is the column name in CSV file. Set optional to "false" for all key element fields. Click save after set configuring each line.
- First Row is Heade : Set true since the data file we use has header

Click save button
![Import Extension](createApp12.png)
![Import Extension](createApp14.png)

Switch to Inputs and map current input data to upstream output data
- CSVString -> $flow.FileContent
- SequenceNumber -> $flow.LineNumber

Click save when finish it
![Import Extension](createApp13.png)

Now the data has bean transform to the object which could be recognized by the system. The next step is to convert data to graph entities (nodes, edges and attributes). We use the core activity "Build Graph" to perform this tranformation.
Let's select GraphBuilder -> Bruild Graph and configue it.

![Import Extension](createApp15.png)

Filling setting
- Graph Model : Select "Northwind" which we just created. The Northwind graph model now associated with this activity. You would see this when we setup Inputs.
- Allow Null Key : Will Generate node, edge even their primary key will null element.
- Batch Mode : Set false since we process one data each time.
- Pass Through Fields : (leave it empty)
- Modify Size of Instances : (leave it empty, will use it in Employee data setup)

Click save button
![Import Extension](createApp16.png)

Before we can map the input data let's take look of the output of CSVParser (upstream data of current Build Graph activity). Since CSVParser has ability to handle multiple CSV rows, the output of it is an array not just a single object.
![Import Extension](createApp16-9.png)

For procees the incoming array type of data, we need to turn on the iterator to iterate through upstream output data. Even there is only one element in the array for the current case. Following screenshot showing how to do it.
![Import Extension](createApp17.png)

For mapping the input data you may notice that 1. the Northwind graph model has been brought to this activity as input schema, 2. The mapping target is not to CSVParser but the local interation. For the data coming from customers.csv you can populate more than one type of nodes which are deinfed in Northwind graph. Here is the nodes data which will be set.

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

You don't have to setup mapping for edge if you don't have attribute need to be setup for them (we don't configue "label" since TGDB doesn't need it). BuildGraph activity is going to use the edge definded in graph model to create edge between nodes automatically. 

![Import Extension](createApp18.png)
![Import Extension](createApp17-5.png)
![Import Extension](createApp19.png)

After we convert data to graph entities we can insert them to TGDB. Let's create TGDB connection first. In Connections tab select Add Connection -> TGDB Connector
![Import Extension](createModel03.png)

In the diaog box filling following information
- Connection name (for example "TGDB")
- TGDB Server URL
- Username
- Password
- Keep Connection Alive : select "true"

Click connect

![Import Extension](createModel04.png)

Now back to application "Customer Data" flow to add TGDB activity. Slect GraphBuilder_TGDB -> TGDBUpsert.
![Import Extension](createApp20.png)

Filling Setting for
- TGDB connection : Select the "TGDB" Connection we just created
- Set Allow empty sting key to true (so node with empty string key still get inserted)

Click save

![Import Extension](createApp21.png)

Map input data
- Graph{} - $activity[BuildGraph].Graph{}

Since the Graph object is immutable, you are not allow to access the internal structure detail. 

![Import Extension](createApp22.png)

You can insert built in log activity by following steps:
Make room for logger activity by shifting activities one position to the right.

![Import Extension](createApp23.png)

Add log activity by select Default -> Log

![Import Extension](createApp24.png)

Setup message for priniting (you can apply built-in fuctions to incoming data fields)
- message : string.concat(string.tostring($flow.LineNumber), " - ", $flow.FileContent)

![Import Extension](createApp25.png)

You can add write the entities which generate be BuildGraph activity by adding GraphBuilder -> GraphtoFile activity

![Import Extension](createApp26.png)

Specify the output folder and filename for writing the graph entities

![Import Extension](createApp27.png)

Input data is and only can be Graph. The input setup same as TGDBUpsert 

![Import Extension](createApp28.png)

Congradulations you have finish the first data flow for the application

![Import Extension](createApp29.png)

Now you can follow the same steps to finish all the rest of flows

![Import Extension](createApp30.png)

When you work on employee flow, please pay attention on following steps.  

In employee data there are two fields EmployeeID and ReportTo which represent one indivisual employee. It implies that by the infomation in the employee data we can populate two employee nodes. One for empoyee himself/herself and one for his/her manager. We have to incresae the employee instance for the data mapping.
- Modify size of instances : Select "Employee" node and make the number of instances to 2

Click save

![Import Extension](createApp31.png)

Switch to Inputs you will see two employee nodes appears (Employee0 and Employee1). Let's make employee0 the employee (not manager) so all data can be populate to this node.

![Import Extension](createApp32.png)

We make the emplyee1 node represent the manager of employee0 node so the only information we have for it (in the data) is "ReportTo" which will populate employee1's EmployeeID.

![Import Extension](createApp33.png)

Then we need to tell BuildGrap activity the relation between employee0 and employee1.

![Import Extension](createApp34.png)


