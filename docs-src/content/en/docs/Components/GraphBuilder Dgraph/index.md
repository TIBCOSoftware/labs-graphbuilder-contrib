---
title: "Dgraph"
linkTitle: "Dgraph"
weight: 3
description: >
  This extension contains activities to perform upsert operation against Dgraph database
---

* [Connector](https://github.com/TIBCOSoftware/labs-graphbuilder-contrib/tree/master/dgraph/connector/dgraph/)
	: A Dgraph connector is the component to host your Dgraph server connecting information. Activities which connect to same Dgraph connector would connect to same Dgraph server instance
* [DgraphUpsert](https://github.com/TIBCOSoftware/labs-graphbuilder-contrib/tree/master/dgraph/activity/dgraphupsert/)
	: A DgraphUpsert activity consumes the graph data from BuildGraph activity and insert/update to Dgraph server
