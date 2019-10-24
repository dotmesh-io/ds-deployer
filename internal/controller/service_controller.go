package controller

import (
	"context"
	"sync"

	deployer_v1 "github.com/dotmesh-io/ds-deployer/apis/deployer/v1"
	"github.com/dotmesh-io/ds-deployer/pkg/status"

	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/intstr"
)

func (c *Controller) synchronizeServices() error {
	var wg sync.WaitGroup

	for _, modelDeployment := range c.cache.ModelDeployments() {
		existing, ok := c.cache.GetService(modelDeployment.Namespace, getDeploymentName(modelDeployment))
		if !ok {
			if c.statusCache.Get(modelDeployment.Id).Service != status.StatusConfiguring {
				c.statusCache.Set(modelDeployment.Id, status.ModuleService, status.StatusConfiguring)
			}

			wg.Add(1)
			// creating new deployment
			go func(modelDeployment *deployer_v1.Deployment) {
				err := c.createService(modelDeployment)
				if err != nil {
					c.logger.Errorw("failed to create service",
						"error", err,
						"deployment_id", modelDeployment.GetId(),
					)
				}
				wg.Done()
			}(modelDeployment)
			// should get created if it doesn't exist
			continue
		}

		c.logger.Debugf("service %s/%s found, checking for updates", existing.Namespace, existing.Name)

		if !servicesEqual(toKubernetesService(modelDeployment, c.controllerIdentifier), existing) {
			if c.statusCache.Get(modelDeployment.Id).Service != status.StatusConfiguring {
				c.statusCache.Set(modelDeployment.Id, status.ModuleService, status.StatusConfiguring)
			}

			updatedService := updateService(existing, modelDeployment)

			wg.Add(1)
			go func(updatedService *corev1.Service) {
				err := c.client.Update(context.Background(), updatedService)
				if err != nil {
					c.logger.Errorw("failed to update service",
						"error", err,
						"service_namespace", modelDeployment.Namespace,
						"service_name", updatedService.GetName(),
						"deployment_id", modelDeployment.GetId(),
					)
				}
				wg.Done()
			}(updatedService)
		} else {
			if c.statusCache.Get(modelDeployment.Id).Service != status.StatusReady {
				c.statusCache.Set(modelDeployment.Id, status.ModuleService, status.StatusReady)
			}
		}
	}

	// going through existing services to see which ones should
	// be removed
	for meta, service := range c.cache.services {

		if service.GetAnnotations() == nil {
			continue
		}

		_, ok := c.cache.modelDeployments[Meta{namespace: meta.namespace, name: service.GetAnnotations()["name"]}]
		if !ok {
			// not found in model deployments, should delete
			c.logger.Infof("service %s/%s not found in model deployments, deleting", service.GetNamespace(), service.GetName())
			err := c.client.Delete(context.Background(), service)
			if err != nil {
				c.logger.Errorw("failed to delete service",
					"error", err,
					"name", service.GetName(),
					"namespace", service.GetNamespace(),
				)
			}
		}
	}

	wg.Wait()

	return nil
}

func (c *Controller) createService(modelDeployment *deployer_v1.Deployment) error {
	return c.client.Create(context.Background(), toKubernetesService(modelDeployment, c.controllerIdentifier))
}

func toKubernetesService(md *deployer_v1.Deployment, controllerIdentifier string) *corev1.Service {

	service := &corev1.Service{
		TypeMeta: metav1.TypeMeta{},
		ObjectMeta: metav1.ObjectMeta{
			Name:      getDeploymentName(md),
			Namespace: md.GetNamespace(),
			Labels: map[string]string{
				"owner": "ds-deployer",
			},
			Annotations: map[string]string{
				AnnControllerIdentifier: controllerIdentifier,
				// based on model deployment name we will need this later
				// to ensure we delete what's not needed anymore
				"name": md.GetName(),
			},
		},
		Spec: corev1.ServiceSpec{
			Ports: getServicePorts(md),
			Selector: map[string]string{
				"deployment": md.GetId(),
			},
			Type: getServiceType(md),
		},
	}

	return service
}

func servicesEqual(desired, existing *corev1.Service) bool {
	if desired.Spec.Type != existing.Spec.Type {
		return false
	}

	if desired.Spec.LoadBalancerIP != existing.Spec.LoadBalancerIP {
		return false
	}

	if len(desired.Spec.Ports) != len(existing.Spec.Ports) {
		return false
	}

	for i := range desired.Spec.Ports {
		if desired.Spec.Ports[i].Name != existing.Spec.Ports[i].Name {
			return false
		}
		if desired.Spec.Ports[i].Port != existing.Spec.Ports[i].Port {
			return false
		}
		if desired.Spec.Ports[i].TargetPort.IntVal != existing.Spec.Ports[i].TargetPort.IntVal {
			return false
		}

	}

	return true
}

func updateService(existing *corev1.Service, md *deployer_v1.Deployment) *corev1.Service {
	updated := existing.DeepCopy()

	updated.Spec.Ports = getServicePorts(md)
	updated.Spec.Type = getServiceType(md)

	return updated
}

func getServicePorts(md *deployer_v1.Deployment) []corev1.ServicePort {
	var servicePorts []corev1.ServicePort

	for _, p := range md.Service.GetPorts() {
		servicePorts = append(servicePorts, corev1.ServicePort{
			Name: p.GetName(),
			Port: p.GetPort(),
			TargetPort: intstr.IntOrString{
				IntVal: p.GetTargetPort(),
				Type:   intstr.Int,
			},
		})
	}

	return servicePorts
}

func getServiceType(md *deployer_v1.Deployment) corev1.ServiceType {

	var serviceType corev1.ServiceType
	switch md.Service.GetType() {
	case "loadbalancer":
		serviceType = corev1.ServiceTypeLoadBalancer
	default:
		serviceType = corev1.ServiceTypeNodePort
	}
	return serviceType
}
