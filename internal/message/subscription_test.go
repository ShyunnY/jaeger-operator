package message

import (
	"context"
	"fmt"
	"github.com/telepresenceio/watchable"
	"testing"
	"time"
)

func TestSubscription(t *testing.T) {

	m := watchable.Map[string, string]{}

	go SubscriptionIR(m.Subscribe(context.TODO()),
		func(update watchable.Update[string, string], errCh chan error) {
			if update.Delete {
				return
			}

			fmt.Printf("handler: %+v\n", update.Key)
		})
	time.Sleep(time.Second)

	m.Store("z3", "20")
	m.Store("z6", "23")
	time.Sleep(time.Second)
	m.Store("s7", "24")
	m.Store("l4", "21")
	m.Store("w5", "22")
	m.Delete("z3")
	time.Sleep(time.Second)
}
