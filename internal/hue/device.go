package hue

import (
	"encoding/json"
	"fmt"
)

type Device struct {
	conn *HueConnection
	id   string
	on   bool
	name string
}

func MakeNewDevice(conn *HueConnection, id string, on bool, name string) *Device {
	return &Device{conn, id, on, name}
}

func (l Device) Id() string {
	return l.id
}

func (l Device) IsPoweredOn() bool {
	return l.on
}

func (l Device) Name() string {
	return l.name
}

func (l *Device) ChangePowerState(desiredState bool) {
	payload := struct {
		On struct {
			On bool `json:"on"`
		} `json:"on"`
	}{
		On: struct {
			On bool `json:"on"`
		}{On: desiredState},
	}
	bytes, err := json.Marshal(payload)
	if err != nil {
		panic(err.Error())
	}

	url := fmt.Sprintf("/clip/v2/resource/light/%s", l.id)
	body := l.conn.MakeRequest(PutRequest, url, bytes)

	var resp map[string]interface{}
	err = json.Unmarshal(body, &resp)
	if err != nil {
		panic("Failed to unmarshal response. Panicing.")
	}

	// assuming no error occurred, the power state has successfully been changed
	l.on = desiredState
}
