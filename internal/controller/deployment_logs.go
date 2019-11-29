package controller

import (
	"errors"
	"io"
	"strings"

	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/labels"

	deployer_v1 "github.com/dotmesh-io/ds-deployer/apis/deployer/v1"
)

var (
	ErrDeploymentNotFound = errors.New("deployment not found")
)

func (c *Controller) Logs(request *deployer_v1.LogsRequest) (io.ReadCloser, error) {

	md, ok := c.cache.GetModelDeployment(request.DeploymentId)
	if !ok {
		return nil, ErrDeploymentNotFound
	}

	podLogOpts := &corev1.PodLogOptions{}

	var podName string

	set := labels.Set(map[string]string{
		"deployment": md.GetId(),
	})

	if pods, err := c.clientSet.CoreV1().Pods("").List(metav1.ListOptions{
		LabelSelector: set.AsSelector().String(),
	}); err != nil {
		c.logger.Errorf("list pods of deployment %s error:%v", md.Id, err)
		return nil, err
	} else {

		if len(pods.Items) == 1 {
			podName = pods.Items[0].GetName()

			if len(pods.Items[0].Spec.Containers) > 1 {

				switch request.Container {
				case deployer_v1.LogsRequest_MODEL:
					podLogOpts.Container = lookupModelContainerName(pods.Items[0].Spec.Containers)
				case deployer_v1.LogsRequest_PROXY:
					podLogOpts.Container = lookupModelProxyContainerName(pods.Items[0].Spec.Containers)
				}
			}
		}
	}

	req := c.clientSet.CoreV1().Pods(md.GetNamespace()).GetLogs(podName, podLogOpts)
	return req.Stream()

}

func lookupModelContainerName(containers []corev1.Container) string {
	for _, c := range containers {
		if strings.HasPrefix(c.Name, "ds-md") {
			return c.Name
		}
	}
	return ""
}

func lookupModelProxyContainerName(containers []corev1.Container) string {
	for _, c := range containers {
		if strings.HasPrefix(c.Name, "ds-mx") {
			return c.Name
		}
	}
	return ""
}
