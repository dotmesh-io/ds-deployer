# Dotscience Deployer

Model deployment to Kubernetes.

## Design

Deployer subscribes to a deployment stream from Gateway and also does periodic queries to ensure that we have the latest desired configuration.

Once deployment reconciler detects that there are deployment requests that are not yet spawned in the cluster - it creates them. 

Deployer configures model-proxy as a sidecar, creates a service and an ingress. Webhook Relay ingress controller will then provision a tunnel. 