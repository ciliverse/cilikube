# CiliKube Desktop

Electron shell that launches the Go sidecar and loads the UI from `http://127.0.0.1:<port>`.

## Artifacts (CI)

| Platform | Package |
| --- | --- |
| Windows x64 | `CiliKube-*-win-x64-setup.exe` / portable `.exe` |
| macOS Apple Silicon | `CiliKube-*-mac-arm64.dmg` |
| Linux x64 | `CiliKube-*-linux-x64.AppImage` |

Triggered by tags like `v0.9.2-desktop.1` via `.github/workflows/release-desktop.yml`.

Packaging must run on the target OS (or matching GitHub runner) because `gorm.io/driver/sqlite` → `mattn/go-sqlite3` needs CGO.

## Local (developers)

```bash
# from repo root
pnpm --dir web install && pnpm --dir web build

# build native sidecar on this OS
mkdir -p build/desktop
# Windows:
#   go build -ldflags "-w -s" -o build/desktop/cilikube.exe ./cmd/server
# macOS / Linux:
go build -ldflags "-w -s" -o build/desktop/cilikube ./cmd/server

pnpm --dir desktop install
node desktop/scripts/prepare-sidecar.mjs

# pick one:
pnpm --dir desktop pack:win    # Windows only
pnpm --dir desktop pack:mac    # macOS arm64 only
pnpm --dir desktop pack:linux  # Linux AppImage only
```

Outputs under `desktop/release/`.

## Runtime env (set by Electron)

| Env | Purpose |
| --- | --- |
| `CILIKUBE_DESKTOP=1` | Desktop defaults (loopback, OAuth off) |
| `CILIKUBE_ADDR` | `127.0.0.1:<port>` |
| `CILIKUBE_CONFIG` | Config under Electron `userData` |
| `CILIKUBE_WEB_ROOT` | Packaged `resources/web` |
| `CILIKUBE_CASBIN_MODEL` | Packaged `resources/pkg/auth/model.conf` |

Default first login: `admin` / `12345678` (forced password change on first login).
