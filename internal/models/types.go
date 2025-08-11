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
	Name      string    `json:"name"`
	Active    bool      `json:"active"`
	Traffic   int       `json:"traffic"`
	CreatedAt time.Time `json:"createdAt"`
	Status    string    `json:"status"`
}

type Container struct {
	Name    string
	Image   string
	Command []string
	Args    []string
	// add Ports, Env etc. as you like
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
