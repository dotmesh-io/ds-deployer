package controller

import (
	"testing"

	"github.com/stretchr/testify/assert"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	deployer_v1 "github.com/dotmesh-io/ds-deployer/apis/deployer/v1"
	"go.uber.org/zap"
)

// test that we are receiving copies
func TestGetCachedDeployments(t *testing.T) {
	c := NewKubernetesCache("123", zap.S())

	dep := &deployer_v1.Deployment{
		Id:        "1111-11-11-11-1111",
		Namespace: "original",
		Name:      "cats",
		Deployment: &deployer_v1.DeploymentSpec{
			Replicas: 1,
			Image:    "quay.io/image:tag",
			Ports:    []int32{8080},
		},
		Service: &deployer_v1.ServiceSpec{
			Ports: []*deployer_v1.ServicePort{
				{
					Name:       "foo",
					Port:       8080,
					TargetPort: 8080,
				},
			},
		},
		Ingress: &deployer_v1.IngressSpec{
			Class: "nginx",
			Host:  "foo.bar.com",
		},
		Labels: make(map[string]string),
	}

	inserted := c.Insert(dep)
	assert.True(t, inserted)

	cached := c.ModelDeployments()
	assert.Equal(t, len(cached), 1)

	dep.Namespace = "updated"
	c.Insert(dep)

	assert.Equal(t, "original", cached[0].Namespace)
}

//  test to ensure that we get copies of objcets
func TestGetCachedService(t *testing.T) {
	c := NewKubernetesCache("123", zap.S())

	inserted := c.Insert(&corev1.Service{
		ObjectMeta: metav1.ObjectMeta{
			Namespace: "foo",
			Name:      "bar",
			Annotations: map[string]string{
				AnnControllerIdentifier: "123",
			},
		},
	})
	assert.True(t, inserted)

	c.Insert(inserted)

	svc, ok := c.GetService("foo", "bar")
	assert.Equal(t, ok, true)

	svc.Namespace = "new-namespace"

	assert.Equal(t, "foo", c.services[Meta{namespace: "foo", name: "bar"}].Namespace)
}
