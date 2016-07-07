package gstats

import (
	"errors"
	"fmt"
	"os"
	"time"
	"unicode"
	"unicode/utf8"

	"github.com/cactus/go-statsd-client/statsd"
	"github.com/etgryphon/stringUp"
)

type incrementer func(stat string, incrementBy int64) error

type Statser interface {
	End(string, time.Time, int64)
	BufferedEnd(string, time.Time, int64)
	Inc(string) error
	IncErr(string, error) error
	IncrementBy(string, int64) error
	BufferedIncrementBy(string, int64) error
	Gauge(string, int64) error
}

// wrapper/adapter around the cactus statsd client
type Statistics struct {
	client            statsd.Statter
	IncrementBuffers  map[string]int64
	BufferFlushPeriod time.Duration
	mu                chan bool
}

func CreateStatsdClient() (*Statistics, error) {
	return _CreateStatsdClient(time.Second)
}

func _CreateStatsdClient(bufferFlushPeriod time.Duration) (*Statistics, error) {
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
	wrapper := Statistics{client, make(map[string]int64), bufferFlushPeriod, make(chan bool, 1)}
	go wrapper.autoFlushBufferedStats()
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

func (s *Statistics) autoFlushBufferedStats() {
	for {
		<-time.After(s.BufferFlushPeriod)
		s.flushBufferedStats()
	}
}

func (s *Statistics) flushBufferedStats() {
	for stat, incValue := range s.IncrementBuffers {
		if incValue > 0 {
			s.IncrementBuffers[stat] = 0
			s.client.Inc(stat, incValue, 1.0)
		}
	}
}

func (s *Statistics) _End(traceIdentifier string, timestamp time.Time, incrementBy int64, incFunc incrementer) {
	endingTimestamp := time.Now()
	duration := int64(endingTimestamp.Sub(timestamp) / time.Millisecond)
	if incrementBy > 0 {
		incFunc(traceIdentifier+".count", incrementBy)
	}
	s.client.Timing(traceIdentifier, duration, 1)
}

func (s *Statistics) End(traceIdentifier string, timestamp time.Time, incrementBy int64) {
	s._End(traceIdentifier, timestamp, incrementBy, s.IncrementBy)
}

func (s *Statistics) BufferedEnd(traceIdentifier string, timestamp time.Time, incrementBy int64) {
	s._End(traceIdentifier, timestamp, incrementBy, s.BufferedIncrementBy)
}

// stats.Inc("AnEvent")
func (s *Statistics) Inc(stat string) error {
	return s.IncrementBy(stat+".count", 1)
}

// stats.IncErr("This.Event.Records", "my error message") => "This.Event.Records.MyErrorMessage.count"
// expected behavior to strip non-western characters
func (s *Statistics) IncErr(stat string, err error) error {
	stat = fmt.Sprintf(stat+".%s.count", normalize(err))
	return s.IncrementBy(stat, 1)
}

// stats.IncrementBy("Requests", 5) // I got 5 requests!
func (s *Statistics) IncrementBy(stat string, incrementBy int64) error {
	return s.client.Inc(stat, incrementBy, 1.0)
}

// stats.BufferedIncrementBy("Requests", 5) // I got 5 requests!
func (s *Statistics) BufferedIncrementBy(stat string, incrementBy int64) error {
	s.mu <- true
	defer func() { <-s.mu }()
	// if we have never seen this stat before, we simply return 0 for val
	val, _ := s.IncrementBuffers[stat]
	s.IncrementBuffers[stat] = val + incrementBy
	incValue, _ := s.IncrementBuffers[stat]
	if s.IncrementBuffers[stat] >= 100 {
		s.IncrementBuffers[stat] = 0
		return s.client.Inc(stat, incValue, 1.0)
	}
	return nil
}

func (s *Statistics) Gauge(stat string, value int64) error {
	return s.client.Gauge(stat, value, 1.0)
}

// expected behavior to strip non-western characters
func normalize(toRecord error) string {
	// the stringUp takes out the non-western chars
	// "camelCase"
	cameled := stringUp.CamelCase(toRecord.Error())

	// "CamelCase"
	rune, size := utf8.DecodeRuneInString(cameled)
	return string(unicode.ToUpper(rune)) + cameled[size:]
}
