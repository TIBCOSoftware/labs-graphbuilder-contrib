---
title: "Components"
linkTitle: "Components"
weight: 4
description: >
  The solution is based on four core main components: Graph Connector, BuildGraph Activity, TGDB Connector and TGDBUpsert. This section explains in detail each one of them.
---

- [Graph Connector](https://github.com/TIBCOSoftware/labs-graphbuilder-contrib/blob/master/builder/connector/graph): Graph connector is a component which hosts a graph model for sharing throughout graph construction related activities. The activities that connect to the same Graph connector share the same graph model (data schema)
- [BuildGraph Activity](https://github.com/TIBCOSoftware/labs-graphbuilder-contrib/blob/master/builder/activity/builder): A BuildGraph Activity must connect to a Graph Connector to build the input data schema from the graph model which is hosted in the Graph Connector. BuildGraph Activity transforms the input data into graph entities (nodes, edges and their attributes) based on the graph model.
- [TGDB Connector](https://github.com/TIBCOSoftware/labs-graphbuilder-contrib/blob/master/tgdb/connector/): A TGDB Connector is the component that stores TIBCO® Graph Database server connection information. Activities which connect to the same TGDB Connector are connecting to the same TIBCO® Graph Database server instance
- [TGDBUpsert](https://github.com/TIBCOSoftware/labs-graphbuilder-contrib/blob/master/tgdb/activity/tgdbupsert): A TGDBUpsert activity consumes the graph entities from a BuildGraph Activity and inserts/updates them into TIBCO® Graph Database
