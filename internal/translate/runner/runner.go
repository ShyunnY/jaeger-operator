package runner

import (
	"context"
	"github.com/ShyunnY/jaeger-operator/internal/message"
	"github.com/ShyunnY/jaeger-operator/internal/translate"
	"github.com/telepresenceio/watchable"

	jaegerv1a1 "github.com/ShyunnY/jaeger-operator/api/v1alpha1"
	"github.com/ShyunnY/jaeger-operator/internal/config"
)

type Config struct {
	config.Server

	InfraIRMap  *message.InfraIRMaps
	StatusIRMap *message.StatusIRMaps
	IrMessage   *message.IRMessage
}

type Runner struct {
	Config
}

func New(cfg *Config) Runner {
	return Runner{Config: *cfg}
}

func (r *Runner) Name() string {
	return jaegerv1a1.TranslatorComponent
}

func (r *Runner) Start(ctx context.Context) error {
	r.Logger = r.Logger.WithName(r.Name())
	go r.translateResources(ctx)
	r.Logger.Info("translator started")

	return nil
}

func (r *Runner) translateResources(ctx context.Context) {
	translator := translate.Translator{
		InfraIRMap:  r.Config.InfraIRMap,
		StatusIRMap: r.Config.StatusIRMap,
	}

	message.SubscriptionIR(
		r.IrMessage.Subscribe(ctx),
		func(update watchable.Update[string, *jaegerv1a1.Jaeger], errCh chan error) {
			r.Logger.Info("translator takes the ir instance and handler it", "instance", update.Key)

			if update.Delete {
				// TODO: 跳过资源翻译, 我们删除各个irKey即可
			}

			// 将资源翻译成infraIR,statusIR等
			if err := translator.Translate(update.Value); err != nil {
				errCh <- err
			}
		},
	)
}
