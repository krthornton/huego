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

type eventContainer []struct {
	Creationtime time.Time        `json:"creationtime"`
	Data         []map[string]any `json:"data"`
	ID           string           `json:"id"`
	Type         string           `json:"type"`
}

func (d *Database) startEventListener() error {
	// build the request
	headers := map[string]string{
		"Accept": "text/event-stream",
	}
	var request = d.conn.buildRequest(getRequest, "/eventstream/clip/v2", nil, headers)

	// make the request
	resp, err := d.conn.httpClient.Do(request)
	if err != nil {
		return err
	}

	// fire up event listener to handle events
	go d.eventListener(resp.Body)

	return nil
}

func (d Database) eventListener(body io.ReadCloser) {
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
			var event eventContainer
			err = json.Unmarshal([]byte(matches[1]), &event)
			if err != nil {
				panic(err)
			}
			d.eventChan <- &event
		}
	}
}

func (d Database) processEvents() {
outterLoop:
	for {
		select {
		case containers := <-d.eventChan:
			// at least one event container is waiting to be processed
			for _, container := range *containers {
				// process each container
				if container.Type == "update" {
					for _, event := range container.Data {
						// process each event within container
						for _, res := range d.resources {
							if res.Id() == event["id"] {
								res.HandleUpdate(event)
							}
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
