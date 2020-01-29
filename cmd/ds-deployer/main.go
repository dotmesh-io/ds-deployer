package main

import (
	"context"
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"os"
	"path/filepath"
	"time"

	kingpin "gopkg.in/alecthomas/kingpin.v2"

	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/api/extensions/v1beta1"
	"k8s.io/client-go/kubernetes"
	_ "k8s.io/client-go/plugin/pkg/client/auth/gcp"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
	"sigs.k8s.io/controller-runtime/pkg/client/config"
	"sigs.k8s.io/controller-runtime/pkg/controller"
	"sigs.k8s.io/controller-runtime/pkg/handler"
	"sigs.k8s.io/controller-runtime/pkg/manager"
	"sigs.k8s.io/controller-runtime/pkg/source"

	deploymentController "github.com/dotmesh-io/ds-deployer/internal/controller"
	"github.com/dotmesh-io/ds-deployer/pkg/deployer"
	"github.com/dotmesh-io/ds-deployer/pkg/health"
	"github.com/dotmesh-io/ds-deployer/pkg/logger"
	"github.com/dotmesh-io/ds-deployer/pkg/status"
	"github.com/dotmesh-io/ds-deployer/pkg/workgroup"
)

const EnvGatewayAddress = "GATEWAY_ADDRESS"
const EnvAuthToken = "TOKEN"
const EnvHealthPort = "HEALTH_PORT"

const gatewayServerAddress = "cloud.dotscience.com:8800"
const controllerName = "deployment-controller"

func main() {

	app := kingpin.New("ds-deployer", "Dotscience runner")

	run := app.Command("run", "Start the deployer")
	token := run.Flag("token", "Authentication token (each registered runner gets a token)").Default(os.Getenv(EnvAuthToken)).String()
	requireTLS := run.Flag("require-tls", "Require TLS for connection to the server").Default("true").Bool()
	serverAddr := run.Flag("addr", "Server address").Default(gatewayServerAddress).String()

	gracePeriod := run.Flag("grace-period", "Grace period before starting deployment synchronization").Default("10").Int()
	kubeconfig := run.Flag("kubeconfig", "path to kubeconfig (if not in running inside a cluster)").Default(filepath.Join(os.Getenv("HOME"), ".kube", "config")).String()
	inCluster := run.Flag("incluster", "use in cluster configuration.").Default("true").Bool()

	healthServerPort := run.Flag("health-port", "Health server port").Default("9300").OverrideDefaultFromEnvar(EnvHealthPort).String()
	metricsServerUser := run.Flag("metrics-user", "Metrics server username").Default("admin").String()
	metricsServerPassword := run.Flag("metrics-password", "Metrics server password").Default("metrics").String()

	logger := logger.GetInstance().Sugar()

	args := os.Args[1:]
	switch kingpin.MustParse(app.Parse(args)) {
	default:
		app.Usage(args)
		os.Exit(2)
	case run.FullCommand():

		if *token == "" {
			logger.Errorf("token not supplied, use --token <YOUR TOKEN> or environment variable '%s' to specify the token", EnvAuthToken)
			os.Exit(1)
		}
		// Setup a Manager
		logger.Info("setting up manager...!")
		cfg := config.GetConfigOrDie()
		logger.Info("got to 0")

		mgr, err := manager.New(cfg, manager.Options{
			Port:               7777,
			MetricsBindAddress: "0",
		})
		if err != nil {
			logger.Errorw("unable to set up overall controller manager",
				"error", err,
			)
			os.Exit(1)
		}
		logger.Info("got to 1")
		controllerIdentifier := getMD5Hash(*token)
		logger.Info("got to 2")

		statusCache := status.New()
		logger.Info("got to 3")
		cache := deploymentController.NewKubernetesCache(controllerIdentifier, logger.With("module", "cache"))
		logger.Info("got to 4")

		healthServer := health.NewServer(&health.Opts{
			Port:        *healthServerPort,
			Logger:      logger,
			Username:    *metricsServerUser,
			Password:    *metricsServerPassword,
			ObjectCache: cache,
		})
		logger.Info("got to 5")

		gatewayAddress := *serverAddr
		if os.Getenv(EnvGatewayAddress) != "" {
			gatewayAddress = os.Getenv(EnvGatewayAddress)
		}
		logger.Info("got to 6")

		kubeClient := newClient(*kubeconfig, *inCluster)
		logger.Info("got to 7")

		controllerOptions := []deploymentController.Option{
			deploymentController.WithIdentifier(controllerIdentifier),
			deploymentController.WithClient(mgr.GetClient()),
			deploymentController.WithCache(cache),
			deploymentController.WithLogger(logger.With("module", "deployment-reconciler")),
			deploymentController.WithStatusCache(statusCache),
			deploymentController.WithGatewayModule(healthServer),
			deploymentController.WithClientSet(kubeClient),
		}
		logger.Info("got to 8")

		deploymentReconciler, err := deploymentController.New(controllerOptions...)
		if err != nil {
			logger.Errorw("unable to set up dotscience deployment controller",
				"error", err,
			)
			os.Exit(1)
		}
		logger.Info("got to 9")

		gatewayClient := deployer.New(&deployer.Opts{
			Addr:          gatewayAddress,
			Token:         *token,
			RequireTLS:    *requireTLS,
			ObjectCache:   cache,
			StatusCache:   statusCache,
			PodLogsGetter: deploymentReconciler,
			Logger:        logger,
		})
		logger.Info("got to 10")

		healthServer.RegisterModule("gateway_conn", gatewayClient)
		logger.Info("got to 11")

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

		var g workgroup.Group

		buf := deploymentController.NewBuffer(&g, cache, logger, 128)

		deploymentController.WatchServices(&g, kubeClient, logger, buf)
		deploymentController.WatchDeployments(&g, kubeClient, logger, buf)
		deploymentController.WatchIngress(&g, kubeClient, logger, buf)

		g.Add(func(stop <-chan struct{}) error {
			ctx, cancel := context.WithCancel(context.Background())

			go func() {
				<-stop
				cancel()
			}()

			return gatewayClient.StartDeployer(ctx)
		})

		g.Add(func(stop <-chan struct{}) error {
			logger.Info("health server starting...")
			defer logger.Info("health server stopped")
			go func() {
				<-stop
				err := healthServer.Stop()
				if err != nil {
					logger.Errorw("error while stopping health server",
						"error", err,
					)
				}
			}()

			return healthServer.Start()
		})

		// start controller
		g.Add(func(stop <-chan struct{}) error {
			logger.Info("starting controller")
			defer logger.Info("controller stopped")
			ctx, cancel := context.WithCancel(context.Background())

			go func() {
				<-stop
				logger.Info("signal received, stopping controller")
				cancel()
			}()

			time.Sleep(time.Duration(*gracePeriod) * time.Second)

			return deploymentReconciler.Start(ctx)
		})

		err = g.Run()
		if err != nil {
			logger.Errorf("deployer stopped with an error: %s", err)
			os.Exit(1)
		}
	}

}

func newClient(kubeconfig string, inCluster bool) *kubernetes.Clientset {
	var err error
	var config *rest.Config
	if kubeconfig != "" && !inCluster {
		config, err = clientcmd.BuildConfigFromFlags("", kubeconfig)
		check(err)
	} else {
		config, err = rest.InClusterConfig()
		check(err)
	}

	client, err := kubernetes.NewForConfig(config)
	check(err)
	return client
}

func check(err error) {
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func getMD5Hash(text string) string {
	hasher := md5.New()
	hasher.Write([]byte(text))
	return hex.EncodeToString(hasher.Sum(nil))
}
