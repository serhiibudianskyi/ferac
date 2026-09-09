#!/usr/bin/env bash
set -euo pipefail

PROJECT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"
VUE_DIR="$PROJECT_DIR/vue"
PREVIEW_PORT="${VUE_PREVIEW_PORT:-4173}"
LOG_FILE="$VUE_DIR/preview.log"

if [[ ! -f "$VUE_DIR/package.json" ]]; then
  echo "Vue project not found: $VUE_DIR" >&2
  exit 1
fi

cd "$VUE_DIR"
pnpm build

mapfile -t preview_pids < <(pgrep -f '/vite/bin/vite.js preview' || true)
if ((${#preview_pids[@]} > 0)); then
  kill "${preview_pids[@]}"
  for _ in {1..50}; do
    still_running=false
    for pid in "${preview_pids[@]}"; do
      if kill -0 "$pid" 2>/dev/null; then
        still_running=true
        break
      fi
    done
    [[ "$still_running" == false ]] && break
    sleep 0.1
  done
fi

nohup pnpm preview --host --port "$PREVIEW_PORT" > "$LOG_FILE" 2>&1 &
preview_pid=$!

for _ in {1..50}; do
  if curl -fsS -o /dev/null "http://127.0.0.1:${PREVIEW_PORT}/ferac/"; then
    echo "Vue rebuilt and preview started."
    echo "PID: $preview_pid"
    echo "URL: https://budianskyi-s.com/ferac/"
    exit 0
  fi
  sleep 0.2
done

echo "Preview did not become ready. Check $LOG_FILE" >&2
exit 1
