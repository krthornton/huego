package hue

import (
	"bytes"
	"container/list"
	"fmt"
	"io"
	"net/http"
	"os"
	"time"
)

type hueRequest struct {
	httpRequest     *http.Request
	reqType         string
	responseChannel chan *[]byte
	submissionTime  *time.Time
}

func (c HueConnection) SubmitHueRequest(reqType string, reqMethod string, path string, payload []byte, headers map[string]string) chan *[]byte {
	httpRequest := c.buildHttpRequest(reqMethod, path, payload, headers)
	respChan := make(chan *[]byte, 1)
	hueRequest := &hueRequest{
		httpRequest:     httpRequest,
		reqType:         reqType,
		responseChannel: respChan,
	}
	c.requestChannel <- hueRequest

	return hueRequest.responseChannel
}

func (c HueConnection) StartRequestHandler() {
	go c.requestHandler()
}

func (c HueConnection) requestHandler() {
	requests := list.New()

	for {
		// receive next request from channel
		request := <-c.requestChannel

		// check  that we have not hit the throttle limit
		timeElapsed := time.Duration(0)
		currentTime := time.Now()
		oneSecond := time.Duration(1 * time.Second)
		for reqCount, item := 1, requests.Back(); timeElapsed < oneSecond && item != nil; reqCount, item = reqCount+1, item.Prev() {
			prevReq, ok := item.Value.(*hueRequest)
			if !ok {
				fmt.Println("Unexpected, non-request object in request stack.")
				os.Exit(1)
			}

			if reqCount == 10 {
				// threshold reached; pause and wait before submitting next request
				time.Sleep(timeElapsed - oneSecond)
				break
			} else {
				timeElapsed += currentTime.Sub(*prevReq.submissionTime)
			}
		}

		// push new request onto stack for tracking and remove stale ones
		requests.PushBack(request)
		// TODO: remove stale items

		// make the request and log submission time
		resp, err := c.httpClient.Do(request.httpRequest)
		if err != nil {
			fmt.Println(err.Error())
			os.Exit(1)
		}
		request.submissionTime = &currentTime

		// capture response and return as bytes
		body, err := io.ReadAll(resp.Body)
		if err != nil {
			fmt.Println(err.Error())
			os.Exit(1)
		}

		// return response to requester
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
