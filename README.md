# PhoneInfoga Desktop (Go + Wails)

**Goal:** Cross-platform desktop GUI (Windows/macOS/Linux) that orchestrates **PhoneInfoga** using **serve (REST)** preferred and **CLI** fallback. Includes scans (single/batch), history/cases (SQLite via ent), evidence, and export/reporting.

> **You are viewing a bootstrap repo.** It contains folder structure, Go interfaces/stubs, ent schemas, and Wails binding contracts. You will still run `wails init` to generate the actual web UI shell, then wire these packages. Follow the steps below.

---

## Quick Start

1. **Prereqs**
   - Go 1.22+
   - Node 18+ / PNPM or NPM
   - Wails v3: `go install github.com/wailsapp/wails/v3/cmd/wails@latest`
   - Ent: `go install entgo.io/ent/cmd/ent@latest`
   - wkhtmltopdf (for nicer PDF reports) — optional initially
   - PhoneInfoga binaries (we can bundle later), or install system-wide

2. **Create UI shell with Wails**
   ```bash
   wails init -n frontend-shell -t svelte
   ```
   This creates a new folder (e.g. `frontend-shell`). Copy its `frontend/` contents into the `frontend/` folder of this repo, and copy its `wails.json` to the repo root. Adjust app name/ID as you prefer.

3. **Wire Go backend**
   - In `cmd/desktop/main.go`, import Wails and point to `internal/app` bootstrap.
   - Implement TODOs in `internal/phoneinfoga/manager.go` to actually run/parse PhoneInfoga.
   - Fill `pkg/uiapi` methods and bind them via Wails in `internal/app/bootstrap.go`.

4. **Install deps & generate code**
   ```bash
   make tools
   make generate
   ```

5. **Run in dev**
   ```bash
   make dev
   ```

6. **Build installers**
   ```bash
   make build
   make package
   ```

---

## Project layout

```
cmd/desktop/               # Wails bootstrap; app entry
internal/app/              # App init, DI wiring, versioning, logging
internal/cfg/              # Config manager (Viper), paths
internal/secrets/          # OS keyring wrapper
internal/phoneinfoga/      # Serve/CLI manager + parsers
internal/scan/             # Job orchestrator (queue, cancel, retries)
internal/store/ent/schema/ # Ent schemas (Scan, Result, Case, Artifact, Setting)
internal/report/           # HTML templates, PDF export, hashing
internal/evidence/         # Artifact ingestion & hashing
internal/netx/             # HTTP client factory (proxy, timeouts, UA)
internal/audit/            # Provenance stamps & version info
pkg/dto/                   # Shared DTOs (Go) mirrored to TS
pkg/uiapi/                 # Wails bindings (Go <-> JS)
frontend/                  # Wails Svelte UI (populate after wails init)
scripts/                   # helper scripts
```

---

## Notes

- **Binary strategy:** during packaging, embed per-OS PhoneInfoga binaries under `resources/bin/{{win,mac,linux}}` and copy on first run to an app-data dir (`~/.local/share`, `%APPDATA%`, etc.). Provide a "Use system binary" toggle in Settings.
- **Secrets:** use OS keychain (macOS Keychain, Windows Credential Manager, libsecret/KWallet on Linux) via `go-keyring`.
- **Compliance:** keep raw JSON, timestamps, hashes, and PhoneInfoga version for chain-of-custody.

© 2025 Your Company. MIT for GUI code; PhoneInfoga licensed per upstream and distributed unmodified.
