package controller

import (
	"sync"

	"github.com/jinzhu/copier"
	"go.uber.org/zap"
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/api/extensions/v1beta1"
	"k8s.io/client-go/tools/cache"

	deployer_v1 "github.com/dotmesh-io/ds-deployer/apis/deployer/v1"

	"github.com/dotmesh-io/ds-deployer/pkg/cond"
)

// A KubernetesCache holds Kubernetes objects and associated configuration and produces
// DAG values.
type KubernetesCache struct {
	ingressesMu *sync.RWMutex
	ingresses   map[Meta]*v1beta1.Ingress

	deploymentsMu *sync.RWMutex
	deployments   map[Meta]*appsv1.Deployment

	servicesMu *sync.RWMutex
	services   map[Meta]*corev1.Service

	modelDeploymentsMu *sync.RWMutex
	modelDeployments   map[Meta]*deployer_v1.Deployment

	cond.Cond

	// hash of the api key
	controllerIdentifier string

	logger *zap.SugaredLogger
}

func NewKubernetesCache(controllerIdentifier string, logger *zap.SugaredLogger) *KubernetesCache {
	return &KubernetesCache{
		ingressesMu:          &sync.RWMutex{},
		deploymentsMu:        &sync.RWMutex{},
		servicesMu:           &sync.RWMutex{},
		controllerIdentifier: controllerIdentifier,
		modelDeploymentsMu:   &sync.RWMutex{},
		logger:               logger,
	}
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
		kc.servicesMu.Lock()
		if kc.services == nil {
			kc.services = make(map[Meta]*corev1.Service)
		}
		kc.services[m] = obj
		kc.servicesMu.Unlock()
		// return kc.serviceTriggersRebuild(obj)
		return true
	case *v1beta1.Ingress:
		if getDeployerID(obj.GetAnnotations()) != kc.controllerIdentifier {
			return false
		}

		m := Meta{name: obj.Name, namespace: obj.Namespace}
		kc.ingressesMu.Lock()
		if kc.ingresses == nil {
			kc.ingresses = make(map[Meta]*v1beta1.Ingress)
		}
		kc.ingresses[m] = obj
		kc.ingressesMu.Unlock()
		return true
	case *appsv1.Deployment:
		if getDeployerID(obj.GetAnnotations()) != kc.controllerIdentifier {
			return false
		}

		m := Meta{name: obj.Name, namespace: obj.Namespace}

		kc.deploymentsMu.Lock()
		if kc.deployments == nil {
			kc.deployments = make(map[Meta]*appsv1.Deployment)
		}
		kc.deployments[m] = obj
		kc.deploymentsMu.Unlock()
		return true
	case *deployer_v1.Deployment:
		kc.modelDeploymentsMu.Lock()
		m := Meta{name: obj.Name, namespace: obj.Namespace}
		if kc.modelDeployments == nil {
			kc.modelDeployments = make(map[Meta]*deployer_v1.Deployment)
		}
		kc.modelDeployments[m] = obj
		kc.modelDeploymentsMu.Unlock()
		return true
	default:
		// not an interesting object
		return false
	}
}

// Remove removes obj from the KubernetesCache.
// Remove returns a boolean indiciating if the cache changed after the remove operation.
func (kc *KubernetesCache) Remove(obj interface{}) bool {
	defer kc.Notify()
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
		kc.modelDeploymentsMu.Lock()
		m := Meta{name: obj.Name, namespace: obj.Namespace}
		_, ok := kc.modelDeployments[m]
		delete(kc.modelDeployments, m)
		kc.modelDeploymentsMu.Unlock()
		return ok
	default:
		// not interesting
		return false
	}
}

// ModelDeployments returns model deployments
func (kc *KubernetesCache) ModelDeployments() []*deployer_v1.Deployment {
	var deployments []*deployer_v1.Deployment

	kc.modelDeploymentsMu.RLock()
	for _, v := range kc.modelDeployments {

		var cp deployer_v1.Deployment

		err := copier.Copy(&cp, v)
		if err != nil {
			kc.logger.Errorw("failed to copy deployment",
				"id", v.GetId(),
				"error", zap.Error(err),
			)
			continue
		}

		deployments = append(deployments, &cp)
	}
	kc.modelDeploymentsMu.RUnlock()

	return deployments
}
