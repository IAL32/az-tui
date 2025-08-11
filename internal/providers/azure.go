package providers

import (
	"context"

	"github.com/IAL32/az-tui/internal/azure"
	"github.com/IAL32/az-tui/internal/models"
)

// AzureProvider implements the DataProvider interface using Azure CLI
type AzureProvider struct{}

// NewAzureProvider creates a new Azure CLI data provider
func NewAzureProvider() *AzureProvider {
	return &AzureProvider{}
}

// ListResourceGroups returns all available resource groups
func (p *AzureProvider) ListResourceGroups(ctx context.Context) ([]models.ResourceGroup, error) {
	return azure.ListResourceGroups(ctx)
}

// ListContainerApps returns container apps, optionally filtered by resource group
func (p *AzureProvider) ListContainerApps(ctx context.Context, resourceGroup string) ([]models.ContainerApp, error) {
	return azure.ListContainerApps(ctx, resourceGroup)
}

// GetAppDetails returns detailed JSON information for a specific app
func (p *AzureProvider) GetAppDetails(ctx context.Context, name, resourceGroup string) (string, error) {
	return azure.GetAppDetails(ctx, name, resourceGroup)
}

// ListRevisions returns all revisions for a specific container app
func (p *AzureProvider) ListRevisions(ctx context.Context, appName, resourceGroup string) ([]models.Revision, error) {
	return azure.ListRevisions(ctx, appName, resourceGroup)
}

// ListContainers returns all containers for a specific app revision
func (p *AzureProvider) ListContainers(ctx context.Context, app models.ContainerApp, revisionName string) ([]models.Container, error) {
	return azure.ListContainersCmd(ctx, app, revisionName)
}
