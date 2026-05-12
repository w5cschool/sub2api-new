#!/usr/bin/env bash
set -euo pipefail

ROOT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"
FRONTEND_DIR="$ROOT_DIR/frontend"

FRONTEND_HOST="${FRONTEND_HOST:-127.0.0.1}"
FRONTEND_PORT="${FRONTEND_PORT:-3000}"
MOCK_BACKEND_HOST="${MOCK_BACKEND_HOST:-127.0.0.1}"
MOCK_BACKEND_PORT="${MOCK_BACKEND_PORT:-18080}"
OPEN_BROWSER="${OPEN_BROWSER:-1}"
USE_REAL_BACKEND="${USE_REAL_BACKEND:-0}"
REAL_BACKEND_URL="${REAL_BACKEND_URL:-http://127.0.0.1:8080}"

MOCK_BACKEND_PID=""

cleanup() {
  if [[ -n "$MOCK_BACKEND_PID" ]] && kill -0 "$MOCK_BACKEND_PID" >/dev/null 2>&1; then
    kill "$MOCK_BACKEND_PID" >/dev/null 2>&1 || true
  fi
}
trap cleanup EXIT INT TERM

require_command() {
  if ! command -v "$1" >/dev/null 2>&1; then
    echo "Missing required command: $1" >&2
    exit 1
  fi
}

wait_for_url() {
  local url="$1"
  local tries="${2:-40}"
  local delay="${3:-0.25}"

  for _ in $(seq 1 "$tries"); do
    if curl -fsS "$url" >/dev/null 2>&1; then
      return 0
    fi
    sleep "$delay"
  done

  echo "Timed out waiting for $url" >&2
  return 1
}

start_mock_backend() {
  local mock_file
  mock_file="$(mktemp "${TMPDIR:-/tmp}/sub2api-local-mock.XXXXXX.mjs")"

  cat >"$mock_file" <<'NODE'
import http from 'node:http'

const host = process.env.MOCK_BACKEND_HOST || '127.0.0.1'
const port = Number(process.env.MOCK_BACKEND_PORT || '18080')

const publicSettings = {
  site_name: 'tocodex',
  site_subtitle: 'GPT Codex Gateway',
  site_logo: '',
  site_version: 'local',
  api_base_url: 'https://tocodex.cc',
  contact_info: 'Icanmeetu',
  doc_url: '',
  home_content: '',
  payment_enabled: true,
  backend_mode_enabled: false,
  registration_enabled: true,
  show_request_prices: false,
  custom_menu_items: []
}

function sendJson(res, status, body) {
  const payload = JSON.stringify(body)
  res.writeHead(status, {
    'content-type': 'application/json; charset=utf-8',
    'cache-control': 'no-store',
    'access-control-allow-origin': '*',
    'access-control-allow-methods': 'GET,POST,OPTIONS',
    'access-control-allow-headers': 'content-type,authorization'
  })
  res.end(payload)
}

const server = http.createServer((req, res) => {
  if (req.method === 'OPTIONS') {
    res.writeHead(204, {
      'access-control-allow-origin': '*',
      'access-control-allow-methods': 'GET,POST,OPTIONS',
      'access-control-allow-headers': 'content-type,authorization'
    })
    res.end()
    return
  }

  const url = new URL(req.url || '/', `http://${req.headers.host || `${host}:${port}`}`)

  if (url.pathname === '/api/v1/settings/public') {
    sendJson(res, 200, { code: 0, message: 'ok', data: publicSettings })
    return
  }

  if (url.pathname === '/setup/status') {
    sendJson(res, 200, { code: 0, message: 'ok', data: { initialized: true, setup_required: false } })
    return
  }

  if (url.pathname === '/healthz') {
    sendJson(res, 200, { ok: true })
    return
  }

  sendJson(res, 404, { code: 404, message: `mock endpoint not implemented: ${url.pathname}` })
})

server.listen(port, host, () => {
  console.log(`[local-verify] mock backend listening at http://${host}:${port}`)
})

process.on('SIGTERM', () => server.close(() => process.exit(0)))
process.on('SIGINT', () => server.close(() => process.exit(0)))
NODE

  MOCK_BACKEND_HOST="$MOCK_BACKEND_HOST" MOCK_BACKEND_PORT="$MOCK_BACKEND_PORT" node "$mock_file" &
  MOCK_BACKEND_PID="$!"
  wait_for_url "http://$MOCK_BACKEND_HOST:$MOCK_BACKEND_PORT/healthz"
}

require_command node
require_command pnpm
require_command curl

if [[ ! -d "$FRONTEND_DIR/node_modules" ]]; then
  echo "[local-verify] frontend dependencies missing; running pnpm install --frozen-lockfile"
  pnpm --dir "$FRONTEND_DIR" install --frozen-lockfile
fi

if [[ "$USE_REAL_BACKEND" == "1" ]]; then
  BACKEND_URL="$REAL_BACKEND_URL"
  echo "[local-verify] using real backend: $BACKEND_URL"
else
  start_mock_backend
  BACKEND_URL="http://$MOCK_BACKEND_HOST:$MOCK_BACKEND_PORT"
fi

FRONTEND_URL="http://$FRONTEND_HOST:$FRONTEND_PORT/home"

echo "[local-verify] starting Vite frontend"
echo "[local-verify] frontend: $FRONTEND_URL"
echo "[local-verify] proxy target: $BACKEND_URL"

if [[ "$OPEN_BROWSER" == "1" ]] && command -v open >/dev/null 2>&1; then
  (sleep 1.5 && open "$FRONTEND_URL") &
fi

cd "$FRONTEND_DIR"
VITE_DEV_PORT="$FRONTEND_PORT" VITE_DEV_PROXY_TARGET="$BACKEND_URL" pnpm dev --host "$FRONTEND_HOST"
