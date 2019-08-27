package main

import (
	"os"

	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/api/extensions/v1beta1"
	_ "k8s.io/client-go/plugin/pkg/client/auth/gcp"
	"sigs.k8s.io/controller-runtime/pkg/client/config"
	"sigs.k8s.io/controller-runtime/pkg/controller"
	"sigs.k8s.io/controller-runtime/pkg/handler"
	"sigs.k8s.io/controller-runtime/pkg/manager"
	"sigs.k8s.io/controller-runtime/pkg/manager/signals"
	"sigs.k8s.io/controller-runtime/pkg/source"

	deploymentController "github.com/dotmesh-io/ds-deployer/internal/controller"
	"github.com/dotmesh-io/ds-deployer/pkg/logger"
)

const controllerName = "deployment-controller"

func main() {

	logger := logger.GetInstance().Sugar()

	// Setup a Manager
	logger.Info("setting up manager")
	mgr, err := manager.New(config.GetConfigOrDie(), manager.Options{})
	if err != nil {
		logger.Errorw("unable to set up overall controller manager",
			"error", err,
		)
		os.Exit(1)
	}

	controllerOptions := []deploymentController.Option{
		deploymentController.WithClient(mgr.GetClient()),
		deploymentController.WithLogger(logger.With("module", "deployment-reconciler")),
	}

	deploymentReconciler, err := deploymentController.New(controllerOptions...)
	if err != nil {
		logger.Errorw("unable to set up dotscience deployment controller",
			"error", err,
		)
		os.Exit(1)
	}

	// Setup a new controller to reconcile dotscience deployments
	logger.Info("Setting up controller")
	c, err := controller.New(controllerName, mgr, controller.Options{
		Reconciler: deploymentReconciler,
	})
	if err != nil {
		logger.Errorw("unable to set up individual controller",
			"error", err,
		)
		os.Exit(1)
	}

	// Watch Deployment and enqueue ReplicaSet object key
	if err := c.Watch(&source.Kind{Type: &appsv1.Deployment{}}, &handler.EnqueueRequestForObject{}); err != nil {
		logger.Error("unable to watch Deployment",
			"error", err,
		)
		os.Exit(1)
	}

	if err := c.Watch(&source.Kind{Type: &v1beta1.Ingress{}}, &handler.EnqueueRequestForObject{}); err != nil {
		logger.Error("unable to watch Ingress",
			"error", err,
		)
		os.Exit(1)
	}

	if err := c.Watch(&source.Kind{Type: &corev1.Service{}}, &handler.EnqueueRequestForObject{}); err != nil {

		logger.Error("unable to watch Service",
			"error", err,
		)
		os.Exit(1)
	}

	// // Watch Pods and enqueue owning ReplicaSet key
	// if err := c.Watch(&source.Kind{Type: &corev1.Pod{}},
	// 	&handler.EnqueueRequestForOwner{OwnerType: &appsv1.ReplicaSet{}, IsController: true}); err != nil {
	// 	entryLog.Error(err, "unable to watch Pods")
	// 	os.Exit(1)
	// }

	logger.Info("starting manager")
	if err := mgr.Start(signals.SetupSignalHandler()); err != nil {
		logger.Error("unable to run manager",
			"error", err,
		)
		os.Exit(1)
	}
}
