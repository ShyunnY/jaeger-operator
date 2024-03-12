package message

import (
	"github.com/ShyunnY/jaeger-operator/internal/logging"
	"github.com/telepresenceio/watchable"
)

func SubscriptionIR[K comparable, V any](
	subscription <-chan watchable.Snapshot[K, V],
	handlerFunc func(update watchable.Update[K, V], errCh chan error),
) {

	errChan := make(chan error, 20)
	go func() {
		for err := range errChan {
			logging.DefaultLogger().WithName("message").Error(err, "observed an error")
		}
	}()

	if snapshot, ok := <-subscription; ok {
		for k, v := range snapshot.State {
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
			handlerFunc(
				update,
				errChan,
			)
		}
	}

}
