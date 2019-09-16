// Package ping provides a facility for pinging hosts that uses the installed ping command
// rather than Go's ICMP support, enabling basic pinging without needing administrator/root
// privileges.
package ping

import (
	"bufio"
	"bytes"
	"io"
	"os/exec"
	"strconv"

	"github.com/getlantern/errors"
)

const (
	// DefaultCount is 1
	DefaultCount = 1

	// DefaultPayloadSize is 1472 bytes
	DefaultPayloadSize = 1472
)

// Opts specifies options for pinging
type Opts struct {
	// Count is the number of times to ping (defaults to 1)
	Count int

	// PayloadSize is the size of the data field in the ICMP packet, defaults to 1472 (for a 1500 byte total packet)
	PayloadSize int
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

// Run pings the specified host the specified number of times
func Run(host string, opts *Opts) (*Stats, error) {
	if opts == nil {
		opts = &Opts{}
	}
	if opts.Count <= 0 {
		opts.Count = DefaultCount
	}
	if opts.PayloadSize <= 0 {
		opts.PayloadSize = DefaultPayloadSize
	}

	out, err := exec.Command("ping", args(host, opts)...).Output()
	if err != nil {
		return nil, errors.New("Unable to run ping command. Output: %v: %v", string(out), err)
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
					return nil, errors.New("Unable to parse RTT %v: %v", matches[rttMinIdx], err)
				}
				stats.RTTAvg, err = strconv.ParseFloat(matches[rttAvgIdx], 64)
				if err != nil {
					return nil, errors.New("Unable to parse RTT %v: %v", matches[rttAvgIdx], err)
				}
				stats.RTTMax, err = strconv.ParseFloat(matches[rttMaxIdx], 64)
				if err != nil {
					return nil, errors.New("Unable to parse RTT %v: %v", matches[rttMaxIdx], err)
				}
			} else if matches := packetLossRegex.FindStringSubmatch(line); matches != nil {
				foundPLR = true
				stats.PLR, err = strconv.ParseFloat(matches[1], 64)
				if err != nil {
					return nil, errors.New("Unable to parse packet loss %v: %v", matches[1], err)
				}
			}
		}
		if err != nil {
			if err != io.EOF {
				return nil, errors.New("Error reading output: %v", err)
			}
			break
		}
	}

	if !foundRTT {
		return nil, errors.New("PING result did not include RTT information")
	}
	if !foundPLR {
		return nil, errors.New("PING result did not include packet loss information")
	}
	return stats, nil
}
