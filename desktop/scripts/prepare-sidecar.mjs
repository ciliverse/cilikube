#!/usr/bin/env node
/**
 * Stage Go binary + web dist + casbin model for electron-builder extraResources.
 *
 * Env:
 *   SIDECAR_BIN   path to cilikube / cilikube.exe (required)
 *   WEB_DIST      path to web/dist (default: ../web/dist)
 *   REPO_ROOT     cilikube repo root (default: parent of desktop/)
 */
import fs from 'node:fs'
import path from 'node:path'
import { fileURLToPath } from 'node:url'

const __dirname = path.dirname(fileURLToPath(import.meta.url))
const desktopRoot = path.resolve(__dirname, '..')
const repoRoot = process.env.REPO_ROOT
  ? path.resolve(process.env.REPO_ROOT)
  : path.resolve(desktopRoot, '..')

const staging = path.join(desktopRoot, 'resources-staging')
const sidecarBin = process.env.SIDECAR_BIN
  ? path.resolve(process.env.SIDECAR_BIN)
  : path.join(repoRoot, 'build', 'desktop', process.platform === 'win32' ? 'cilikube.exe' : 'cilikube')
const webDist = process.env.WEB_DIST
  ? path.resolve(process.env.WEB_DIST)
  : path.join(repoRoot, 'web', 'dist')
const modelSrc = path.join(repoRoot, 'pkg', 'auth', 'model.conf')

function rmrf(p) {
  fs.rmSync(p, { recursive: true, force: true })
}

function copyDir(src, dest) {
  fs.mkdirSync(dest, { recursive: true })
  for (const ent of fs.readdirSync(src, { withFileTypes: true })) {
    const s = path.join(src, ent.name)
    const d = path.join(dest, ent.name)
    if (ent.isDirectory()) copyDir(s, d)
    else fs.copyFileSync(s, d)
  }
}

if (!fs.existsSync(sidecarBin)) {
  console.error(`sidecar binary missing: ${sidecarBin}`)
  console.error('Build first: go build -o build/desktop/cilikube.exe ./cmd/server')
  process.exit(1)
}
if (!fs.existsSync(path.join(webDist, 'index.html'))) {
  console.error(`web dist missing index.html: ${webDist}`)
  console.error('Build first: pnpm --dir web build')
  process.exit(1)
}
if (!fs.existsSync(modelSrc)) {
  console.error(`casbin model missing: ${modelSrc}`)
  process.exit(1)
}

rmrf(staging)
fs.mkdirSync(path.join(staging, 'bin'), { recursive: true })
const binName = path.basename(sidecarBin).endsWith('.exe') ? 'cilikube.exe' : 'cilikube'
fs.copyFileSync(sidecarBin, path.join(staging, 'bin', binName))
if (process.platform !== 'win32') {
  fs.chmodSync(path.join(staging, 'bin', binName), 0o755)
}

copyDir(webDist, path.join(staging, 'web'))
fs.mkdirSync(path.join(staging, 'pkg', 'auth'), { recursive: true })
fs.copyFileSync(modelSrc, path.join(staging, 'pkg', 'auth', 'model.conf'))

console.log('staged resources →', staging)
console.log('  bin/', binName)
console.log('  web/', '(from', webDist + ')')
console.log('  pkg/auth/model.conf')
