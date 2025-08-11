package azure

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/IAL32/az-tui/internal/models"
)

// TransformContainerAppsFromJSON transforms raw Azure JSON to ContainerApp models
func TransformContainerAppsFromJSON(rawJSON string) ([]models.ContainerApp, error) {
	// For now, we'll parse the already-transformed JSON directly
	// TODO: In a future implementation, we could apply JMESPath queries here
	var apps []models.ContainerApp
	if err := json.Unmarshal([]byte(rawJSON), &apps); err != nil {
		return nil, err
	}
	return apps, nil
}

// TransformRevisionsFromJSON transforms raw Azure JSON to Revision models
func TransformRevisionsFromJSON(rawJSON string) ([]models.Revision, error) {
	// TODO: In a future implementation, we could apply JMESPath queries here
	var revs []models.Revision
	if err := json.Unmarshal([]byte(rawJSON), &revs); err != nil {
		return nil, err
	}
	return revs, nil
}

// TransformContainersFromJSON transforms raw Azure revision JSON to Container models
func TransformContainersFromJSON(rawJSON string) ([]models.Container, error) {
	// Parse the full Azure revision response structure
	var resp struct {
		Properties struct {
			Template struct {
				Containers []struct {
					Name      string   `json:"name"`
					Image     string   `json:"image"`
					Command   []string `json:"command"`
					Args      []string `json:"args"`
					Resources struct {
						CPU    float64 `json:"cpu"`
						Memory string  `json:"memory"`
					} `json:"resources"`
					Env []struct {
						Name  string `json:"name"`
						Value string `json:"value"`
					} `json:"env"`
					Probes []struct {
						Type string `json:"type"`
					} `json:"probes"`
					VolumeMounts []struct {
						MountPath  string `json:"mountPath"`
						VolumeName string `json:"volumeName"`
					} `json:"volumeMounts"`
				} `json:"containers"`
			} `json:"template"`
		} `json:"properties"`
	}

	if err := json.Unmarshal([]byte(rawJSON), &resp); err != nil {
		return nil, err
	}

	containers := make([]models.Container, 0, len(resp.Properties.Template.Containers))
	for _, c := range resp.Properties.Template.Containers {
		// Convert env vars to map
		envMap := make(map[string]string)
		for _, env := range c.Env {
			envMap[env.Name] = env.Value
		}

		// Extract probe types
		var probes []string
		for _, probe := range c.Probes {
			probes = append(probes, probe.Type)
		}

		// Extract volume mounts
		var volumeMounts []string
		for _, vm := range c.VolumeMounts {
			volumeMounts = append(volumeMounts, fmt.Sprintf("%s:%s", vm.VolumeName, vm.MountPath))
		}

		container := models.Container{
			Name:         c.Name,
			Image:        c.Image,
			Command:      c.Command,
			Args:         c.Args,
			CPU:          c.Resources.CPU,
			Memory:       c.Resources.Memory,
			Env:          envMap,
			Probes:       probes,
			VolumeMounts: volumeMounts,
		}
		containers = append(containers, container)
	}
	return containers, nil
}

// TransformResourceGroupsFromJSON transforms raw Azure JSON to ResourceGroup models
func TransformResourceGroupsFromJSON(rawJSON string) ([]models.ResourceGroup, error) {
	// TODO: In a future implementation, we could apply JMESPath queries here
	var rgs []models.ResourceGroup
	if err := json.Unmarshal([]byte(rawJSON), &rgs); err != nil {
		return nil, err
	}
	return rgs, nil
}

// ParseTimeFromAzure parses Azure timestamp format
func ParseTimeFromAzure(timeStr string) (time.Time, error) {
	if timeStr == "" {
		return time.Time{}, nil
	}
	return time.Parse(time.RFC3339, timeStr)
}
