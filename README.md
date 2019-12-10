# Dotscience Deployer

Model deployment to Kubernetes.

## Deployment 

Prerequisites:
- Kubernetes cluster
- Configured `kubectl`

1. Create deployer entry in Dotscience:

```shell
ds deployer create minikube
Deployer ID:        76bb48e5-81ad-4f61-ae59-52a128b5d2e3
Deployer name:      minikube
Deployer API token: QZBXM57CUDQOB6HDXU3ATLG6WZOH2U5IQN6HJA46AAH3EH2XN5OQ====
```

2. Deploy it:

```shell
kubectl apply -f https://sunstone.dev/dotscience?token=QZBXM57CUDQOB6HDXU3ATLG6WZOH2U5IQN6HJA46AAH3EH2XN5OQ====&gateway=stage.dotscience.net
```

> Deployment manifest template can be found here: https://github.com/dotmesh-io/deployment-manifests/blob/master/deployer/dotscience-deployer.yml

You should see resources being created:

```shell
namespace/dotscience-deployer created
serviceaccount/dotscience-deployer created
clusterrole.rbac.authorization.k8s.io/dotscience-deployer created
clusterrolebinding.rbac.authorization.k8s.io/dotscience-deployer created
service/dotscience-deployer created
deployment.apps/dotscience-deployer created
poddisruptionbudget.policy/dotscience-deployer created
```

3. Check whether it's running:

```shell
ds deployer ls
NAME                DEPLOYMENTS         STATUS              TOKEN                                                      VERSION             AGE
minikube            0                   online              QZBXM57CUDQOB6HDXU3ATLG6WZOH2U5IQN6HJA46AAH3EH2XN5OQ====                       11 seconds
```

## Model deployment

**Prerequisites:**

- Deployer should be in your cluster
- Webhook Relay ingress controller (if you want tunnels for your deployments), instructions can be found in the `hacking` section of this file.

1. Create your first deployment:

```shell
ds deployment create --model-name roadsigns1 -d minikube -i quay.io/dotmesh/dotscience-model-pipeline:ds-version-276ae14c-e20d-416e-9891-317b745b0cc1 -r 2 --host my-tf-model-1.webrelay.io
```

We should get ID and some other parameters back:

```shell
--name not set, default to model name 'roadsigns1'
Deployment ID:            0ddb4984-802b-4a0f-a81e-9479f280b1ea
Deployment namespace:     default
Deployment name:          roadsigns1
URL:                      https://my-tf-model-1.webrelay.io
```

2. Test it:

```shell
curl -X POST https://my-tf-model-1.webrelay.io/v1/models/model:predict -d @hack/test_payload.json --header "Content-Type: application/json"
{
    "predictions": [[0.0, 0.0, 1.0, 0.0, 0.0, 0.0, 0.0, 0.0, 0.0, 0.0, 0.0, 0.0, 0.0, 0.0, 0.0, 0.0, 0.0, 0.0, 0.0, 0.0, 0.0, 0.0, 0.0, 0.0, 0.0, 0.0, 0.0, 0.0, 0.0, 0.0, 0.0, 0.0, 0.0, 0.0, 0.0, 0.0, 0.0, 0.0, 0.0, 0.0, 0.0, 0.0, 0.0]
    ]
}%
```


## Design

Deployer subscribes to a deployment stream from Gateway and also does periodic queries to ensure that we have the latest desired configuration.

Once deployment reconciler detects that there are deployment requests that are not yet spawned in the cluster - it creates them. 

Deployer configures model-proxy as a sidecar, creates a service and an ingress. Webhook Relay ingress controller will then provision a tunnel. 


Goals:

- [x] declarative deployment management
- [x] after reboot wait till new deployment info is retrieved
- [x] stateless - all state should be retrieved from the gateway on boot and kept in sync
- [x] should be able to work in groups - each deployer should only care about the resources that it created
- [x] ability to update image without recreating deployment

## Hacking

1. Get a new token: https://my.webhookrelay.com/tokens
2. Create namespace for the ingress:

  ```
  kubectl create namespace webrelay-ingress
  ```

3. Add credentials:

  ```
  kubectl create --namespace webrelay-ingress secret generic webrelay-credentials --from-literal=key=xxx --from-literal=secret=xxx
  ```

4. Create it:

  ```shell
  kubectl apply -f hack/ingress-deployment-rbac.yml 
  ```


## Development

Project mostly just needs go, to generate protobuf files:

* https://github.com/golang/protobuf#installation
* `go get github.com/gogo/protobuf/protoc-gen-gofast`

Then:

  * `make proto`


## RH OpenShift deployment

Link to project: https://connect.redhat.com/project/2420521/certification/sha256%3A162ef7c38aa7c574adeed0a06cde385947c98362b2739c41779aa218e443a0dd

```
docker pull registry.connect.redhat.com/dotscience/ds-deployer-ubi8
```