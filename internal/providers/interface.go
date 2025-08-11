package providers

import (
	"context"

	"github.com/IAL32/az-tui/internal/models"
)

// DataProvider defines the interface for fetching Azure Container Apps data
type DataProvider interface {
	// ListResourceGroups returns all available resource groups
	ListResourceGroups(ctx context.Context) ([]models.ResourceGroup, error)

	// ListContainerApps returns container apps, optionally filtered by resource group
	ListContainerApps(ctx context.Context, resourceGroup string) ([]models.ContainerApp, error)

	// GetAppDetails returns detailed JSON information for a specific app
	GetAppDetails(ctx context.Context, name, resourceGroup string) (string, error)

	// ListRevisions returns all revisions for a specific container app
	ListRevisions(ctx context.Context, appName, resourceGroup string) ([]models.Revision, error)

	// ListContainers returns all containers for a specific app revision
	ListContainers(ctx context.Context, app models.ContainerApp, revisionName string) ([]models.Container, error)
}
