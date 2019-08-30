package controller

import (
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/api/extensions/v1beta1"
	"k8s.io/client-go/tools/cache"

	deployer_v1 "github.com/dotmesh-io/ds-deployer/apis/deployer/v1"

	"github.com/dotmesh-io/ds-deployer/pkg/cond"
)

const DEFAULT_INGRESS_CLASS = "ds-deployer"

// A KubernetesCache holds Kubernetes objects and associated configuration and produces
// DAG values.
type KubernetesCache struct {

	// Deployer's IngressClass.
	// If not set, defaults to DEFAULT_INGRESS_CLASS.
	IngressClass string

	ingresses        map[Meta]*v1beta1.Ingress
	deployments      map[Meta]*appsv1.Deployment
	services         map[Meta]*corev1.Service
	modelDeployments map[Meta]*deployer_v1.Deployment

	cond.Cond

	// hash of the api key
	controllerIdentifier string
}

func NewKubernetesCache(controllerIdentifier string) *KubernetesCache {
	return &KubernetesCache{controllerIdentifier: controllerIdentifier}
}

// Meta holds the name and namespace of a Kubernetes object.
type Meta struct {
	name, namespace string
}

func (kc *KubernetesCache) OnAdd(obj interface{}) {
	_ = kc.Insert(obj)
}

func (kc *KubernetesCache) OnUpdate(_, newObj interface{}) {
	_ = kc.Insert(newObj)
}

func (kc *KubernetesCache) OnDelete(obj interface{}) {
	_ = kc.Remove(obj)
}

// Insert inserts obj into the KubernetesCache.
// Insert returns true if the cache accepted the object, or false if the value
// is not interesting to the cache. If an object with a matching type, name,
// and namespace exists, it will be overwritten.
func (kc *KubernetesCache) Insert(obj interface{}) bool {
	// notify subscribers about cache changes
	defer kc.Notify()

	switch obj := obj.(type) {
	case *corev1.Service:

		if getDeployerID(obj.GetAnnotations()) != kc.controllerIdentifier {
			return false
		}

		m := Meta{name: obj.Name, namespace: obj.Namespace}
		if kc.services == nil {
			kc.services = make(map[Meta]*corev1.Service)
		}
		kc.services[m] = obj
		// return kc.serviceTriggersRebuild(obj)
		return true
	case *v1beta1.Ingress:

		if getDeployerID(obj.GetAnnotations()) != kc.controllerIdentifier {
			return false
		}

		class := getIngressClassAnnotation(obj.Annotations)
		if class != "" && class != kc.ingressClass() {
			return false
		}

		m := Meta{name: obj.Name, namespace: obj.Namespace}
		if kc.ingresses == nil {
			kc.ingresses = make(map[Meta]*v1beta1.Ingress)
		}
		kc.ingresses[m] = obj
		return true
	case *appsv1.Deployment:

		if getDeployerID(obj.GetAnnotations()) != kc.controllerIdentifier {
			return false
		}

		m := Meta{name: obj.Name, namespace: obj.Namespace}
		if kc.deployments == nil {
			kc.deployments = make(map[Meta]*appsv1.Deployment)
		}
		kc.deployments[m] = obj
		return true
	case *deployer_v1.Deployment:
		m := Meta{name: obj.Name, namespace: obj.Namespace}
		if kc.modelDeployments == nil {
			kc.modelDeployments = make(map[Meta]*deployer_v1.Deployment)
		}
		kc.modelDeployments[m] = obj
		return true
	default:
		// not an interesting object
		return false
	}
}

// Remove removes obj from the KubernetesCache.
// Remove returns a boolean indiciating if the cache changed after the remove operation.
func (kc *KubernetesCache) Remove(obj interface{}) bool {
	switch obj := obj.(type) {
	default:
		return kc.remove(obj)
	case cache.DeletedFinalStateUnknown:
		return kc.Remove(obj.Obj) // recurse into ourselves with the tombstoned value
	}
}

func (kc *KubernetesCache) remove(obj interface{}) bool {
	switch obj := obj.(type) {
	case *corev1.Service:
		m := Meta{name: obj.Name, namespace: obj.Namespace}
		_, ok := kc.services[m]
		delete(kc.services, m)
		return ok
	case *v1beta1.Ingress:
		m := Meta{name: obj.Name, namespace: obj.Namespace}
		_, ok := kc.ingresses[m]
		delete(kc.ingresses, m)
		return ok
	case *appsv1.Deployment:
		m := Meta{name: obj.Name, namespace: obj.Namespace}
		_, ok := kc.deployments[m]
		delete(kc.deployments, m)
		return ok
	case *deployer_v1.Deployment:
		m := Meta{name: obj.Name, namespace: obj.Namespace}
		_, ok := kc.modelDeployments[m]
		delete(kc.modelDeployments, m)
		return ok
	default:
		// not interesting
		return false
	}
}

// ModelDeployments returns model deployments
func (kc *KubernetesCache) ModelDeployments() []*deployer_v1.Deployment {
	var deployments []*deployer_v1.Deployment

	for _, v := range kc.modelDeployments {
		deployments = append(deployments, v)
	}

	return deployments
}

// ingressClass returns the IngressClass
// or DEFAULT_INGRESS_CLASS if not configured.
func (kc *KubernetesCache) ingressClass() string {
	return stringOrDefault(kc.IngressClass, DEFAULT_INGRESS_CLASS)
}

func stringOrDefault(s, def string) string {
	if s == "" {
		return def
	}
	return s
}

// getIngressClassAnnotation checks for the acceptable ingress class annotations
// 1. ds-deployer.dotscience.com/ingress.class
// 2. kubernetes.io/ingress.class
//
// it returns the first matching ingress annotation (in the above order) with test
func getIngressClassAnnotation(annotations map[string]string) string {
	class, ok := annotations["ds-deployer.dotscience.com/ingress.class"]
	if ok {
		return class
	}

	class, ok = annotations["kubernetes.io/ingress.class"]
	if ok {
		return class
	}

	return ""
}
