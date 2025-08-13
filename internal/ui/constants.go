package ui

import (
	"github.com/IAL32/az-tui/internal/ui/core"
)

// Table column keys for Apps mode
const (
	columnKeyAppName      = "name"
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

// Message type aliases for backward compatibility
// These now use the core message types
type loadedAppsMsg = core.LoadedAppsMsg
type loadedRevsMsg = core.LoadedRevisionsMsg
type loadedContainersMsg = core.LoadedContainersMsg
type revisionRestartedMsg = core.RevisionRestartedMsg
type loadedResourceGroupsMsg = core.LoadedResourceGroupsMsg
type leaveEnvVarsMsg = core.LeaveEnvVarsMsg
