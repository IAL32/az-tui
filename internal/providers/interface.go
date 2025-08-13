package providers

import (
	"context"

	"github.com/IAL32/az-tui/internal/models"
)

// DataProvider defines the interface for fetching Azure Container Apps data
type DataProvider interface {
	ListResourceGroups(ctx context.Context) ([]models.ResourceGroup, error)
	ListContainerApps(ctx context.Context, resourceGroup string) ([]models.ContainerApp, error)
	GetAppDetails(ctx context.Context, name, resourceGroup string) (string, error)
	ListRevisions(ctx context.Context, appName, resourceGroup string) ([]models.Revision, error)
	ListContainers(ctx context.Context, app models.ContainerApp, revisionName string) ([]models.Container, error)
}
