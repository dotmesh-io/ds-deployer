package deployer

import (
	"context"
	"fmt"
	"io"

	deployer_v1 "github.com/dotmesh-io/ds-deployer/apis/deployer/v1"
)

type PodLogsGetter interface {
	Logs(md *deployer_v1.Deployment, request *deployer_v1.LogsRequest) (io.ReadCloser, error)
}

func (c *DefaultClient) getLogRequests(ctx context.Context, filter *deployer_v1.LogsFilter) error {

	// calling the streaming API
	stream, err := c.client.StreamLogRequests(ctx, filter)
	if err != nil {
		return fmt.Errorf("error while getting log requests: %s", err)
	}

	c.logger.Info("listening for log requests...")

	for {
		select {
		case <-ctx.Done():
			return nil
		default:
			logsRequest, err := stream.Recv()
			if err == io.EOF {
				break
			}
			if err != nil {
				// log.WithFields(log.Fields{
				// 	"error":   err,
				// 	"address": c.opts.Addr,
				// }).Error("failed to get stream from server")
				c.logger.Errorw("failed to establish log requests stream",
					"error", err,
					"addr", c.opts.Addr,
				)
				return err
			}

			c.logger.Infow("new logs request received",
				"id", logsRequest.GetDeploymentId(),
				"tx_id", logsRequest.GetTxId(),
			)
			// c.objectCache.Insert(deployment)
		}
	}
}
