package welford

import (
	"fmt"
	"sync"
)

type WelfordAggregate interface {
	Reset() WelfordAggregate
	Update(...float64) WelfordAggregate
	Count() int
	Mean() float64
	Variance() float64
	SampleVariance() float64
	Results() (int, float64, float64, float64)
	String() string
}

// Aggregate is an opaque struct which holds the current status of the online
// calculation of mean and standard deviation of the corresponding random
// variable.
type Aggregate struct {
	count int
	mean  float64
	m2    float64
}

// ConcurrentAggregate is an extension of the basic Aggregate where it allows
// for concurrent access by multiple go-routines or producers.
type ConcurrentAggregate struct {
	Aggregate
	mu *sync.RWMutex
}

// NewAggregate
//
// This function creates and initializes a new Aggregate value, ready to be used.
func NewAggregate() *Aggregate {
	return &Aggregate{}
}

// Reset
//
// This function resets an existing Aggregate back to a zero-value, readying it
// to be used on a new random sequence.
func (A *Aggregate) Reset() *Aggregate {

	A.count = 0
	A.mean = 0
	A.m2 = 0

	return A
}

// Update
//
// This function accepts a new random sample, and updates the internal state of
// the Aggregate to account for the newly provided sample value.
func (A *Aggregate) Update(Values ...float64) *Aggregate {

	for _, v := range Values {
		A = A.update(v)
	}

	return A
}

func (A *Aggregate) update(Value float64) *Aggregate {

	A.count++
	Delta := Value - A.mean
	A.mean += Delta / float64(A.count)
	Delta2 := Value - A.mean
	A.m2 += Delta * Delta2

	return A
}

func (A *Aggregate) Count() int {
	return A.count
}

func (A *Aggregate) Mean() float64 {
	return A.mean
}

func (A *Aggregate) Variance() float64 {

	if A.count == 0 {
		return 0
	}

	return A.m2 / float64(A.count)
}

func (A *Aggregate) SampleVariance() float64 {

	if A.count <= 1 {
		return 0
	}

	return A.m2 / float64(A.count-1)
}

func (A *Aggregate) Results() (int, float64, float64, float64) {
	return A.Count(), A.Mean(), A.Variance(), A.SampleVariance()
}

func (A *Aggregate) String() string {

	Count, Mean, Variance, SampleVariance := A.Results()

	return fmt.Sprintf("Count: %d, Mean: %f, Variance: %f, Sample Variance: %f", Count, Mean, Variance, SampleVariance)
}

// NewConcurrentAggregate
//
// This function creates and initializes a new ConcurrentAggregate value, ready to be used.
func NewConcurrentAggregate() *ConcurrentAggregate {
	return &ConcurrentAggregate{
		mu: &sync.RWMutex{},
	}
}

// Reset
//
// This function resets an existing ConcurrentAggregate back to a zero-value, readying it
// to be used on a new random sequence.
func (A *ConcurrentAggregate) Reset() *ConcurrentAggregate {

	A.mu.Lock()
	defer A.mu.Unlock()

	A.Aggregate.Reset()

	return A
}

// Update
//
// This function accepts a new random sample, and updates the internal state of
// the ConcurrentAggregate to account for the newly provided sample value.
func (A *ConcurrentAggregate) Update(Values ...float64) *ConcurrentAggregate {

	A.mu.Lock()
	defer A.mu.Unlock()

	A.Aggregate.Update(Values...)
	return A

}

func (A *ConcurrentAggregate) Variance() float64 {

	A.mu.RLock()
	defer A.mu.RUnlock()

	return A.Aggregate.Variance()
}

func (A *ConcurrentAggregate) SampleVariance() float64 {

	A.mu.RLock()
	defer A.mu.RUnlock()

	return A.Aggregate.SampleVariance()
}

func (A *ConcurrentAggregate) Results() (int, float64, float64, float64) {

	A.mu.RLock()
	defer A.mu.RUnlock()

	return A.Aggregate.Results()
}

func (A *ConcurrentAggregate) String() string {

	Count, Mean, Variance, SampleVariance := A.Results()

	return fmt.Sprintf("Count: %d, Mean: %f, Variance: %f, Sample Variance: %f", Count, Mean, Variance, SampleVariance)
}
