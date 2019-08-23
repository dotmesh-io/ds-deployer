package controller

import (
	"github.com/go-logr/logr"

	"sigs.k8s.io/controller-runtime/pkg/client"
)

// Option configures a controller option.
type Option func(*Reconciler)

func WithClient(client client.Client) Option {
	return func(r *Reconciler) {
		r.client = client
	}
}

func WithLogger(log logr.Logger) Option {
	return func(r *Reconciler) {
		r.log = log
	}
}
