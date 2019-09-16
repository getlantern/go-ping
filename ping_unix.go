// +build !windows

package ping

import (
	"fmt"
	"regexp"
)

var (
	rttRegex        = regexp.MustCompile(`min/avg/max/.+dev = ([0-9\.]+)/([0-9\.]+)/([0-9\.]+)`)
	rttMinIdx       = 1
	rttAvgIdx       = 2
	rttMaxIdx       = 3
	packetLossRegex = regexp.MustCompile(`, ([0-9\.]+)% packet loss`)
)

func args(host string, opts *Opts) []string {
	return []string{
		"-q",
		"-c", fmt.Sprintf("%d", opts.Count),
		"-s", fmt.Sprintf("%d", opts.PayloadSize),
		host,
	}
}
