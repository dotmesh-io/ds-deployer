package controller

import (
	"context"
	"errors"
	"fmt"
	"time"

	"go.uber.org/zap"

	appsv1 "k8s.io/api/apps/v1"
	apiMachineryErrors "k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/client-go/kubernetes"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"

	"github.com/dotmesh-io/ds-deployer/pkg/health"
	"github.com/dotmesh-io/ds-deployer/pkg/logger"
	"github.com/dotmesh-io/ds-deployer/pkg/status"
)

var (
	ErrGatewayConnNotOK = errors.New("gateway module is not healthy, cannot sync")
)

// Controller reconciles Deployments
type Controller struct {
	// client can be used to retrieve objects from the APIServer.
	// client client.Client
	client client.Client

	clientSet *kubernetes.Clientset

	logger *zap.SugaredLogger
	cache  *KubernetesCache
	// status cache
	statusCache status.Cache

	gatewayConnModule health.Module

	// hash of the api key
	controllerIdentifier string
}

func New(opts ...Option) (*Controller, error) {
	c := new(Controller)
	for _, opt := range opts {
		opt(c)
	}

	if c.client == nil {
		return nil, fmt.Errorf("Kubernetes client is missing")
	}
	if c.logger == nil {
		c.logger = logger.GetInstance().Sugar()
	}

	return c, nil
}

func (c *Controller) Start(ctx context.Context) error {
	ticker := time.NewTicker(30 * time.Second)
	defer ticker.Stop()
	v := event{}
	ch := make(chan event, 1)
	c.register(ctx, ch)

	for {
		select {
		case <-ctx.Done():
			return nil
		case v = <-ch:
			err := c.sync()
			if err != nil {
				c.logger.Errorw("failed to process deployment diff",
					"error", err,
				)
				continue
			}
			c.logger.Debugw("cache changes detected, synchronizing",
				"version", v.Version,
				"resource", v.Resource,
			)
		case <-ticker.C:

			err := c.sync()
			if err != nil {
				c.logger.Errorw("failed to process cache diff",
					"error", err,
				)
			}
		}
	}
}

type event struct {
	Version  int
	Resource string
}

func (c *Controller) register(ctx context.Context, ch chan event) {
	go func() {
		bch := make(chan int, 1)
		bLast := 0
		for {
			c.cache.Register(bch, bLast)
			select {
			case <-ctx.Done():
				return
			case bLast = <-bch:
				ch <- event{Version: bLast, Resource: "cache"}
			}
		}
	}()
}

func (c *Controller) sync() error {

	if !c.gatewayConnModule.OK() {
		return ErrGatewayConnNotOK
	}

	// check deployments
	err := c.synchronizeDeployments()
	if err != nil {
		c.logger.Errorw("failed to synchronize deployments",
			"error", err,
		)
	}

	// check services
	err = c.synchronizeServices()
	if err != nil {
		c.logger.Errorw("failed to synchronize services",
			"error", err,
		)
	}

	// check ingresses
	err = c.synchronizeIngresses()
	if err != nil {
		c.logger.Errorw("failed to synchronize ingresses",
			"error", err,
		)
	}

	return nil
}

func (c *Controller) Reconcile(request reconcile.Request) (reconcile.Result, error) {
	// set up a convenient log object so we don't have to type request over and over again
	log := c.logger.With("request", request)

	log.Infow("Reconciling object", "namespaced_name", request.NamespacedName)

	// Fetch the Deployment from the cache
	deployment := &appsv1.Deployment{}
	err := c.client.Get(context.TODO(), request.NamespacedName, deployment)
	if apiMachineryErrors.IsNotFound(err) {
		log.Warnw("Could not find Deployment",
			"name", request.NamespacedName,
		)
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
	err = c.client.Update(context.TODO(), deployment)
	if err != nil {
		log.Error(err, "Could not write Deployment")
		return reconcile.Result{}, err
	}

	return reconcile.Result{}, nil
}
