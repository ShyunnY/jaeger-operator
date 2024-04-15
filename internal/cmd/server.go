package cmd

import (
	"github.com/spf13/cobra"
	ctrl "sigs.k8s.io/controller-runtime"

	"github.com/ShyunnY/jaeger-operator/internal/config"
	"github.com/ShyunnY/jaeger-operator/internal/consts"
	infrarunner "github.com/ShyunnY/jaeger-operator/internal/infra/runner"
	kubernetesrunner "github.com/ShyunnY/jaeger-operator/internal/kubernetes/runner"
	"github.com/ShyunnY/jaeger-operator/internal/logging"
	"github.com/ShyunnY/jaeger-operator/internal/message"
	"github.com/ShyunnY/jaeger-operator/internal/metrics"
	statusrunner "github.com/ShyunnY/jaeger-operator/internal/status/runner"
	"github.com/ShyunnY/jaeger-operator/internal/tracing"
	translaterunner "github.com/ShyunnY/jaeger-operator/internal/translate/runner"
	"github.com/ShyunnY/jaeger-operator/internal/utils"
)

var (
	logLevel          string
	namespace         string
	testTraceEndpoint string
)

func GetServerCommand() *cobra.Command {

	sreCmd := &cobra.Command{
		Use:     "server",
		Aliases: []string{"serve"},
		Short:   "Serve Jaeger Controller",
		RunE: func(*cobra.Command, []string) error {
			return server()
		},
	}

	sreCmd.Flags().StringVarP(&logLevel, "log-level", "v", "info", "config log output level")
	sreCmd.Flags().StringVarP(&namespace, "namespaces", "n", "", "config watch namespace, use ',' to separate multiple namespaces, empty means watch all (e.g. prod, dev)")
	sreCmd.Flags().StringVarP(&testTraceEndpoint, "endpoint", "e", "", "test trace endpoint")

	return sreCmd
}

// server serves jaeger operator
func server() error {

	cfg := &config.Server{
		JaegerOperatorName: consts.OperatorName,
		Logger:             logging.NewLogger(consts.LogLevel(logLevel)).WithName(Name()),
		NamespaceSet:       utils.ExtractNamespace(namespace),
		Observability: config.Observability{
			Metric: &config.Metrics{},
			Trace: &config.Traces{
				Endpoint: &testTraceEndpoint,
			},
		},
	}

	// running admin serve
	if err := NewAdmin(cfg); err != nil {
		return err
	}

	// tracing serve
	if err := tracing.New(cfg); err != nil {
		return err
	}

	// metrics serve
	if err := metrics.New(cfg); err != nil {
		return err
	}

	// starting runners
	if err := setRunners(cfg); err != nil {
		return err
	}

	return nil
}

// setRunners Start all runners required for the Jaeger Operator service to running
func setRunners(cfg *config.Server) error {

	ctx := ctrl.SetupSignalHandler()

	// init IRs
	irMessage := new(message.IRMessage)
	infraIRs := new(message.InfraIRMaps)
	statusIRs := new(message.StatusIRMaps)

	// 1. kubernetes controller runner
	k8sRunner := kubernetesrunner.New(&kubernetesrunner.Config{
		Server:    *cfg,
		IrMessage: irMessage,
	})
	if err := k8sRunner.Start(ctx); err != nil {
		return err
	}

	// 2. translator runner
	translatorRunner := translaterunner.New(&translaterunner.Config{
		Server:      *cfg,
		InfraIRMap:  infraIRs,
		StatusIRMap: statusIRs,
		IrMessage:   irMessage,
	})
	if err := translatorRunner.Start(ctx); err != nil {
		return err
	}

	// 3. infrastructure runner
	infraRunner := infrarunner.New(&infrarunner.Config{
		Server:    *cfg,
		InfraMap:  infraIRs,
		StatusMap: statusIRs,
	})
	if err := infraRunner.Start(ctx); err != nil {
		return err
	}

	// 4. status handler runner
	statusRunner := statusrunner.New(&statusrunner.Config{
		Server:      *cfg,
		StatusIRMap: statusIRs,
	})
	if err := statusRunner.Start(ctx); err != nil {
		return err
	}

	// wait until ctx done...
	<-ctx.Done()

	// TODO: post clean work

	cfg.Logger.Info("jaeger operator shutting down")

	return nil
}
