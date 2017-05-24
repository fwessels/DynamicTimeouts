package dynamictimeouts

import (
	"sync/atomic"
	"time"
)

const (
	DynTimeOutIncreaseThresholdPct = 0.33 // Upper threshold for failures in order to increase timeout
	DynTimeOutDecreaseThresholdPct = 0.10 // Lower threshold for failures in order to decrease timeout
	DynTimeOutLogSize              = 16
	maxDuration                    = time.Duration(1<<63 - 1)
)

type DynamicTimout struct {
	timeout int64
	entries int64
	log     [DynTimeOutLogSize]time.Duration
}

func NewDynamicTimeout(timeout time.Duration) *DynamicTimout {
	return &DynamicTimout{timeout: int64(timeout)}
}

// Timeout returns the current timeout value
func (dt *DynamicTimout) Timeout() time.Duration {
	return time.Duration(atomic.LoadInt64(&dt.timeout))
}

// LogSuccess logs the duration of a successful action that
// did not hit the timeout
func (dt *DynamicTimout) LogSuccess(duration time.Duration) {
	dt.logEntry(duration)
}

// LogFailure logs an action that hit the timeout
func (dt *DynamicTimout) LogFailure() {
	dt.logEntry(maxDuration)
}

func (dt *DynamicTimout) logEntry(duration time.Duration) {
	entries := int(atomic.AddInt64(&dt.entries, 1))
	index := entries - 1
	if index < DynTimeOutLogSize {
		dt.log[index] = duration
	}
	if entries == DynTimeOutLogSize {
		dt.adjust(entries)
	}
}

func (dt *DynamicTimout) adjust(entries int) {

	failures, average := 0, 0
	for i := 0; i < entries; i++ {
		if dt.log[i] == maxDuration {
			failures++
		} else {
			average += int(dt.log[i])
		}
	}
	if failures < entries {
		average /= entries - failures
	}

	timeOutHitPct := float64(failures) / float64(entries)

	if timeOutHitPct > DynTimeOutIncreaseThresholdPct {
		// We are hitting the timeout too often, so increase the timeout by 25%
		timeout := atomic.LoadInt64(&dt.timeout) * 125 / 100
		atomic.StoreInt64(&dt.timeout, timeout)
	} else if timeOutHitPct < DynTimeOutDecreaseThresholdPct {
		// We are hitting the timeout relatively few times, so let's decrease the timeout
		average = average * 125 / 100 // Add buffer of 25% on top of average

		timeout := (atomic.LoadInt64(&dt.timeout) + int64(average)) / 2 // Middle between current timeout and average success
		atomic.StoreInt64(&dt.timeout, timeout)
	}

	copy
	// reset log entries
	atomic.StoreInt64(&dt.entries, 0)
}
