package hue

import (
	"encoding/json"
	"fmt"
)

type Device struct {
	conn       *HueConnection
	id         string
	on         bool
	brightness float64
	name       string
}

func MakeNewDevice(conn *HueConnection, id string, on bool, brightness float64, name string) *Device {
	return &Device{conn, id, on, brightness, name}
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

func (d Device) Brightness() float64 {
	return d.brightness
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

func (d *Device) ChangeBrightness(desiredBrightness float64) {
	payload := struct {
		Dimming struct {
			Brightness float64 `json:"brightness"`
		} `json:"dimming"`
	}{
		Dimming: struct {
			Brightness float64 `json:"brightness"`
		}{Brightness: desiredBrightness},
	}
	bytes, err := json.Marshal(payload)
	if err != nil {
		panic(err.Error())
	}

	url := fmt.Sprintf("/clip/v2/resource/light/%s", d.id)
	body := d.conn.MakeRequest(PutRequest, url, bytes)

	var resp map[string]interface{}
	err = json.Unmarshal(body, &resp)
	if err != nil {
		panic("Failed to unmarshal response. Panicing.")
	}

	// assuming no error occurred, tthe brightness has successfully been changed
	d.brightness = desiredBrightness
}
