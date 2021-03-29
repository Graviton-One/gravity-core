#!/bin/bash 

set -e

override_config() {
	local relevant_path=$1
	local override_path=$2

  temp_config="/etc/gravity/config_tmp.json"

	cat "$relevant_path" | jq ". += $(echo $(cat $override_path))" > "$temp_config"

  cat $temp_config > "$relevant_path"
  rm $temp_config

}

if [ "$INIT_CONFIG" -eq 1 ]
then
  # Config folder is empty, generating keys
  if [ -z "$(ls -A /etc/gravity/)" ]; then
    /bin/gravity ledger --home=/etc/gravity/ init --network="$GRAVITY_NETWORK"
  fi
fi

if [ ! -z $ADAPTERS_CFG_PATH ]
then
	override_config '/etc/gravity/config.json' "$ADAPTERS_CFG_PATH"
fi

if [ ! -z $GENESIS_CFG_PATH ]
then
	override_config '/etc/gravity/genesis.json' "$GENESIS_CFG_PATH"
fi

if [ -n "$GRAVITY_BOOTSTRAP" ]
then
  /bin/gravity ledger --home=/etc/gravity start --rpc="$GRAVITY_PRIVATE_RPC" --bootstrap="$GRAVITY_BOOTSTRAP"
else
  /bin/gravity ledger --home=/etc/gravity start --rpc="$GRAVITY_PRIVATE_RPC" --bootstrap=""
fi
