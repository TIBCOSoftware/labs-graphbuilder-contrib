---
title: "Neo4j"
linkTitle: "Neo4j"
weight: 4
description: >
  Extension that contains activities to perform upsert operation against Neo4j database
---

* [Connector](https://github.com/TIBCOSoftware/labs-graphbuilder-contrib/tree/master/neo4j/connector/neo4j/)
	: A Neo4j connector is the component to store your Neo4j server connection information. Activities which connect to the same Neo4j connector would connect to the same Neo4j server instance
* [Neo4jUpsert](https://github.com/TIBCOSoftware/labs-graphbuilder-contrib/tree/master/neo4j/activity/neo4jupsert/)
	: A Neo4jUpsert activity consumes the graph entities from BuildGraph activity and inserts/updates them to Neo4j server
