package controller

import (
	"context"
	"strings"
	"sync"

	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	deployer_v1 "github.com/dotmesh-io/ds-deployer/apis/deployer/v1"
)

func (c *Controller) synchronizeDeployments() error {

	var wg sync.WaitGroup

	for meta, modelDeployment := range c.cache.modelDeployments {
		// checking if we have this deployment
		existing, ok := c.cache.deployments[Meta{
			namespace: meta.namespace,
			name:      getDeploymentName(modelDeployment),
		}]
		if !ok {
			wg.Add(1)
			// creating new deployment
			go func(modelDeployment *deployer_v1.Deployment) {
				err := c.createDeployment(modelDeployment)
				if err != nil {
					c.logger.Errorw("failed to create deployment",
						"error", err,
						"deployment_id", modelDeployment.GetId(),
					)
				}
				wg.Done()
			}(modelDeployment)
			// should get created if it doesn't exist
			continue
		}

		c.logger.Debugf("deployment %s/%s found, checking for updates", existing.Namespace, existing.Name)

		if !deploymentsEqual(toKubernetesDeployment(modelDeployment, c.controllerIdentifier), existing) {
			updatedDeployment := updateDeployment(existing, modelDeployment)

			wg.Add(1)
			go func(updatedDeployment *appsv1.Deployment) {
				err := c.client.Update(context.Background(), updatedDeployment)
				if err != nil {
					c.logger.Errorw("failed to update deployment",
						"error", err,
						"deployment_id", modelDeployment.GetId(),
					)
				}
				wg.Done()
			}(updatedDeployment)
		}
	}

	// going through existing deployments to see which ones should
	// be removed
	for meta, deployment := range c.cache.deployments {

		if deployment.GetAnnotations() == nil {
			continue
		}

		_, ok := c.cache.modelDeployments[Meta{namespace: meta.namespace, name: deployment.GetAnnotations()["name"]}]
		if !ok {
			// not found in model deployments, should delete
			c.logger.Infof("deployment %s/%s not found in model deployments, deleting", deployment.GetNamespace(), deployment.GetName())
			err := c.client.Delete(context.Background(), deployment)
			if err != nil {
				c.logger.Errorw("failed to delete deployment",
					"error", err,
					"name", deployment.GetName(),
					"namespace", deployment.GetNamespace(),
				)
			}
		}
	}

	wg.Wait()

	return nil
}

func getDeploymentName(d *deployer_v1.Deployment) string {
	return "ds-" + d.GetName() + "-" + shortUUID(d.GetId())
}

func getPodName(d *deployer_v1.Deployment) string {
	return "ds-" + d.GetName() + "-" + shortUUID(d.GetId())
}

func shortUUID(u string) string {
	return strings.Split(u, "-")[0]
}

func (c *Controller) createDeployment(modelDeployment *deployer_v1.Deployment) error {
	return c.client.Create(context.Background(), toKubernetesDeployment(modelDeployment, c.controllerIdentifier))
}

func toKubernetesDeployment(modelDeployment *deployer_v1.Deployment, controllerIdentifier string) *appsv1.Deployment {

	cp := []corev1.ContainerPort{}

	for _, p := range modelDeployment.Deployment.GetPorts() {
		cp = append(cp, corev1.ContainerPort{
			ContainerPort: int32(p),
		})
	}

	deployment := &appsv1.Deployment{
		TypeMeta: metav1.TypeMeta{},
		ObjectMeta: metav1.ObjectMeta{
			Name:      getDeploymentName(modelDeployment),
			Namespace: modelDeployment.GetNamespace(),
			Labels: map[string]string{
				"owner": "ds-deployer",
			},
			Annotations: map[string]string{
				AnnControllerIdentifier: controllerIdentifier,
				// based on model deployment name we will need this later
				// to ensure we delete what's not needed anymore
				"name": modelDeployment.GetName(),
			},
		},
		Spec: appsv1.DeploymentSpec{
			Replicas: toInt32(int(modelDeployment.Deployment.GetReplicas())),
			Selector: &metav1.LabelSelector{
				MatchLabels: map[string]string{
					"deployment": modelDeployment.GetId(),
				},
			},
			Template: corev1.PodTemplateSpec{
				ObjectMeta: metav1.ObjectMeta{
					Labels: map[string]string{
						"deployment": modelDeployment.GetId(),
					},
				},
				Spec: corev1.PodSpec{
					ImagePullSecrets: []corev1.LocalObjectReference{
						// TODO: pass in secrets
					},
					Containers: []corev1.Container{
						corev1.Container{
							Name:  getPodName(modelDeployment),
							Image: modelDeployment.Deployment.GetImage(),
							Ports: cp,
						},
					},
				},
			},
		},
	}

	return deployment
}

// compares replicas, image, port, image pull secrets
func deploymentsEqual(desired, existing *appsv1.Deployment) bool {
	if desired.Spec.Replicas != existing.Spec.Replicas {
		return false
	}

	if desired.Spec.Template.GetLabels()["deployment"] != existing.Spec.Template.GetLabels()["deployment"] {
		return false
	}

	if len(desired.Spec.Template.Spec.ImagePullSecrets) != len(existing.Spec.Template.Spec.ImagePullSecrets) {
		return false
	}

	for i := range desired.Spec.Template.Spec.ImagePullSecrets {
		if existing.Spec.Template.Spec.ImagePullSecrets[i] != desired.Spec.Template.Spec.ImagePullSecrets[i] {
			return false
		}
	}

	// comparing images
	if len(desired.Spec.Template.Spec.Containers) != len(existing.Spec.Template.Spec.Containers) {
		return false
	}

	existingContainers := make(map[string]corev1.Container)

	for _, container := range existing.Spec.Template.Spec.Containers {
		existingContainers[container.Name] = container
	}

	for _, container := range desired.Spec.Template.Spec.Containers {
		existingContainer, ok := existingContainers[container.Name]
		if !ok {
			return false
		}
		if existingContainer.Image != container.Name {
			return false
		}
		if len(existingContainer.Ports) != len(container.Ports) {
			return false
		}
		for i := range container.Ports {
			if container.Ports[i] != existingContainer.Ports[i] {
				return false
			}
		}
	}

	return true
}

func updateDeployment(existing *appsv1.Deployment, md *deployer_v1.Deployment) *appsv1.Deployment {
	updated := existing.DeepCopy()

	cp := []corev1.ContainerPort{}
	for _, p := range md.Deployment.GetPorts() {
		cp = append(cp, corev1.ContainerPort{
			ContainerPort: int32(p),
		})
	}

	// ensuring deployment ID
	labels := updated.Spec.Template.GetLabels()
	labels["deployment"] = md.GetId()
	updated.SetLabels(labels)

	// updating spec
	if updated.Spec.Selector.MatchLabels == nil {
		updated.Spec.Selector.MatchLabels = map[string]string{
			"deployment": md.GetId(),
		}
	} else {
		updated.Spec.Selector.MatchLabels["deployment"] = md.GetId()
	}

	if updated.Spec.Template.Labels == nil {
		updated.Spec.Template.Labels = map[string]string{
			"deployment": md.GetId(),
		}
	} else {
		updated.Spec.Template.Labels["deployment"] = md.GetId()
	}

	updated.Spec.Replicas = toInt32(int(md.Deployment.Replicas))
	modelPodName := getPodName(md)
	for idx, c := range updated.Spec.Template.Spec.Containers {
		if c.Name == modelPodName {
			updated.Spec.Template.Spec.Containers[idx].Image = md.Deployment.GetImage()
			updated.Spec.Template.Spec.Containers[idx].Ports = cp
		}
	}

	return updated
}

func toInt32(v int) *int32 {
	i32 := int32(v)
	return &i32
}
