---
title: "Labs"
linkTitle: "Labs"
weight: 3
description: >
  In following hands on labs I will step by step guide you through the process of building an Flogo application which transforms csv data then inserts/updates data to TIBCO® Graph Database. 
---

We are going to use the Northwind dataset to create an Flogo application. The Northwind data is a sample dataset used by Microsoft to demonstrate the features of Microsoft's relational database. We will demonstrate how to use GraphBuilder to convert relational data to graph then insert into TIBCO® Graph Database.

You can check out or download Project GraphBuilder from <a href="https://github.com/TIBCOSoftware/labs-graphbuilder-contrib" target="_blank">here</a> and the artifacts which you need for your labs project

- Download GraphBuilders user extensions <a href="https://github.com/TIBCOSoftware/labs-graphbuilder-contrib/blob/master/dist" target="_blank">here</a>
- Download project data, graph model and TGDB configuration <a href="https://github.com/TIBCOSoftware/labs-graphbuilder-contrib/blob/master/sample-applications/Northwind/" target="_blank">here</a>
- Sownload GUI utilities <a href="https://github.com/TIBCOSoftware/labs-graphbuilder-contrib/blob/master/sample-applications/utilities/" target="_blank">here</a>

In the labs we use TIBCO Flogo® Enterprise studio to configure the lab applications. You need to have it installed before you can start building the application. You can get TIBCO Flogo® Enterprise studio from <a href="https://edelivery.tibco.com/storefront/en/eval/tibco-flogo-enterprise/prod11810.html" target="_blank">here</a>

![Import Extension](upload01.png)

After installed studio import all required user extensions files (builder.zip, tgdb.zip, tools.zip and sse.zip)
1. In "Extensions" tab click "Upload" button
2. Click "From a Zip file"
3. Each time select one user extension (for example builder.zip) from your download folder
4. Click "Upload and compiling"

![Import Extension](upload02.png)

Click "Done" when extension get uploaded and compiled

![Import Extension](upload03.png)

Uploaded extension will be display on left panel

![Import Extension](upload04.png)

Keep uploading all other required extensions. Here is all four required user extensions
- GraphBuilder
- GraphBuilder_TGDB
- GraphBuilder_Tools
- GraphBuilder_SSE

![Import Extension](upload05.png)

Now you are good to go for the upcoming labs

