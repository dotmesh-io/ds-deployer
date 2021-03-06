package controller

import (
	"context"
	"encoding/base64"
	"strconv"
	"strings"
	"sync"

	"github.com/dotmesh-io/ds-deployer/pkg/status"

	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/intstr"

	deployer_v1 "github.com/dotmesh-io/ds-deployer/apis/deployer/v1"
)

const (
	ModelProxyContainerPort int32 = 9501
	ModelProxyAPIPort       int32 = 9502
)

func (c *Controller) synchronizeDeployments() error {

	var wg sync.WaitGroup

	for _, modelDeployment := range c.cache.ModelDeployments() {
		existing, ok := c.cache.GetDeployment(modelDeployment.Namespace, getDeploymentName(modelDeployment))
		if !ok {
			// updating status cache
			if c.statusCache.Get(modelDeployment.Id).Deployment == status.StatusConfiguring {
				c.statusCache.Set(modelDeployment.Id, status.ModuleDeployment, status.StatusConfiguring)
			}

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
			c.logger.Debugf("deployment %s/%s needs updating", existing.Namespace, existing.Name)
			// updating status cache
			deploymentStatus := c.statusCache.Get(modelDeployment.Id)
			if deploymentStatus.Deployment != status.StatusConfiguring {
				c.statusCache.Set(modelDeployment.Id, status.ModuleDeployment, status.StatusConfiguring)
			}

			if deploymentStatus.AvailableReplicas != int64(existing.Status.AvailableReplicas) {
				c.statusCache.SetAvailableReplicas(modelDeployment.GetId(), int64(existing.Status.AvailableReplicas))
			}

			updatedDeployment := updateDeployment(existing, modelDeployment)

			wg.Add(1)
			go func(updatedDeployment *appsv1.Deployment) {
				err := c.client.Update(context.Background(), updatedDeployment)
				if err != nil {
					c.logger.Errorw("failed to update deployment via k8s client",
						"error", err,
						"deployment_id", modelDeployment.GetId(),
					)
				}
				wg.Done()
			}(updatedDeployment)
		} else {
			// deployment is in sync
			deploymentStatus := c.statusCache.Get(modelDeployment.Id)
			if deploymentStatus.Deployment != status.StatusReady {
				c.statusCache.Set(modelDeployment.Id, status.ModuleDeployment, status.StatusReady)
			}

			if deploymentStatus.AvailableReplicas != int64(existing.Status.AvailableReplicas) {
				c.statusCache.SetAvailableReplicas(modelDeployment.GetId(), int64(existing.Status.AvailableReplicas))
			}
		}
	}

	// going through existing deployments to see which ones should
	// be removed
	for _, deployment := range c.cache.deployments {

		if deployment.GetAnnotations() == nil {
			continue
		}

		_, ok := c.cache.GetModelDeployment(deployment.GetAnnotations()["deployment"])

		if !ok {
			c.statusCache.Delete(deployment.GetAnnotations()["deployment"])
			// not found in model deployments, should delete
			c.logger.Infof("deployment %s/%s (dep ID: %s) not found in model deployments, deleting", deployment.GetNamespace(), deployment.GetName(), deployment.GetAnnotations()["deployment"])
			err := c.client.Delete(context.Background(), deployment)
			if err != nil {
				if strings.Contains(err.Error(), "not found") {
					// it's fine
					continue
				}
				c.logger.Errorw("failed to delete deployment",
					"error", err,
					"deployment_id", deployment.GetAnnotations()["deployment"],
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

const (
	modelContainerPrefix      = "ds-md-"
	modelProxyContainerPrefix = "ds-mx-"
)

func getModelContainerName(d *deployer_v1.Deployment) string {
	// return modelContainerPrefix + d.GetName() + "-" + shortUUID(d.GetId())
	return "model"
}

func getModelProxyContainerName(d *deployer_v1.Deployment) string {
	return "proxy"
}

func shortUUID(u string) string {
	return strings.Split(u, "-")[0]
}

func (c *Controller) createDeployment(modelDeployment *deployer_v1.Deployment) error {
	return c.client.Create(context.Background(), toKubernetesDeployment(modelDeployment, c.controllerIdentifier))
}

// try decoding model classes
func getModelClasses(val string) string {
	decoded, err := base64.StdEncoding.DecodeString(val)
	if err != nil {
		return val
	}
	return string(decoded)
}

// Configures readiness probe to look at
// https://github.com/dotmesh-io/model-proxy/blob/master/pkg/api/api.go#L89
func getModelProxyReadinessProbe(modelDeployment *deployer_v1.Deployment) *corev1.Probe {
	return &corev1.Probe{
		Handler: corev1.Handler{
			HTTPGet: &corev1.HTTPGetAction{
				Path: "/health",
				Port: intstr.FromInt(int(ModelProxyAPIPort)),
			},
		},
		InitialDelaySeconds: 30,
		TimeoutSeconds:      10,
	}
}

func toKubernetesDeployment(modelDeployment *deployer_v1.Deployment, controllerIdentifier string) *appsv1.Deployment {

	cp := []corev1.ContainerPort{}

	for _, p := range modelDeployment.Deployment.GetPorts() {
		cp = append(cp, corev1.ContainerPort{
			ContainerPort: int32(p),
		})
	}

	containers := []corev1.Container{
		corev1.Container{
			Name:  getModelContainerName(modelDeployment),
			Image: modelDeployment.Deployment.GetImage(),
			Ports: cp,
		},
	}

	annotations := map[string]string{
		AnnControllerIdentifier: controllerIdentifier,
		// based on model deployment name we will need this later
		// to ensure we delete what's not needed anymore
		"name":       modelDeployment.GetName(),
		"deployment": modelDeployment.GetId(),
	}

	podAnnotations := map[string]string{}

	if modelDeployment.ModelProxyEnabled() && len(cp) > 0 {
		// configuration example can be found here:
		// https://github.com/dotmesh-io/k8s-manifests/blob/master/e2e-demo-prototype/model-dep.yaml
		containers = append(containers, corev1.Container{
			Name:  getModelProxyContainerName(modelDeployment),
			Image: modelDeployment.Metrics.GetImage(),
			Env: []corev1.EnvVar{
				{
					Name:  "TF_SERVING_ADDR",
					Value: "http://127.0.0.1:" + strconv.Itoa(int(cp[0].ContainerPort)),
				},
				{
					Name:  "TF_SERVING_PROXY_PORT",
					Value: strconv.Itoa(int(ModelProxyContainerPort)),
				},
				{
					Name:  "TF_CLASSES",
					Value: getModelClasses(modelDeployment.Metrics.Classes),
				},
				{
					Name:  "DEPLOYMENT_ID",
					Value: modelDeployment.GetId(),
				},
			},
			Ports: []corev1.ContainerPort{
				{
					ContainerPort: ModelProxyContainerPort, // Model reverse proxy, traffic will go Ingress -> Service -> 9501 -> 8501
				},
				{
					ContainerPort: ModelProxyAPIPort, // Proxy API
				},
			},
			LivenessProbe: getModelProxyReadinessProbe(modelDeployment),
		})

		// add prometheus scraping configuration
		// prometheus.io/scrape: "true"
		// prometheus.io/port: "9502"
		podAnnotations["prometheus.io/scrape"] = "true"
		podAnnotations["prometheus.io/path"] = "/api/metrics"
		podAnnotations["prometheus.io/port"] = strconv.Itoa(int(ModelProxyAPIPort))

	}

	deployment := &appsv1.Deployment{
		TypeMeta: metav1.TypeMeta{},
		ObjectMeta: metav1.ObjectMeta{
			Name:      getDeploymentName(modelDeployment),
			Namespace: modelDeployment.GetNamespace(),
			Labels: map[string]string{
				"owner": "ds-deployer",
			},
			Annotations: annotations,
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
					Annotations: podAnnotations,
				},
				Spec: corev1.PodSpec{
					ImagePullSecrets: []corev1.LocalObjectReference{
						// TODO: pass in secrets
					},
					Containers: containers,
				},
			},
		},
	}

	return deployment
}

// compares replicas, image, port, image pull secrets
func deploymentsEqual(desired, existing *appsv1.Deployment) bool {
	if *desired.Spec.Replicas != *existing.Spec.Replicas {
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

	for idx, container := range desired.Spec.Template.Spec.Containers {
		existingContainer := existing.Spec.Template.Spec.Containers[idx]

		if existingContainer.Name != container.Name {
			return false
		}

		if existingContainer.Image != container.Image {
			return false
		}
		if len(existingContainer.Ports) != len(container.Ports) {
			return false
		}
		for i := range container.Ports {
			if container.Ports[i].ContainerPort != existingContainer.Ports[i].ContainerPort {
				return false
			}
		}

		if !envEqual(existingContainer.Env, container.Env) {
			return false
		}

		if !livelinessProbeEqual(existingContainer.ReadinessProbe, container.ReadinessProbe) {
			return false
		}
	}

	return true
}

func livelinessProbeEqual(l, r *corev1.Probe) bool {
	if l == nil && r == nil {
		return true
	}
	// if one is nil and another is not, need to recreate
	if l == nil || r == nil {
		return false
	}

	if l.Handler.HTTPGet != nil && r.Handler.HTTPGet != nil {
		if l.Handler.HTTPGet.Path != r.Handler.HTTPGet.Path {
			return false
		}

		if l.Handler.HTTPGet.Port != r.Handler.HTTPGet.Port {
			return false
		}

		if l.InitialDelaySeconds != r.InitialDelaySeconds {
			return false
		}

		if l.TimeoutSeconds != r.TimeoutSeconds {
			return false
		}
	}

	return true
}

func envEqual(l, r []corev1.EnvVar) bool {
	if len(l) != len(r) {
		return false
	}

	for idx, val := range l {
		if r[idx].Name != val.Name {
			return false
		}
		if r[idx].Value != val.Value {
			return false
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
	modelPodName := getModelContainerName(md)
	proxyPodName := getModelProxyContainerName(md)
	for idx, c := range updated.Spec.Template.Spec.Containers {
		if c.Name == modelPodName {
			updated.Spec.Template.Spec.Containers[idx].Image = md.Deployment.GetImage()
			updated.Spec.Template.Spec.Containers[idx].Ports = cp
		}
		if c.Name == proxyPodName {
			updated.Spec.Template.Spec.Containers[idx].ReadinessProbe = getModelProxyReadinessProbe(md)
		}
	}

	return updated
}

func toInt32(v int) *int32 {
	i32 := int32(v)
	return &i32
}
