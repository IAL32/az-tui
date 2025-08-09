package models

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
	Name    string `json:"name"`
	Active  bool   `json:"active"`
	Traffic int    `json:"traffic"`
}
