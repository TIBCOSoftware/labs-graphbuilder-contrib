---
title: "Gremlin"
linkTitle: "Gremlin"
weight: 5
description: >
  This extension contains activities to perform upsert operation against Gremlin Janusgraph server
---

* [Connector](https://github.com/TIBCOSoftware/labs-graphbuilder-contrib/tree/master/gremlin/connector/gremlin/)
	: A gremlin connector is the component to store your Janusgraph server connection information. Activities which connect to the same gremlin connector would connect to the same Janusgraph server instance

* [JanusgraphUpsert](https://github.com/TIBCOSoftware/labs-graphbuilder-contrib/tree/master/gremlin/activity/janusgraphupsert/)
	: A JanusgraphUpsert activity consumes the graph entities from BuildGraph activity and inserts/updates them to Janusgraph server
