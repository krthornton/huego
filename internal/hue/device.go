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

func (c *HueConnection) FetchDevices() {
	body := c.MakeRequest(GetRequest, "/clip/v2/resource/light", nil)

	var resp map[string]interface{}
	err := json.Unmarshal(body, &resp)
	if err != nil {
		panic("Failed to unmarshal response. Panicing.")
	}

	var devices []*Device
	resources := resp["data"].([]interface{})
	for _, resc := range resources {
		resMap := resc.(map[string]interface{})
		device := MakeNewDevice(
			c,
			resMap["id"].(string),
			resMap["on"].(map[string]interface{})["on"].(bool),
			resMap["dimming"].(map[string]interface{})["brightness"].(float64),
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

func (c HueConnection) HandleDeviceEvent(event Event) {
	device := c.getDevice(event.ID)
	if device == nil {
		// event is for a device we aren't aware of
		return
	}

	// update device power state from event
	if event.On != nil {
		device.on = event.On.On
	}

	// update device brightness from event
	if event.Dimming != nil {
		device.brightness = event.Dimming.Brightness
	}
}
