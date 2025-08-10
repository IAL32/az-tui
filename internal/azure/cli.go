package azure

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"os/exec"
	"strings"

	m "az-tui/internal/models"
)

func RunAz(ctx context.Context, args ...string) (string, error) {
	cmd := exec.CommandContext(ctx, "az", args...)
	var out, errb bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &errb
	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("az %s: %w%s", strings.Join(args, " "), err, errb.String())
	}
	return out.String(), nil
}

func ListContainerApps(ctx context.Context, rg string) ([]m.ContainerApp, error) {
	q := "[].{name:name, resourceGroup:resourceGroup, environmentId:managedEnvironmentId, location:location, latestRevisionName:latestRevisionName, ingressFqdn:properties.configuration.ingress.fqdn}"
	args := []string{"containerapp", "list", "-o", "json", "--query", q}
	if rg != "" {
		args = append(args, "-g", rg)
	}
	raw, err := RunAz(ctx, args...)
	if err != nil {
		return nil, err
	}
	var apps []m.ContainerApp
	if err := json.Unmarshal([]byte(raw), &apps); err != nil {
		return nil, err
	}
	return apps, nil
}

func GetAppDetails(ctx context.Context, name, rg string) (string, error) {
	raw, err := RunAz(ctx, "containerapp", "show", "-n", name, "-g", rg, "-o", "json")
	if err != nil {
		return "", err
	}
	var buf bytes.Buffer
	if err := json.Indent(&buf, []byte(raw), "", "  "); err != nil {
		return raw, nil
	}
	return buf.String(), nil
}

func ListRevisions(ctx context.Context, name string, rg string) ([]m.Revision, error) {
	q := "[].{name:name, active:properties.active, traffic:(properties.trafficWeight||`0`)}"
	raw, err := RunAz(ctx, "containerapp", "revision", "list", "-n", name, "-g", rg, "-o", "json", "--query", q)
	if err != nil {
		return nil, err
	}
	var revs []m.Revision
	if err := json.Unmarshal([]byte(raw), &revs); err != nil {
		return nil, err
	}
	return revs, nil
}
func ListContainersCmd(ctx context.Context, ct m.ContainerApp, revName string) ([]m.Container, error) {
	// az containerapp revision show --name <app> --resource-group <rg> --revision <rev>
	raw, err := RunAz(ctx, "containerapp", "revision", "show",
		"-n", ct.Name, "-g", ct.ResourceGroup, "--revision", revName)
	if err != nil {
		return nil, err
	}

	// Parse containers
	var resp struct {
		Properties struct {
			Template struct {
				Containers []struct {
					Name    string   `json:"name"`
					Image   string   `json:"image"`
					Command []string `json:"command"`
					Args    []string `json:"args"`
				} `json:"containers"`
			} `json:"template"`
		} `json:"properties"`
	}
	if jerr := json.Unmarshal([]byte(raw), &resp); jerr != nil {
		return nil, jerr
	}
	cs := make([]m.Container, 0, len(resp.Properties.Template.Containers))
	for _, c := range resp.Properties.Template.Containers {
		cs = append(cs, m.Container{Name: c.Name, Image: c.Image, Command: c.Command, Args: c.Args})
	}
	return cs, nil
}
