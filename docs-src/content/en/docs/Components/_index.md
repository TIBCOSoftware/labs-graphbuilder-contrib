---
title: "Components"
linkTitle: "Components"
weight: 4
description: >
  The solution is based on four core components
---

- [Graph Connector](https://github.com/TIBCOSoftware/labs-graphbuilder-contrib/blob/master/builder/connector/graph)
	: A Graph connector is a component which hosts your graph model for sharing graph model among graph construction related activity. Activities which connect to the same Graph connector would share same graph model (data schema)
- [BuildGraph Activity](https://github.com/TIBCOSoftware/labs-graphbuilder-contrib/blob/master/builder/activity/builder)
	: A BuildGraph activity must connect to a Graph connector so it can build its input data schema from the graph model which is hosted in that Graph connector. BuildGraph activity transform the input data to graph entities (nodes, edges and their attributes) based on the graph model.
- [TGDB Connector](https://github.com/TIBCOSoftware/labs-graphbuilder-contrib/blob/master/tgdb/connector/)
	: A TGDB connector is a component to store your TIBCO® Graph Database server connection information. Activities which connect to the same TGDB connector are actually connecting to the same TIBCO® Graph Database server instance
- [TGDBUpsert](https://github.com/TIBCOSoftware/labs-graphbuilder-contrib/blob/master/tgdb/activity/tgdbupsert)
	: A TGDBUpsert activity consumes the graph entities from BuildGraph activity and inserts/updates them to TIBCO® Graph Database