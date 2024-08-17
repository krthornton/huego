package hue

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
	"os"
	"time"
)

type hueRequest struct {
	httpRequest     *http.Request
	responseChannel chan *[]byte
}

func (c HueConnection) SubmitHueRequest(reqType string, path string, payload []byte, headers map[string]string) chan *[]byte {
	httpRequest := c.buildHttpRequest(reqType, path, payload, headers)
	respChan := make(chan *[]byte, 1)
	hueRequest := &hueRequest{
		httpRequest:     httpRequest,
		responseChannel: respChan,
	}
	c.requestChannel <- hueRequest

	return hueRequest.responseChannel
}

func (c HueConnection) StartRequestHandler() {
	go c.requestHandler()
}

func (c HueConnection) requestHandler() {
	for {
		// receive next request from channel
		request := <-c.requestChannel

		// make the request
		resp, err := c.httpClient.Do(request.httpRequest)
		if err != nil {
			fmt.Println(err.Error())
			os.Exit(1)
		}

		// capture response and return as bytes
		body, err := io.ReadAll(resp.Body)
		if err != nil {
			fmt.Println(err.Error())
			os.Exit(1)
		}

		// return response to requestee
		request.responseChannel <- &body
		close(request.responseChannel)
	}
}

func (c HueConnection) buildHttpRequest(reqType string, path string, payload []byte, headers map[string]string) *http.Request {
	// create url from provided path
	url := fmt.Sprintf("https://%s%s", c.ipAddr, path)

	// build body of request (if provided)
	var buf bytes.Buffer
	if payload != nil {
		buf = *bytes.NewBuffer(payload)
	}

	// build request
	request, err := http.NewRequest(reqType, url, &buf)
	if err != nil {
		panic(err.Error())
	}
	request.Header.Add("hue-application-key", c.apiKey)
	for key, val := range headers {
		request.Header.Add(key, val)
	}

	return request
}

func (c *HueConnection) StartRequestTimer() {
	c.requestTimer = time.AfterFunc(100*time.Millisecond, func() {
		// remove timer lock after timer expires
		c.requestTimer = nil
	})
}
