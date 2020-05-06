---
title: "Gremlin"
linkTitle: "Gremlin"
weight: 5
description: >
  This extension contains activities to perform upsert operation against Gremlin Janusgraph server
---

* [Connector](https://github.com/TIBCOSoftware/labs-graphbuilder-contrib/tree/master/gremlin/connector/gremlin/)
	: A gremlin connector is the component to host your Janusgraph server connecting information. Activities which connect to same gremlin connector would connect to same Janusgraph server instance

* [JanusgraphUpsert](https://github.com/TIBCOSoftware/labs-graphbuilder-contrib/tree/master/gremlin/activity/janusgraphupsert/)
	: A JanusgraphUpsert activity consumes the graph data from BuildGraph activity and insert/update to Janusgraph server
