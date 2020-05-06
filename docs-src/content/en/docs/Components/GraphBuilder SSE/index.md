---
title: "SSE"
linkTitle: "SSE"
weight: 6
description: >
  SSE user extension contains activities for implemeting Server Sent Event (HTTP based streaming event) client and server. 
---

* [Connector](https://github.com/TIBCOSoftware/labs-graphbuilder-contrib/tree/master/sse/connector/sse/)
	: A SSE connector is the component to store your sse server (Outbound = false) or connection (Outbound = true) information. Activities connect to same SSE connector are connecting to same SSE service.

* [SSESubscriber](https://github.com/TIBCOSoftware/labs-graphbuilder-contrib/tree/master/sse/trigger/ssesub/)
	: A SSESubscriber trigger subscribes to remote sse server then consumes streaming events. The SSE Connector for a subscriber need to be configured as Outbound = true.

* [SSEServer](https://github.com/TIBCOSoftware/labs-graphbuilder-contrib/tree/master/sse/trigger/sseserver/)
	: A SSEServer trigger works as an SSE server which serves streamming events. It maitains the incoming connection and requests but generate any data. The streaming data comes from another activity called SSEEndPoint. The SSE Connector for a subscriber need to be configured as Outbound = false
	
* [SSEEndpoint](https://github.com/TIBCOSoftware/labs-graphbuilder-contrib/tree/master/sse/activity/sseendpoint/)
	: A SSEEndpoint activity sits on different flow from SSEServer. It takes input event and streams it to SSEServer. The link between a SSEServer and a SSEEndPoint is that both of them need to connect to the same SSE connector (Outbound = false)
