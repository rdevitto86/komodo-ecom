set shell := ["bash", "-cu"]

SERVICES_FILE := "infra/local/services.jsonc"
COMPOSE       := "docker compose -f infra/local/docker-compose.yml"
COMPOSE_PROXY := "docker compose -f infra/local/docker-compose.yml -f infra/local/docker-compose.proxy.yml"

# Show available commands
help:
    @just --list --unsorted

# Install repo toolchain via Homebrew
bootstrap:
    #!/usr/bin/env bash
    which brew > /dev/null || (echo "Install Homebrew first: https://brew.sh" && exit 1)
    brew bundle

# ── Local Dev ──────────────────────────────────────────────────────────────
#
#   just up                    everything enabled in services.jsonc (local)
#   just up api                api services only (local)
#   just up ui                 ui only (local)
#   just up api ui             api + ui (local)
#   just up dev                everything enabled (proxied to AWS dev)
#   just up api dev            api only (proxied to AWS dev)
#   just up order-api          single service by profile name (local)
#   just up order-api dev      single service (proxied to AWS dev)
#
#   Toggle persistent services in: infra/local/services.jsonc

# Start services — args: any combo of [api] [ui] [support] [dev] or a profile name
up +args="all":
    #!/usr/bin/env bash
    set -euo pipefail

    ARGS=" {{args}} "
    COMPOSE="{{COMPOSE}}"
    PROFILES="--profile infra"

    [[ "$ARGS" == *" dev "* ]] && COMPOSE="{{COMPOSE_PROXY}}"

    # Return enabled --profile flags for a given section of services.jsonc
    section_profiles() {
        sed 's|//.*||g' "{{SERVICES_FILE}}" \
            | jq -r ".${1} | to_entries[] | select(.value == true) | .key" \
            | sed 's/^/--profile /' | tr '\n' ' '
    }

    # "all" or bare "dev" = every enabled section
    if [[ "$ARGS" == *" all "* ]] || \
       [[ "$ARGS" == *" dev "* && "$ARGS" != *" api "* && "$ARGS" != *" ui "* && "$ARGS" != *" support "* ]]; then
        PROFILES="$PROFILES $(section_profiles api) $(section_profiles ui) $(section_profiles support)"
    else
        [[ "$ARGS" == *" api "*     ]] && PROFILES="$PROFILES $(section_profiles api)"
        [[ "$ARGS" == *" ui "*      ]] && PROFILES="$PROFILES $(section_profiles ui)"
        [[ "$ARGS" == *" support "* ]] && PROFILES="$PROFILES $(section_profiles support)"
    fi

    # Any unrecognised arg is treated as a raw profile name
    for arg in {{args}}; do
        case "$arg" in
            all|dev|api|ui|support) ;;
            *) PROFILES="$PROFILES --profile $arg" ;;
        esac
    done

    echo "Starting: $PROFILES"
    $COMPOSE $PROFILES up -d --build

# Stop all running services
down:
    #!/usr/bin/env bash
    {{COMPOSE}} --profile infra --profile api --profile ui --profile support down --remove-orphans

# Stop all services and remove volumes
down-clean:
    #!/usr/bin/env bash
    {{COMPOSE}} --profile infra --profile api --profile ui --profile support down --remove-orphans --volumes

# Follow logs for all running services
logs:
    #!/usr/bin/env bash
    {{COMPOSE}} --profile infra --profile api --profile ui --profile support logs -f

# Show running containers
ps:
    #!/usr/bin/env bash
    {{COMPOSE}} --profile infra --profile api --profile ui --profile support ps
