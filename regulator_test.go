package regulator

import (
	"errors"
	"sync/atomic"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestConcurrency(t *testing.T) {
	r := NewRegulator(2)
	count := 10
	collectTimes := make(chan int, count)

	for i := 0; i < count; i++ {
		r.Execute(func() error {
			collectTimes <- timestamp()
			time.Sleep(time.Millisecond * 20)
			return nil
		})
	}
	r.Wait()

	times := []int{}
	for i := 0; i < count; i++ {
		times = append(times, <-collectTimes)
		if i%2 == 1 {
			assert.True(t, times[i]-times[i-1] < 5, "go routines did not execute in groups of 2")
		}
	}
}

func TestErrorHandling(t *testing.T) {
	r := NewRegulator(2)
	count := 10
	executeCount := int32(0)

	for i := 0; i < count; i++ {
		r.Execute(func() error {
			atomic.AddInt32(&executeCount, 1)
			if executeCount >= 5 {
				return errors.New("ut oh")
			}
			return nil
		})
	}
	err := r.Wait()
	regErr := err.(RegulatorError)
	assert.Equal(t, "ut oh", regErr.Error())
	assert.True(t, executeCount <= 6, "regulator should have not executed more than 6 go routines")
}

func timestamp() int {
	return int(time.Now().UnixNano() / int64(time.Millisecond))
}
