package controller

import (
	// "github.com/go-logr/logr"
	"go.uber.org/zap"

	"sigs.k8s.io/controller-runtime/pkg/client"
)

// Option configures a controller option.
type Option func(*Controller)

func WithClient(client client.Client) Option {
	return func(r *Controller) {
		r.client = client
	}
}

func WithLogger(logger *zap.SugaredLogger) Option {
	return func(r *Controller) {
		r.logger = logger
	}
}

func WithCache(cache *KubernetesCache) Option {
	return func(r *Controller) {
		r.cache = cache
	}
}
