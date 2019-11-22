// +build !windows

package ping

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"os/exec"
	"regexp"
	"strconv"

	"github.com/getlantern/errors"
)

var (
	rttRegex        = regexp.MustCompile(`min/avg/max/.+dev = ([0-9\.]+)/([0-9\.]+)/([0-9\.]+)`)
	rttMinIdx       = 1
	rttAvgIdx       = 2
	rttMaxIdx       = 3
	packetLossRegex = regexp.MustCompile(`, ([0-9\.]+)% packet loss`)
)

// Run pings the specified host with the given options
func Run(host string, opts *Opts) (*Stats, error) {
	opts = opts.withDefaults()

	out, err := exec.Command("ping", args(host, opts)...).Output()
	if err != nil {
		return nil, errors.New("unable to run ping command. Output: %v: %v", string(out), err)
	}

	stats := &Stats{}
	reader := bufio.NewReader(bytes.NewReader(out))
	foundRTT := false
	foundPLR := false
	for {
		_line, _, err := reader.ReadLine()
		if _line != nil {
			line := string(_line)
			if matches := rttRegex.FindStringSubmatch(line); matches != nil {
				foundRTT = true
				stats.RTTMin, err = strconv.ParseFloat(matches[rttMinIdx], 64)
				if err != nil {
					return nil, errors.New("unable to parse RTT %v: %v", matches[rttMinIdx], err)
				}
				stats.RTTAvg, err = strconv.ParseFloat(matches[rttAvgIdx], 64)
				if err != nil {
					return nil, errors.New("unable to parse RTT %v: %v", matches[rttAvgIdx], err)
				}
				stats.RTTMax, err = strconv.ParseFloat(matches[rttMaxIdx], 64)
				if err != nil {
					return nil, errors.New("unable to parse RTT %v: %v", matches[rttMaxIdx], err)
				}
			} else if matches := packetLossRegex.FindStringSubmatch(line); matches != nil {
				foundPLR = true
				stats.PLR, err = strconv.ParseFloat(matches[1], 64)
				if err != nil {
					return nil, errors.New("unable to parse packet loss %v: %v", matches[1], err)
				}
			}
		}
		if err != nil {
			if err != io.EOF {
				return nil, errors.New("error reading output: %v", err)
			}
			break
		}
	}

	if !foundRTT {
		return nil, errors.New("ping result did not include RTT information")
	}
	if !foundPLR {
		return nil, errors.New("ping result did not include packet loss information")
	}
	return stats, nil
}

func args(host string, opts *Opts) []string {
	return []string{
		"-q",
		"-c", fmt.Sprintf("%d", opts.Count),
		"-s", fmt.Sprintf("%d", opts.PayloadSize),
		host,
	}
}
