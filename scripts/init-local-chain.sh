#!/usr/bin/env bash
set -euo pipefail

CHAIN_ID="ferac-1"
DENOM="uferac"
NODE_HOME="${HOME}/.ferac"
PROJECT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"
KEYRING_BACKEND="test"
FERACD="${FERACD_BIN:-$PROJECT_DIR/feracd}"

if [[ ! -x "$FERACD" ]]; then
  echo "feracd binary not found: $FERACD" >&2
  exit 1
fi

for command_name in python3; do
  if ! command -v "$command_name" >/dev/null 2>&1; then
    echo "$command_name is required" >&2
    exit 1
  fi
done

if [[ "${1:-}" == "--reset" ]]; then
  pkill -x feracd 2>/dev/null || true
  for _ in {1..50}; do
    pgrep -x feracd >/dev/null || break
    sleep 0.1
  done
  pkill -9 -x feracd 2>/dev/null || true
  rm -rf "$NODE_HOME"
elif [[ -e "$NODE_HOME/config/genesis.json" ]]; then
  echo "${NODE_HOME} already exists. Use --reset to recreate the chain." >&2
  exit 1
fi

"$FERACD" init mynode \
  --chain-id "$CHAIN_ID" \
  --default-denom "$DENOM" \
  --home "$NODE_HOME"

add_key() {
  local name="$1"
  local mnemonic="$2"
  local mnemonic_file

  mnemonic_file="$(mktemp)"
  printf '%s' "$mnemonic" > "$mnemonic_file"

  "$FERACD" keys add "$name" \
    --recover \
    --source "$mnemonic_file" \
    --home "$NODE_HOME" \
    --keyring-backend "$KEYRING_BACKEND"

  rm -f "$mnemonic_file"
}

add_key creator "absurd amount doctor acoustic avoid letter advice cage absurd amount doctor acoustic avoid letter advice cage absurd amount doctor acoustic avoid letter advice comic"
add_key team "acoustic avoid letter advice cage absurd amount doctor acoustic avoid letter advice cage absurd amount doctor acoustic avoid letter advice cage absurd amount exchange"
add_key validator_reserve "adapt blossom school alcohol coral light army gather adapt blossom school alcohol coral light army gather adapt blossom school alcohol coral light army hold"
add_key community "advice cage absurd amount doctor acoustic avoid letter advice cage absurd amount doctor acoustic avoid letter advice cage absurd amount doctor acoustic avoid negative"
add_key dex_liquidity "agree choice donor anxiety expect little beef pass agree choice donor anxiety expect little beef pass agree choice donor anxiety expect little beef sample"
add_key treasury "alcohol coral light army gather adapt blossom school alcohol coral light army gather adapt blossom school alcohol coral light army gather adapt blossom tackle"
add_key ecosystem "alpha deal scrub asthma idea logic bright thought alpha deal scrub asthma idea logic bright thought alpha deal scrub asthma idea logic bright truly"
add_key val1 "amount doctor acoustic avoid letter advice cage absurd amount doctor acoustic avoid letter advice cage absurd amount doctor acoustic avoid letter advice cage again"
add_key val2 "animal embark drastic bamboo mountain loyal category cancel animal embark drastic bamboo mountain loyal category cancel animal embark drastic bamboo mountain loyal category cover"
add_key val3 "anxiety expect little beef pass agree choice donor anxiety expect little beef pass agree choice donor anxiety expect little beef pass agree choice fire"
add_key val4 "arch flame security bid radar machine club gesture arch flame security bid radar machine club gesture arch flame security bid radar machine club lawn"

add_genesis_account() {
  local name="$1"
  local amount="$2"

  "$FERACD" genesis add-genesis-account "$name" "${amount}${DENOM}" \
    --home "$NODE_HOME" \
    --keyring-backend "$KEYRING_BACKEND"
}

add_genesis_account creator 8888888800000
add_genesis_account team 2666666640000
add_genesis_account validator_reserve 3555555520000
add_genesis_account community 23111110800000
add_genesis_account dex_liquidity 17777777600000
add_genesis_account treasury 13333333200000
add_genesis_account ecosystem 11555555440000
add_genesis_account val1 2000000000000
add_genesis_account val2 2000000000000
add_genesis_account val3 2000000000000
add_genesis_account val4 2000000000000

python3 - "$PROJECT_DIR/config.yml" "$NODE_HOME/config/genesis.json" \
  > /tmp/ferac-genesis.json <<'PY'
import json
import sys

import yaml


def deep_merge(left, right):
    merged = dict(left)
    for key, value in right.items():
        if isinstance(merged.get(key), dict) and isinstance(value, dict):
            merged[key] = deep_merge(merged[key], value)
        else:
            merged[key] = value
    return merged


config_path, genesis_path = sys.argv[1:]
with open(config_path, encoding="utf-8") as config_file:
    config = yaml.safe_load(config_file)
with open(genesis_path, encoding="utf-8") as genesis_file:
    genesis = json.load(genesis_file)

genesis["app_state"] = deep_merge(genesis["app_state"], config["genesis"]["app_state"])
genesis["consensus"] = deep_merge(genesis["consensus"], config["genesis"]["consensus"])
json.dump(genesis, sys.stdout, indent=2)
print()
PY
mv /tmp/ferac-genesis.json "$NODE_HOME/config/genesis.json"

mkdir -p "$NODE_HOME/config/gentx"

VALIDATOR_HOMES="$(mktemp -d)"
trap 'rm -rf "$VALIDATOR_HOMES"' EXIT

for validator in val1 val2 val3 val4; do
  validator_home="$VALIDATOR_HOMES/$validator"
  "$FERACD" init "$validator" \
    --chain-id "$CHAIN_ID" \
    --default-denom "$DENOM" \
    --home "$validator_home" \
    >/dev/null
  validator_pubkey="$($FERACD tendermint show-validator --home "$validator_home")"

  "$FERACD" genesis gentx "$validator" 1000000000000${DENOM} \
    --home "$NODE_HOME" \
    --chain-id "$CHAIN_ID" \
    --moniker "$validator" \
    --pubkey "$validator_pubkey" \
    --output-document "$NODE_HOME/config/gentx/gentx-${validator}.json" \
    --commission-rate 0.05 \
    --commission-max-rate 0.05 \
    --commission-max-change-rate 0.00 \
    --min-self-delegation 1000000000 \
    --keyring-backend "$KEYRING_BACKEND" \
    --yes \
    >/dev/null

  if [[ "$validator" == "val1" ]]; then
    cp "$validator_home/config/priv_validator_key.json" "$NODE_HOME/config/priv_validator_key.json"
  fi
done

"$FERACD" genesis collect-gentxs --home "$NODE_HOME"

sed -i 's/^minimum-gas-prices = .*/minimum-gas-prices = "0uferac"/' \
  "$NODE_HOME/config/app.toml"
sed -i 's/^enabled-unsafe-cors = false/enabled-unsafe-cors = true/' \
  "$NODE_HOME/config/app.toml"
sed -i '/^\[api\]/,/^\[/ {
  s/^enable = false$/enable = true/
  s|^address = "tcp://localhost:1317"$|address = "tcp://0.0.0.0:1317"|
}' "$NODE_HOME/config/app.toml"
sed -i 's/^timeout_commit = .*/timeout_commit = "8s"/' \
  "$NODE_HOME/config/config.toml"
sed -i 's/^timeout_propose = .*/timeout_propose = "3s"/' \
  "$NODE_HOME/config/config.toml"

MULTI_NODE_ROOT="${HOME}/.ferac-nodes"
rm -rf "$MULTI_NODE_ROOT"
mkdir -p "$MULTI_NODE_ROOT"

declare -A NODE_HOMES=(
  [val1]="$NODE_HOME"
  [val2]="$MULTI_NODE_ROOT/val2"
  [val3]="$MULTI_NODE_ROOT/val3"
  [val4]="$MULTI_NODE_ROOT/val4"
)
declare -A RPC_PORTS=(
  [val1]=26657
  [val2]=26658
  [val3]=26659
  [val4]=26660
)
declare -A P2P_PORTS=(
  [val1]=26756
  [val2]=26757
  [val3]=26758
  [val4]=26759
)

for validator in val1 val2 val3 val4; do
  node_home="${NODE_HOMES[$validator]}"

  if [[ "$validator" != "val1" ]]; then
    "$FERACD" init "$validator" \
      --chain-id "$CHAIN_ID" \
      --default-denom "$DENOM" \
      --home "$node_home" \
      >/dev/null
    cp "$NODE_HOME/config/genesis.json" "$node_home/config/genesis.json"
    cp "$VALIDATOR_HOMES/$validator/config/priv_validator_key.json" \
      "$node_home/config/priv_validator_key.json"
    cp "$VALIDATOR_HOMES/$validator/config/node_key.json" \
      "$node_home/config/node_key.json"
    cp "$NODE_HOME/config/app.toml" "$node_home/config/app.toml"
    cp "$NODE_HOME/config/config.toml" "$node_home/config/config.toml"
    sed -i '/^\[api\]/,/^\[/ s/^enable = true$/enable = false/' \
      "$node_home/config/app.toml"
    sed -i '/^\[grpc\]/,/^\[/ s/^enable = true$/enable = false/' \
      "$node_home/config/app.toml"
  fi

  rpc_port="${RPC_PORTS[$validator]}"
  p2p_port="${P2P_PORTS[$validator]}"
  sed -i "s|^laddr = \"tcp://127.0.0.1:[0-9]*\"|laddr = \"tcp://127.0.0.1:${rpc_port}\"|" \
    "$node_home/config/config.toml"
  sed -i "s|^laddr = \"tcp://0.0.0.0:[0-9]*\"|laddr = \"tcp://0.0.0.0:${p2p_port}\"|" \
    "$node_home/config/config.toml"
  sed -i 's/^addr_book_strict = true/addr_book_strict = false/' \
    "$node_home/config/config.toml"
  sed -i 's/^allow_duplicate_ip = false/allow_duplicate_ip = true/' \
    "$node_home/config/config.toml"
  sed -i 's/^pex = true/pex = false/' "$node_home/config/config.toml"
done

ALL_NODE_IDS=""
for validator in val1 val2 val3 val4; do
  node_id="$($FERACD comet show-node-id --home "${NODE_HOMES[$validator]}")"
  if [[ -n "$ALL_NODE_IDS" ]]; then
    ALL_NODE_IDS+=","
  fi
  ALL_NODE_IDS+="$node_id"
done

for current_validator in val1 val2 val3 val4; do
  persistent_peers=""
  for validator in val1 val2 val3 val4; do
    if [[ "$validator" == "$current_validator" ]]; then
      continue
    fi
    if [[ "$current_validator" != "val1" && "$validator" != "val1" ]]; then
      continue
    fi
    node_home="${NODE_HOMES[$validator]}"
    node_id="$($FERACD comet show-node-id --home "$node_home")"
    p2p_port="${P2P_PORTS[$validator]}"
    if [[ -n "$persistent_peers" ]]; then
      persistent_peers+=","
    fi
    persistent_peers+="${node_id}@127.0.0.1:${p2p_port}"
  done

  node_home="${NODE_HOMES[$current_validator]}"
  sed -i "s|^persistent_peers = .*|persistent_peers = \"${persistent_peers}\"|" \
    "$node_home/config/config.toml"
  sed -i "s|^unconditional_peer_ids = .*|unconditional_peer_ids = \"${ALL_NODE_IDS}\"|" \
    "$node_home/config/config.toml"
done

for validator in val1 val2 val3 val4; do
  node_home="${NODE_HOMES[$validator]}"
  nohup "$FERACD" start --home "$node_home" \
    > "$PROJECT_DIR/ferac-${validator}.log" 2>&1 &
done

echo "Four-node chain started. Logs: $PROJECT_DIR/ferac-val*.log"
echo "REST: http://127.0.0.1:1317"
