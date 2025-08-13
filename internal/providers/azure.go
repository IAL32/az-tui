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

func (p *AzureProvider) ListResourceGroups(ctx context.Context) ([]models.ResourceGroup, error) {
	return azure.ListResourceGroups(ctx)
}

func (p *AzureProvider) ListContainerApps(ctx context.Context, resourceGroup string) ([]models.ContainerApp, error) {
	return azure.ListContainerApps(ctx, resourceGroup)
}

func (p *AzureProvider) GetAppDetails(ctx context.Context, name, resourceGroup string) (string, error) {
	return azure.GetAppDetails(ctx, name, resourceGroup)
}

func (p *AzureProvider) ListRevisions(ctx context.Context, appName, resourceGroup string) ([]models.Revision, error) {
	return azure.ListRevisions(ctx, appName, resourceGroup)
}

func (p *AzureProvider) ListContainers(ctx context.Context, app models.ContainerApp, revisionName string) ([]models.Container, error) {
	return azure.ListContainersCmd(ctx, app, revisionName)
}
