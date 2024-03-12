package infra

import (
	"context"
	"fmt"

	"k8s.io/client-go/util/retry"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

type Client struct {
	client.Client
}

type DeltaObjectFunc func() *InventoryObject

func New(cli client.Client) *Client {
	return &Client{Client: cli}
}

func (c *Client) CreateOrUpdateOrDelete(ctx context.Context, deltaHandler DeltaObjectFunc) error {

	inventoryObject := deltaHandler()
	createObjs, updateObjs, deleteObjs := inventoryObject.CreateObjects, inventoryObject.UpdateObjects, inventoryObject.DeleteObjects

	return retry.RetryOnConflict(retry.DefaultRetry, func() error {

		if createObjs != nil && len(createObjs) != 0 {
			for i := range createObjs {
				createObj := createObjs[i]
				if err := c.Client.Create(ctx, createObj); err != nil {
					return fmt.Errorf("create obj error: %w", err)
				}
			}
		}

		if updateObjs != nil && len(updateObjs) != 0 {
			for i := range updateObjs {
				updateObject := updateObjs[i]
				if err := c.Update(ctx, updateObject); err != nil {
					return fmt.Errorf("update obj error: %w", err)
				}
			}
		}

		if deleteObjs != nil && len(deleteObjs) != 0 {
			for i := range deleteObjs {
				deleteObj := deleteObjs[i]
				if err := c.Delete(ctx, deleteObj); err != nil {
					return fmt.Errorf("delete obj error: %w", err)
				}
			}
		}

		return nil
	})
}
