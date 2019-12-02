package controller

import (
	"reflect"
	"testing"

	deployer_v1 "github.com/dotmesh-io/ds-deployer/apis/deployer/v1"
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func Test_toKubernetesDeployment(t *testing.T) {
	type args struct {
		modelDeployment      *deployer_v1.Deployment
		controllerIdentifier string
	}
	tests := []struct {
		name string
		args args
		want *appsv1.Deployment
	}{
		{
			name: "standard deployment",
			args: args{
				modelDeployment: &deployer_v1.Deployment{
					Id:        "1111-11-11-11-1111",
					Namespace: "ns",
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
				},
				controllerIdentifier: "5555",
			},
			want: &appsv1.Deployment{
				TypeMeta: metav1.TypeMeta{},
				ObjectMeta: metav1.ObjectMeta{
					Name:      "ds-cats-1111",
					Namespace: "ns",
					Labels: map[string]string{
						"owner": "ds-deployer",
					},
					Annotations: map[string]string{
						AnnControllerIdentifier: "5555",
						// based on model deployment name we will need this later
						// to ensure we delete what's not needed anymore
						"name":       "cats",
						"deployment": "1111-11-11-11-1111",
					},
				},
				Spec: appsv1.DeploymentSpec{
					Replicas: toInt32(1),
					Selector: &metav1.LabelSelector{
						MatchLabels: map[string]string{
							"deployment": "1111-11-11-11-1111",
						},
					},
					Template: corev1.PodTemplateSpec{
						ObjectMeta: metav1.ObjectMeta{
							Labels: map[string]string{
								"deployment": "1111-11-11-11-1111",
							},
							Annotations: map[string]string{},
						},
						Spec: corev1.PodSpec{
							ImagePullSecrets: []corev1.LocalObjectReference{
								// TODO: pass in secrets
							},
							Containers: []corev1.Container{
								corev1.Container{
									Name:  "ds-md-cats-1111",
									Image: "quay.io/image:tag",
									Ports: []corev1.ContainerPort{
										{
											ContainerPort: int32(8080),
										},
									},
								},
							},
						},
					},
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := toKubernetesDeployment(tt.args.modelDeployment, tt.args.controllerIdentifier); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("toKubernetesDeployment() = %v, want %v", got, tt.want)
			}
		})
	}
}
