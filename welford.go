package welford

type Aggregate struct {
	count int
	mean  float64
	m2    float64
}

func NewAggregate() *Aggregate {
	return &Aggregate{}
}

func (A *Aggregate) Update(NewValue float64) *Aggregate {

	A.count++
	Delta := NewValue - A.mean
	A.mean += Delta / float64(A.count)
	Delta2 := NewValue - A.mean
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
	return A.m2 / float64(A.count)
}

func (A *Aggregate) SampleVariance() float64 {
	return A.m2 / float64(A.count-1)
}
