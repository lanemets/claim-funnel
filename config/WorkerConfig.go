package config

import "time"

type Worker struct {
	Retries            int
	RetryTimeoutMillis int

	LockDuration              time.Duration
	MaxTasks                  int
	MaxParallelTaskPerHandler int
	LongPollingTimeout        time.Duration
}