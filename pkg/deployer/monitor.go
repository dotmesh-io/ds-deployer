package deployer

import (
	"context"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/health/grpc_health_v1"
)

func (c *DefaultClient) monitorHealth(ctx context.Context, cc *grpc.ClientConn, cancelConn func()) {
	defer cancelConn()
	defer cc.Close()

	ticker := time.NewTicker(800 * time.Millisecond)
	defer ticker.Stop()
	healthClient := grpc_health_v1.NewHealthClient(cc)

	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			checkCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
			_, err := healthClient.Check(checkCtx, &grpc_health_v1.HealthCheckRequest{})
			cancel()
			if err != nil {
				c.logger.Errorf("gateway healthcheck failed: %s", err)
				return
			}
		}
	}
}
