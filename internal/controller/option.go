package controller

import (
	"go.uber.org/zap"
	"sigs.k8s.io/controller-runtime/pkg/client"

	"github.com/dotmesh-io/ds-deployer/pkg/status"
)

// Option configures a controller option.
type Option func(*Controller)

func WithClient(client client.Client) Option {
	return func(c *Controller) {
		c.client = client
	}
}

func WithLogger(logger *zap.SugaredLogger) Option {
	return func(c *Controller) {
		c.logger = logger
	}
}

func WithCache(cache *KubernetesCache) Option {
	return func(c *Controller) {
		c.cache = cache
	}
}

func WithStatusCache(cache status.Cache) Option {
	return func(c *Controller) {
		c.statusCache = cache
	}
}

func WithIdentifier(identifier string) Option {
	return func(c *Controller) {
		c.controllerIdentifier = identifier
	}
}
