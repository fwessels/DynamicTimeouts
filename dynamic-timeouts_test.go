package dynamictimeouts

import (
	"testing"
	"time"
	"fmt"
	"math/rand"
)

func TestSingleIncreaseTimeout(t *testing.T) {

	timeout := NewDynamicTimeout(time.Minute)

	initial := timeout.Timeout()
	fmt.Println(initial)

	for i := 0; i < DynTimeOutLogSize; i++ {
		timeout.LogFailure()
	}

	adjusted := timeout.Timeout()
	fmt.Println(adjusted)
}

func TestDualIncreaseTimeout(t *testing.T) {

	timeout := NewDynamicTimeout(time.Minute)

	initial := timeout.Timeout()
	fmt.Println(initial)

	for i := 0; i < DynTimeOutLogSize; i++ {
		timeout.LogFailure()
	}

	adjusted := timeout.Timeout()
	fmt.Println(adjusted)

	for i := 0; i < DynTimeOutLogSize; i++ {
		timeout.LogFailure()
	}

	adjustedAgain := timeout.Timeout()
	fmt.Println(adjustedAgain)
}

func TestSingleDecreaseTimeout(t *testing.T) {

	timeout := NewDynamicTimeout(time.Minute)

	initial := timeout.Timeout()
	fmt.Println(initial)

	for i := 0; i < DynTimeOutLogSize; i++ {
		timeout.LogSuccess(20 * time.Second)
	}

	adjusted := timeout.Timeout()
	fmt.Println(adjusted)
}

func TestDualDecreaseTimeout(t *testing.T) {

	timeout := NewDynamicTimeout(time.Minute)

	initial := timeout.Timeout()
	fmt.Println(initial)

	for i := 0; i < DynTimeOutLogSize; i++ {
		timeout.LogSuccess(20 * time.Second)
	}

	adjusted := timeout.Timeout()
	fmt.Println(adjusted)

	for i := 0; i < DynTimeOutLogSize; i++ {
		timeout.LogSuccess(20 * time.Second)
	}

	adjustedAgain := timeout.Timeout()
	fmt.Println(adjustedAgain)
}

func TestInfiniteDecreaseTimeout(t *testing.T) {

	timeout := NewDynamicTimeout(time.Minute)

	initial := timeout.Timeout()
	fmt.Println(initial)

	for l := 0; l < 100; l++ {
		for i := 0; i < DynTimeOutLogSize; i++ {
			timeout.LogSuccess(20 * time.Second)
		}

		adjusted := timeout.Timeout()
		if l == 99 {
			fmt.Println(adjusted)
		}
	}
}

func testAdjustTimeout(t *testing.T, timeout *DynamicTimout, f func() float64) {

	for i := 0; i < DynTimeOutLogSize; i++ {

		rnd := f()
		duration := time.Duration(float64(20 * time.Second) * rnd)

		if duration < 100 * time.Millisecond {
			duration = 100 * time.Millisecond
		}
		if duration >= time.Minute {
			timeout.LogFailure()
		} else {
			timeout.LogSuccess(duration)
		}
	}
}

func TestAdjustTimeoutExponential(t *testing.T) {

	timeout := NewDynamicTimeout(time.Minute)

	rand.Seed(time.Now().UTC().UnixNano())

	initial := timeout.Timeout()
	fmt.Println(initial)

	for try := 0; try < 10; try++ {

		testAdjustTimeout(t, timeout, rand.ExpFloat64)

		adjusted := timeout.Timeout()
		fmt.Println(adjusted)
	}
}

func TestAdjustTimeoutNormalized(t *testing.T) {

	timeout := NewDynamicTimeout(time.Minute)

	rand.Seed(time.Now().UTC().UnixNano())

	initial := timeout.Timeout()
	fmt.Println(initial)

	for try := 0; try < 10; try++ {

		testAdjustTimeout(t, timeout, func() float64 {
			return 1.0 + rand.NormFloat64()
		})

		adjusted := timeout.Timeout()
		fmt.Println(adjusted)
	}
}
