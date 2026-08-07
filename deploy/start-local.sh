#!/usr/bin/env bash
set -euo pipefail

# Build and run the local Docker development stack after releasing its HTTP port.
#
# Usage:
#   ./deploy/start-local.sh            # foreground logs (recommended while debugging)
#   ./deploy/start-local.sh --detach   # run in the background
#
# Environment overrides:
#   COMPOSE_FILE=docker-compose.dev.yml PROJECT_NAME=sub2api-dev SERVER_PORT=8080

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
COMPOSE_FILE="${COMPOSE_FILE:-docker-compose.dev.yml}"
PROJECT_NAME="${PROJECT_NAME:-sub2api-dev}"
DETACH=false

case "${1:-}" in
  "") ;;
  --detach|-d) DETACH=true ;;
  -h|--help)
    sed -n '2,12p' "$0"
    exit 0
    ;;
  *)
    echo "Unknown option: $1" >&2
    exit 2
    ;;
esac

cd "$SCRIPT_DIR"

if ! command -v docker >/dev/null 2>&1; then
  echo "Docker is required." >&2
  exit 1
fi

if ! docker compose version >/dev/null 2>&1; then
  echo "Docker Compose is required." >&2
  exit 1
fi

ensure_docker_daemon() {
  if docker info >/dev/null 2>&1; then
    return 0
  fi

  if [[ "$(uname -s)" == "Darwin" ]] && command -v open >/dev/null 2>&1; then
    echo "Docker daemon is not running; opening Docker Desktop..."
    open -a Docker >/dev/null 2>&1 || true
    for _ in {1..60}; do
      sleep 1
      if docker info >/dev/null 2>&1; then
        echo "Docker daemon is ready."
        return 0
      fi
    done
  fi

  echo "Docker daemon is unavailable. Start Docker Desktop (or your Docker daemon) and run this script again." >&2
  exit 1
}

ensure_docker_daemon

if ! command -v lsof >/dev/null 2>&1; then
  echo "lsof is required to release an occupied port." >&2
  exit 1
fi

if [[ ! -f "$COMPOSE_FILE" ]]; then
  echo "Compose file not found: $SCRIPT_DIR/$COMPOSE_FILE" >&2
  exit 1
fi

read_env_value() {
  local key="$1"
  [[ -f .env ]] || return 0
  awk -F= -v key="$key" '
    $0 !~ /^[[:space:]]*#/ && $1 == key {
      sub(/^[^=]*=/, "")
      value = $0
    }
    END { if (value != "") print value }
  ' .env
}

PORT="${SERVER_PORT:-$(read_env_value SERVER_PORT)}"
PORT="${PORT:-8080}"

if [[ ! "$PORT" =~ ^[0-9]{1,5}$ ]] || (( PORT < 1 || PORT > 65535 )); then
  echo "SERVER_PORT must be a valid TCP port; got: $PORT" >&2
  exit 2
fi

compose() {
  docker compose -p "$PROJECT_NAME" -f "$COMPOSE_FILE" "$@"
}

release_port() {
  local pids
  pids="$(lsof -tiTCP:"$PORT" -sTCP:LISTEN 2>/dev/null || true)"
  [[ -n "$pids" ]] || return 0

  echo "Releasing TCP port $PORT (PID(s): ${pids//$'\n'/, })"
  kill $pids 2>/dev/null || true

  for _ in {1..20}; do
    sleep 0.25
    pids="$(lsof -tiTCP:"$PORT" -sTCP:LISTEN 2>/dev/null || true)"
    [[ -z "$pids" ]] && return 0
  done

  echo "Force-stopping remaining listener(s) on TCP port $PORT: ${pids//$'\n'/, }"
  kill -9 $pids 2>/dev/null || true
  sleep 0.25

  if lsof -tiTCP:"$PORT" -sTCP:LISTEN >/dev/null 2>&1; then
    echo "Could not release TCP port $PORT." >&2
    exit 1
  fi
}

echo "Stopping the previous local stack, if any..."
compose down --remove-orphans >/dev/null 2>&1 || true
release_port

echo "Starting local development stack on http://127.0.0.1:$PORT"
if [[ "$DETACH" == true ]]; then
  compose up --build -d
  echo "Running in the background. Follow logs with:"
  echo "  docker compose -p $PROJECT_NAME -f $COMPOSE_FILE logs -f sub2api"
else
  exec docker compose -p "$PROJECT_NAME" -f "$COMPOSE_FILE" up --build
fi
