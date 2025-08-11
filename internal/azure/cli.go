package azure

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"os/exec"
	"strings"

	m "github.com/IAL32/az-tui/internal/models"
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
	q := `[].{
		name:name,
		resourceGroup:resourceGroup,
		environmentId:managedEnvironmentId,
		location:location,
		latestRevisionName:latestRevisionName,
		ingressFqdn:properties.configuration.ingress.fqdn,
		provisioningState:properties.provisioningState,
		runningStatus:properties.runningStatus,
		minReplicas:properties.template.scale.minReplicas,
		maxReplicas:properties.template.scale.maxReplicas,
		cpu:properties.template.containers[0].resources.cpu,
		memory:properties.template.containers[0].resources.memory,
		ingressExternal:properties.configuration.ingress.external,
		targetPort:properties.configuration.ingress.targetPort,
		identityType:identity.type,
		workloadProfile:properties.workloadProfileName,
		createdAt:systemData.createdAt,
		lastModifiedAt:systemData.lastModifiedAt
	}`
	args := []string{"containerapp", "list", "-o", "json", "--query", q}
	if rg != "" {
		args = append(args, "-g", rg)
	}
	raw, err := RunAz(ctx, args...)
	if err != nil {
		return nil, err
	}
	return TransformContainerAppsFromJSON(raw)
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
	q := `[].{
		name:name,
		active:properties.active,
		traffic:(properties.trafficWeight||` + "`0`" + `),
		createdAt:properties.createdTime,
		fqdn:properties.fqdn,
		replicas:properties.replicas,
		healthState:properties.healthState,
		provisioningState:properties.provisioningState,
		runningState:properties.runningState,
		minReplicas:properties.template.scale.minReplicas,
		maxReplicas:properties.template.scale.maxReplicas,
		cpu:properties.template.containers[0].resources.cpu,
		memory:properties.template.containers[0].resources.memory
	}`
	raw, err := RunAz(ctx, "containerapp", "revision", "list", "-n", name, "-g", rg, "-o", "json", "--query", q)
	if err != nil {
		return nil, err
	}
	return TransformRevisionsFromJSON(raw)
}
func ListContainersCmd(ctx context.Context, ct m.ContainerApp, revName string) ([]m.Container, error) {
	// az containerapp revision show --name <app> --resource-group <rg> --revision <rev>
	raw, err := RunAz(ctx, "containerapp", "revision", "show",
		"-n", ct.Name, "-g", ct.ResourceGroup, "--revision", revName)
	if err != nil {
		return nil, err
	}
	return TransformContainersFromJSON(raw)
}

func ListResourceGroups(ctx context.Context) ([]m.ResourceGroup, error) {
	q := `[].{
		name:name,
		location:location,
		provisioningState:properties.provisioningState,
		tags:tags
	}`
	raw, err := RunAz(ctx, "group", "list", "-o", "json", "--query", q)
	if err != nil {
		return nil, err
	}
	return TransformResourceGroupsFromJSON(raw)
}
