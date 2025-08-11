package ui

import (
	models "github.com/IAL32/az-tui/internal/models"
)

// Table column keys for Apps mode
const (
	columnKeyAppName      = "name"
	columnKeyAppRG        = "rg"
	columnKeyAppLocation  = "location"
	columnKeyAppRevision  = "revision"
	columnKeyAppFQDN      = "fqdn"
	columnKeyAppStatus    = "status"
	columnKeyAppReplicas  = "replicas"
	columnKeyAppResources = "resources"
	columnKeyAppIngress   = "ingress"
	columnKeyAppIdentity  = "identity"
	columnKeyAppWorkload  = "workload"
)

// Table column keys for Revisions mode
const (
	columnKeyRevName      = "name"
	columnKeyRevActive    = "active"
	columnKeyRevTraffic   = "traffic"
	columnKeyRevCreated   = "created"
	columnKeyRevStatus    = "status"
	columnKeyRevReplicas  = "replicas"
	columnKeyRevScaling   = "scaling"
	columnKeyRevResources = "resources"
	columnKeyRevHealth    = "health"
	columnKeyRevRunning   = "running"
	columnKeyRevFQDN      = "fqdn"
)

// Table column keys for Containers mode
const (
	columnKeyCtrName      = "name"
	columnKeyCtrImage     = "image"
	columnKeyCtrCommand   = "command"
	columnKeyCtrArgs      = "args"
	columnKeyCtrResources = "resources"
	columnKeyCtrEnvCount  = "envcount"
	columnKeyCtrProbes    = "probes"
	columnKeyCtrVolumes   = "volumes"
	columnKeyCtrStatus    = "status"
)

// Table column keys for Resource Groups mode
const (
	columnKeyRGName     = "name"
	columnKeyRGLocation = "location"
	columnKeyRGState    = "state"
	columnKeyRGTags     = "tags"
)

// Message types for different operations
type loadedAppsMsg struct {
	apps []models.ContainerApp
	err  error
}

type loadedRevsMsg struct {
	revs []models.Revision
	err  error
}

type loadedContainersMsg struct {
	appID   string
	revName string
	ctrs    []models.Container
	err     error
}

type revisionRestartedMsg struct {
	appID   string
	revName string
	err     error
	out     string
}

type loadedResourceGroupsMsg struct {
	resourceGroups []models.ResourceGroup
	err            error
}
