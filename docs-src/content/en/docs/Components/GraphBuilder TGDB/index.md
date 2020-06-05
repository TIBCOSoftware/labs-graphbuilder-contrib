---
title: "TGDB"
linkTitle: "TGDB"
weight: 2
description: >
  Extension that contains activities that perform CRUD operations against TIBCO® Graph Database
---

* [TGDB Connector](https://github.com/TIBCOSoftware/labs-graphbuilder-contrib/blob/master/tgdb/connector/)
	:  A TGDB connector is a component to store TIBCO® Graph Database server connection information. Activities which connect to the same TGDB connector are actually connecting to the same TIBCO® Graph Database server instance
* [TGDBUpsert](https://github.com/TIBCOSoftware/labs-graphbuilder-contrib/blob/master/tgdb/activity/tgdbupsert)
	: A TGDBUpsert activity consumes the graph entities from BuildGraph activity and inserts/updates them to TIBCO® Graph Database
* [TGDBQuery](https://github.com/TIBCOSoftware/labs-graphbuilder-contrib/tree/master/tgdb/activity/tgdbquery/)
	: With TGDBQuery activity users can build their own application to query against TIBCO® Graph Database. It supports both TGQL and Gremlin query language
* [TGDBDelete](https://github.com/TIBCOSoftware/labs-graphbuilder-contrib/tree/master/tgdb/activity/tgdbdelete/)
	: TGDBDelete activity implements the deletion of graph entities for TIBCO® Graph Database. It takes graph entities (with primary key attributes populated) from BuildGraph then performs the deletion on them
