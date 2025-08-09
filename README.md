# AzulTUI

AzulTUI is a **K9s-like Terminal UI** for managing **Azure Container Apps**, built with [Charm's Bubble Tea](https://github.com/charmbracelet/bubbletea). It wraps the Azure CLI for fast iteration and leverages your existing `az` authentication and context.

---

## Features

- **Browse Container Apps** in your subscription (optional filter by resource group via `ACA_RG`)
- **View detailed JSON** for the selected app (name, resource group, location, FQDN, latest revision)
- **Inspect revisions** with active flags and traffic split
- **Tail logs** for live debugging (`l` key)
- **Exec into running container** (`e` key)
- **Keyboard navigation**:
  - `q` — quit
  - `r` — refresh apps
  - `R` — reload revisions
  - `tab` — switch between details and revisions pane
  - `/` — filter list

---

## Installation

```bash
go mod init az-tui
go get github.com/charmbracelet/bubbletea \
       github.com/charmbracelet/bubbles \
       github.com/charmbracelet/lipgloss
```

You must have:

- **Azure CLI** installed
- **Container Apps CLI extension**:

```bash
az extension add -n containerapp
```

---

## Usage

```bash
# Optional: scope to a resource group
export ACA_RG="my-resource-group"

# Run AzulTUI
go run .
```

---

## Roadmap

- Fuzzy search across apps
- Subscription and environment filters
- Replica health status
- Traffic split editing
- Job run browsing
- Metrics integration

---

## License

MIT
