package metrics

import (
	"context"
	"testing"
	"time"

	ctrl "sigs.k8s.io/controller-runtime"

	"github.com/ShyunnY/jaeger-operator/internal/config"
	"github.com/ShyunnY/jaeger-operator/internal/logging"
	"github.com/stretchr/testify/assert"
)

func TestNewPrometheus(t *testing.T) {

	ctx := ctrl.SetupSignalHandler()

	err := New(&config.Server{
		Logger: logging.NewLogger("debug"),
		Observability: config.Observability{
			Metric: &config.Metrics{
				Protocol: "xxx-ooo",
			},
		},
	})
	assert.NoError(t, err)

	counter, err := meterProvider.Meter("envoy-gateway").Float64Counter("request_total")
	if err != nil {
		panic(err)
	}

	go func() {
		for {
			counter.Add(context.Background(), 1)
			time.Sleep(time.Second)
		}
	}()

	<-ctx.Done()
}
