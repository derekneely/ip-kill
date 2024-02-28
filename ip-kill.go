package main

import (
	"flag"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"os/exec"
	"os/signal"
	"syscall"
	"time"
)

const (
	api_url = "https://api.ipify.org"

	defInterval   = 5
	defNetTimeout = 5
	defNetErrors  = 3

	colorReset  = "\033[0m"
	colorRed    = "\033[31m"
	colorGreen  = "\033[32m"
	colorYellow = "\033[33m"
	colorBlue   = "\033[34m"
	colorPurple = "\033[35m"
	colorCyan   = "\033[36m"
	colorWhite  = "\033[37m"
)

func main() {
	printBanner()

	c := fetchArgs()

	err := start(c)
	if err != nil {
		log.Fatal(err)
	}
}

func start(c config) error {
	// fetch the initial ip
	ip, err := fetchIp(c.netTimeout)
	if err != nil {
		return err
	}

	// start message
	fmt.Println(string(colorGreen), fmt.Sprintf("** Starting ip monitor: %s", ip))
	printConfig(c)

	// sigint channel handling
	sigint := make(chan os.Signal, 1)
	signal.Notify(sigint, os.Interrupt, syscall.SIGTERM)
	go func() {
		<-sigint
		fmt.Println(string(colorYellow), "** Stopping ip monitor")
		os.Exit(0)
	}()

	// Infinite loop
	ec := 0
	for {
		ipChk, err := fetchIp(c.netTimeout)
		if err != nil {
			ec += 1
			if ec >= c.netErr {
				fmt.Println(string(colorRed), "X\n** NETWORK ERRORS EXCEEDED")
				killProc(c.process)
				break
			}

			// set ipChk value to ip to pass checks downstream
			ipChk = ip
			fmt.Print(string(colorRed), "X")
		} else {
			ec = 0
			fmt.Print(string(colorGreen), "*")
		}

		if ipChk != ip {
			fmt.Println(string(colorRed), "\n** IP CHANGED **")
			killProc(c.process)
			break
		}

		time.Sleep(time.Duration(c.interval) * time.Second)
	}

	fmt.Println("** Shutting down **")
	return nil
}

func killProc(p string) {
	fmt.Printf("Aborting process: %s\n", p)

	cmd := exec.Command("pkill", "-f", p)
	// var out bytes.Buffer
	// cmd.Stdout = &out
	err := cmd.Run()
	if err != nil {
		fmt.Printf("Failed to kill process: %s\n", err)
		return
	}

	fmt.Printf("Process: %s has been killed.\n", p)
}

func fetchIp(to time.Duration) (string, error) {
	client := http.Client{
		Timeout: to * time.Second,
	}

	res, err := client.Get(api_url)
	if err != nil {
		return "", err
	}

	ip, err := io.ReadAll(res.Body)
	if err != nil {
		return "", err
	}

	return string(ip), nil
}

type config struct {
	process    string
	interval   int
	netTimeout time.Duration
	netErr     int
}

func fetchArgs() config {
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

	proc := flag.Arg(0)
	if proc == "" {
		printHelp()
		os.Exit(1)
	}

	return config{
		process:    proc,
		interval:   itv,
		netErr:     ne,
		netTimeout: time.Duration(nto),
	}
}

func printConfig(c config) {
	fmt.Print(string(colorGreen))
	fmt.Printf("    --Target process: %s\n", c.process)
	fmt.Printf("    --Interval: %d\n", c.interval)
	fmt.Printf("    --Network Timeout: %d\n", c.netTimeout)
	fmt.Printf("    --Network Error Tolerance: %d\n", c.netErr)
}

func printHelp() {
	fmt.Println(string(colorReset))
	fmt.Println("Todo: Print Help")
}

func printBanner() {
	w := `
._____________           ____  __.__.__  .__   
|   \______   \         |    |/ _|__|  | |  |  
|   ||     ___/  ______ |      < |  |  | |  |  
|   ||    |     /_____/ |    |  \|  |  |_|  |__
|___||____|             |____|__ \__|____/____/
                                \/             
`
	fmt.Println(string(colorBlue))
	fmt.Println(w)
}
