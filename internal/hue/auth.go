package hue

import (
	"encoding/json"
	"fmt"
)

type KeyManager interface {
	GetApiKey(ip string) (string, error)
	SetApiKey(ip string, apiKey string) error
}

type authRequest struct {
	DeviceType string `json:"devicetype"`
	Generate   bool   `json:"generateclientkey"`
}

type authError struct {
	Type        int    `json:"type"`
	Address     string `json:"address"`
	Description string `json:"description"`
}

type authSuccess struct {
	Username  string `json:"username"`
	ClientKey string `json:"clientkey"`
}

type authResponse []struct {
	Error   authError   `json:"error"`
	Success authSuccess `json:"success"`
}

type UnauthenticatedError struct {
	ip string
}

func (e UnauthenticatedError) Error() string {
	return fmt.Sprintf("button press required for hub with IP %s", e.ip)
}

type AuthenticationFailureError struct {
	msg string
}

func (e AuthenticationFailureError) Error() string {
	return e.msg
}

func (c *BridgeConnection) Authenticate(keyMan KeyManager) error {
	// determine if we already have an API key for the given IP
	apiKey, err := keyMan.GetApiKey(c.ipAddr)
	if err != nil {
		return AuthenticationFailureError{msg: err.Error()}
	}

	if apiKey != "" {
		// API key already exists; assume authentication successful
		c.apiKey = apiKey
		return nil
	}

	// otherwise, attempt to obtain new API key from Hue hub
	id := fmt.Sprintf("huego#%s", "REPLACE_ME")
	payload := authRequest{id, true}
	bytes, err := json.Marshal(payload)
	if err != nil {
		return AuthenticationFailureError{msg: err.Error()}
	}

	respBytes := c.makeRequest(postRequest, "/api", bytes)
	var resp authResponse
	err = json.Unmarshal(respBytes, &resp)
	if err != nil {
		return AuthenticationFailureError{msg: "failed to unmarshal auth check response"}
	}

	if len(resp) != 1 {
		return AuthenticationFailureError{msg: "unexpected response length in auth"}
	}
	authContent := resp[0]

	apiKey = authContent.Success.Username
	if apiKey == "" {
		return UnauthenticatedError{ip: c.ipAddr}
	}

	// authentication successful; save API key
	err = keyMan.SetApiKey(c.ipAddr, apiKey)
	if err != nil {
		return AuthenticationFailureError{msg: err.Error()}
	}

	c.apiKey = apiKey
	return nil
}
