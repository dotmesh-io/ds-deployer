package status

import (
	"fmt"
	"sync"

	"github.com/dotmesh-io/ds-deployer/pkg/cond"
)

// static check
var _ Cache = &StatusCache{}

type StatusCache struct {
	deployments map[string]DeploymentStatus
	mu          *sync.RWMutex
	cond.Cond
}

type Cache interface {
	Set(deploymentID string, module Module, status Status)
	SetAvailableReplicas(deploymentID string, replicas int64)
	Get(deploymentID string) DeploymentStatus
	Delete(deploymentID string)
	List() map[string]DeploymentStatus
}

func New() *StatusCache {
	return &StatusCache{
		deployments: make(map[string]DeploymentStatus),
		mu:          &sync.RWMutex{},
	}
}

type Status int

const (
	StatusNone Status = iota
	StatusConfiguring
	StatusReady
	StatusError
)

func (s Status) String() string {
	switch s {
	case StatusNone:
		return "none"
	case StatusConfiguring:
		return "configuring"
	case StatusReady:
		return "ready"
	case StatusError:
		return "error"
	}
	return "unknown"
}

type DeploymentStatus struct {
	Deployment        Status
	AvailableReplicas int64
	Service           Status
	Ingress           Status
}

// Status returns overall status of the deployment
func (s *DeploymentStatus) Status() string {
	return fmt.Sprintf("Deployment: %s Service: %s Ingress: %s", s.Deployment, s.Service, s.Ingress)
}

type Module int

const (
	ModuleDeployment Module = iota
	ModuleService
	ModuleIngress
)

func (c *StatusCache) Set(deploymentID string, module Module, status Status) {
	defer c.Notify()
	c.mu.Lock()
	deploymentStatus, ok := c.deployments[deploymentID]
	if !ok {
		deploymentStatus = DeploymentStatus{}
	}

	switch module {
	case ModuleDeployment:
		deploymentStatus.Deployment = status
	case ModuleService:
		deploymentStatus.Service = status
	case ModuleIngress:
		deploymentStatus.Ingress = status
	}

	c.deployments[deploymentID] = deploymentStatus

	c.mu.Unlock()
}

func (c *StatusCache) SetAvailableReplicas(deploymentID string, replicas int64) {
	defer c.Notify()
	c.mu.Lock()
	deploymentStatus, ok := c.deployments[deploymentID]
	if !ok {
		deploymentStatus = DeploymentStatus{}
	}

	deploymentStatus.AvailableReplicas = replicas

	c.deployments[deploymentID] = deploymentStatus

	c.mu.Unlock()
}

func (c *StatusCache) Get(deploymentID string) DeploymentStatus {
	c.mu.RLock()
	defer c.mu.RUnlock()
	status, ok := c.deployments[deploymentID]
	if !ok {
		return DeploymentStatus{}
	}

	copy := new(DeploymentStatus)
	*copy = status

	return *copy
}

func (c *StatusCache) Delete(deploymentID string) {
	c.mu.Lock()
	delete(c.deployments, deploymentID)
	c.mu.Unlock()
}

func (c *StatusCache) List() map[string]DeploymentStatus {
	result := make(map[string]DeploymentStatus)
	c.mu.RLock()

	for k, v := range c.deployments {
		ds := new(DeploymentStatus)
		*ds = v

		result[k] = *ds
	}

	c.mu.RUnlock()
	return result
}
