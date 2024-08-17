package hue

import (
	"encoding/json"
	"fmt"
	"time"

	tea "github.com/charmbracelet/bubbletea"
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

func (c *HueConnection) checkAuthResponse(body []byte) bool {
	var resp authResponse
	err := json.Unmarshal(body, &resp)
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

func (c *HueConnection) Authenticate() tea.Msg {
	id := fmt.Sprintf("huego#%s", "REPLACE_ME")
	payload := authRequest{id, true}
	bytes, err := json.Marshal(payload)
	if err != nil {
		panic(err.Error())
	}

	respChan := c.SubmitHueRequest("POST", "/api", bytes, nil)
	resp := *<-respChan
	if c.checkAuthResponse(resp) {
		return "Success"
	}

	now := time.Now()
	timeout, _ := time.ParseDuration("1m")
	end := now.Add(timeout)

	sleepTime, _ := time.ParseDuration("2s")
	for now.Before(end) {
		time.Sleep(sleepTime)
		respChan = c.SubmitHueRequest("POST", "/api", bytes, nil)
		resp = *<-respChan
		if c.checkAuthResponse(resp) {
			return "Success"
		}
		now = time.Now()
	}

	panic("failed to authenticate")
}
