package models

import "time"

// ----------------------------- Data -----------------------------

type ContainerApp struct {
	Name              string  `json:"name"`
	ResourceGroup     string  `json:"resourceGroup"`
	Location          string  `json:"location"`
	EnvironmentID     string  `json:"environmentId"`
	LatestRevision    string  `json:"latestRevisionName"`
	IngressFQDN       string  `json:"ingressFqdn"`
	ProvisioningState string  `json:"provisioningState"`
	RunningStatus     string  `json:"runningStatus"`
	MinReplicas       int     `json:"minReplicas"`
	MaxReplicas       int     `json:"maxReplicas"`
	CPU               float64 `json:"cpu"`
	Memory            string  `json:"memory"`
	IngressExternal   bool    `json:"ingressExternal"`
	TargetPort        int     `json:"targetPort"`
	IdentityType      string  `json:"identityType"`
	WorkloadProfile   string  `json:"workloadProfile"`
	CreatedAt         string  `json:"createdAt"`
	LastModifiedAt    string  `json:"lastModifiedAt"`
}

type Revision struct {
	Name              string    `json:"name"`
	Active            bool      `json:"active"`
	Traffic           int       `json:"traffic"`
	CreatedAt         time.Time `json:"createdAt"`
	Status            string    `json:"status"`
	FQDN              string    `json:"fqdn"`
	Replicas          int       `json:"replicas"`
	HealthState       string    `json:"healthState"`
	ProvisioningState string    `json:"provisioningState"`
	RunningState      string    `json:"runningState"`
	MinReplicas       int       `json:"minReplicas"`
	MaxReplicas       int       `json:"maxReplicas"`
	CPU               float64   `json:"cpu"`
	Memory            string    `json:"memory"`
}

type Container struct {
	Name         string            `json:"name"`
	Image        string            `json:"image"`
	Command      []string          `json:"command"`
	Args         []string          `json:"args"`
	CPU          float64           `json:"cpu"`
	Memory       string            `json:"memory"`
	Env          map[string]string `json:"env"`
	Ports        []ContainerPort   `json:"ports"`
	Probes       []string          `json:"probes"`
	VolumeMounts []string          `json:"volumeMounts"`
}

type ContainerPort struct {
	Name     string `json:"name"`
	Port     int    `json:"port"`
	Protocol string `json:"protocol"`
}

type RevItem struct{ Revision }

func (ri RevItem) Title() string { return ri.Name }
func (ri RevItem) Description() string {
	if ri.Active {
		return "active"
	}
	return ""
}
func (ri RevItem) FilterValue() string { return ri.Name }

type ResourceGroup struct {
	Name     string            `json:"name"`
	Location string            `json:"location"`
	State    string            `json:"provisioningState"`
	Tags     map[string]string `json:"tags"`
}
