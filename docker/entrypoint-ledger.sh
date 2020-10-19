#!/bin/bash

update_config_field() {
  local key=$1
  local value=$2 

  temp_config=/etc/gravity/config_tmp.json
  # shellcheck disable=SC2002
  cat /etc/gravity/config.json | jq "$key = \"$value\"" > "$temp_config"
  # cat /etc/gravity/config.json | jq ".Adapters.ethereum.NodeUrl = \"$ETH_NODE_URL\"" > "$temp_config"
  cat $temp_config > /etc/gravity/config.json 
  rm $temp_config
}

if [ "$INIT_CONFIG" -eq 1 ]
then
  # Config folder is empty, generating keys
  if [ -z "$(ls -A /etc/gravity/)" ]; then
    gravity ledger --home=/etc/gravity/ init --network="$GRAVITY_NETWORK"
  fi
fi


env_keys=('.Adapters.ethereum.NodeUrl' '.Adapters.ethereum.GravityContractAddress' '.Adapters.waves.GravityContractAddress' '.Adapters.waves.ChainId' '.Adapters.waves.NodeUrl')
env_values=("$ETH_NODE_URL" "$GRAVITY_ETH_ADDRESS" "$GRAVITY_WAVES_ADDRESS" "$GRAVITY_WAVES_CHAINID" "$WAVES_NODE_URL")
list_len=${#env_keys[@]}


for (( i=0; i < list_len; i++ ))
do
  env_key="${env_keys[i]}"
  value="${env_values[i]}"

  echo "KEY: $env_key; VALUE: $value"

  if [ -n "$value" ]
  then
    update_config_field "$env_key" "$value"
  fi

done

if [ -n "$GRAVITY_BOOTSTRAP" ]
then
  gravity ledger --home=/etc/gravity start --rpc="$GRAVITY_RPC" --bootstrap="$GRAVITY_BOOTSTRAP"
else
  gravity ledger --home=/etc/gravity start --rpc="$GRAVITY_RPC"
fi
