package config

import (
	"encoding/json"
	"huego/internal/hue"
)

type DeviceDb struct {
	conn    *hue.HueConnection
	devices *[]*hue.Device
}

func MakeNewDeviceDb(conn *hue.HueConnection) *DeviceDb {
	devices := make([]*hue.Device, 0)

	return &DeviceDb{
		conn:    conn,
		devices: &devices,
	}
}

func (d *DeviceDb) FetchDevices() {
	body := d.conn.MakeRequest(hue.GetRequest, "/clip/v2/resource/light", nil)

	var resp map[string]interface{}
	err := json.Unmarshal(body, &resp)
	if err != nil {
		panic("Failed to unmarshal response. Panicing.")
	}

	var devices []*hue.Device
	resources := resp["data"].([]interface{})
	for _, resc := range resources {
		resMap := resc.(map[string]interface{})
		device := hue.MakeNewDevice(
			d.conn,
			resMap["id"].(string),
			resMap["on"].(map[string]interface{})["on"].(bool),
			resMap["dimming"].(map[string]interface{})["brightness"].(float64),
			resMap["metadata"].(map[string]interface{})["name"].(string),
		)
		devices = append(devices, device)
	}

	d.devices = &devices
}

func (d DeviceDb) GetDevices() []*hue.Device {
	return *d.devices
}

func (d DeviceDb) GetDevice(index int) *hue.Device {
	return (*d.devices)[index]
}
