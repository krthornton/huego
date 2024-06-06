package hue

import (
	"bytes"
	"crypto/tls"
	"fmt"
	"io"
	"net/http"
	"os"
)

type RequestType int

const (
	GetRequest  RequestType = 0
	PostRequest RequestType = 1
	PutRequest  RequestType = 3
)

func getRequestTypeString(reqType RequestType) string {
	switch reqType {
	case PostRequest:
		return "POST"
	case PutRequest:
		return "PUT"
	default:
		return "GET"
	}
}

type HueConnection struct {
	httpClient    *http.Client
	httpTransport *http.Transport
	tlsConfig     *tls.Config
	ipAddr        string
	apiKey        string
}

func NewHueConnection() *HueConnection {
	tlsConfig := &tls.Config{InsecureSkipVerify: true}
	httpTransport := &http.Transport{
		TLSClientConfig: tlsConfig,
	}
	httpClient := &http.Client{
		Transport: httpTransport,
	}

	conn := &HueConnection{
		httpClient:    httpClient,
		httpTransport: httpTransport,
		tlsConfig:     tlsConfig,
	}

	return conn
}

func (c HueConnection) GetIpAddress() string {
	return c.ipAddr
}

func (c *HueConnection) SetIpAddress(ipAddr string) {
	c.ipAddr = ipAddr
}

func (c *HueConnection) SetApiKey(apiKey string) {
	c.apiKey = apiKey
}

func (c HueConnection) MakeRequest(reqType RequestType, path string, payload []byte) []byte {
	// create url from provided path
	url := fmt.Sprintf("https://%s%s", c.ipAddr, path)

	// build body of request (if provided)
	var buf bytes.Buffer
	if payload != nil {
		buf = *bytes.NewBuffer(payload)
	}

	// build request
	reqString := getRequestTypeString(reqType)
	request, err := http.NewRequest(reqString, url, &buf)
	if err != nil {
		panic(err.Error())
	}
	request.Header.Add("hue-application-key", c.apiKey)
	if payload != nil {
		request.Header.Add("Content-Type", "application/json")
	}

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
