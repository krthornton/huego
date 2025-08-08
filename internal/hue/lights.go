package hue

import (
	"encoding/json"
	"errors"
	"fmt"
	"math"
)

type Light struct {
	id               string
	name             string
	powerLocal       bool
	powerRemote      bool
	brightnessLocal  int
	brightnessRemote int
}

func newLightFromJsonMap(jsonMap map[string]interface{}) (*Light, error) {
	id := jsonMap["id"].(string)
	name := jsonMap["metadata"].(map[string]interface{})["name"].(string)
	power := jsonMap["on"].(map[string]interface{})["on"].(bool)

	brightnessFloat := jsonMap["dimming"].(map[string]interface{})["brightness"].(float64)
	brightnessInt := int(math.Round(brightnessFloat))

	light := &Light{
		id:               id,
		name:             name,
		powerLocal:       power,
		powerRemote:      power,
		brightnessLocal:  brightnessInt,
		brightnessRemote: brightnessInt,
	}

	return light, nil
}

func (l Light) Type() string { return "light" }
func (l Light) Id() string   { return l.id }
func (l Light) Name() string { return l.name }

func (l Light) GetPowerState() bool       { return l.powerLocal }
func (l *Light) SetPowerState(state bool) { l.powerLocal = state }

func (l Light) GetBrightnessLevel() int       { return l.brightnessLocal }
func (l *Light) SetBrightnessLevel(level int) { l.brightnessLocal = level }

func (l *Light) SubmitChanges(conn *BridgeConnection) error {
	payload := make(map[string]any)

	if l.powerLocal != l.powerRemote {
		onMap := make(map[string]bool)
		onMap["on"] = l.powerLocal
		payload["on"] = onMap
	}

	if l.brightnessLocal != l.brightnessRemote {
		dimmingMap := make(map[string]float32)
		dimmingMap["brightness"] = float32(l.brightnessLocal)
		payload["dimming"] = dimmingMap
	}

	if len(payload) == 0 {
		// no changes to submit
		return nil
	}

	// submit changes
	bytes, err := json.Marshal(payload)
	if err != nil {
		panic(err.Error())
	}
	url := fmt.Sprintf("/clip/v2/resource/light/%s", l.id)
	body := conn.makeRequest(putRequest, url, bytes)

	var resp map[string]interface{}
	err = json.Unmarshal(body, &resp)
	if err != nil {
		return errors.New("failed to unmarshal response")
	}

	l.powerRemote = l.powerLocal
	l.brightnessRemote = l.brightnessLocal

	return nil
}

func (l *Light) HandleUpdate(updateMap map[string]any) error {
	onData, ok := updateMap["on"]
	if ok {
		l.powerLocal = onData.(map[string]any)["on"].(bool)
		l.powerRemote = l.powerLocal
	}

	dimmingData, ok := updateMap["dimming"]
	if ok {
		brightnessFloat := dimmingData.(map[string]any)["brightness"].(float64)
		brightnessFloat = math.Round(brightnessFloat/10) * 10
		l.brightnessLocal = int(brightnessFloat)
		l.brightnessRemote = l.brightnessLocal
	}

	return nil
}
