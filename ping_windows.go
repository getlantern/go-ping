package ping

import (
	"fmt"
	"regexp"
)

var (
	rttRegex        = regexp.MustCompile(`.+Minimum = ([0-9\.]+)ms, Maximum = ([0-9\.]+)ms, Average = ([0-9\.]+)ms.*`)
	rttMinIdx       = 1
	rttAvgIdx       = 3
	rttMaxIdx       = 2
	packetLossRegex = regexp.MustCompile(`.+\(([0-9\.]+)% loss\).*`)
)

func args(host string, opts *Opts) []string {
	return []string{
		"/n", fmt.Sprintf("%d", opts.Count),
		"/l", fmt.Sprintf("%d", opts.PayloadSize),
		host,
	}
}
