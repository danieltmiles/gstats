package gstats

import (
	"errors"
	"fmt"
	"os"
	"time"

	"github.com/cactus/go-statsd-client/statsd"
)

type statsToBeIncremented struct {
	Stat  string
	Count int64
	Rate  float32
}
type statsToBeTraced struct {
	Stat  string
	Start time.Time
}
type StatsClientWrapper struct {
	Client             *statsd.Client
	IncrementableStats map[string]statsToBeIncremented
	TimeableStats      map[string]statsToBeTraced
}

func CreateStatsdClient() (*StatsClientWrapper, error) {
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
	wrapper := StatsClientWrapper{client, make(map[string]statsToBeIncremented), make(map[string]statsToBeTraced)}
	return &wrapper, err
}
func (w *StatsClientWrapper) StartTrace(traceIdentifier string) error {
	timestamp := time.Now()
	w.TimeableStats[traceIdentifier] = statsToBeTraced{traceIdentifier, timestamp}
	return nil
}

func (w *StatsClientWrapper) EndTrace(traceIdentifier string) error {
	timestamp := time.Now()
	_, keyExists := w.TimeableStats[traceIdentifier]
	if !keyExists {
		return errors.New(fmt.Sprintf("Error, must start trace \"%s\" before ending it", traceIdentifier))
	}
	startTime := w.TimeableStats[traceIdentifier].Start
	duration := int64(timestamp.Sub(startTime))
	w.Client.Timing(traceIdentifier, duration, 1)
	delete(w.TimeableStats, traceIdentifier)
	return nil
}

func (w *StatsClientWrapper) Inc(stat string) error {
	_, keyExists := w.IncrementableStats[stat]
	if !keyExists {
		w.IncrementableStats[stat] = statsToBeIncremented{stat, 0, 1.0}
	}
	count := w.IncrementableStats[stat].Count + 1
	oldRate := w.IncrementableStats[stat].Rate
	// BOZO: cannot assign to w.IncrementableStats[stat].Count, so replace with
	//       totally new object and hope the garbage collector does its thing
	w.IncrementableStats[stat] = statsToBeIncremented{stat, count, oldRate}
	w.Client.Inc(stat, count, oldRate)
	return nil
}
