---
title: "Airline"
linkTitle: "Airline"
weight: 1
description: >
  Using graph model to construct the relation between passenger, destination data and flight events
---

#### Implementation Source

Download application artifacts from [here](https://github.com/TIBCOSoftware/labs-graphbuilder-contrib/tree/master/sample-applications/Airline).

Download GraphBuilder user extensions from [here](https://github.com/TIBCOSoftware/labs-graphbuilder-contrib/blob/master/dist)

#### Installation

Open TIBCO Flogo® Enterprise 2.8.1 studio and upload required user extensions (builder.zip, tgdb.zip and tools.zip)

![Import Extension](user_extensions.png)

Create an empty application for Airline 

![Import Extension](create_app.png)

Import application from the pre-configured Airline application descriptor

![Import Extension](import_app.png)

Find and select descriptor (airline.json) from download folder

![Import Extension](import_app2.png)

Ignore waning just click "Import All"

![Import Extension](import_app3.png)

Airline application with two data flow (flight and pnr) is imported

![Import Extension](import_app4.png)

Check connection tab to see two connections to be fixed

![Import Extension](fix_conn.png)

Select TGDB connection and edit it

![Import Extension](fix_conn1.png)

Make configuration meet your TIBCO® Graph Database setup then click "Connect" to save it

![Import Extension](fix_conn2.png)

Same to the Graph connection but just click "Connect" since graph model has been correctly set

![Import Extension](fix_conn3.png)

> This example is created in TIBCO Flogo® Enterprise 2.8.1 studio.