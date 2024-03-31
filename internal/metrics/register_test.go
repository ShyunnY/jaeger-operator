package metrics

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"go.opentelemetry.io/otel/attribute"
	"k8s.io/apimachinery/pkg/util/rand"
	ctrl "sigs.k8s.io/controller-runtime"

	"github.com/ShyunnY/jaeger-operator/internal/config"
	"github.com/ShyunnY/jaeger-operator/internal/logging"
)

func TestNoopMetrics(t *testing.T) {
	err := New(&config.Server{})
	assert.NoError(t, err)
}

func TestNewPrometheus(t *testing.T) {

	ctx := ctrl.SetupSignalHandler()

	err := New(&config.Server{
		Logger: logging.NewLogger("debug"),
		Observability: config.Observability{
			Metric: &config.Metrics{},
		},
	})
	assert.NoError(t, err)

	counter := NewCounter("request_total", "handler request count", WithAttribute(attribute.String("backend", "login")))
	gauge := NewGauge("message_use_rate", "use rate in message handler of disk", WithAttribute(attribute.String("message", "disk")))

	go func() {
		for {
			counter.Increment()
			gauge.Record(float64(rand.Int()))
			time.Sleep(time.Second)
		}
	}()

	<-ctx.Done()
}
