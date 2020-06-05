---
title: "Dgraph"
linkTitle: "Dgraph"
weight: 3
description: >
  Extension that contains activities to perform upsert operations against Dgraph database
---

* [Connector](https://github.com/TIBCOSoftware/labs-graphbuilder-contrib/tree/master/dgraph/connector/dgraph/)
	: A Dgraph connector is a component to store your Dgraph server connection information. Activities which connect to the same Dgraph connector would connect to the same Dgraph server instance
* [DgraphUpsert](https://github.com/TIBCOSoftware/labs-graphbuilder-contrib/tree/master/dgraph/activity/dgraphupsert/)
	:  A DgraphUpsert activity consumes the graph entities from BuildGraph activity and inserts/updates them to Dgraph server
