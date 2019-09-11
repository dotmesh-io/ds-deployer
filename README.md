# Dotscience Deployer

Model deployment to Kubernetes.

## Design

Deployer subscribes to a deployment stream from Gateway and also does periodic queries to ensure that we have the latest desired configuration.

Once deployment reconciler detects that there are deployment requests that are not yet spawned in the cluster - it creates them. 

Deployer configures model-proxy as a sidecar, creates a service and an ingress. Webhook Relay ingress controller will then provision a tunnel. 


Goals:

- [x] declarative deployment management
- [ ] after reboot wait till new deployment info is retrieved
- [x] stateless - all state should be retrieved from the gateway on boot and kept in sync
- [x] should be able to work in groups - each deployer should only care about the resources that it created
- [x] ability to update image without recreating deployment 


Nice to have:
- [ ] Self update functionality, however, can be outsourced to Keel


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