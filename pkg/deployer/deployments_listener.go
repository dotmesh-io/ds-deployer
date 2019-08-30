package deployer

import (
	"context"
	"fmt"
	"io"

	deployer_v1 "github.com/dotmesh-io/ds-deployer/apis/deployer/v1"
)

func (c *DefaultClient) getDeployments(ctx context.Context, filter *deployer_v1.DeploymentFilter) error {

	// calling the streaming API
	stream, err := c.client.StreamDeployments(ctx, filter)
	if err != nil {
		return fmt.Errorf("error while getting deployments: %s", err)
	}

	c.logger.Info("listening for deployments...")

	for {
		select {
		case <-ctx.Done():
			return nil
		default:
			deployment, err := stream.Recv()
			if err == io.EOF {
				break
			}
			if err != nil {
				// log.WithFields(log.Fields{
				// 	"error":   err,
				// 	"address": c.opts.Addr,
				// }).Error("failed to get stream from server")
				c.logger.Errorw("failed to establish deployments stream",
					"error", err,
					"addr", c.opts.Addr,
				)
				return err
			}

			c.logger.Infow("new deployment received",
				"name", deployment.GetName(),
				"namespace", deployment.GetNamespace(),
				"id", deployment.GetId(),
			)
			c.objectCache.Insert(deployment)
		}
	}
}
