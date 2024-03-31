package metrics

import (
	"context"
	"net"
	"net/http"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/otlp/otlpmetric/otlpmetricgrpc"
	"go.opentelemetry.io/otel/exporters/otlp/otlpmetric/otlpmetrichttp"
	otelprom "go.opentelemetry.io/otel/exporters/prometheus"
	"go.opentelemetry.io/otel/sdk/metric"
	ctrlmetrics "sigs.k8s.io/controller-runtime/pkg/metrics"

	"github.com/ShyunnY/jaeger-operator/internal/config"
	"github.com/ShyunnY/jaeger-operator/internal/consts"
	"github.com/ShyunnY/jaeger-operator/internal/logging"
)

var (
	meterProvider   = otel.GetMeterProvider().Meter("jaeger-operator")
	metricLogger    = logging.NewLogger(consts.LogLevelDInfo).WithName("metrics")
	metricsEndpoint = "/metrics"
)

// TODO: 需要将controller-runtime中的metrics端点改变到此server上
func init() {
	otel.SetLogger(metricLogger.Logger)
}

type Options struct {
	address string

	sink struct {
		host     string
		port     string
		protocol string
	}

	prometheus struct {
		enableProm bool
		registry   prometheus.Registerer
		gatherer   prometheus.Gatherer
	}
}

func (o *Options) enableSink() bool {
	if len(o.sink.host) != 0 &&
		len(o.sink.port) != 0 {
		return true
	}

	return false
}

// New Create the metrics service and create the corresponding metrics Exporter.
// By default, only the Prometheus Exporter is enabled (of course, this can be explicitly disabled).
// If you have a Sink configured, you can export metrics to Collectors that support the OTLP protocol.
func New(cfg *config.Server) error {

	opts := applyConfig(cfg)
	if err := registerProvider(opts); err != nil {
		return err
	}

	// if Prometheus is not disabled, we start a web server to provide to prom for pull
	if !cfg.DisablePrometheus() {
		return promServer(opts)
	}

	return nil
}

func applyConfig(cfg *config.Server) *Options {
	opts := &Options{}

	// we only configure the Sink if it is configured.
	// otherwise we will not register the otel exporter
	if cfg.Metric.Sink != nil {
		opts.sink.host = consts.OtelHost
		opts.sink.port = consts.OtelPort
		opts.sink.protocol = consts.OtelProtol

		if cfg.Metric != nil && len(cfg.Metric.Sink.Host) != 0 {
			opts.sink.host = cfg.Metric.Sink.Host
		}

		if cfg.Metric != nil && len(cfg.Metric.Sink.Port) != 0 {
			opts.sink.port = cfg.Metric.Sink.Port
		}

		if cfg.Metric != nil && len(cfg.Metric.Sink.Protocol) != 0 {
			opts.sink.protocol = cfg.Metric.Sink.Protocol
		}
	}

	// if Prometheus is not disabled, we build the values required for prom
	if !cfg.DisablePrometheus() {
		opts.address = net.JoinHostPort(consts.MetricsHost, consts.MetricsPort)
		opts.prometheus.enableProm = true
		opts.prometheus.registry = ctrlmetrics.Registry
		opts.prometheus.gatherer = ctrlmetrics.Registry
	}

	return opts
}

func registerProvider(opts *Options) error {

	metricOpts := []metric.Option{}
	if opts.enableSink() {
		switch opts.sink.protocol {
		case "http", "https":
			if err := registerOTELHTTPExporter(opts, &metricOpts); err != nil {
				return err
			}
		case "grpc":
			if err := registerOTELGRPCExporter(opts, &metricOpts); err != nil {
				return err
			}
		}
	}

	if opts.prometheus.enableProm {
		if err := registerPrometheus(opts, &metricOpts); err != nil {
			return err
		}
	}
	meterProvider := metric.NewMeterProvider(metricOpts...)
	otel.SetMeterProvider(meterProvider)

	return nil
}

// promServer Run the prometheus service
func promServer(opts *Options) error {

	metricLogger.Info("start prom metrics server", "address", opts.address)

	handler := promhttp.HandlerFor(
		opts.prometheus.gatherer,
		promhttp.HandlerOpts{},
	)

	mux := http.NewServeMux()
	mux.Handle(metricsEndpoint, handler)

	sre := http.Server{
		Handler:           mux,
		Addr:              opts.address,
		ReadTimeout:       5 * time.Second,
		ReadHeaderTimeout: 5 * time.Second,
		WriteTimeout:      10 * time.Second,
		IdleTimeout:       15 * time.Second,
	}

	go func() {
		if err := sre.ListenAndServe(); err != nil {
			metricLogger.Error(err, "failed to start the prom metrics server")
		}
	}()

	return nil
}

// registerPrometheus Register the prometheus metrics exporter
func registerPrometheus(opts *Options, metricOpts *[]metric.Option) error {

	promExporter, err := otelprom.New(
		otelprom.WithoutTargetInfo(),
		otelprom.WithoutScopeInfo(),
		otelprom.WithoutCounterSuffixes(),
		otelprom.WithRegisterer(opts.prometheus.registry),
	)
	if err != nil {
		return err
	}

	*metricOpts = append(*metricOpts, metric.WithReader(promExporter))
	metricLogger.Info("build the otel metrics pull endpoint", "address", opts.address)

	return nil
}

// registerOTELHTTPExporter Register the otel metrics http exporter
func registerOTELHTTPExporter(opts *Options, metricOpts *[]metric.Option) error {

	address := net.JoinHostPort(opts.sink.host, opts.sink.port)
	httpExporter, err := otlpmetrichttp.New(
		context.Background(),
		otlpmetrichttp.WithEndpoint(address),
		otlpmetrichttp.WithInsecure(),
	)
	if err != nil {
		return err
	}

	reader := metric.NewPeriodicReader(httpExporter)
	*metricOpts = append(*metricOpts, metric.WithReader(reader))
	metricLogger.Info("build the otel metrics http push endpoint", "address", address)

	return nil
}

// registerOTELGRPCExporter Register the otel metrics gRPC exporter
func registerOTELGRPCExporter(opts *Options, metricOpts *[]metric.Option) error {

	address := net.JoinHostPort(opts.sink.host, opts.sink.port)
	gRPCExporter, err := otlpmetricgrpc.New(
		context.Background(),
		otlpmetricgrpc.WithEndpoint(address),
		otlpmetricgrpc.WithInsecure(),
	)
	if err != nil {
		return err
	}

	reader := metric.NewPeriodicReader(gRPCExporter)
	*metricOpts = append(*metricOpts, metric.WithReader(reader))
	metricLogger.Info("build the otel metrics gRPC push endpoint", "address", address)

	return nil
}
