package hue

import (
	"io"
	"log"

	"github.com/hashicorp/mdns"
)

func DiscoverIpAddress() string {
	// Make a channel for results and start listening
	entriesCh := make(chan *mdns.ServiceEntry, 4)
	ipAddrCh := make(chan string, 1)
	go func() {
		for entry := range entriesCh {
			ipAddrCh <- entry.AddrV4.String()
		}
	}()

	// Start the lookup using custom query param to disable verbose logging
	log.SetOutput(io.Discard)
	mdns.Lookup("_hue._tcp", entriesCh)

	// wait for discovery of local hue bridge
	ipAddr := <-ipAddrCh

	// close search by shutting channel
	close(entriesCh)

	return ipAddr
}
