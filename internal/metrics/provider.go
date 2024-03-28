package metrics

import "go.opentelemetry.io/otel"

var (
	meterProvider = otel.GetMeterProvider()
)

// TODO: 包装meter的三种类型, 以及添加对应的元数据Metadata
