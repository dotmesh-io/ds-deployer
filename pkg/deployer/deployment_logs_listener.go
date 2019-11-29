package deployer

import (
	"bufio"
	"context"
	"fmt"
	"io"

	deployer_v1 "github.com/dotmesh-io/ds-deployer/apis/deployer/v1"
)

type PodLogsGetter interface {
	Logs(request *deployer_v1.LogsRequest) (io.ReadCloser, error)
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

			go func() {
				err := c.processLogRequest(ctx, logsRequest)
				if err != nil {
					c.logger.Errorw("failed to process logs request",
						"id", logsRequest.GetDeploymentId(),
						"tx_id", logsRequest.GetTxId(),
						"error", err,
					)

				}
			}()
		}
	}
}

func (c *DefaultClient) processLogRequest(ctx context.Context, request *deployer_v1.LogsRequest) error {

	logStream, err := c.podLogsGetter.Logs(request)
	if err != nil {
		return err
	}
	defer logStream.Close()

	sendStream, err := c.client.SendLogs(ctx)
	if err != nil {
		return err
	}

	var createdIndex int64

	rd := bufio.NewReader(logStream)

	for {
		str, err := rd.ReadString('\n')
		if err != nil {
			if err == io.EOF {
				// sending last line
				return sendStream.Send(&deployer_v1.Logs{
					CreatedIndex: createdIndex,
					TxId:         request.TxId,
					Line:         str,
					Eof:          true,
				})
			}
			return err
		}
		err = sendStream.Send(&deployer_v1.Logs{
			TxId: request.TxId,
			Line: str,
		})
		if err != nil {
			return err
		}
		createdIndex++
	}
}
