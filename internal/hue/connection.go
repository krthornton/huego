package hue

import (
	"crypto/tls"
	"net/http"
	"time"
)

type HueConnection struct {
	httpClient    *http.Client
	httpTransport *http.Transport
	tlsConfig     *tls.Config
	eventsChannel chan EventContainer
	requestTimer  *time.Timer
	devices       *[]*Device
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
	devices := make([]*Device, 0)

	// make a buffered event listener channel to accept multiple before block
	eventsChannel := make(chan EventContainer, 25)

	conn := &HueConnection{
		httpClient:    httpClient,
		httpTransport: httpTransport,
		tlsConfig:     tlsConfig,
		eventsChannel: eventsChannel,
		devices:       &devices,
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
