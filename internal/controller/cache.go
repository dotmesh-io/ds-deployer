package controller

import (
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/api/extensions/v1beta1"

	deployer_v1 "github.com/dotmesh-io/ds-deployer/apis/deployer/v1"
)

// A KubernetesCache holds Kubernetes objects and associated configuration and produces
// DAG values.
type KubernetesCache struct {
	ingresses        map[Meta]*v1beta1.Ingress
	deployments      map[Meta]*appsv1.Deployment
	services         map[Meta]*corev1.Service
	modelDeployments map[Meta]*deployer_v1.Deployment
}

// Meta holds the name and namespace of a Kubernetes object.
type Meta struct {
	name, namespace string
}
