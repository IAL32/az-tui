package mock

import (
	"context"
	"embed"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/IAL32/az-tui/internal/azure"
	"github.com/IAL32/az-tui/internal/models"
)

//go:embed testdata/*.json
var testDataFS embed.FS

// Provider implements the DataProvider interface using mock data
type Provider struct{}

// NewProvider creates a new mock data provider
func NewProvider() (*Provider, error) {
	return &Provider{}, nil
}

// ListResourceGroups returns all available resource groups
func (p *Provider) ListResourceGroups(ctx context.Context) ([]models.ResourceGroup, error) {
	// Simulate some processing time
	select {
	case <-ctx.Done():
		return nil, ctx.Err()
	default:
	}

	// Load raw JSON data and transform it using shared helpers
	rgData, err := testDataFS.ReadFile("testdata/resource_groups.json")
	if err != nil {
		return nil, fmt.Errorf("failed to read resource groups: %w", err)
	}

	return azure.TransformResourceGroupsFromJSON(string(rgData))
}

// ListContainerApps returns container apps, optionally filtered by resource group
func (p *Provider) ListContainerApps(ctx context.Context, resourceGroup string) ([]models.ContainerApp, error) {
	// Simulate some processing time
	select {
	case <-ctx.Done():
		return nil, ctx.Err()
	default:
	}

	// Load raw JSON data and transform it using shared helpers
	appsData, err := testDataFS.ReadFile("testdata/container_apps.json")
	if err != nil {
		return nil, fmt.Errorf("failed to read container apps: %w", err)
	}

	apps, err := azure.TransformContainerAppsFromJSON(string(appsData))
	if err != nil {
		return nil, fmt.Errorf("failed to transform container apps: %w", err)
	}

	if resourceGroup == "" {
		return apps, nil
	}

	// Filter by resource group
	var filtered []models.ContainerApp
	for _, app := range apps {
		if app.ResourceGroup == resourceGroup {
			filtered = append(filtered, app)
		}
	}

	return filtered, nil
}

// GetAppDetails returns detailed JSON information for a specific app
func (p *Provider) GetAppDetails(ctx context.Context, name, resourceGroup string) (string, error) {
	// Simulate some processing time
	select {
	case <-ctx.Done():
		return "", ctx.Err()
	default:
	}

	// Load and find the app
	apps, err := p.ListContainerApps(ctx, resourceGroup)
	if err != nil {
		return "", err
	}

	for _, app := range apps {
		if app.Name == name && app.ResourceGroup == resourceGroup {
			// Return a formatted JSON representation
			data, err := json.MarshalIndent(app, "", "  ")
			if err != nil {
				return "", fmt.Errorf("failed to marshal app details: %w", err)
			}
			return string(data), nil
		}
	}

	return "", fmt.Errorf("app %s not found in resource group %s", name, resourceGroup)
}

// ListRevisions returns all revisions for a specific container app
func (p *Provider) ListRevisions(ctx context.Context, appName, resourceGroup string) ([]models.Revision, error) {
	// Simulate some processing time
	select {
	case <-ctx.Done():
		return nil, ctx.Err()
	default:
	}

	// Load raw JSON data and transform it using shared helpers
	revsData, err := testDataFS.ReadFile("testdata/revisions.json")
	if err != nil {
		return nil, fmt.Errorf("failed to read revisions: %w", err)
	}

	// Parse the JSON structure to get revisions for the specific app
	var allRevisions map[string]json.RawMessage
	if err := json.Unmarshal(revsData, &allRevisions); err != nil {
		return nil, fmt.Errorf("failed to unmarshal revisions: %w", err)
	}

	appRevisions, exists := allRevisions[appName]
	if !exists {
		return []models.Revision{}, nil
	}

	return azure.TransformRevisionsFromJSON(string(appRevisions))
}

// ListContainers returns all containers for a specific app revision
func (p *Provider) ListContainers(ctx context.Context, app models.ContainerApp, revisionName string) ([]models.Container, error) {
	// Simulate some processing time
	select {
	case <-ctx.Done():
		return nil, ctx.Err()
	default:
	}

	// Load raw JSON data for the specific revision
	containersData, err := testDataFS.ReadFile("testdata/revision_details.json")
	if err != nil {
		return nil, fmt.Errorf("failed to read revision details: %w", err)
	}

	// Parse the JSON structure to get containers for the specific revision
	var allRevisionDetails map[string]json.RawMessage
	if err := json.Unmarshal(containersData, &allRevisionDetails); err != nil {
		return nil, fmt.Errorf("failed to unmarshal revision details: %w", err)
	}

	key := fmt.Sprintf("%s-%s", app.Name, revisionName)
	revisionDetail, exists := allRevisionDetails[key]
	if !exists {
		return []models.Container{}, nil
	}

	return azure.TransformContainersFromJSON(string(revisionDetail))
}

// GetContainerKey returns the key used to store containers in the mock data
func GetContainerKey(appName, revisionName string) string {
	return fmt.Sprintf("%s-%s", appName, revisionName)
}

// FilterAppsByResourceGroup is a helper function to filter apps by resource group
func FilterAppsByResourceGroup(apps []models.ContainerApp, resourceGroup string) []models.ContainerApp {
	if resourceGroup == "" {
		return apps
	}

	var filtered []models.ContainerApp
	for _, app := range apps {
		if strings.EqualFold(app.ResourceGroup, resourceGroup) {
			filtered = append(filtered, app)
		}
	}
	return filtered
}
