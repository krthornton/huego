package hue

import (
	"encoding/json"
	"errors"
	"fmt"
	"math"
	"time"
)

type Light struct {
	id                     string
	name                   string
	powerLocal             bool
	powerLocalChanged      bool
	powerRemote            bool
	brightnessLocal        int
	brightnessLocalChanged bool
	brightnessRemote       int
	lastUpdate             *time.Time
}

func newLightFromJsonMap(jsonMap map[string]any) (*Light, error) {
	id := jsonMap["id"].(string)
	name := jsonMap["metadata"].(map[string]any)["name"].(string)
	power := jsonMap["on"].(map[string]any)["on"].(bool)

	brightnessFloat := jsonMap["dimming"].(map[string]any)["brightness"].(float64)
	brightnessInt := int(math.Round(brightnessFloat))

	light := &Light{
		id:                     id,
		name:                   name,
		powerLocal:             power,
		powerLocalChanged:      false,
		powerRemote:            power,
		brightnessLocal:        brightnessInt,
		brightnessLocalChanged: false,
		brightnessRemote:       brightnessInt,
		lastUpdate:             nil,
	}

	return light, nil
}

func (l Light) Type() string { return "light" }
func (l Light) Id() string   { return l.id }
func (l Light) Name() string { return l.name }

func (l Light) GetPowerState() bool { return l.powerLocal }
func (l *Light) SetPowerState(state bool) {
	l.powerLocal = state
	l.powerLocalChanged = true
}

func (l Light) GetBrightnessLevel() int { return l.brightnessLocal }
func (l *Light) SetBrightnessLevel(level int) {
	l.brightnessLocal = level
	l.brightnessLocalChanged = true
}

func (l *Light) SubmitChanges(conn *BridgeConnection) error {
	payload := make(map[string]any)

	if l.powerLocalChanged {
		onMap := make(map[string]bool)
		onMap["on"] = l.powerLocal
		payload["on"] = onMap
	}

	if l.brightnessLocalChanged {
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

	var resp map[string]any
	err = json.Unmarshal(body, &resp)
	if err != nil {
		return errors.New("failed to unmarshal response")
	}

	l.powerLocalChanged = false
	l.brightnessLocalChanged = false

	currentTime := time.Now()
	l.lastUpdate = &currentTime

	return nil
}

func (l *Light) HandleUpdate(updateMap map[string]any) error {
	onData, ok := updateMap["on"]
	if ok {
		l.powerRemote = onData.(map[string]any)["on"].(bool)
	}

	dimmingData, ok := updateMap["dimming"]
	if ok {
		brightnessFloat := dimmingData.(map[string]any)["brightness"].(float64)
		brightnessFloat = math.Round(brightnessFloat/10) * 10
		l.brightnessRemote = int(brightnessFloat)
	}

	return nil
}

func (l *Light) Sync(delay time.Duration) error {
	if l.lastUpdate != nil {
		delta := time.Since(*l.lastUpdate)
		if delta.Seconds() < delay.Seconds() {
			// wait time has not elapsed; skip sync
			return nil
		}
	}

	l.powerLocal = l.powerRemote
	l.brightnessLocal = l.brightnessRemote
	return nil
}
