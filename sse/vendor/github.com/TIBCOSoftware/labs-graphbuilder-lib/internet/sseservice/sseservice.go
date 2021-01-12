/*
 * Copyright © 2020. TIBCO Software Inc.
 * This file is subject to the license terms contained
 * in the license file that is distributed with this file.
 */
package sse

import (
	"bufio"
	"fmt"
	"net/http"
)

type SSEService interface {
	SetEventListener(listener SSEListener)
	Start() error
	Stop() error
}

type SSEListener interface {
	ProcessEvent(event string) error
}

//type UpdateFunc func(current []byte) (updated []byte, err error)

type SSEServiceFactory struct {
}

func NewSSEServiceFactory() SSEServiceFactory {
	return SSEServiceFactory{}
}

func (this SSEServiceFactory) GetService(properties map[string]interface{}) SSEService {
	sseService := SSEServiceImpl{}
	sseService.url = properties["url"].(string)
	sseService.resource = properties["resource"].(string)
	sseService.accessToken = properties["accessToken"].(string)
	return &sseService
}

type SSEServiceImpl struct {
	url         string
	resource    string
	accessToken string
	response    *http.Response
	listener    SSEListener
}

func (this *SSEServiceImpl) SetEventListener(listener SSEListener) {
	this.listener = listener
}

func (this *SSEServiceImpl) Start() error {

	url := this.url + this.resource
	if "" != this.accessToken {
		url = url + "?access_token=" + this.accessToken
	}

	res, err := http.Get(url)
	this.response = res

	if err != nil {
		return err
	}

	if res.StatusCode != 200 {
		return fmt.Errorf("got response status code %d\n", res.StatusCode)
	}

	reader := bufio.NewReader(res.Body)

	for {
		line, isPrefix, err := reader.ReadLine()

		if nil != err {
			return err
		}

		if nil == line {
			//break
			fmt.Println("\n\n\n********** nil line ************\n\n\n")
		}

		if isPrefix {
			continue
		}

		this.listener.ProcessEvent(string(line))
	}
	return nil
}

func (this *SSEServiceImpl) Stop() error {
	closed := this.response.Close
	if !closed {
		return fmt.Errorf("SSE service is unable to be closed!")
	}
	return nil
}
