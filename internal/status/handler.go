package status

import (
	"context"
	"github.com/ShyunnY/jaeger-operator/internal/logging"
	kerrors "k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/client-go/util/retry"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

type Update struct {
	NamespacedName types.NamespacedName
	Object         client.Object
	Mutator        MutatorFunc
}
type MutatorFunc func(client.Object) client.Object

type UpdateHandler struct {
	client   client.Client
	logger   logging.Logger
	updateCh chan Update
	enabled  chan struct{}
}

func NewUpdateHandler(client client.Client, logger logging.Logger) *UpdateHandler {
	return &UpdateHandler{
		client:   client,
		logger:   logger,
		updateCh: make(chan Update),
		enabled:  make(chan struct{}),
	}
}

func (u *UpdateHandler) apply(upd Update) {
	if err := retry.RetryOnConflict(retry.DefaultRetry, func() error {
		obj := upd.Object

		if err := u.client.Get(context.TODO(), upd.NamespacedName, obj); err != nil {
			if kerrors.IsNotFound(err) {
				return nil
			}

			return err
		}

		newObj := upd.Mutator(obj)
		newObj.SetUID(obj.GetUID())

		// TODO:
		err := u.client.Status().Update(context.TODO(), newObj)
		if err != nil {
			return err
		}
		return nil
	}); err != nil {
		u.logger.Error(err, "failed to update object condition",
			"kind", upd.Object.GetObjectKind().GroupVersionKind().Kind,
			"name", upd.Object.GetName(),
		)
	}
}

func (u *UpdateHandler) Start(ctx context.Context) {

	u.logger.Info("status update handler started")
	defer u.logger.Info("status update handler shutdown")

	close(u.enabled)

	for {
		select {
		case <-ctx.Done():
			return
		case update := <-u.updateCh:
			u.apply(update)
		}
	}

}

func (u *UpdateHandler) Write(upd Update) {
	select {
	case <-u.enabled:
		u.updateCh <- upd
	default:

	}
}
