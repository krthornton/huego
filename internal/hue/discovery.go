package hue

import (
	"io"
	"log"
	"net"

	"github.com/hashicorp/mdns"
)

func DiscoverHueBridges() chan *net.IP {
	// discover available interfaces
	infs, err := net.Interfaces()
	if err != nil {
		panic("Failed to retrieve available interfaces.")
	}

	// start query routines for each interface
	ipCh := make(chan *net.IP, 10)
	for _, inf := range infs {
		// ensure inf object outlives loop
		inf := inf

		go func() {
			// start a routine to receive detected Hubs
			entriesCh := make(chan *mdns.ServiceEntry)
			go func() {
				for entry := range entriesCh {
					ipCh <- &entry.AddrV4
				}
			}()

			// build out query
			params := mdns.DefaultParams("_hue._tcp")
			params.Interface = &inf
			params.DisableIPv6 = true
			params.Entries = entriesCh
			logger := &log.Logger{}
			logger.SetOutput(io.Discard)
			params.Logger = logger

			// start the search for Hubs
			mdns.Query(params)
		}()
	}

	return ipCh
}
