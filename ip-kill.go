package main

import (
	"flag"
	"log"
	"os"
	"time"

	"githug.com/derekneely/ip-kill/monitor"
)

const (
	defInterval   = 5
	defNetTimeout = 5
	defNetErrors  = 3
)

func main() {
	flag.Usage = monitor.PrintUsage

	monitor.PrintBanner()

	c := fetchArgs()
	m := monitor.NewMonitor(&c)

	err := m.Start()
	if err != nil {
		log.Fatal(err)
	}
}

func fetchArgs() monitor.Config {
	// fetch the inteval value
	var itv int
	flag.IntVar(&itv, "interval", defInterval, "how often to run the ip check (default 5)")

	// fetch the network error value
	var nto int
	flag.IntVar(&nto, "timeout", defNetTimeout, "the number of seconds before network ip fetch timeout (default f)")

	// fetch the network error value
	var ne int
	flag.IntVar(&ne, "neterr", defNetErrors, "the number of consecutive network errors before aborting (default 3)")

	// parse flags
	flag.Parse()

	// validate flags
	if itv <= 0 {
		itv = defInterval
	}

	// fetch the process path
	proc := flag.Arg(0)
	if proc == "" {
		monitor.PrintUsage()
		os.Exit(1)
	}

	// return monitor config based on
	return monitor.Config{
		Process:    proc,
		Interval:   itv,
		NetErr:     ne,
		NetTimeout: time.Duration(nto),
	}
}
