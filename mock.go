package gstats

/*
 * provides a MockStatser struct with accompaning functions to implement the
 * statser interface that will record what gets called. See mock_test.go
 * for usage examples, but in general, it looks like this:
 * > mockStatserInstance.Inc("some stat")
 * > fmt.Printf("%v\n", mockStatserInstance.CallsToInc[0].IncVal)
 * some stat
 */

import (
	"time"
)

type IncSignature struct {
	IncVal string
}

type IncErrSignature struct {
	IncVal string
	Err    error
}

type EndSignature struct {
	Str string
	Tim time.Time
	Num int64
}

type IncrementBySignature struct {
	Str string
	Num int64
}

type GaugeSignature struct {
	Str string
	Num int64
}

type MockStatser struct {
	CallsToInc                 []IncSignature
	CallsToIncErr              []IncErrSignature
	CallsToEnd                 []EndSignature
	CallsToBufferedEnd         []EndSignature
	CallsToIncrementBy         []IncrementBySignature
	CallsToBufferedIncrementBy []IncrementBySignature
	CallsToGauge               []GaugeSignature
}

func (t *MockStatser) Inc(incVal string) error {
	t.CallsToInc = append(t.CallsToInc, IncSignature{incVal})
	return nil
}

func (t *MockStatser) IncErr(incVal string, err error) error {
	t.CallsToIncErr = append(t.CallsToIncErr, IncErrSignature{incVal, err})
	return nil
}

func (t *MockStatser) End(str string, tim time.Time, num int64) {
	t.CallsToEnd = append(t.CallsToEnd, EndSignature{str, tim, num})
}

func (t *MockStatser) BufferedEnd(str string, tim time.Time, num int64) {
	t.CallsToBufferedEnd = append(t.CallsToBufferedEnd, EndSignature{str, tim, num})
}

func (t *MockStatser) IncrementBy(str string, num int64) error {
	t.CallsToIncrementBy = append(t.CallsToIncrementBy, IncrementBySignature{str, num})
	return nil
}

func (t *MockStatser) BufferedIncrementBy(str string, num int64) error {
	t.CallsToBufferedIncrementBy = append(t.CallsToBufferedIncrementBy, IncrementBySignature{str, num})
	return nil
}

func (t *MockStatser) Gauge(str string, num int64) error {
	t.CallsToGauge = append(t.CallsToGauge, GaugeSignature{str, num})
	return nil
}
