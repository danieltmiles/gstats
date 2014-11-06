package gstats

import (
	"time"
)

type DummyStats struct {
}

func (s *DummyStats) End(string, time.Time, int64) {
}

func (s *DummyStats) BufferedEnd(str string, tim time.Time, num int64) {
}

func (s *DummyStats) Inc(string) error {
	return nil
}

func (s *DummyStats) IncErr(string, error) error {
	return nil
}

func (s *DummyStats) IncrementBy(string, int64) error {
	return nil
}

func (s *DummyStats) BufferedIncrementBy(string, int64) error {
	return nil
}

func (s *DummyStats) Gauge(string, int64) error {
	return nil
}
