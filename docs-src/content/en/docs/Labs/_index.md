---
title: "Labs"
linkTitle: "Labs"
weight: 3
description: >
  Build an application to insert/update real-time data to TIBCO® Graph Database 
---

The following section provides step-by-step, hands-on exercises that show how to build a Flogo application by parsing data coming in CSV files and inserting it to TIBCO® Graph Database.

The Labs leverage the Northwind dataset (sample dataset used by Microsoft to demonstrate the features of their relational database). The exercises illustrate how GraphBuilder can be used to convert relational data into graph and then insert it into TIBCO® Graph Database.

- Project GraphBuilder and the artifacts needed for the exercises can be downloaded <a href="https://github.com/TIBCOSoftware/labs-graphbuilder-contrib" target="_blank">here</a> 
- Click <a href="https://github.com/TIBCOSoftware/labs-graphbuilder-contrib/blob/master/dist" target="_blank">here</a> to download GraphBuilder user extensions 
- Click  <a href="https://github.com/TIBCOSoftware/labs-graphbuilder-contrib/blob/master/sample-applications/Northwind/" target="_blank">here</a> to download the project data, graph model and TGDB configuration
- Click <a href="https://github.com/TIBCOSoftware/labs-graphbuilder-contrib/blob/master/sample-applications/utilities/" target="_blank">here</a> to download the GUI utilities 

The Labs use TIBCO Flogo® Enterprise studio to configure the applications. It is required to have it locally installed before starting building the application. Click <a href="https://edelivery.tibco.com/storefront/en/eval/tibco-flogo-enterprise/prod11810.html" target="_blank">here</a> to download TIBCO Flogo® Enterprise studio

![Import Extension](upload01.png)

After the installation of TIBCO Flogo® Enterprise studio has been completed, import all required user extensions files (builder.zip, tgdb.zip, tools.zip and sse.zip) as shown in the image below

- In “Extensions” tab click “Upload” button
- Click “From a Zip file”
- Select one user extension (for example builder.zip) from your download folder at a time
- Click “Upload and compiling”

![Import Extension](upload02.png)

Click “Done” when extensions are uploaded and compiled

![Import Extension](upload03.png)

Uploaded extension will be display on left panel

![Import Extension](upload04.png)

Keep uploading all other required extensions. Here are required user extensions
- GraphBuilder
- GraphBuilder_TGDB
- GraphBuilder_Tools
- GraphBuilder_SSE

![Import Extension](upload05.png)

This completes the set up for the Labs 

