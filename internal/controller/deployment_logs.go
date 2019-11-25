package controller

import (
	"errors"
	"io"

	deployer_v1 "github.com/dotmesh-io/ds-deployer/apis/deployer/v1"

	corev1 "k8s.io/api/core/v1"
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

	switch request.Container {
	case deployer_v1.LogsRequest_MODEL:
		podName = getPodName(md)
	case deployer_v1.LogsRequest_PROXY:
		podName = getModelProxyPodName(md)
	}

	req := c.clientSet.CoreV1().Pods(md.GetNamespace()).GetLogs(podName, podLogOpts)
	return req.Stream()

}
