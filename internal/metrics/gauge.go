package metrics

import (
	"context"
	"sync"

	"go.opentelemetry.io/otel/attribute"
	otelmetrics "go.opentelemetry.io/otel/metric"
)

const GaugeMetricType = "gauge"

type Gauge struct {
	name  string
	val   float64
	attr  []attribute.KeyValue
	gauge otelmetrics.Float64ObservableGauge
	mux   *sync.RWMutex
}

func NewGauge(name, description string, opts ...MetricOptions) *Gauge {
	option := applyOption(opts...)

	gau := &Gauge{
		name: name,
		attr: option.Attr,
		mux:  &sync.RWMutex{},
	}

	g, err := meterProvider.Float64ObservableGauge(
		name,
		otelmetrics.WithDescription(description),
		otelmetrics.WithUnit(option.Unit),
		otelmetrics.WithFloat64Callback(func(c context.Context, o otelmetrics.Float64Observer) error {
			gau.mux.RLock()
			defer gau.mux.RUnlock()

			o.Observe(gau.val, otelmetrics.WithAttributes(gau.attr...))
			return nil
		}),
	)
	if err != nil {
		metricLogger.Error(err, "failed to create new metric", "type", GaugeMetricType, "name", name)
	}
	gau.gauge = g

	return gau
}

func (g *Gauge) Record(val float64) {
	g.mux.Lock()
	defer g.mux.Unlock()

	g.val = val
}
