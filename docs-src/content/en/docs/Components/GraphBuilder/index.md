---
title: "GraphBuilder"
linkTitle: "GraphBuilder"
weight: 1
description: >
  The core extension which provides BuildGraph activity to construct graph from input data based on user predefined graph model
---

- [Graph Connector](https://github.com/TIBCOSoftware/labs-graphbuilder-contrib/blob/master/builder/connector/graph)
	: A Graph connector is the component to host your graph model. Activities which connect to same graph connector would share same graph model (data schema)

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
	: A BuildGraph activity must connect to a Graph connector so it can build the input data structure from the graph model which associated with connector. BuildGraph activity transform the input data to graph structure (nodes and edges) based on the graph model.
- [GraphToFile](https://github.com/TIBCOSoftware/labs-graphbuilder-contrib/blob/master/builder/activity/graphtofile)
	: A GraphToFile activity takes graph entities (nodes and edges) from BuildGraph and writes them to file. It's an useful utility for troubleshooting