// Package ping provides a facility for pinging hosts that uses the installed ping command
// rather than Go's ICMP support, enabling basic pinging without needing administrator/root
// privileges.
package ping

import (
	"fmt"
)

const (
	// DefaultCount is 1
	DefaultCount = 1

	// DefaultPayloadSize is 56 bytes
	DefaultPayloadSize = 56
)

// Opts specifies options for pinging
type Opts struct {
	// Count is the number of times to ping (defaults to 1)
	Count int

	// PayloadSize is the size of the data field in the ICMP packet, defaults to 1472 (for a 1500 byte total packet)
	PayloadSize int
}

func (opts *Opts) withDefaults() *Opts {
	if opts == nil {
		opts = &Opts{}
	}
	if opts.Count <= 0 {
		opts.Count = DefaultCount
	}
	if opts.PayloadSize <= 0 {
		opts.PayloadSize = DefaultPayloadSize
	}
	return opts
}

// Stats represents the statistics of pinging a host
type Stats struct {
	// RTT Minimum milliseconds
	RTTMin float64
	// RTT Average milliseconds
	RTTAvg float64
	// RTT Maximum milliseconds
	RTTMax float64
	// Packet Loss Rate (ratio, not percent)
	PLR float64
}

func (s *Stats) String() string {
	return fmt.Sprintf("rtt min/avg/max (%vms / %vms / %vms)    packet loss rate (%v%%)", s.RTTMin, s.RTTAvg, s.RTTMax, s.PLR)
}
