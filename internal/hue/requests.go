package hue

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
	"os"
)

type requestMethod int

const (
	getRequest  requestMethod = 0
	postRequest requestMethod = 1
	putRequest  requestMethod = 3
)

func getRequestTypeString(method requestMethod) string {
	switch method {
	case postRequest:
		return "POST"
	case putRequest:
		return "PUT"
	default:
		return "GET"
	}
}

func (c BridgeConnection) buildRequest(method requestMethod, path string, payload []byte, headers map[string]string) *http.Request {
	// create url from provided path
	url := fmt.Sprintf("https://%s%s", c.ipAddr, path)

	// build body of request (if provided)
	var buf bytes.Buffer
	if payload != nil {
		buf = *bytes.NewBuffer(payload)
	}

	// build request
	reqString := getRequestTypeString(method)
	request, err := http.NewRequest(reqString, url, &buf)
	if err != nil {
		panic(err.Error())
	}
	request.Header.Add("hue-application-key", c.apiKey)
	for key, val := range headers {
		request.Header.Add(key, val)
	}

	return request
}

func (c BridgeConnection) makeRequest(method requestMethod, path string, payload []byte) []byte {
	// build request
	headers := map[string]string{
		"Content-Type": "application/json",
	}
	var request = c.buildRequest(method, path, payload, headers)

	// make the request
	resp, err := c.httpClient.Do(request)
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}
	defer resp.Body.Close()

	// capture response and return as bytes
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}

	return body
}
