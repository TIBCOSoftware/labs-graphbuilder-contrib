---
title: "Components"
linkTitle: "Components"
weight: 4
description: >
  The solution is based on four core components
---

- [Graph Connector](https://github.com/TIBCOSoftware/labs-graphbuilder-contrib/blob/master/builder/connector/graph)
	: A Graph connector is the component to host your graph model. Activities which connect to same Graph connector would share same graph model (data schema)
- [BuildGraph Activity](https://github.com/TIBCOSoftware/labs-graphbuilder-contrib/blob/master/builder/activity/builder)
	: A BuildGraph activity must connect to a Graph connector so it can build the input data structure from the graph model which associated with connector. BuildGraph activity transform the input data to graph structure (nodes and edges) based on the graph model.
- [TGDB Connector](https://github.com/TIBCOSoftware/labs-graphbuilder-contrib/blob/master/tgdb/connector/)
	: A TGDB connector is a component to store your TGDB server connecting information. Activities which connect to same TGDB connector are actually connect to the same TGDB server instance
- [TGDBUpsert](https://github.com/TIBCOSoftware/labs-graphbuilder-contrib/blob/master/tgdb/activity/tgdbupsert)
	: A TGDBUpsert activity consumes the graph data from BuildGraph activity and insert/update to TGDB.