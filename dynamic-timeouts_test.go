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
