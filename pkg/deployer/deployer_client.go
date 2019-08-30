// package deployer is responsible for subscribing to a gateway via GRPC connection
// and updating the internal cache of what should be deployed
package deployer

import (
	"context"
	"fmt"
	"strings"
	"sync"
	"time"

	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/connectivity"
	"google.golang.org/grpc/credentials"

	deployer_v1 "github.com/dotmesh-io/ds-deployer/apis/deployer/v1"
	"github.com/dotmesh-io/ds-deployer/pkg/stopper"
	"github.com/dotmesh-io/ds-deployer/pkg/timeutil"
	"github.com/dotmesh-io/ds-deployer/pkg/version"
)

const (
	maxBackOff = 10 * time.Second
)

type ObjectCache interface {
	Insert(obj interface{}) bool
}

type Opts struct {
	Addr        string
	Token       string
	RequireTLS  bool
	ObjectCache ObjectCache
	Logger      *zap.SugaredLogger
}

// TODO:
// 1. delete deployments from cache
// 2. signal readiness when should operator start managing k8s resources
// so we don't start deleting them on boot
type DefaultClient struct {
	opts *Opts

	conn *grpc.ClientConn

	dialOpts []grpc.DialOption

	client deployer_v1.DeployerClient

	connectedMu sync.Mutex
	connected   bool

	objectCache ObjectCache

	logger *zap.SugaredLogger
}

func New(opts *Opts) *DefaultClient {

	dialOpts := []grpc.DialOption{
		grpc.WithPerRPCCredentials(&LoginCreds{
			Token:      opts.Token,
			RequireTLS: opts.RequireTLS,
		}),
		grpc.WithUserAgent(fmt.Sprintf("client/ds-deployer-%s", version.GetVersion().Version)),
		grpc.WithBackoffMaxDelay(30 * time.Second),
		WithKeepAliveDialer(),
	}

	if opts.RequireTLS {
		opts.Logger.Info("TLS connection requirement set")
		dialOpts = append(dialOpts, grpc.WithTransportCredentials(credentials.NewClientTLSFromCert(nil, "")))
	} else {
		dialOpts = append(dialOpts, grpc.WithInsecure())
	}

	return &DefaultClient{
		opts:        opts,
		dialOpts:    dialOpts,
		objectCache: opts.ObjectCache,
		logger:      opts.Logger,
	}
}

func (c *DefaultClient) startPeriodicSync(ctx context.Context) {
	ticker := time.NewTicker(10 * time.Second)
	defer ticker.Stop()

	for range ticker.C {
		deploymentsResp, err := c.client.ListDeployments(ctx, &deployer_v1.DeploymentFilter{})
		if err != nil {
			c.logger.Errorw("failed to retrieve deployments",
				"error", err,
			)
			continue
		}

		for _, d := range deploymentsResp.Deployments {
			// c.logger.Infow("configured deployment detected",
			// 	"image_name", d.GetImageName(),
			// 	"id", d.GetId(),
			// 	"name", d.GetName(),
			// 	"namespace", d.GetNamespace(),
			// 	"ingress_host", d.GetIngress().GetHost(),
			// 	"ingress_class", d.GetIngress().GetClass(),
			// )
			c.objectCache.Insert(d)
		}

		// TODO: check what's in the cache and remove anything that shouldn't be there anymore

	}
}

func (c *DefaultClient) StartDeployer(ctx context.Context) error {

	attempts := 0
	var backOff time.Duration

	stp := stopper.NewStopper(ctx)

RECONNECT:

	c.connectedMu.Lock()
	c.connected = false
	c.connectedMu.Unlock()

	err := c.dial(ctx)
	if err != nil {
		c.logger.Errorw("dial failed",
			"error", err,
		)
		return err
	}

	go c.startPeriodicSync(ctx)

	for {
		select {
		case <-ctx.Done():
			return err
		default:
			c.connectedMu.Lock()
			c.connected = true
			c.connectedMu.Unlock()

			fl := &deployer_v1.DeploymentFilter{}
			err = c.getDeployments(ctx, fl)
			if err != nil {
				if strings.Contains(err.Error(), "unauthorized") {
					return fmt.Errorf("deployer authentication failed, check your deployer token")
				}

				c.logger.Errorw("can't open connection to stream deployment requests, retrying...",
					"error", err,
					"attempts", attempts,
					"address", c.opts.Addr,
				)

				backOff = timeutil.ExpBackoff(backOff, maxBackOff)
				attempts++

				stp.Sleep(backOff)
				goto RECONNECT
			}
		}
	}
}

func (c *DefaultClient) dial(ctx context.Context) error {
	dialCtx, dialCancel := context.WithTimeout(ctx, time.Second*10)
	defer dialCancel()
	conn, err := grpc.DialContext(dialCtx, c.opts.Addr, c.dialOpts...)
	if err != nil {
		return err
	}
	if c.conn != nil {
		c.conn.Close()
	}

	c.conn = conn

	c.client = deployer_v1.NewDeployerClient(c.conn)

	cancel := func() {
		err := conn.Close()
		if err != nil {
			c.logger.Errorw("failed to close connection",
				"error", err,
			)
		}
	}

	go c.monitorHealth(ctx, conn, cancel)

	return nil
}

// OK is used by the health check to determine whether the runner is working
func (c *DefaultClient) OK() bool {
	if c.conn != nil {
		return c.conn.GetState() == connectivity.Ready
	}
	return false
}
