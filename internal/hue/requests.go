package hue

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
	"os"
	"time"
)

type HueRequestType int

type HueRequest struct {
	RequestType HueRequestType
	LimitId     string
	HttpRequest *http.Request
}

func (c HueConnection) requestHandler() {

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

func (c HueConnection) MakeRequest(reqType string, path string, payload []byte) []byte {
	// build request
	headers := map[string]string{
		"Content-Type": "application/json",
	}
	var request = c.buildHttpRequest(reqType, path, payload, headers)

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

func (c *HueConnection) StartRequestTimer() {
	c.requestTimer = time.AfterFunc(100*time.Millisecond, func() {
		// remove timer lock after timer expires
		c.requestTimer = nil
	})
}
