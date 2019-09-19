package status

import (
	"sync"

	"github.com/dotmesh-io/ds-deployer/pkg/cond"
)

type StatusCache struct {
	deployments map[string]DeploymentStatus
	mu          *sync.RWMutex
	cond.Cond
}

type Cache interface {
	Set(deploymentID string, module Module, status Status)
	Get(deploymentID string) DeploymentStatus
	Delete(deploymentID string)
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

type DeploymentStatus struct {
	Deployment Status
	Service    Status
	Ingress    Status
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
