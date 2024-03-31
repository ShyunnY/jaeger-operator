package message

import (
	"github.com/ShyunnY/jaeger-operator/internal/logging"
	"github.com/ShyunnY/jaeger-operator/internal/metrics"
	"github.com/telepresenceio/watchable"
)

var (
	messageCounter    = metrics.NewCounter("message_ir_count", "the number of ir accepted by message")
	messageErrCounter = metrics.NewCounter("message_err_count", "the number of errors handled by the handler after message accepts ir")
)

func SubscriptionIR[K comparable, V any](
	subscription <-chan watchable.Snapshot[K, V],
	handlerFunc func(update watchable.Update[K, V], errCh chan error),
) {

	errChan := make(chan error, 20)
	go func() {
		for err := range errChan {
			messageErrCounter.Increment()
			logging.DefaultLogger().WithName("message").Error(err, "observed an error")
		}
	}()

	if snapshot, ok := <-subscription; ok {
		for k, v := range snapshot.State {
			messageCounter.Increment()
			handlerFunc(
				watchable.Update[K, V]{
					Key:   k,
					Value: v,
				},
				errChan,
			)
		}
	}

	for snapshot := range subscription {
		for _, update := range snapshot.Updates {
			messageCounter.Increment()
			handlerFunc(
				update,
				errChan,
			)
		}
	}

}
