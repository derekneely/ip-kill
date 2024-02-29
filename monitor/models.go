package monitor

import "time"

type Config struct {
	Process    string
	Interval   int
	NetTimeout time.Duration
	NetErr     int
}

type Status struct {
	ProgressBar string
	SuccessCnt  int
	ErrCnt      int
	ConsErrCnt  int
}
