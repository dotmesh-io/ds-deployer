package controller

import (
	"context"
	"strings"
	"sync"

	deployer_v1 "github.com/dotmesh-io/ds-deployer/apis/deployer/v1"
	"github.com/dotmesh-io/ds-deployer/pkg/status"
	"k8s.io/api/extensions/v1beta1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/intstr"
)

const (
	kubernetesIngressClassAnnotation = "kubernetes.io/ingress.class"
)

func (c *Controller) synchronizeIngresses() error {

	var wg sync.WaitGroup

	for _, modelDeployment := range c.cache.ModelDeployments() {
		existing, ok := c.cache.GetIngress(modelDeployment.Namespace, getDeploymentName(modelDeployment))
		if !ok {
			if c.statusCache.Get(modelDeployment.Id).Ingress != status.StatusConfiguring {
				c.statusCache.Set(modelDeployment.Id, status.ModuleIngress, status.StatusConfiguring)
			}

			wg.Add(1)
			// creating new deployment
			go func(modelDeployment *deployer_v1.Deployment) {
				err := c.createIngress(modelDeployment)
				if err != nil {
					c.logger.Errorw("failed to create ingress",
						"error", err,
						"deployment_id", modelDeployment.GetId(),
					)
				}
				wg.Done()
			}(modelDeployment)
			// should get created if it doesn't exist
			continue
		}

		c.logger.Debugf("ingress %s/%s found, checking for updates", existing.Namespace, existing.Name)

		if !ingressesEqual(toKubernetesIngress(modelDeployment, c.controllerIdentifier), existing) {
			if c.statusCache.Get(modelDeployment.Id).Ingress != status.StatusConfiguring {
				c.statusCache.Set(modelDeployment.Id, status.ModuleIngress, status.StatusConfiguring)
			}

			updatedIngress := updateIngress(existing, modelDeployment)

			wg.Add(1)
			go func(updatedIngress *v1beta1.Ingress) {
				err := c.client.Update(context.Background(), updatedIngress)
				if err != nil {
					c.logger.Errorw("failed to update service",
						"error", err,
						"service_namespace", modelDeployment.Namespace,
						"service_name", updatedIngress.GetName(),
						"deployment_id", modelDeployment.GetId(),
					)
				}
				wg.Done()
			}(updatedIngress)
		} else {
			if c.statusCache.Get(modelDeployment.Id).Ingress != status.StatusReady {
				c.statusCache.Set(modelDeployment.Id, status.ModuleIngress, status.StatusReady)
			}
		}
	}

	// going through existing ingresses to see which ones should
	// be removed
	for _, ingress := range c.cache.ingresses {

		if ingress.GetAnnotations() == nil {
			continue
		}

		_, ok := c.cache.modelDeployments[ingress.GetAnnotations()["deployment"]]
		if !ok {
			// not found in model deployments, should delete
			c.logger.Infof("ingress %s/%s not found in model deployments, deleting", ingress.GetNamespace(), ingress.GetName())
			err := c.client.Delete(context.Background(), ingress)
			if err != nil {
				if strings.Contains(err.Error(), "not found") {
					// it's fine
					continue
				}

				c.logger.Errorw("failed to delete ingress",
					"error", err,
					"name", ingress.GetName(),
					"namespace", ingress.GetNamespace(),
				)
			}
		}
	}

	wg.Wait()

	return nil
}

func (c *Controller) createIngress(modelDeployment *deployer_v1.Deployment) error {
	return c.client.Create(context.Background(), toKubernetesIngress(modelDeployment, c.controllerIdentifier))
}

func getIngressSpec(md *deployer_v1.Deployment) v1beta1.IngressSpec {
	var modelPort intstr.IntOrString

	for _, p := range md.Service.GetPorts() {
		if p.GetName() == "model-http" {
			modelPort = intstr.IntOrString{
				Type:   intstr.Int,
				IntVal: p.GetPort(),
			}
		}
	}
	return v1beta1.IngressSpec{
		Rules: []v1beta1.IngressRule{
			{
				Host: md.Ingress.GetHost(),
				IngressRuleValue: v1beta1.IngressRuleValue{
					HTTP: &v1beta1.HTTPIngressRuleValue{
						Paths: []v1beta1.HTTPIngressPath{
							{
								Path: "/",
								Backend: v1beta1.IngressBackend{
									ServiceName: getDeploymentName(md),
									ServicePort: modelPort,
								},
							},
						},
					},
				},
			},
		},
	}
}

func ingressesEqual(desired, existing *v1beta1.Ingress) bool {
	if len(desired.Spec.Rules) != len(existing.Spec.Rules) {
		return false
	}

	if desired.GetAnnotations() == nil || existing.GetAnnotations() == nil {
		return false
	}

	if desired.GetAnnotations()[kubernetesIngressClassAnnotation] != existing.GetAnnotations()[kubernetesIngressClassAnnotation] {
		return false
	}

	for i := range desired.Spec.Rules {
		if desired.Spec.Rules[i].Host != existing.Spec.Rules[i].Host {
			return false
		}
		if len(desired.Spec.Rules[i].HTTP.Paths) != len(existing.Spec.Rules[i].HTTP.Paths) {
			return false
		}

		for idx := range desired.Spec.Rules[i].HTTP.Paths {
			if desired.Spec.Rules[i].HTTP.Paths[idx].Path != existing.Spec.Rules[i].HTTP.Paths[idx].Path {
				return false
			}
			if desired.Spec.Rules[i].HTTP.Paths[idx].Backend.ServiceName != existing.Spec.Rules[i].HTTP.Paths[idx].Backend.ServiceName {
				return false
			}
			if desired.Spec.Rules[i].HTTP.Paths[idx].Backend.ServicePort != existing.Spec.Rules[i].HTTP.Paths[idx].Backend.ServicePort {
				return false
			}
		}
	}

	return true
}

func updateIngress(existing *v1beta1.Ingress, md *deployer_v1.Deployment) *v1beta1.Ingress {
	updated := existing.DeepCopy()
	updated.Annotations[kubernetesIngressClassAnnotation] = md.Ingress.GetClass()
	updated.Spec = getIngressSpec(md)

	return updated
}

func toKubernetesIngress(md *deployer_v1.Deployment, controllerIdentifier string) *v1beta1.Ingress {

	annotations := map[string]string{
		AnnControllerIdentifier: controllerIdentifier,
		// based on model deployment name we will need this later
		// to ensure we delete what's not needed anymore
		"name":                           md.GetName(),
		"deployment":                     md.GetId(),
		kubernetesIngressClassAnnotation: md.Ingress.GetClass(),
	}

	if md.Ingress.GetClass() == "nginx" {
		annotations["nginx.ingress.kubernetes.io/proxy-body-size"] = "100m"
	}

	ingress := &v1beta1.Ingress{
		TypeMeta: metav1.TypeMeta{},
		ObjectMeta: metav1.ObjectMeta{
			Name:      getDeploymentName(md),
			Namespace: md.GetNamespace(),
			Labels: map[string]string{
				"owner": "ds-deployer",
			},
			Annotations: annotations,
		},
		Spec: getIngressSpec(md),
	}

	return ingress
}
