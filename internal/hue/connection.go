package hue

import (
	"crypto/tls"
	"net/http"
	"time"
)

type HueConnection struct {
	httpClient     *http.Client
	httpTransport  *http.Transport
	tlsConfig      *tls.Config
	requestChannel chan *hueRequest
	requestTimer   *time.Timer
	devices        *[]*Device
	ipAddr         string
	apiKey         string
}

func NewHueConnection() *HueConnection {
	tlsConfig := &tls.Config{InsecureSkipVerify: true}
	httpTransport := &http.Transport{
		TLSClientConfig: tlsConfig,
	}
	httpClient := &http.Client{
		Transport: httpTransport,
	}
	devices := make([]*Device, 0)

	requestChannel := make(chan *hueRequest, 10)

	conn := &HueConnection{
		httpClient:     httpClient,
		httpTransport:  httpTransport,
		requestChannel: requestChannel,
		tlsConfig:      tlsConfig,
		devices:        &devices,
	}

	return conn
}

func (c HueConnection) GetIpAddress() string {
	return c.ipAddr
}

func (c HueConnection) GetApiKey() string {
	return c.apiKey
}

func (c *HueConnection) SetIpAddress(ipAddr string) {
	c.ipAddr = ipAddr
}

func (c *HueConnection) SetApiKey(apiKey string) {
	c.apiKey = apiKey
}
