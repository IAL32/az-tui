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

func ListRevisions(ctx context.Context, name, rg string) ([]m.Revision, error) {
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
