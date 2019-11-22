package ping

import (
	"time"

	"github.com/sparrc/go-ping"
)

// Run pings the specified host with the given options
func Run(host string, opts *Opts) (*Stats, error) {
	opts = opts.withDefaults()

	pinger, err := ping.NewPinger(host)
	if err != nil {
		return nil, err
	}
	pinger.SetPrivileged(true)

	pinger.Count = opts.Count
	pinger.Size = opts.PayloadSize
	// On UNIX, the ping command won't wait forever for each packet. This keeps Windows from waiting forever as well.
	// If the RTT is higher than 10 seconds then there's worse problems anyway.
	pinger.Timeout = 10 * time.Duration(opts.Count) * time.Second
	pinger.Run()
	stats := pinger.Statistics()

	return &Stats{
		RTTMin: millis(stats.MinRtt),
		RTTAvg: millis(stats.AvgRtt),
		RTTMax: millis(stats.MaxRtt),
		PLR:    stats.PacketLoss,
	}, nil
}

func millis(d time.Duration) float64 {
	return float64(d) / 1000000
}
