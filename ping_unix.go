package ping

import (
	"fmt"
	"regexp"
)

var (
	rttRegex        = regexp.MustCompile(`.+min/avg/max/.+dev = ([0-9\.]+)/([0-9\.]+)/([0-9\.]+)/([0-9\.]+).+`)
	rttMinIdx       = 1
	rttAvgIdx       = 2
	rttMaxIdx       = 3
	rttDevIdx       = 4
	packetLossRegex = regexp.MustCompile(`.+, ([0-9\.]+)% packet loss.*`)
)

func args(host string, opts *Opts) []string {
	return []string{
		host, "-q",
		fmt.Sprintf("-c %d", opts.Count),
		fmt.Sprintf("-s %d", opts.PayloadSize),
	}
}
