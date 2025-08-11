# Az-TUI

Az-TUI is a **terminal-based user interface (TUI)** for managing Azure Container Apps, inspired by tools like K9s. It allows you to browse your Container Apps, inspect details and revisions, and access container logs or shells, all from within a single terminal UI.

## Features

- **Browse Azure Container Apps** across your subscription (or limit to a resource group via `ACA_RG`).
- **View detailed app information** (JSON) including name, resource group, location, ingress FQDN, and latest revision.
- **Inspect revisions** with active indicators and traffic percentages.
- **Tail logs** for apps, revisions, or containers.
- **Exec into running containers** for debugging.
- **Keyboard-driven navigation** with familiar shortcuts.

## Key Bindings

### Global

- `q` / `Ctrl+C` – Quit
- `/` – Filter list
- `Esc` – Go back (in nested views)
- `Enter` – Select/drill down
- `?` – Toggle help
- `Shift+←` / `Shift+→` – Scroll table left/right

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

(Optional) restrict to a resource group:

```bash
export ACA_RG="my-resource-group"
```

Run:

```bash
./az-tui
```

Navigate with arrow keys or `j`/`k`, drill down with `Enter`, and use the key bindings above for actions.

## Architecture

Az-TUI uses the [Bubble Tea](https://github.com/charmbracelet/bubbletea) framework:

- **Modes:** `resource groups` → `apps` → `revisions` → `containers` → `environment variables`
- **Azure CLI integration:** Fetches data using `az containerapp` and `az group` commands
- **UI Components:** Bubble Table for data display with filtering and navigation
- **Asynchronous updates:** Commands run in background and update the model via messages
- **Help system:** Built-in help with `?` key showing context-sensitive keybindings

## Roadmap

- Fuzzy search
- Subscription/environment switching
- Show container replica health
- Edit traffic split allocations
- Browse Azure Container Apps Jobs
- Integrate metrics (CPU/memory, HTTP rates)

## License

MIT License. See [LICENSE](./LICENSE).
