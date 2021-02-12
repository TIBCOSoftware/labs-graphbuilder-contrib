/*
 * Copyright © 2020. TIBCO Software Inc.
 * This file is subject to the license terms contained
 * in the license file that is distributed with this file.
 */
package httpclient

import (
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"
)

type HttpClient struct {
	url         string
	resource    string
	accessToken string
	parameters  map[string]interface{}
	headers     map[string]interface{}
	response    *http.Response
}

func (this *HttpClient) FetchData() (string, error) {

	url := this.url + this.resource

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return "", errors.New(fmt.Sprintf("Error create request : %v ", err))
	}

	query := req.URL.Query()
	for para, value := range this.parameters {
		query.Add(para, value.(string))
	}
	req.URL.RawQuery = query.Encode()

	for header, value := range this.headers {
		req.Header.Set(header, value.(string))
	}

	client := &http.Client{Timeout: time.Second * 10}

	this.response, err = client.Do(req)
	if err != nil {
		return "", errors.New(fmt.Sprintf("Error reading response : %v ", err))
	}
	defer this.response.Body.Close()

	if this.response.StatusCode != 200 {
		this.response.Body.Close()
		return "", fmt.Errorf("got response status code %d\n", this.response.StatusCode)
	}

	data, err := ioutil.ReadAll(this.response.Body)
	return string(data), err
}
