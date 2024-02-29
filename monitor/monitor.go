package monitor

import (
	"flag"
	"fmt"
	"io"
	"net/http"
	"os"
	"os/exec"
	"os/signal"
	"strings"
	"syscall"
	"time"
)

const (
	api_url = "https://api.ipify.org"

	sChar = "*"
	eChar = "X"

	colorReset  = "\033[0m"
	colorRed    = "\033[31m"
	colorGreen  = "\033[32m"
	colorYellow = "\033[33m"
	colorBlue   = "\033[34m"
)

type monitor struct {
	config *Config
	status *Status
}

func NewMonitor(c *Config) *monitor {
	return &monitor{
		config: c,
		status: &Status{},
	}
}

func (m *monitor) Start() error {
	// fetch the initial ip
	ip, err := m.fetchIp()
	if err != nil {
		return err
	}

	// start message
	m.printConfig(ip)

	// sigint channel handling
	sigint := make(chan os.Signal, 1)
	signal.Notify(sigint, os.Interrupt, syscall.SIGTERM)
	go func() {
		<-sigint
		fmt.Println(colorYellow, "\n** Stopping monitor", colorReset)
		os.Exit(0)
	}()

	// processing loop
	for {
		ipChk, err := m.fetchIp()
		if err != nil {
			m.updateStatus(false)
			m.printStatus()
			if m.checkStatus() {
				fmt.Println(colorRed, "\n** NETWORK ERRORS EXCEEDED **", colorReset)
				m.killProc()
				break
			}

			// set ipChk value to ip to pass checks downstream
			ipChk = ip
		} else {
			m.updateStatus(true)
			m.printStatus()
		}

		if ipChk != ip {
			fmt.Println(colorRed, "\n** IP CHANGED **", colorReset)
			m.killProc()
			break
		}

		time.Sleep(time.Duration(m.config.Interval) * time.Second)
	}

	fmt.Println("** Shutting down **")
	return nil
}

func (m *monitor) killProc() {
	fmt.Printf("Aborting process: %s\n", m.config.Process)
	cmd := exec.Command("pkill", "-f", m.config.Process)
	err := cmd.Run()
	if err != nil {
		fmt.Printf("Failed to kill process: %s\n", err)
		return
	}

	fmt.Printf("Process: %s has been killed.\n", m.config.Process)
}

func (m *monitor) fetchIp() (string, error) {
	client := http.Client{
		Timeout: m.config.NetTimeout * time.Second,
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

func (m *monitor) updateStatus(success bool) {
	// reset the progress bar text
	if len(m.status.ProgressBar) >= 50 {
		m.status.ProgressBar = ""
	}

	// set raw progress bar character sequence
	if success {
		m.status.ProgressBar += sChar
		m.status.SuccessCnt += 1

		// reset consecutive errs
		m.status.ConsErrCnt = 0
	} else {
		m.status.ProgressBar += eChar
		m.status.ErrCnt += 1
		m.status.ConsErrCnt += 1
	}
}

func (m *monitor) checkStatus() bool {
	return m.status.ConsErrCnt >= m.config.NetErr
}

func (m *monitor) printStatus() {
	const pLength = 50
	const ansiLength = 9

	// update the probess bar with some colors
	// todo: replace color strings with const
	pBar := strings.ReplaceAll(m.status.ProgressBar, sChar, fmt.Sprintf("\033[32m%s\033[0m", sChar))
	pBar = strings.ReplaceAll(pBar, eChar, fmt.Sprintf("\033[31m%s\033[0m", eChar))
	sCount := strings.Count(pBar, "\033[32m")
	eCount := strings.Count(pBar, "\033[31m")
	vLen := len(pBar) - ((eCount + sCount) * ansiLength)
	pBar = pBar + strings.Repeat(" ", pLength-vLen)

	// success string formatting
	sucStr := fmt.Sprintf("Success: %s %5d%s", colorGreen, m.status.SuccessCnt, colorReset)

	// error count formatting
	errColor := colorReset
	if m.status.ErrCnt > 0 {
		errColor = colorRed
	}
	errStr := fmt.Sprintf("Total Errs: %s %5d %s", errColor, m.status.ErrCnt, colorReset)

	// consecutive error formatting
	conErrColor := colorReset
	if m.status.ConsErrCnt > 0 {
		conErrColor = colorRed
	}
	conErrStr := fmt.Sprintf("%sCons. Errs: %5d/%d %s", conErrColor, m.status.ConsErrCnt, m.config.NetErr, colorReset)

	// fmt.Printf("\r%s[%-50s] Success: %5d | Total Errs: %5d | Cons. Errs: %5d   ", colorReset, pBar, m.status.SuccessCnt, m.status.ErrCnt, m.status.ConsErrCnt)
	fmt.Printf("\r%s[%-50s] %s | %s | %s       ", colorReset, pBar, sucStr, errStr, conErrStr)
}

func (m *monitor) printConfig(ip string) {
	fmt.Print(colorBlue)
	fmt.Printf("** Starting ip monitor: %s\n", ip)
	fmt.Printf("    --Target process: %s\n", m.config.Process)
	fmt.Printf("    --Interval: %d\n", m.config.Interval)
	fmt.Printf("    --Network Timeout: %d\n", m.config.NetTimeout)
	fmt.Printf("    --Network Error Tolerance: %d\n", m.config.NetErr)
	fmt.Println(colorReset)
}

func PrintUsage() {
	fmt.Print(colorGreen)
	fmt.Println("Usage: ./ip-kill [-interval=5] [-timeout=5] [-neterr=3] process_to_kill")
	flag.PrintDefaults()
	fmt.Println(colorReset)
}

func PrintBanner() {
	w := `
._____________           ____  __.__.__  .__   
|   \______   \         |    |/ _|__|  | |  |  
|   ||     ___/  ______ |      < |  |  | |  |  
|   ||    |     /_____/ |    |  \|  |  |_|  |__
|___||____|             |____|__ \__|____/____/
                                \/             
`
	fmt.Println(colorGreen, w, colorReset)
}
