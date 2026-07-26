import {
  app,
  BrowserWindow,
  dialog,
  shell,
} from 'electron'
import { ChildProcessWithoutNullStreams, spawn } from 'node:child_process'
import { createServer } from 'node:net'
import fs from 'node:fs'
import path from 'node:path'
import crypto from 'node:crypto'
import http from 'node:http'
import os from 'node:os'

let mainWindow: BrowserWindow | null = null
let sidecar: ChildProcessWithoutNullStreams | null = null
let sidecarAddr = ''
const sidecarLog: string[] = []

function pushLog(line: string) {
  sidecarLog.push(line)
  if (sidecarLog.length > 200) sidecarLog.shift()
}

function resourcesRoot(): string {
  if (app.isPackaged) return process.resourcesPath
  // Dev: desktop/resources-staging
  return path.join(__dirname, '..', 'resources-staging')
}

function sidecarBinary(): string {
  const name = process.platform === 'win32' ? 'cilikube.exe' : 'cilikube'
  return path.join(resourcesRoot(), 'bin', name)
}

function defaultKubeconfig(): string {
  if (process.env.KUBECONFIG) return process.env.KUBECONFIG
  return path.join(os.homedir(), '.kube', 'config')
}

async function getFreePort(): Promise<number> {
  return new Promise((resolve, reject) => {
    const srv = createServer()
    srv.listen(0, '127.0.0.1', () => {
      const addr = srv.address()
      if (!addr || typeof addr === 'string') {
        srv.close()
        reject(new Error('failed to allocate port'))
        return
      }
      const port = addr.port
      srv.close((err) => (err ? reject(err) : resolve(port)))
    })
    srv.on('error', reject)
  })
}

function ensureDesktopConfig(configPath: string, dbPath: string): void {
  if (fs.existsSync(configPath)) return

  const jwtSecret = crypto.randomBytes(32).toString('hex')
  const encKey = crypto.randomBytes(16).toString('hex') + crypto.randomBytes(16).toString('hex')
  const kubeconfig = defaultKubeconfig().replace(/\\/g, '/')
  const db = dbPath.replace(/\\/g, '/')

  const yaml = `server:
  host: "127.0.0.1"
  port: "17880"
  read_timeout: 30
  write_timeout: 30
  mode: release
  activeCluster: ""
  encryptionKey: "${encKey}"
kubernetes:
  kubeconfig: "${kubeconfig}"
database:
  enabled: true
  type: sqlite
  host: ""
  port: 0
  username: ""
  password: ""
  database: "${db}"
  charset: ""
storage:
  type: database
  database:
    enabled: true
    type: sqlite
    host: ""
    port: 0
    username: ""
    password: ""
    database: "${db}"
    charset: ""
jwt:
  secret_key: "${jwtSecret}"
  expire_duration: 24h0m0s
  issuer: cilikube
oauth:
  github:
    client_id: ""
    client_secret: ""
    redirect_url: ""
    enabled: false
  allow_registration: true
  auto_link_accounts: false
security:
  password:
    min_length: 8
    require_uppercase: false
    require_lowercase: true
    require_numbers: true
    require_symbols: false
  session:
    max_idle_time: 1800
    max_session_time: 86400
    max_concurrent_sessions: 5
  login:
    max_attempts: 20
    lockout_duration: 300
    attempt_window: 600
  audit:
    enabled: true
    retention_days: 90
    log_login: true
    log_api: true
    log_admin: true
clusters: []
`
  fs.mkdirSync(path.dirname(configPath), { recursive: true })
  fs.mkdirSync(path.dirname(dbPath), { recursive: true })
  fs.writeFileSync(configPath, yaml, 'utf8')
}

function waitForHealth(addr: string, timeoutMs: number): Promise<void> {
  const deadline = Date.now() + timeoutMs
  return new Promise((resolve, reject) => {
    const tryOnce = () => {
      const req = http.get(`http://${addr}/health`, (res) => {
        res.resume()
        if (res.statusCode && res.statusCode >= 200 && res.statusCode < 300) {
          resolve()
          return
        }
        if (Date.now() > deadline) {
          reject(new Error(`health check failed: HTTP ${res.statusCode}`))
          return
        }
        setTimeout(tryOnce, 300)
      })
      req.on('error', () => {
        if (Date.now() > deadline) {
          reject(new Error('sidecar did not become healthy in time'))
          return
        }
        setTimeout(tryOnce, 300)
      })
      req.setTimeout(2000, () => req.destroy())
    }
    tryOnce()
  })
}

async function startSidecar(): Promise<string> {
  const bin = sidecarBinary()
  if (!fs.existsSync(bin)) {
    throw new Error(`sidecar binary not found: ${bin}`)
  }

  const port = await getFreePort()
  const addr = `127.0.0.1:${port}`
  const userData = app.getPath('userData')
  const configPath = path.join(userData, 'config.yaml')
  const dbPath = path.join(userData, 'cilikube.db')
  ensureDesktopConfig(configPath, dbPath)

  const webRoot = path.join(resourcesRoot(), 'web')
  const casbinModel = path.join(resourcesRoot(), 'pkg', 'auth', 'model.conf')
  const cwd = resourcesRoot()

  const env = {
    ...process.env,
    CILIKUBE_DESKTOP: '1',
    CILIKUBE_ADDR: addr,
    CILIKUBE_CONFIG: configPath,
    CILIKUBE_WEB_ROOT: webRoot,
    CILIKUBE_CASBIN_MODEL: casbinModel,
  }

  pushLog(`starting sidecar: ${bin}`)
  pushLog(`addr=${addr} cwd=${cwd}`)

  sidecar = spawn(bin, ['--config', configPath], {
    cwd,
    env,
    windowsHide: true,
  })

  sidecar.stdout.on('data', (buf) => {
    String(buf)
      .split(/\r?\n/)
      .filter(Boolean)
      .forEach((line) => pushLog(line))
  })
  sidecar.stderr.on('data', (buf) => {
    String(buf)
      .split(/\r?\n/)
      .filter(Boolean)
      .forEach((line) => pushLog(line))
  })
  sidecar.on('exit', (code, signal) => {
    pushLog(`sidecar exited code=${code} signal=${signal}`)
    sidecar = null
  })

  try {
    await waitForHealth(addr, 45000)
  } catch (err) {
    stopSidecar()
    const tail = sidecarLog.slice(-40).join('\n')
    throw new Error(`${(err as Error).message}\n\n--- sidecar log ---\n${tail}`)
  }

  sidecarAddr = addr
  return addr
}

function stopSidecar() {
  if (!sidecar) return
  const child = sidecar
  sidecar = null
  try {
    if (process.platform === 'win32') {
      spawn('taskkill', ['/pid', String(child.pid), '/f', '/t'], { windowsHide: true })
    } else {
      child.kill('SIGTERM')
      setTimeout(() => {
        try {
          child.kill('SIGKILL')
        } catch {
          /* ignore */
        }
      }, 2000)
    }
  } catch {
    /* ignore */
  }
}

async function createWindow(addr: string) {
  mainWindow = new BrowserWindow({
    width: 1440,
    height: 900,
    minWidth: 1100,
    minHeight: 700,
    title: 'CiliKube',
    webPreferences: {
      preload: path.join(__dirname, 'preload.js'),
      contextIsolation: true,
      nodeIntegration: false,
      sandbox: true,
    },
  })

  mainWindow.webContents.setWindowOpenHandler(({ url }) => {
    shell.openExternal(url)
    return { action: 'deny' }
  })

  await mainWindow.loadURL(`http://${addr}/`)
  mainWindow.on('closed', () => {
    mainWindow = null
  })
}

const gotLock = app.requestSingleInstanceLock()
if (!gotLock) {
  app.quit()
} else {
  app.on('second-instance', () => {
    if (mainWindow) {
      if (mainWindow.isMinimized()) mainWindow.restore()
      mainWindow.focus()
    }
  })

  app.whenReady().then(async () => {
    try {
      const addr = await startSidecar()
      await createWindow(addr)
    } catch (err) {
      dialog.showErrorBox('CiliKube failed to start', String((err as Error).message || err))
      app.quit()
    }
  })

  app.on('before-quit', () => {
    stopSidecar()
  })

  app.on('window-all-closed', () => {
    stopSidecar()
    app.quit()
  })
}

// silence unused in some builds
void sidecarAddr
