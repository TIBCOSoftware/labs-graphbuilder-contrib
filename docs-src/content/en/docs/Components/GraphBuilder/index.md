---
title: "GraphBuilder"
linkTitle: "GraphBuilder"
weight: 1
description: >
  The core extension which provides BuildGraph activity to construct graph based on input data coming from user’s predefined graph model
---

- [Graph Connector](https://github.com/TIBCOSoftware/labs-graphbuilder-contrib/blob/master/builder/connector/graph)
	: Graph Connector is a component which hosts your graph model for sharing graph model among graph construction related activity. Activities which connect to the same Graph connector would share same graph model (data schema)

Here is the schema of graph model

```
{
    "$schema": "http://json-schema.org/draft-04/schema#",
    "type": "object",
    "properties": {
        "nodes": {
            "type": "array",
            "items": {
                "type": "object",
                "properties": {
                    "name": {
                        "type": "string"
                    },
                    "key": {
                        "type": "array",
                        "items": {
                            "type": "string"
                        }
                    },
                    "attributes": {
                        "type": "array",
                        "items": {
                            "type": "object",
                            "properties": {
                                "name": {
                                    "type": "string"
                                },
                                "type": {
                                    "type": "string"
                                }
                            },
                            "required": [
                                "name",
                                "type"
                            ]
                        }
                    }
                },
                "required": [
                    "name",
                    "key",
                    "attributes"
                ]
            }
        },
        "edges": {
            "type": "array",
            "items": {
                "type": "object",
                "properties": {
                    "to": {
                        "type": "string"
                    },
                    "name": {
                        "type": "string"
                    },
                    "from": {
                        "type": "string"
                    },
                    "attributes": {
                        "type": "array",
                        "items": {
                            "type": "object",
                            "properties": {
                                "name": {
                                    "type": "string"
                                },
                                "type": {
                                    "type": "string"
                                }
                            }
                        }
                    }
                },
                "required": [
                    "from",
                    "name",
                    "to",
                    "attributes"
                ]
            }
        }
    }
}
```

- [BuildGraph Activity](https://github.com/TIBCOSoftware/labs-graphbuilder-contrib/blob/master/builder/activity/builder)
  : BuildGraph Activity must connect to a Graph Connector so it can build its input data schema from the graph model which is hosted in that Graph connector. BuildGraph activity transform the input data to graph entities (nodes, edges and their attributes) based on the graph model
- [GraphToFile](https://github.com/TIBCOSoftware/labs-graphbuilder-contrib/blob/master/builder/activity/graphtofile)
  : GraphToFile activity takes graph entities (nodes and edges) from BuildGraph and writes them to a file. It's a useful utility for troubleshooting