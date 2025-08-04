package hue

import (
	"encoding/json"
	"fmt"
)

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

func (c *HueConnection) Authenticate() bool {
	id := fmt.Sprintf("huego#%s", "REPLACE_ME")
	payload := authRequest{id, true}
	bytes, err := json.Marshal(payload)
	if err != nil {
		panic(err.Error())
	}

	respBytes := c.MakeRequest(PostRequest, "/api", bytes)
	var resp authResponse
	err = json.Unmarshal(respBytes, &resp)
	if err != nil {
		panic("failed to unmarshal auth check response")
	}

	if len(resp) != 1 {
		panic("unexpected response length in auth")
	}
	authContent := resp[0]

	if authContent.Success.Username != "" {
		c.apiKey = authContent.Success.Username
		return true
	}

	return false
}
