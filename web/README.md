# CiliKube Web (React)

Mainstream React control plane for the CiliKube monorepo.

## Stack (current)

- React 19 + TypeScript 7 + Vite 8
- Tailwind CSS v4 (`@tailwindcss/vite`)
- TanStack Query 5 + React Router 7
- Framer Motion + Recharts + Lucide
- Oxlint + pnpm (`packageManager` pinned)

Visual language is aligned with [CiliTerm](https://github.com/cillianxtech/ciliterm) (TRON / HUD), adapted for a Kubernetes admin shell.

## Develop

```bash
# backend (repo root)
go run ./cmd/server --config configs/config.yaml

# frontend
cd web
pnpm install
pnpm dev
```

App: http://localhost:8888  
API: `VITE_BASE_API` (default `http://localhost:8080`)

```bash
# from monorepo root
pnpm --dir web dev
pnpm --dir web build
pnpm --dir web lint
```

## Note on Vue

The previous Vue UI lives in the separate [`cilikube-web`](https://github.com/cillianxtech/cilikube-web) repository.
This monorepo’s frontend source of truth is `web/` (React).
A sync patch for recent Vue pages is under `../patches/`.
