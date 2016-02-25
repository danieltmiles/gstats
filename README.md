# Gstats -- A statsd wrapper for go

# Usage
```go
import "github.com/monsooncommerce/gstats"
// initialize a statser
var stats gstats.Statser
// make sure that STATSD_ADDRESS and STATSD_PREFIX are defined in your environment
// the address should be your statsd server indluding port number, the prefix is the
// first element of your statistic path
stats, err = gstats.CreateStatsdClient()
if err != nil {
	// you know what to do
}
```
Now let's use it
```go
stats.Inc("statistic") // increments a counter called stats.counters.$STATSD_PREFIX.statistic
stats.IncrementBy("statistic", 5) // increments stats.counters.$STATSD_PREFIX.statistic counter by 5
stats.BufferedIncrementBy("statistic", 5) // Sums up all calls to BufferedIncrementBy and
                                          // sends them off when it gets 100 of them or
                                          // when it's been saving them for 1 second.
                                          // This can be important if your stating is causing
                                          // performance issues
stats.Gauge("statistic", 5) // sets stats.gauges.$STATSD_PREFIX.statistic gauge value to 5
```
Now let's use the package to stat how long a given function took to execute
```go
func MyFunc(arg string) string {
	//start your function with a defer statement. As with all Go defer statements,
	// gstats.TraceAndIncrement will execute as the first statement in the function
	// and stats.BufferedEnd will execute as the last thing done before returning.
	defer stats.BufferedEnd(gstats.TraceAndIncrement("MyFunc"))
	value, err := someFunctionCall()
	if err != nil {
		return "we had an error"
	}
	return "everything ran fine"
}
```
this way we know how long your function took to execute, no matter which exit-point it finished at.
