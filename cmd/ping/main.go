package main

import (
	"flag"
	"fmt"
	"log"

	"github.com/getlantern/ping"
)

func main() {
	opts := &ping.Opts{}
	flag.IntVar(&opts.Count, "c", ping.DefaultCount, "number of times to ping")
	flag.IntVar(&opts.PayloadSize, "s", ping.DefaultPayloadSize, "number of bytes payload in each ICMP echo request")
	flag.Parse()

	if len(flag.Args()) == 0 {
		log.Fatal("Please specify a host")
	}

	host := flag.Args()[0]
	stats, err := ping.Run(host, opts)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("\"%v\",%v,%v,%v,%v,%v\n", host, stats.RTTMin, stats.RTTAvg, stats.RTTMax, stats.RTTDev, stats.PLR)
}
