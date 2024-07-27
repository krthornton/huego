package hue

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"regexp"
	"time"
)

type Event struct {
	Dimming *struct {
		Brightness float64 `json:"brightness"`
	} `json:"dimming"`
	ID   string `json:"id"`
	IDV1 string `json:"id_v1"`
	On   *struct {
		On bool `json:"on"`
	} `json:"on"`
	Owner *struct {
		Rid   string `json:"rid"`
		Rtype string `json:"rtype"`
	} `json:"owner"`
	Type string `json:"type"`
}

// initially created using https://mholt.github.io/json-to-go/
type EventContainer []struct {
	Creationtime time.Time `json:"creationtime"`
	Data         []Event   `json:"data"`
	ID           string    `json:"id"`
	Type         string    `json:"type"`
}

func (c HueConnection) StartEventListener() {
	// build the request
	headers := map[string]string{
		"Accept": "text/event-stream",
	}
	var request = c.buildRequest(GetRequest, "/eventstream/clip/v2", nil, headers)

	// make the request
	resp, err := c.httpClient.Do(request)
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}

	// fire up event listener to handle events
	go c.eventListener(resp.Body)
}

func (c HueConnection) eventListener(body io.ReadCloser) {
	// listen and print out responses
	var reader = bufio.NewReader(body)
	re := regexp.MustCompile("data: ?(.*)")
	for {
		line, err := reader.ReadBytes('\n')
		if err != nil {
			fmt.Println(err.Error())
			os.Exit(1)
		}
		str := string(line)
		if matches := re.FindStringSubmatch(str); len(matches) > 0 {
			var event EventContainer
			err = json.Unmarshal([]byte(matches[1]), &event)
			if err != nil {
				fmt.Println(err.Error())
				os.Exit(1)
			}
			c.eventsChannel <- event
		}
	}
}

func (c HueConnection) ProcessEvents() {
outterLoop:
	for {
		select {
		case eventContainer := <-c.eventsChannel:
			// at least one event container is waiting to be processed
			for _, container := range eventContainer {
				// process each container
				if container.Type == "update" {
					for _, event := range container.Data {
						// process each event
						if event.Type == "light" {
							c.HandleDeviceEvent(event)
						}
					}
				}
			}
		default:
			// no events are waiting to be processed; break and continue
			break outterLoop
		}
	}
}
