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

func toKubernetesDeployments(deployments map[Meta]*deployer_v1.Deployment) map[Meta]*appsv1.Deployment {

	computed := make(map[Meta]*appsv1.Deployment)

	// for _, d := range deployments {

	// }

	return computed
}

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
			go func(modelDeployment *deployer_v1.Deployment) {
				err := c.createDeployment(modelDeployment)
				if err != nil {
					c.logger.Errorw("failed to create deployment",
						"error", err,
						"deployment_id", modelDeployment.GetId(),
					)
				}
			}(modelDeployment)
			// should get created if it doesn't exist
			continue
		}

		c.logger.Infof("deployment %s/%s found, checking for updates", existing.Namespace, existing.Name)

		existing.GetLabels()
	}

	wg.Wait()

	return nil
}

func getDeploymentName(d *deployer_v1.Deployment) string {
	return "ds-" + d.GetName() + shortUUID(d.GetId())
}

func getPodName(d *deployer_v1.Deployment) string {
	return "ds-" + d.GetName() + shortUUID(d.GetId())
}

func shortUUID(u string) string {
	return strings.Split(u, "-")[0]
}

func (c *Controller) createDeployment(modelDeployment *deployer_v1.Deployment) error {
	return c.client.Create(context.Background(), toKubernetesDeployment(modelDeployment))

}

func toKubernetesDeployment(modelDeployment *deployer_v1.Deployment) *appsv1.Deployment {

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
			Annotations: map[string]string{},
		},
		Spec: appsv1.DeploymentSpec{
			Replicas: toInt32(1),
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

func toInt32(v int) *int32 {
	i32 := int32(v)
	return &i32
}
