#!/usr/bin/env bash
set -euo pipefail

PROJECT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"

cd "$PROJECT_DIR"

echo "Building feracd..."
go build -o "$PROJECT_DIR/feracd" ./cmd/feracd

echo "Creating a clean chain..."
"$PROJECT_DIR/scripts/init-local-chain.sh" --reset

echo "Rebuilding Vue and restarting preview..."
"$PROJECT_DIR/scripts/rebuild-vue.sh"

for _ in {1..100}; do
  if curl -fsS -o /dev/null http://127.0.0.1:26657/status; then
    echo "All services restarted."
    echo "App: https://budianskyi-s.com/ferac/"
    echo "Logs: $PROJECT_DIR/ferac-val*.log"
    exit 0
  fi
  sleep 0.2
done

echo "Chain RPC did not become ready. Check $PROJECT_DIR/ferac-val*.log" >&2
exit 1
