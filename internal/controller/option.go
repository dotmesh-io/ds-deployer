package controller

import (
	// "github.com/go-logr/logr"
	"go.uber.org/zap"

	"sigs.k8s.io/controller-runtime/pkg/client"
)

// Option configures a controller option.
type Option func(*Reconciler)

func WithClient(client client.Client) Option {
	return func(r *Reconciler) {
		r.client = client
	}
}

func WithLogger(logger *zap.SugaredLogger) Option {
	return func(r *Reconciler) {
		r.logger = logger
	}
}

func WithCache(cache *KubernetesCache) Option {
	return func(r *Reconciler) {
		r.cache = cache
	}
}
