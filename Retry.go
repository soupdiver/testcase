package testcase

import (
	"testing"
	"time"

	"github.com/adamluzsi/testcase/internal"
)

// Retry Automatically retries operations whose failure is expected under certain defined conditions.
// This pattern enables fault-tolerance.
//
// A common scenario where using Retry will benefit you is testing concurrent operations.
// Due to the nature of async operations, one might need to wait
// and observe the system with multiple tries before the outcome can be seen.
type Retry struct{ Strategy RetryStrategy }

type RetryStrategy interface {
	// While implements the retry strategy looping part.
	// Depending on the outcome of the condition,
	// the RetryStrategy can decide whether further iterations can be done or not
	While(condition func() bool)
}

type RetryStrategyFunc func(condition func() bool)

func (fn RetryStrategyFunc) While(condition func() bool) { fn(condition) }

// Assert will attempt to assert with the assertion function block multiple times until the expectations in the function body met.
// In case expectations are failed, it will retry the assertion block using the RetryStrategy.
// The last failed assertion results would be published to the received testing.TB.
// Calling multiple times the assertion function block content should be a safe and repeatable operation.
func (r Retry) Assert(tb testing.TB, blk func(testing.TB)) {
	tb.Helper()
	var lastRecorder *internal.RecorderTB

	r.Strategy.While(func() bool {
		tb.Helper()
		lastRecorder = &internal.RecorderTB{TB: tb}
		internal.RecoverExceptGoexit(func() {
			tb.Helper()
			blk(lastRecorder)
		})
		if lastRecorder.IsFailed {
			lastRecorder.CleanupNow()
		}
		return lastRecorder.IsFailed
	})

	if lastRecorder != nil {
		lastRecorder.Forward()
	}
}

//func (r Retry) setup(s *Spec) {
//	s.flaky = &r
//}

func RetryCount(times int) RetryStrategy {
	return RetryStrategyFunc(func(condition func() bool) {
		for i := 0; i < times+1; i++ {
			if ok := condition(); !ok {
				return
			}
		}
	})
}

func makeRetry(i interface{}) (Retry, bool) {
	switch n := i.(type) {
	case time.Duration:
		return Retry{Strategy: Waiter{WaitTimeout: n}}, true
	case int:
		return Retry{Strategy: RetryCount(n)}, true
	case RetryStrategy:
		return Retry{Strategy: n}, true
	case Retry:
		return n, true
	default:
		return Retry{}, false
	}
}
