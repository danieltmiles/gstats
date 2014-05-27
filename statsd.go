package gstats

import (
	"errors"
	"os"
	"time"

	"github.com/cactus/go-statsd-client/statsd"
)

// wrapper/adapter around the cactus statsd client
type Statistics struct {
	client *statsd.Client
}

func CreateStatsdClient() (*Statistics, error) {
	address := os.Getenv("STATSD_ADDRESS")
	if address == "" {
		return nil, errors.New("environment variable STATSD_ADDRESS not defined, cannot continue")
	}
	prefix := os.Getenv("STATSD_PREFIX")
	if prefix == "" {
		return nil, errors.New("environment variable STATSD_PREFIX not defined, cannot continue")
	}
	client, err := statsd.New(address, prefix)
	if err != nil {
		return nil, errors.New("Couldn't initialize statsd.  StatsdInitError=\"" + err.Error() + "\"")
	}
	wrapper := Statistics{client}
	return &wrapper, err
}
// defer stats.End(Trace("foobar"))
func Trace(traceIdentifier string) (string, time.Time, int64) {
	timestamp := time.Now()
	return traceIdentifier, timestamp, 0
}

// defer stats.End(TraceAndIncrement("foobar"))
func TraceAndIncrement(traceIdentifier string) (string, time.Time, int64) {
	timestamp := time.Now()
	return traceIdentifier, timestamp, 1
}

func (s *Statistics) End(traceIdentifier string, timestamp time.Time, incrementBy int64) {
	endingTimestamp := time.Now()
	duration := int64(endingTimestamp.Sub(timestamp))
	if incrementBy > 0 {
		s.IncrementBy(traceIdentifier, incrementBy)
	}
	s.client.Timing(traceIdentifier, duration, 1)
}

func (s *Statistics) Inc(stat string) error {
	return s.IncrementBy(stat, 1)
}

func (s *Statistics) IncrementBy(stat string, incrementBy int64) error {
	return s.client.Inc(stat, incrementBy, 1.0)
}
