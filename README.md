# Az-TUI

Az-TUI is a **terminal-based user interface (TUI)** for managing Azure Container Apps, inspired by tools like K9s. It allows you to browse your Container Apps, inspect details and revisions, and access container logs or shells, all from within a single terminal UI.

![Demo](demo.gif)

## Features

- **Browse Azure Container Apps** across your subscription (or limit to a resource group via `ACA_RG`).
- **View detailed app information** (JSON) including name, resource group, location, ingress FQDN, and latest revision.
- **Inspect revisions** with active indicators and traffic percentages.
- **Tail logs** for apps, revisions, or containers.
- **Exec into running containers** for debugging.
- **Keyboard-driven navigation** with familiar shortcuts.
- **Mock data mode** for development and testing without Azure CLI dependencies.

## Key Bindings

### Global

- `q` / `Ctrl+C` – Quit
- `/` – Filter list
- `Esc` – Go back (in nested views)
- `Enter` – Select/drill down
- `?` – Toggle help
- `:` – Context switching (VIM/k9s-like navigation)
- `Shift+←` / `Shift+→` – Scroll table left/right

### Context Switching

Az-TUI now supports **VIM/k9s-like context switching** for quick navigation between different views:

- Press `:` to open the context selection menu
- Use `↑`/`k` and `↓`/`j` to navigate between available contexts
- Press `Enter` to switch to the selected context
- Press `Esc` to cancel context switching

**Context behavior:**
- **From Resource Groups**: Switch to Container Apps or Jobs (preserves or clears resource group selection)
- **From Container Apps**: Switch between Apps and Jobs (preserves current resource group)
- **From Revisions**: Stay in Revisions view (preserves resource group and app selection)
- **From Containers**: Stay in Containers view (preserves all current selections)
- **From Environment Variables**: Stay in Env Vars view (preserves all selections)

The context menu shows only relevant navigation options for your current mode and automatically preserves your selection state when switching contexts.

### Resource Groups Mode

- `r` – Refresh resource groups
- `Enter` – Select resource group and view apps

### Apps Mode

- `r` – Refresh apps
- `R` – Restart revision
- `l` – Logs for app
- `s` – Exec into app
- `v` – View environment variables
- `Enter` – View revisions for app

### Revisions Mode

- `r` – Refresh revisions
- `R` – Restart revision
- `l` – Logs for revision
- `s` – Exec into revision
- `Enter` – View containers in revision

### Containers Mode

- `r` – Refresh containers
- `l` – Logs for container
- `s` – Exec into container
- `v` – View environment variables
- `Enter` – View environment variables for container

### Environment Variables Mode

- `r` – Refresh environment variables
- `Esc` – Go back to previous mode

## Installation

**Prerequisites:**

- Go 1.20+
- Azure CLI (`az login`)
- Azure Container Apps extension: `az extension add -n containerapp`

**Build:**

```bash
git clone https://github.com/IAL32/az-tui.git
cd az-tui
go build -o az-tui cmd/az-tui/main.go
```

Or install directly:

```bash
go install github.com/IAL32/az-tui/cmd/az-tui@latest
```

## Usage

### Standard Mode (Azure CLI)

(Optional) restrict to a resource group:

```bash
export ACA_RG="my-resource-group"
```

Run:

```bash
./az-tui
```

### Mock Data Mode

For development, testing, or demonstration purposes, you can run Az-TUI with mock data instead of connecting to Azure:

```bash
./az-tui --mock
# or
./az-tui -m
```

Mock mode provides a comprehensive dataset including:
- 4 resource groups (production, staging, development, shared services)
- 8 container apps across different environments
- Multiple revisions per app with realistic configurations
- Containers with environment variables, probes, and volume mounts
- Realistic Azure Container Apps scenarios for testing UI functionality

Navigate with arrow keys or `j`/`k`, drill down with `Enter`, and use the key bindings above for actions.

## Architecture

Az-TUI uses the [Bubble Tea](https://github.com/charmbracelet/bubbletea) framework:

- **Modes:** `resource groups` → `apps` → `revisions` → `containers` → `environment variables`
- **Context switching:** VIM/k9s-like navigation system with `:` key for quick mode switching
- **Data providers:** Pluggable architecture supporting both Azure CLI and mock data sources
- **Azure CLI integration:** Fetches data using `az containerapp` and `az group` commands
- **Mock data system:** JSON-based mock data for development and testing
- **UI Components:** Bubble Table for data display with filtering and navigation
- **Asynchronous updates:** Commands run in background and update the model via messages
- **Help system:** Built-in help with `?` key showing context-sensitive keybindings
- **State preservation:** Context switching maintains current selections across mode changes

## Roadmap

- [x] Fuzzy search
- [ ] Subscription/environment switching
- [ ] Show container replica health
- [ ] Edit traffic split allocations
- [ ] Browse Azure Container Apps Jobs
- [ ] Integrate metrics (CPU/memory, HTTP rates)

## License

MIT License. See [LICENSE](./LICENSE).
