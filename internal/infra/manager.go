package infra

import (
	"context"
	"github.com/ShyunnY/jaeger-operator/internal/logging"
	"github.com/ShyunnY/jaeger-operator/internal/message"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

type Manager struct {
	logger      logging.Logger
	cli         client.Client
	infraClient Client
}

func NewManager(cli client.Client, logger logging.Logger) *Manager {
	return &Manager{
		cli:         cli,
		logger:      logger,
		infraClient: Client{Client: cli},
	}
}

func (m *Manager) BuildInfraResources(ctx context.Context, infraIR *message.InfraIR) error {

	// TODO: 如果计算出错, 我们需要进行condition发布

	ic := InventoryComputer{
		namespace:    infraIR.InstanceNamespace,
		instanceName: infraIR.InstanceName,
		cli:          m.cli,
	}

	// create service account
	if saObj, err := ic.ComputeServiceAccount(ctx, infraIR.ServiceAccount); err != nil {
		m.logger.Error(err, "failed to compute ServiceAccount")

		return err
	} else {
		if err := m.infraClient.CreateOrUpdateOrDelete(ctx, func() *InventoryObject {
			return saObj
		}); err != nil {
			m.logger.Error(err, "failed to create or update ServiceAccount")

			return err
		}
	}

	// create deployment
	if deployObj, err := ic.ComputeDeployment(ctx, infraIR.Deployment); err != nil {
		m.logger.Error(err, "failed to compute Deployment")

		return err
	} else {
		if err := m.infraClient.CreateOrUpdateOrDelete(ctx, func() *InventoryObject {
			return deployObj
		}); err != nil {
			m.logger.Error(err, "failed to create or update Deployment")

			return err
		}
	}

	// create service
	if servicesObj, err := ic.ComputeService(ctx, infraIR.Service); err != nil {
		m.logger.Error(err, "failed to compute Service")

		return err
	} else {
		if err := m.infraClient.CreateOrUpdateOrDelete(ctx, func() *InventoryObject {
			return servicesObj
		}); err != nil {
			m.logger.Error(err, "failed to create or update Service")

			return err
		}
	}

	return nil
}
