package metrics

import "go.opentelemetry.io/otel/attribute"

type MetricOption struct {
	Unit string
	Attr []attribute.KeyValue
}

type MetricOptions func(*MetricOption)

func WithUnit(unit string) MetricOptions {
	return func(o *MetricOption) {
		o.Unit = unit
	}
}

func WithAttribute(attr ...attribute.KeyValue) MetricOptions {
	return func(o *MetricOption) {
		if o.Attr != nil && len(o.Attr) > 0 {
			o.Attr = append(o.Attr, attr...)
		} else {
			o.Attr = attr
		}
	}
}

func applyOption(opts ...MetricOptions) *MetricOption {
	metricOpt := &MetricOption{
		Attr: []attribute.KeyValue{},
	}
	for _, opt := range opts {
		opt(metricOpt)
	}

	// TODO: do we need check empty value?

	return metricOpt
}
