package controller

import (
	"context"
	"fmt"

	"go.uber.org/zap"

	appsv1 "k8s.io/api/apps/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"

	"github.com/dotmesh-io/ds-deployer/pkg/logger"
)

// Reconciler reconciles Deployments
type Reconciler struct {
	// client can be used to retrieve objects from the APIServer.
	client client.Client
	logger *zap.SugaredLogger
	cache  *KubernetesCache
}

func New(opts ...Option) (*Reconciler, error) {
	r := new(Reconciler)
	for _, opt := range opts {
		opt(r)
	}

	if r.client == nil {
		return nil, fmt.Errorf("Kubernetes client is missing")
	}
	if r.logger == nil {
		r.logger = logger.GetInstance().Sugar()
	}

	return r, nil
}

// Implement reconcile.Reconciler so the controller can reconcile objects
var _ reconcile.Reconciler = &Reconciler{}

func (r *Reconciler) Reconcile(request reconcile.Request) (reconcile.Result, error) {
	// set up a convenient log object so we don't have to type request over and over again
	log := r.logger.With("request", request)

	// Fetch the Deployment from the cache
	deployment := &appsv1.Deployment{}
	err := r.client.Get(context.TODO(), request.NamespacedName, deployment)
	if errors.IsNotFound(err) {
		log.Error("Could not find Deployment")
		return reconcile.Result{}, nil
	}

	if err != nil {
		log.Errorf("Could not fetch Deployment: %s", err)
		return reconcile.Result{}, err
	}

	// Print the Deployment
	log.Infow("Reconciling Deployment", "container name", deployment.Spec.Template.Spec.Containers[0].Name)

	// Set the label if it is missing
	if deployment.Labels == nil {
		deployment.Labels = map[string]string{}
	}
	if deployment.Labels["heritage"] == "deployer.dotscience.com" {
		return reconcile.Result{}, nil
	}

	// Update the Deployment
	deployment.Labels["heritage"] = "deployer.dotscience.com"
	err = r.client.Update(context.TODO(), deployment)
	if err != nil {
		log.Error(err, "Could not write Deployment")
		return reconcile.Result{}, err
	}

	return reconcile.Result{}, nil
}
