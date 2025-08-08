package hue

import (
	"crypto/tls"
	"net/http"
)

type BridgeConnection struct {
	httpClient    *http.Client
	httpTransport *http.Transport
	tlsConfig     *tls.Config
	ipAddr        string
	apiKey        string
}

func NewBridgeConnection() *BridgeConnection {
	tlsConfig := &tls.Config{InsecureSkipVerify: true}
	httpTransport := &http.Transport{
		TLSClientConfig: tlsConfig,
	}
	httpClient := &http.Client{
		Transport: httpTransport,
	}

	conn := &BridgeConnection{
		httpClient:    httpClient,
		httpTransport: httpTransport,
		tlsConfig:     tlsConfig,
	}

	return conn
}

func (c BridgeConnection) GetIpAddress() string {
	return c.ipAddr
}

func (c BridgeConnection) GetApiKey() string {
	return c.apiKey
}

func (c *BridgeConnection) SetIpAddress(ipAddr string) {
	c.ipAddr = ipAddr
}

func (c *BridgeConnection) SetApiKey(apiKey string) {
	c.apiKey = apiKey
}
