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
	metricsEndpoint = "/metrics"
	metricsLogger   = "metrics"
)

type Options struct {
	logger  logging.Logger
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

func New(cfg *config.Server) error {

	opts := applyConfig(cfg)
	if err := registerProvider(opts); err != nil {
		return err
	}

	if !cfg.DisablePrometheus() {
		return promServer(opts)
	}

	return nil
}

func applyConfig(cfg *config.Server) *Options {
	opts := &Options{
		logger:  cfg.Logger.WithName(metricsLogger),
		address: net.JoinHostPort(consts.MetricsHost, consts.MetricsPort),
	}

	opts.sink.host = consts.OtelHost
	opts.sink.port = consts.OtelPort
	opts.sink.protocol = consts.OtelProtol

	if cfg.Metric != nil && len(cfg.Metric.Host) != 0 {
		opts.sink.host = cfg.Metric.Host
	}

	if cfg.Metric != nil && len(cfg.Metric.Port) != 0 {
		opts.sink.port = cfg.Metric.Port
	}

	if cfg.Metric != nil && len(cfg.Metric.Protocol) != 0 {
		opts.sink.protocol = cfg.Metric.Protocol
	}

	if !cfg.DisablePrometheus() {
		opts.prometheus.enableProm = true
		opts.prometheus.registry = ctrlmetrics.Registry
		opts.prometheus.gatherer = ctrlmetrics.Registry
	}

	return opts
}

func registerProvider(opts *Options) error {

	metricOpts := []metric.Option{}
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

	opts.logger.Info("start prom metrics server", "address", opts.address)

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
			opts.logger.Error(err, "failed to start the prom metrics server")
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
	opts.logger.Info("build the otel metrics pull endpoint", "address", opts.address)

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
	opts.logger.Info("build the otel metrics http push endpoint", "address", address)

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
	opts.logger.Info("build the otel metrics gRPC push endpoint", "address", address)

	return nil
}
