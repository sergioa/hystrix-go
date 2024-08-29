package plugins

import (
	"github.com/DataDog/datadog-go/statsd"
	metricCollector "github.com/sergioa/hystrix-go/hystrix/metric_collector"
)

// These metrics are constants because we're leveraging the Datadog tagging
// extension to statsd.
//
// They only apply to the DatadogCollector and are only useful if providing your
// own implementation of DatadogClient
const (
	// DM = Datadog Metric
	DmCircuitOpen       = "hystrix.circuitOpen"
	DmAttempts          = "hystrix.attempts"
	DmErrors            = "hystrix.errors"
	DmSuccesses         = "hystrix.successes"
	DmFailures          = "hystrix.failures"
	DmRejects           = "hystrix.rejects"
	DmShortCircuits     = "hystrix.shortCircuits"
	DmTimeouts          = "hystrix.timeouts"
	DmFallbackSuccesses = "hystrix.fallbackSuccesses"
	DmFallbackFailures  = "hystrix.fallbackFailures"
	DmTotalDuration     = "hystrix.totalDuration"
	DmRunDuration       = "hystrix.runDuration"
)

type (
	// DatadogClient is the minimum interface needed by
	// NewDatadogCollectorWithClient
	DatadogClient interface {
		Count(name string, value int64, tags []string, rate float64) error
		Gauge(name string, value float64, tags []string, rate float64) error
		TimeInMilliseconds(name string, value float64, tags []string, rate float64) error
	}

	// DatadogCollector fulfills the metricCollector interface allowing users to
	// ship circuit stats to Datadog.
	//
	// This Collector, by default, uses github.com/DataDog/datadog-go/statsd for
	// transport. The main advantage of this over statsd is building graphs and
	// multi-alert monitors around single metrics (constantized above) and
	// adding tag dimensions. You can set up a single monitor to rule them all
	// across services and geographies. Graphs become much simpler to set up by
	// allowing you to create queries like the following
	//
	//   {
	//     "viz": "timeseries",
	//     "requests": [
	//       {
	//         "q": "max:hystrix.runDuration.95percentile{$region} by {hystrixcircuit}",
	//         "type": "line"
	//       }
	//     ]
	//   }
	//
	// As new circuits come online you get graphing and monitoring "for free".
	DatadogCollector struct {
		client DatadogClient
		tags   []string
	}
)

// NewDatadogCollector creates a collector for a specific circuit with a
// "github.com/DataDog/datadog-go/statsd".(*Client).
//
// addr is in the format "<host>:<port>" (e.g. "localhost:8125")
//
// prefix may be an empty string
//
// Example use
//
//	package main
//
//	import (
//		"github.com/afex/hystrix-go/plugins"
//		"github.com/afex/hystrix-go/hystrix/metric_collector"
//	)
//
//	func main() {
//		collector, err := plugins.NewDatadogCollector("localhost:8125", "")
//		if err != nil {
//			panic(err)
//		}
//		metricCollector.Registry.Register(collector)
//	}
func NewDatadogCollector(addr, prefix string) (func(string) metricCollector.MetricCollector, error) {

	c, err := statsd.NewBuffered(addr, 100)
	if err != nil {
		return nil, err
	}

	// Prefix every metric with the app name
	c.Namespace = prefix

	return NewDatadogCollectorWithClient(c), nil
}

// NewDatadogCollectorWithClient accepts an interface which allows you to
// provide your own implementation of a statsd client, alter configuration on
// "github.com/DataDog/datadog-go/statsd".(*Client), provide additional tags per
// circuit-metric tuple, and add logging if you need it.
func NewDatadogCollectorWithClient(client DatadogClient) func(string) metricCollector.MetricCollector {

	return func(name string) metricCollector.MetricCollector {

		return &DatadogCollector{
			client: client,
			tags:   []string{"hystrixcircuit:" + name},
		}
	}
}

func (dc *DatadogCollector) Update(r metricCollector.MetricResult) {
	if r.Attempts > 0 {
		_ = dc.client.Count(DmAttempts, int64(r.Attempts), dc.tags, 1.0)
	}
	if r.Errors > 0 {
		_ = dc.client.Count(DmErrors, int64(r.Errors), dc.tags, 1.0)
	}
	if r.Successes > 0 {
		_ = dc.client.Gauge(DmCircuitOpen, 0, dc.tags, 1.0)
		_ = dc.client.Count(DmSuccesses, int64(r.Successes), dc.tags, 1.0)
	}
	if r.Failures > 0 {
		_ = dc.client.Count(DmFailures, int64(r.Failures), dc.tags, 1.0)
	}
	if r.Rejects > 0 {
		_ = dc.client.Count(DmRejects, int64(r.Rejects), dc.tags, 1.0)
	}
	if r.ShortCircuits > 0 {
		_ = dc.client.Gauge(DmCircuitOpen, 1, dc.tags, 1.0)
		_ = dc.client.Count(DmShortCircuits, int64(r.ShortCircuits), dc.tags, 1.0)
	}
	if r.Timeouts > 0 {
		_ = dc.client.Count(DmTimeouts, int64(r.Timeouts), dc.tags, 1.0)
	}
	if r.FallbackSuccesses > 0 {
		_ = dc.client.Count(DmFallbackSuccesses, int64(r.FallbackSuccesses), dc.tags, 1.0)
	}
	if r.FallbackFailures > 0 {
		_ = dc.client.Count(DmFallbackFailures, int64(r.FallbackFailures), dc.tags, 1.0)
	}

	ms := float64(r.TotalDuration.Nanoseconds() / 1000000)
	_ = dc.client.TimeInMilliseconds(DmTotalDuration, ms, dc.tags, 1.0)

	ms = float64(r.RunDuration.Nanoseconds() / 1000000)
	_ = dc.client.TimeInMilliseconds(DmRunDuration, ms, dc.tags, 1.0)
}

// Reset is a noop operation in this collector.
func (dc *DatadogCollector) Reset() {}
