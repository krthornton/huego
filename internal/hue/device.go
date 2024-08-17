package hue

import (
	"encoding/json"
	"fmt"
	"math"
)

type Device struct {
	conn       *HueConnection
	id         string
	on         bool
	brightness int
	name       string
}

func MakeNewDevice(conn *HueConnection, id string, on bool, brightness int, name string) *Device {
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

func (d Device) Brightness() int {
	return d.brightness
}

func (l *Device) ChangePowerState(desiredState bool) {
	if l.conn.requestTimer != nil {
		return
	} else {
		l.conn.StartRequestTimer()
	}

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
	body := l.conn.MakeRequest("PUT", url, bytes)

	var resp map[string]interface{}
	err = json.Unmarshal(body, &resp)
	if err != nil {
		panic("Failed to unmarshal response. Panicing.")
	}
}

func (d *Device) ChangeBrightness(desiredBrightness float64) {
	if d.conn.requestTimer != nil {
		return
	} else {
		d.conn.StartRequestTimer()
	}

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
	body := d.conn.MakeRequest("PUT", url, bytes)

	var resp map[string]interface{}
	err = json.Unmarshal(body, &resp)
	if err != nil {
		panic("Failed to unmarshal response. Panicing.")
	}
}

func (c *HueConnection) FetchDevices() {
	body := c.MakeRequest("GET", "/clip/v2/resource/light", nil)

	var resp map[string]interface{}
	err := json.Unmarshal(body, &resp)
	if err != nil {
		panic("Failed to unmarshal response. Panicing.")
	}

	var devices []*Device
	resources := resp["data"].([]interface{})
	for _, resc := range resources {
		resMap := resc.(map[string]interface{})
		brightnessFloat := resMap["dimming"].(map[string]interface{})["brightness"].(float64)
		brightnessInt := int(math.Round(brightnessFloat))
		device := MakeNewDevice(
			c,
			resMap["id"].(string),
			resMap["on"].(map[string]interface{})["on"].(bool),
			brightnessInt,
			resMap["metadata"].(map[string]interface{})["name"].(string),
		)
		devices = append(devices, device)
	}

	c.devices = &devices
}

func (c HueConnection) GetDevices() []*Device {
	return *c.devices
}

func (c HueConnection) GetDevice(index int) *Device {
	return (*c.devices)[index]
}

func (c HueConnection) getDevice(id string) *Device {
	for _, device := range *c.devices {
		if device.id == id {
			return device
		}
	}

	return nil
}
