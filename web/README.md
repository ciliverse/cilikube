# CiliKube Web (React)

Mainstream React control plane for CiliKube.

## Stack

- React 19 + TypeScript + Vite
- Tailwind CSS v4
- TanStack Query
- React Router
- Framer Motion
- Recharts + Lucide

## Develop

```bash
# backend
go run cmd/server/main.go

# frontend
cd web
pnpm install
pnpm dev
```

App runs at http://localhost:8888 and talks to `VITE_BASE_API` (default `http://localhost:8080`).

## Note on Vue

The previous Vue UI lives in the separate [`cilikube-web`](https://github.com/cillianxtech/cilikube-web) repository.
A sync patch for recent Vue pages is under `../patches/`.
