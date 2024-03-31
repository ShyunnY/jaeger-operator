package metrics

import (
	"context"

	"go.opentelemetry.io/otel/attribute"
	otelmetrics "go.opentelemetry.io/otel/metric"
)

const CounterMetricType = "counter"

type Counter struct {
	name    string
	attr    []attribute.KeyValue
	counter otelmetrics.Float64Counter
}

func NewCounter(name, description string, opts ...MetricOptions) *Counter {
	option := applyOption(opts...)
	c, err := meterProvider.Float64Counter(
		name,
		otelmetrics.WithDescription(description),
		otelmetrics.WithUnit(option.Unit),
	)
	if err != nil {
		metricLogger.Error(err, "failed to create new metric", "type", CounterMetricType, "name", name)
	}

	return &Counter{
		name:    name,
		attr:    option.Attr,
		counter: c,
	}
}

func (c *Counter) Add(val float64) {
	if c.attr != nil && len(c.attr) != 0 {
		c.counter.Add(
			context.Background(),
			val,
			otelmetrics.WithAttributes(c.attr...),
		)
	} else {
		c.counter.Add(
			context.Background(),
			val,
		)
	}
}

func (c *Counter) Increment() {
	c.Add(1)
}

func (c *Counter) Decrement() {
	c.Add(-1)
}
