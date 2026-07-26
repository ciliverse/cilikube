# CiliKube Desktop

Electron shell that launches the Go sidecar and loads the UI from `http://127.0.0.1:<port>`.

## Local (developers)

Windows packaging must run on Windows (or CI `windows-latest`) because `mattn/go-sqlite3` needs CGO.

```bash
# from repo root
pnpm --dir web install && pnpm --dir web build
# on Windows:
go build -ldflags "-w -s" -o build/desktop/cilikube.exe ./cmd/server
pnpm --dir desktop install
node desktop/scripts/prepare-sidecar.mjs
pnpm --dir desktop pack:win
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

Default first login: `admin` / `12345678`.
