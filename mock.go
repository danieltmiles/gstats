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
	"fmt"
	"sync"
	"time"
)

var (
	mu = sync.Mutex{}
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

func NewMock() MockStatser {
	m := MockStatser{
		CallsToInc:                 []IncSignature{},
		CallsToIncErr:              []IncErrSignature{},
		CallsToEnd:                 []EndSignature{},
		CallsToBufferedEnd:         []EndSignature{},
		CallsToIncrementBy:         []IncrementBySignature{},
		CallsToBufferedIncrementBy: []IncrementBySignature{},
		CallsToGauge:               []GaugeSignature{},
	}
	return m
}

func (t *MockStatser) Inc(incVal string) error {
	mu.Lock()
	defer func() { mu.Unlock() }()
	t.CallsToInc = append(t.CallsToInc, IncSignature{incVal})
	return nil
}

func (t *MockStatser) IncErr(incVal string, err error) error {
	mu.Lock()
	defer func() { mu.Unlock() }()
	t.CallsToIncErr = append(t.CallsToIncErr, IncErrSignature{incVal, err})
	return nil
}

func (t *MockStatser) End(str string, tim time.Time, num int64) {
	mu.Lock()
	defer func() { mu.Unlock() }()
	t.CallsToEnd = append(t.CallsToEnd, EndSignature{str, tim, num})
}

func (t *MockStatser) BufferedEnd(str string, tim time.Time, num int64) {
	mu.Lock()
	defer func() { mu.Unlock() }()
	if t.CallsToBufferedEnd == nil {
		fmt.Printf("callsToBufferedEnd nil\n")
	}
	t.CallsToBufferedEnd = append(t.CallsToBufferedEnd, EndSignature{str, tim, num})
}

func (t *MockStatser) IncrementBy(str string, num int64) error {
	mu.Lock()
	defer func() { mu.Unlock() }()
	t.CallsToIncrementBy = append(t.CallsToIncrementBy, IncrementBySignature{str, num})
	return nil
}

func (t *MockStatser) BufferedIncrementBy(str string, num int64) error {
	mu.Lock()
	defer func() { mu.Unlock() }()
	t.CallsToBufferedIncrementBy = append(t.CallsToBufferedIncrementBy, IncrementBySignature{str, num})
	return nil
}

func (t *MockStatser) Gauge(str string, num int64) error {
	mu.Lock()
	defer func() { mu.Unlock() }()
	if t.CallsToGauge == nil {
		fmt.Printf("calls to gauge is nill\n")
	}
	t.CallsToGauge = append(t.CallsToGauge, GaugeSignature{str, num})
	return nil
}
