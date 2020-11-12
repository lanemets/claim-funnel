package bpm

import "time"

type WorkerConfig struct {
	Retries            int
	RetryTimeoutMillis int

	LockDuration              time.Duration
	MaxTasks                  int
	MaxParallelTaskPerHandler int
	LongPollingTimeout        time.Duration
}