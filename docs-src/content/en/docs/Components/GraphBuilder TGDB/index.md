---
title: "TGDB"
linkTitle: "TGDB"
weight: 2
description: >
  This extension contains activities to perform CRUD operations against TIBCO® Graph Database
---

* [TGDB Connector](https://github.com/TIBCOSoftware/labs-graphbuilder-contrib/blob/master/tgdb/connector/)
	: A TGDB connector is the component to host your TGDB server connecting information. Activities which connect to same TGDB connector would connect to same TGDB server instance
* [TGDBUpsert](https://github.com/TIBCOSoftware/labs-graphbuilder-contrib/blob/master/tgdb/activity/tgdbupsert)
	: A TGDBUpsert activity consumes the graph data from BuildGraph activity and insert/update to TGDB.
* [TGDBQuery](https://github.com/TIBCOSoftware/labs-graphbuilder-contrib/tree/master/tgdb/activity/tgdbquery/)
	: With TGDBQuery activity users can build their own application to query against TGDB. It support both TGQL and Gremlin query language
* [TGDBDelete](https://github.com/TIBCOSoftware/labs-graphbuilder-contrib/tree/master/tgdb/activity/tgdbdelete/)
	: TGDBDelete activity implemeting graph entities deletion for TGDB. It takes graph entities from BuildGraph then performs deletion of them
