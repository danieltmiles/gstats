package gstats

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

type GagueSignature struct {
	Str string
	Num int64
}

type CallRecordingMockStatser struct {
	CallsToInc                 []IncSignature
	CallsToIncErr              []IncErrSignature
	CallsToEnd                 []EndSignature
	CallsToBufferedEnd         []EndSignature
	CallsToIncrementBy         []IncrementBySignature
	CallsToBufferedIncrementBy []IncrementBySignature
	CallsToGauge               []GagueSignature
}

func (t *CallRecordingMockStatser) Inc(incVal string) error {
	t.CallsToInc = append(t.CallsToInc, IncSignature{incVal})
	return nil
}

func (t *CallRecordingMockStatser) IncErr(incVal string, err error) error {
	t.CallsToIncErr = append(t.CallsToIncErr, IncErrSignature{incVal, err})
	return nil
}

func (t *CallRecordingMockStatser) End(str string, tim time.Time, num int64) {
	t.CallsToEnd = append(t.CallsToEnd, EndSignature{str, tim, num})
}

func (t *CallRecordingMockStatser) BufferedEnd(str string, tim time.Time, num int64) {
	t.CallsToBufferedEnd = append(t.CallsToBufferedEnd, EndSignature{str, tim, num})
}

func (t *CallRecordingMockStatser) IncrementBy(str string, num int64) error {
	t.CallsToIncrementBy = append(t.CallsToIncrementBy, IncrementBySignature{str, num})
	return nil
}

func (t *CallRecordingMockStatser) BufferedIncrementBy(str string, num int64) error {
	t.CallsToBufferedIncrementBy = append(t.CallsToBufferedIncrementBy, IncrementBySignature{str, num})
	return nil
}

func (t *CallRecordingMockStatser) Gauge(str string, num int64) error {
	t.CallsToGauge = append(t.CallsToGauge, GagueSignature{str, num})
	return nil
}
