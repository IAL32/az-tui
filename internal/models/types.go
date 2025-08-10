package models

import "time"

// ----------------------------- Data -----------------------------

type ContainerApp struct {
	Name           string `json:"name"`
	ResourceGroup  string `json:"resourceGroup"`
	Location       string `json:"location"`
	EnvironmentID  string `json:"environmentId"`
	LatestRevision string `json:"latestRevisionName"`
	IngressFQDN    string `json:"ingressFqdn"`
}

type Revision struct {
	Name      string    `json:"name"`
	Active    bool      `json:"active"`
	Traffic   int       `json:"traffic"`
	CreatedAt time.Time `json:"createdAt"`
	Status    string    `json:"status"`
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
