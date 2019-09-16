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
	// RTT Deviation milliseconds
	RTTDev float64
	// Packet Loss Rate (ratio, not percent)
	PLR float64
}

// Run pings the specified host the specified number of times
func Run(host string, opts *Opts) (*Stats, error) {
	if opts == nil {
		opts = &Opts{}
	}
	if opts.Count <= 0 {
		opts.Count = 1
	}
	if opts.PayloadSize <= 0 {
		opts.PayloadSize = 1472
	}

	out, err := exec.Command("ping", args(host, opts)...).Output()
	if err != nil {
		return nil, errors.New("Unable to run ping command. Output: %v: %v", string(out), err)
	}

	stats := &Stats{}
	reader := bufio.NewReader(bytes.NewReader(out))
	for {
		_line, _, err := reader.ReadLine()
		if _line != nil {
			line := string(_line)
			if matches := rttRegex.FindStringSubmatch(line); matches != nil {
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
				stats.RTTDev, err = strconv.ParseFloat(matches[rttDevIdx], 64)
				if err != nil {
					return nil, errors.New("Unable to parse RTT %v: %v", matches[rttDevIdx], err)
				}
			} else if matches := packetLossRegex.FindStringSubmatch(line); matches != nil {
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

	return stats, nil
}
