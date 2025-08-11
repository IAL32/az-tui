package azure

import (
	"encoding/json"
	"os"
	"path/filepath"
	"testing"
)

// loadTestData loads test data from the mock testdata files
func loadTestData(filename string) (string, error) {
	path := filepath.Join("../mock/testdata", filename)
	data, err := os.ReadFile(path)
	if err != nil {
		return "", err
	}
	return string(data), nil
}

// TestTransformContainerAppsFromJSON tests the container apps transformation
func TestTransformContainerAppsFromJSON(t *testing.T) {
	t.Run("valid container apps from mock data", func(t *testing.T) {
		data, err := loadTestData("container_apps.json")
		if err != nil {
			t.Fatalf("Failed to load test data: %v", err)
		}

		result, err := TransformContainerAppsFromJSON(data)
		if err != nil {
			t.Fatalf("Unexpected error: %v", err)
		}

		if len(result) == 0 {
			t.Error("Expected at least one container app")
		}

		// Check first app has expected fields from mock data
		firstApp := result[0]
		if firstApp.Name == "" {
			t.Error("Expected app name to be set")
		}
		if firstApp.ResourceGroup == "" {
			t.Error("Expected resource group to be set")
		}
		if firstApp.Location == "" {
			t.Error("Expected location to be set")
		}
		if firstApp.CPU <= 0 {
			t.Error("Expected CPU to be greater than 0")
		}
		if firstApp.Memory == "" {
			t.Error("Expected memory to be set")
		}

		// Verify we have the expected number of apps from mock data
		expectedApps := []string{
			"web-frontend-prod",
			"api-backend-prod",
			"worker-service-prod",
			"web-frontend-staging",
			"api-backend-staging",
			"web-frontend-dev",
			"test-runner-dev",
			"monitoring-dashboard",
		}

		if len(result) != len(expectedApps) {
			t.Errorf("Expected %d apps, got %d", len(expectedApps), len(result))
		}

		// Check that all expected apps are present
		appNames := make(map[string]bool)
		for _, app := range result {
			appNames[app.Name] = true
		}

		for _, expectedName := range expectedApps {
			if !appNames[expectedName] {
				t.Errorf("Expected to find app '%s' in results", expectedName)
			}
		}
	})

	t.Run("invalid JSON", func(t *testing.T) {
		_, err := TransformContainerAppsFromJSON(`{"invalid": json}`)
		if err == nil {
			t.Error("Expected error for invalid JSON")
		}
	})

	t.Run("empty array", func(t *testing.T) {
		result, err := TransformContainerAppsFromJSON(`[]`)
		if err != nil {
			t.Errorf("Unexpected error: %v", err)
		}
		if len(result) != 0 {
			t.Errorf("Expected empty array, got %d items", len(result))
		}
	})
}

// TestTransformResourceGroupsFromJSON tests the resource groups transformation
func TestTransformResourceGroupsFromJSON(t *testing.T) {
	t.Run("valid resource groups from mock data", func(t *testing.T) {
		data, err := loadTestData("resource_groups.json")
		if err != nil {
			t.Fatalf("Failed to load test data: %v", err)
		}

		result, err := TransformResourceGroupsFromJSON(data)
		if err != nil {
			t.Fatalf("Unexpected error: %v", err)
		}

		if len(result) == 0 {
			t.Error("Expected at least one resource group")
		}

		// Check expected resource groups from mock data
		expectedRGs := []string{
			"rg-production-eastus",
			"rg-staging-westus",
			"rg-development-centralus",
			"rg-shared-services",
		}

		if len(result) != len(expectedRGs) {
			t.Errorf("Expected %d resource groups, got %d", len(expectedRGs), len(result))
		}

		// Verify all expected resource groups are present
		rgNames := make(map[string]bool)
		for _, rg := range result {
			rgNames[rg.Name] = true

			// Basic field validation
			if rg.Name == "" {
				t.Error("Expected resource group name to be set")
			}
			if rg.Location == "" {
				t.Error("Expected resource group location to be set")
			}
			if rg.Tags == nil {
				t.Error("Expected resource group tags to be initialized")
			}
		}

		for _, expectedName := range expectedRGs {
			if !rgNames[expectedName] {
				t.Errorf("Expected to find resource group '%s' in results", expectedName)
			}
		}
	})

	t.Run("invalid JSON", func(t *testing.T) {
		_, err := TransformResourceGroupsFromJSON(`{"invalid": json}`)
		if err == nil {
			t.Error("Expected error for invalid JSON")
		}
	})
}

// TestTransformContainersFromJSON tests the containers transformation
func TestTransformContainersFromJSON(t *testing.T) {
	t.Run("valid containers from mock revision details", func(t *testing.T) {
		data, err := loadTestData("revision_details.json")
		if err != nil {
			t.Fatalf("Failed to load test data: %v", err)
		}

		// Parse the revision details JSON to extract a specific revision
		var revisionDetails map[string]interface{}
		if err := json.Unmarshal([]byte(data), &revisionDetails); err != nil {
			t.Fatalf("Failed to parse revision details: %v", err)
		}

		// Test with the first revision in the map
		var firstRevisionJSON string
		for _, revisionData := range revisionDetails {
			revisionBytes, err := json.Marshal(revisionData)
			if err != nil {
				t.Fatalf("Failed to marshal revision data: %v", err)
			}
			firstRevisionJSON = string(revisionBytes)
			break
		}

		result, err := TransformContainersFromJSON(firstRevisionJSON)
		if err != nil {
			t.Fatalf("Unexpected error: %v", err)
		}

		if len(result) == 0 {
			t.Error("Expected at least one container")
		}

		// Basic validation on the first container
		container := result[0]
		if container.Name == "" {
			t.Error("Expected container name to be set")
		}
		if container.Image == "" {
			t.Error("Expected container image to be set")
		}
		if container.CPU <= 0 {
			t.Error("Expected CPU to be greater than 0")
		}
		if container.Memory == "" {
			t.Error("Expected memory to be set")
		}
	})

	t.Run("multi-container scenario from mock data", func(t *testing.T) {
		data, err := loadTestData("revision_details.json")
		if err != nil {
			t.Fatalf("Failed to load test data: %v", err)
		}

		// Parse the revision details JSON
		var revisionDetails map[string]interface{}
		if err := json.Unmarshal([]byte(data), &revisionDetails); err != nil {
			t.Fatalf("Failed to parse revision details: %v", err)
		}

		// Look for a revision with multiple containers (api-backend-prod has 2)
		var multiContainerRevision string
		for key, revisionData := range revisionDetails {
			if key == "api-backend-prod-api-backend-prod--v1-8" {
				revisionBytes, err := json.Marshal(revisionData)
				if err != nil {
					t.Fatalf("Failed to marshal revision data: %v", err)
				}
				multiContainerRevision = string(revisionBytes)
				break
			}
		}

		if multiContainerRevision == "" {
			t.Skip("Multi-container revision not found in mock data")
		}

		result, err := TransformContainersFromJSON(multiContainerRevision)
		if err != nil {
			t.Fatalf("Unexpected error: %v", err)
		}

		if len(result) < 2 {
			t.Errorf("Expected at least 2 containers, got %d", len(result))
		}

		// Verify we have different container names
		names := make(map[string]bool)
		for _, container := range result {
			if names[container.Name] {
				t.Errorf("Duplicate container name: %s", container.Name)
			}
			names[container.Name] = true
		}
	})

	t.Run("invalid JSON", func(t *testing.T) {
		_, err := TransformContainersFromJSON(`{"invalid": json}`)
		if err == nil {
			t.Error("Expected error for invalid JSON")
		}
	})

	t.Run("empty object", func(t *testing.T) {
		result, err := TransformContainersFromJSON(`{}`)
		if err != nil {
			t.Errorf("Unexpected error: %v", err)
		}
		if len(result) != 0 {
			t.Errorf("Expected empty result, got %d containers", len(result))
		}
	})
}

// TestTransformRevisionsFromJSON tests the revisions transformation
func TestTransformRevisionsFromJSON(t *testing.T) {
	t.Run("valid revisions from mock data", func(t *testing.T) {
		data, err := loadTestData("revisions.json")
		if err != nil {
			t.Fatalf("Failed to load test data: %v", err)
		}

		// Parse the revisions JSON to extract revisions for a specific app
		var revisionsMap map[string]json.RawMessage
		if err := json.Unmarshal([]byte(data), &revisionsMap); err != nil {
			t.Fatalf("Failed to parse revisions data: %v", err)
		}

		// Test with the first app's revisions
		var firstAppRevisions string
		for _, revisionData := range revisionsMap {
			firstAppRevisions = string(revisionData)
			break
		}

		if firstAppRevisions == "" {
			t.Fatal("No revision data found in mock file")
		}

		result, err := TransformRevisionsFromJSON(firstAppRevisions)
		if err != nil {
			t.Fatalf("Unexpected error: %v", err)
		}

		if len(result) == 0 {
			t.Error("Expected at least one revision")
		}

		// Basic validation on the first revision
		rev := result[0]
		if rev.Name == "" {
			t.Error("Expected revision name to be set")
		}
		if rev.Status == "" {
			t.Error("Expected revision status to be set")
		}
	})

	t.Run("invalid JSON", func(t *testing.T) {
		_, err := TransformRevisionsFromJSON(`{"invalid": json}`)
		if err == nil {
			t.Error("Expected error for invalid JSON")
		}
	})
}

// TestParseTimeFromAzure tests the Azure time parsing function
func TestParseTimeFromAzure(t *testing.T) {
	tests := []struct {
		name        string
		input       string
		expectError bool
	}{
		{
			name:        "valid RFC3339 time",
			input:       "2024-01-20T10:00:00Z",
			expectError: false,
		},
		{
			name:        "valid RFC3339 time with timezone",
			input:       "2024-01-20T10:00:00+05:00",
			expectError: false,
		},
		{
			name:        "empty string",
			input:       "",
			expectError: false,
		},
		{
			name:        "invalid time format",
			input:       "2024-01-20 10:00:00",
			expectError: true,
		},
		{
			name:        "invalid time string",
			input:       "not-a-time",
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := ParseTimeFromAzure(tt.input)

			if tt.expectError {
				if err == nil {
					t.Errorf("Expected error but got none")
				}
				return
			}

			if err != nil {
				t.Errorf("Unexpected error: %v", err)
				return
			}

			// For empty string, expect zero time
			if tt.input == "" {
				if !result.IsZero() {
					t.Errorf("Expected zero time for empty string, got %v", result)
				}
				return
			}

			// For valid times, just check it's not zero
			if result.IsZero() {
				t.Errorf("Expected non-zero time, got %v", result)
			}
		})
	}
}

// TestTransformContainersFromJSON_VolumeMount tests volume mount formatting using mock data
func TestTransformContainersFromJSON_VolumeMount(t *testing.T) {
	data, err := loadTestData("revision_details.json")
	if err != nil {
		t.Fatalf("Failed to load test data: %v", err)
	}

	// Parse the revision details JSON
	var revisionDetails map[string]interface{}
	if err := json.Unmarshal([]byte(data), &revisionDetails); err != nil {
		t.Fatalf("Failed to parse revision details: %v", err)
	}

	// Look for a revision with volume mounts (web-frontend-prod has volume mounts)
	var revisionWithMounts string
	for key, revisionData := range revisionDetails {
		if key == "web-frontend-prod-web-frontend-prod--v2-3" {
			revisionBytes, err := json.Marshal(revisionData)
			if err != nil {
				t.Fatalf("Failed to marshal revision data: %v", err)
			}
			revisionWithMounts = string(revisionBytes)
			break
		}
	}

	if revisionWithMounts == "" {
		t.Skip("Revision with volume mounts not found in mock data")
	}

	result, err := TransformContainersFromJSON(revisionWithMounts)
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}

	if len(result) == 0 {
		t.Fatal("Expected at least one container")
	}

	// Check that volume mounts are in the expected format "volumeName:mountPath"
	container := result[0]
	if len(container.VolumeMounts) == 0 {
		t.Error("Expected container to have volume mounts from mock data")
	}

	for _, mount := range container.VolumeMounts {
		if mount == "" {
			t.Error("Volume mount should not be empty")
		}
		// Should contain a colon separator
		if len(mount) < 3 {
			t.Errorf("Volume mount should be in format 'volume:path', got %s", mount)
		}
		// Check for colon separator
		hasColon := false
		for _, char := range mount {
			if char == ':' {
				hasColon = true
				break
			}
		}
		if !hasColon {
			t.Errorf("Volume mount should contain ':', got %s", mount)
		}
	}
}

// TestTransformContainersFromJSON_EnvVars tests environment variable handling using mock data
func TestTransformContainersFromJSON_EnvVars(t *testing.T) {
	data, err := loadTestData("revision_details.json")
	if err != nil {
		t.Fatalf("Failed to load test data: %v", err)
	}

	// Parse the revision details JSON
	var revisionDetails map[string]interface{}
	if err := json.Unmarshal([]byte(data), &revisionDetails); err != nil {
		t.Fatalf("Failed to parse revision details: %v", err)
	}

	// Look for a revision with environment variables (web-frontend-prod has many env vars)
	var revisionWithEnvVars string
	for key, revisionData := range revisionDetails {
		if key == "web-frontend-prod-web-frontend-prod--v2-3" {
			revisionBytes, err := json.Marshal(revisionData)
			if err != nil {
				t.Fatalf("Failed to marshal revision data: %v", err)
			}
			revisionWithEnvVars = string(revisionBytes)
			break
		}
	}

	if revisionWithEnvVars == "" {
		t.Skip("Revision with environment variables not found in mock data")
	}

	result, err := TransformContainersFromJSON(revisionWithEnvVars)
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}

	if len(result) == 0 {
		t.Fatal("Expected at least one container")
	}

	container := result[0]
	if len(container.Env) == 0 {
		t.Error("Expected container to have environment variables from mock data")
	}

	// Check that environment variables are properly parsed
	for name, value := range container.Env {
		if name == "" {
			t.Error("Environment variable name should not be empty")
		}
		// Value can be empty (like secrets marked as "***")
		t.Logf("Found env var: %s=%s", name, value)
	}

	// Check for specific env vars that should exist in the mock data
	if _, exists := container.Env["NODE_ENV"]; !exists {
		t.Error("Expected NODE_ENV environment variable from mock data")
	}
}
