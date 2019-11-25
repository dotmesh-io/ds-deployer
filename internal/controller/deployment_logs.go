package controller

import (
	"io"

	deployer_v1 "github.com/dotmesh-io/ds-deployer/apis/deployer/v1"

	corev1 "k8s.io/api/core/v1"
)

func (c *Controller) Logs(md *deployer_v1.Deployment, request *deployer_v1.LogsRequest) (io.ReadCloser, error) {
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
