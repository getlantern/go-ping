package main

import (
	"flag"
	"fmt"
	"log"

	"github.com/getlantern/go-ping"
)

func main() {
	opts := &ping.Opts{}
	flag.IntVar(&opts.Count, "c", ping.DefaultCount, "number of times to ping")
	flag.IntVar(&opts.PayloadSize, "s", ping.DefaultPayloadSize, "number of bytes payload in each ICMP echo request")
	porcelain := flag.Bool("porcelain", false, "print result in machine-readable format")
	flag.Parse()

	if len(flag.Args()) == 0 {
		log.Fatal("Please specify a host")
	}

	host := flag.Args()[0]
	stats, err := ping.Run(host, opts)
	if err != nil {
		log.Fatal(err)
	}
	if *porcelain {
		fmt.Printf("\"%v\",%v,%v,%v,%v\n", host, stats.RTTMin, stats.RTTAvg, stats.RTTMax, stats.PLR)
	} else {
		fmt.Printf("%v - %v\n", host, stats)
	}
}
