package monitor

import (
	"net/http"
	"time"
)

type Config struct {
	Process    string
	Interval   int
	NetTimeout time.Duration
	NetErr     int
	Transport  *http.Transport
}

type Status struct {
	ProgressBar string
	SuccessCnt  int
	ErrCnt      int
	ConsErrCnt  int
}
