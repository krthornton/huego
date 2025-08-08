package hue

import (
	"encoding/json"
	"errors"
)

type Resource interface {
	Type() string
	Id() string
	Name() string
	SubmitChanges(conn *BridgeConnection) error
	HandleUpdate(updateMap map[string]any) error
}

func (d *Database) fetchResources() error {
	body := d.conn.makeRequest(getRequest, "/clip/v2/resource", nil)

	var resp map[string]interface{}
	err := json.Unmarshal(body, &resp)
	if err != nil {
		return errors.New("failed to unmarshal /resource response")
	}

	resources := resp["data"].([]interface{})
	for _, res := range resources {
		resMap := res.(map[string]interface{})
		switch resType := resMap["type"].(string); resType {
		case "light":
			light, err := newLightFromJsonMap(resMap)
			if err != nil {
				return err
			}
			d.resources = append(d.resources, light)
		}
	}

	return nil
}
