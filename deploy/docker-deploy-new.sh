#!/usr/bin/env bash
set -euo pipefail

PROJECT_NAME="${PROJECT_NAME:-sub2api-new}"
COMPOSE_FILE="${COMPOSE_FILE:-docker-compose.new.yml}"
APP_SERVICE="${APP_SERVICE:-sub2api-new}"
DEFAULT_BIND_HOST="${DEFAULT_BIND_HOST:-127.0.0.1}"
DEFAULT_PORT="${DEFAULT_PORT:-18080}"

RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m'

print_info() { echo -e "${BLUE}[INFO]${NC} $1"; }
print_success() { echo -e "${GREEN}[OK]${NC} $1"; }
print_warning() { echo -e "${YELLOW}[WARN]${NC} $1"; }
print_error() { echo -e "${RED}[ERROR]${NC} $1"; }

command_exists() {
    command -v "$1" >/dev/null 2>&1
}

generate_secret() {
    openssl rand -hex 32
}

env_get() {
    local key="$1"
    awk -F= -v key="$key" '
      $1 == key {
        sub(/^[^=]*=/, "")
        value = $0
      }
      END {
        if (value != "") {
          print value
        }
      }
    ' .env
}

env_set() {
    local key="$1"
    local value="$2"
    local tmp
    tmp="$(mktemp)"
    awk -v key="$key" -v value="$value" '
      BEGIN { done = 0 }
      index($0, key "=") == 1 {
        print key "=" value
        done = 1
        next
      }
      { print }
      END {
        if (!done) {
          print key "=" value
        }
      }
    ' .env > "$tmp"
    mv "$tmp" .env
}

ensure_deploy_dir() {
    local script_dir
    script_dir="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
    cd "$script_dir"

    if [[ ! -f "$COMPOSE_FILE" ]]; then
        print_error "Missing ${COMPOSE_FILE}. Run this script from the deploy directory or restore the file."
        exit 1
    fi

    if [[ ! -f ".env.example" ]]; then
        print_error "Missing .env.example."
        exit 1
    fi
}

ensure_env_file() {
    if [[ ! -f ".env" ]]; then
        print_info "Creating .env from .env.example"
        cp .env.example .env
    else
        print_info "Using existing .env"
    fi

    env_set BIND_HOST "$DEFAULT_BIND_HOST"
    env_set SERVER_PORT "$DEFAULT_PORT"

    local postgres_password
    postgres_password="$(env_get POSTGRES_PASSWORD || true)"
    if [[ -z "$postgres_password" || "$postgres_password" == "change_this_secure_password" ]]; then
        env_set POSTGRES_PASSWORD "$(generate_secret)"
        print_info "Generated POSTGRES_PASSWORD"
    fi

    local jwt_secret
    jwt_secret="$(env_get JWT_SECRET || true)"
    if [[ -z "$jwt_secret" ]]; then
        env_set JWT_SECRET "$(generate_secret)"
        print_info "Generated JWT_SECRET"
    fi

    local totp_key
    totp_key="$(env_get TOTP_ENCRYPTION_KEY || true)"
    if [[ -z "$totp_key" ]]; then
        env_set TOTP_ENCRYPTION_KEY "$(generate_secret)"
        print_info "Generated TOTP_ENCRYPTION_KEY"
    fi

    chmod 600 .env
}

check_runtime() {
    if ! command_exists openssl; then
        print_error "openssl is not installed."
        exit 1
    fi

    if ! command_exists docker; then
        print_error "docker is not installed."
        exit 1
    fi

    if ! docker compose version >/dev/null 2>&1; then
        print_error "docker compose is not available."
        exit 1
    fi

    if docker ps -a --format '{{.Names}}' | grep -Eq '^(sub2api-new|sub2api-new-postgres|sub2api-new-redis)$'; then
        print_info "Existing sub2api-new containers found; they will be updated in place."
    fi
}

compose_cmd() {
    docker compose -p "$PROJECT_NAME" -f "$COMPOSE_FILE" "$@"
}

wait_for_health() {
    local bind_host
    local port
    local url
    bind_host="$(env_get BIND_HOST || true)"
    port="$(env_get SERVER_PORT || true)"
    bind_host="${bind_host:-$DEFAULT_BIND_HOST}"
    port="${port:-$DEFAULT_PORT}"
    url="http://${bind_host}:${port}/health"

    print_info "Waiting for Sub2API health endpoint: ${url}"
    for _ in $(seq 1 60); do
        if command_exists curl && curl -fsS "$url" >/dev/null 2>&1; then
            print_success "Health check passed"
            return 0
        fi

        if command_exists wget && wget -q -T 3 -O /dev/null "$url" >/dev/null 2>&1; then
            print_success "Health check passed"
            return 0
        fi

        sleep 2
    done

    print_error "Health check timed out"
    compose_cmd logs --tail=120 "$APP_SERVICE" || true
    exit 1
}

main() {
    echo ""
    echo "=========================================="
    echo "  Sub2API Parallel New Deployment"
    echo "=========================================="
    echo ""

    ensure_deploy_dir
    check_runtime
    ensure_env_file
    mkdir -p data-new postgres_data-new redis_data-new

    print_info "Pulling dependency images"
    compose_cmd pull postgres redis

    print_info "Building current workspace image"
    compose_cmd build "$APP_SERVICE"

    print_info "Starting ${PROJECT_NAME}"
    compose_cmd up -d
    wait_for_health

    echo ""
    print_success "Sub2API new stack is up"
    echo ""
    echo "Deploy directory: $(pwd)"
    echo "Compose file:     ${COMPOSE_FILE}"
    echo "Environment:      $(pwd)/.env"
    echo "Web UI:           http://${DEFAULT_BIND_HOST}:${DEFAULT_PORT}"
    echo ""
    echo "Useful commands:"
    echo "  cd $(pwd)"
    echo "  docker compose -p ${PROJECT_NAME} -f ${COMPOSE_FILE} ps"
    echo "  docker compose -p ${PROJECT_NAME} -f ${COMPOSE_FILE} logs -f ${APP_SERVICE}"
    echo "  docker compose -p ${PROJECT_NAME} -f ${COMPOSE_FILE} down"
    echo ""
    print_info "The existing sub2api stack on 8080 is not stopped or modified."
}

main "$@"
