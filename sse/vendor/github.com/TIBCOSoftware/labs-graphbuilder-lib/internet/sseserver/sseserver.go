/*
 * Copyright © 2020. TIBCO Software Inc.
 * This file is subject to the license terms contained
 * in the license file that is distributed with this file.
 */
package sseserver

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strings"
	"sync"
)

const (
	ServerPort           = "port"
	ConnectionPath       = "path"
	ConnectionTlsEnabled = "tlsEnabled"
	ConnectionUploadCRT  = "uploadCRT"
	ConnectionTlsCRTPath = "tlsCRTPath"
	ConnectionTlsKeyPath = "tlsKeyPath"
	ConnectionTlsCRT     = "tlsCRT"
	ConnectionTlsKey     = "tlsKey"
)

var (
	instance *SSEServerFactory
	once     sync.Once
)

type SSERequestListener interface {
	ProcessRequest(request string) error
}

type SSEServerFactory struct {
	sseServers map[string]*Server
	mux        sync.Mutex
}

func GetFactory() *SSEServerFactory {
	once.Do(func() {
		instance = &SSEServerFactory{sseServers: make(map[string]*Server)}
	})
	return instance
}

func (this *SSEServerFactory) GetServer(serverId string) *Server {
	return this.sseServers[serverId]
}

func (this *SSEServerFactory) CreateServer(
	serverId string,
	properties map[string]interface{},
	listener SSERequestListener) (*Server, error) {

	this.mux.Lock()
	defer this.mux.Unlock()
	server := this.sseServers[serverId]

	if nil == server {
		host := "0.0.0.0"
		if nil != properties["host"] {
			host = properties["host"].(string)
		}

		port := properties[ServerPort].(string)

		path := "/"
		if nil != properties[ConnectionPath] {
			path = properties[ConnectionPath].(string)
		}

		tls := properties[ConnectionTlsEnabled].(bool)
		tlsUploaded := properties[ConnectionUploadCRT].(bool)
		crtLoc := ""
		keyLoc := ""
		if tls {
			crtLoc = properties[ConnectionTlsCRTPath].(string)
			keyLoc = properties[ConnectionTlsKeyPath].(string)
			if tlsUploaded {
				if nil != properties[ConnectionTlsCRT] {
					err := ioutil.WriteFile(crtLoc, properties[ConnectionTlsCRT].([]byte), 0666)
					if err != nil {
						log.Fatal(err)
						return nil, err
					}
				}

				if nil != properties[ConnectionTlsKey] {
					err := ioutil.WriteFile(keyLoc, properties[ConnectionTlsKey].([]byte), 0666)
					if err != nil {
						log.Fatal(err)
						return nil, err
					}
				}
			}
		}

		server = &Server{
			host:     host,
			port:     port,
			path:     path,
			tls:      tls,
			crtLoc:   crtLoc,
			keyLoc:   keyLoc,
			clients:  make(map[string](map[string]*Client)),
			listener: listener,
		}
		this.sseServers[serverId] = server

	}
	return server, nil
}

type Server struct {
	host         string
	port         string
	path         string
	tls          bool
	crtLoc       string
	keyLoc       string
	dataChannels map[string](chan chan []byte)
	clients      map[string](map[string]*Client)
	listener     SSERequestListener
}

func (this *Server) Start() {
	fmt.Println("Start server, path : ", this.path, ", Server : ", this)
	http.Handle(this.path, this)
	if !this.tls {
		http.ListenAndServe(fmt.Sprintf("%s:%s", this.host, this.port), nil)
	} else {
		crt, err := ioutil.ReadFile(this.crtLoc)
		if err != nil {
			log.Fatal(err)
		}
		key, err := ioutil.ReadFile(this.keyLoc)
		if err != nil {
			log.Fatal(err)
		}

		fmt.Println("Got crt from file : ", crt)
		fmt.Println("Got key from file : ", key)
		http.ListenAndServeTLS(fmt.Sprintf("%s:%s", this.host, this.port), this.crtLoc, this.keyLoc, nil)
	}
}

func (this *Server) Stop() {
}

func (this *Server) RegisterClient(streamId string, client *Client) {
	clientsByQuery := this.clients[streamId]
	if nil == clientsByQuery {
		clientsByQuery = make(map[string]*Client)
		this.clients[streamId] = clientsByQuery
	}
	clientsByQuery[client.GetID()] = client
}

func (this *Server) UnRegisterClient(streamId string, client *Client) {
	clientsByQuery := this.clients[streamId]
	if nil != clientsByQuery {
		delete(clientsByQuery, client.GetID())
	}
}

func (this *Server) SendData(streamId string, data []byte) {
	/* Clients which subscribe to all streams */
	clients := this.clients["*"]
	if nil != clients {
		for _, client := range clients {
			client.dataChannel <- data
		}
	}

	clients = this.clients[streamId]
	if nil != clients {
		for _, client := range clients {
			client.dataChannel <- data
		}
	}
}

func (this *Server) ServeHTTP(res http.ResponseWriter, req *http.Request) {
	pos := strings.LastIndex(req.URL.Path, "/")
	path := req.URL.Path[pos+1:]
	id := req.RemoteAddr
	fmt.Println("Client request in, id : ", id, ", request URI : ", path)
	client := NewClient(id)
	this.RegisterClient(path, client)
	fmt.Println("After registered, clients : ", this.clients)
	client.Listening(res)
	fmt.Println("Client quit, id : ", id, ", request URI : ", path)
	this.UnRegisterClient(path, client)
	fmt.Println("After client unregistered, clients : ", this.clients)
}

type Client struct {
	id          string
	dataChannel chan []byte
}

func (this *Client) GetID() string {
	return this.id
}

func (this *Client) Listening(
	res http.ResponseWriter,
) {
	flusher, ok := res.(http.Flusher)

	if !ok {
		http.Error(res, "Streaming unsupported!", http.StatusInternalServerError)
		return
	}

	res.Header().Set("Content-Type", "text/event-stream")
	res.Header().Set("Cache-Control", "no-cache")
	res.Header().Set("Connection", "keep-alive")
	res.Header().Set("Access-Control-Allow-Origin", "*")

	connectionCallback, ok := res.(http.CloseNotifier)
	if !ok {
		http.Error(res, "cannot stream", http.StatusInternalServerError)
		return
	}

	for {
		select {
		case <-connectionCallback.CloseNotify():
			log.Println("done: closed connection")
			return
		case data := <-this.dataChannel:
			fmt.Fprintf(res, "data: %s\n\n", string(data))
			flusher.Flush()
		}
	}
}

func NewClient(id string) *Client {
	client := &Client{
		id:          id,
		dataChannel: make(chan []byte),
	}
	return client
}
