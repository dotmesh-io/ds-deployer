package deployer

import (
	"context"
	"net"
	"time"

	"google.golang.org/grpc"

	"github.com/dotmesh-io/ds-deployer/pkg/version"
)

type LoginCreds struct {
	Token      string
	RequireTLS bool
}

func (c *LoginCreds) GetRequestMetadata(context.Context, ...string) (map[string]string, error) {
	return map[string]string{
		"token":   c.Token,
		"version": version.GetVersion().Version,
	}, nil
}

func (c *LoginCreds) RequireTransportSecurity() bool {
	return c.RequireTLS
}

// WithKeepAliveDialer - required so connections aren't dropped
// more info: https://github.com/grpc/grpc-java/issues/1648
func WithKeepAliveDialer() grpc.DialOption {
	return grpc.WithDialer(func(addr string, timeout time.Duration) (net.Conn, error) {
		d := net.Dialer{Timeout: timeout, KeepAlive: time.Duration(10 * time.Second)}
		return d.Dial("tcp", addr)
	})
}
