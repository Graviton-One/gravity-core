#!/bin/bash

update_config_field() {
  local key=$1
  local value=$2 

  temp_config=/etc/gravity/config_tmp.json
  cat /etc/gravity/config.json | jq "$key = \"$value\"" > "$temp_config"
  # cat /etc/gravity/config.json | jq ".Adapters.ethereum.NodeUrl = \"$ETH_NODE_URL\"" > "$temp_config"
  cat $temp_config > /etc/gravity/config.json 
  rm $temp_config

}

if [ $INIT_CONFIG -eq 1 ]
then
  gravity ledger --home=/etc/gravity/ init --network="$GRAVITY_NETWORK"
fi

env_keys=('.Adapters.ethereum.NodeUrl' '.Adapters.ethereum.GravityContractAddress' '.Adapters.waves.GravityContractAddress' '.Adapters.waves.ChainId')
env_values=("$ETH_NODE_URL" "$GRAVITY_ETH_ADDRESS" "$GRAVITY_WAVES_ADDRESS" "$GRAVITY_WAVES_CHAINID")

for (( i=0; i < 4; i++ ))
do
  env_key="${env_keys[i]}"
  value="${env_values[i]}"

  echo "KEY: $env_key; VALUE: $value"

  if [ ! -z $value ]
  then
    update_config_field "$env_key" "$value"
  fi

done

gravity ledger --home=/etc/gravity start --rpc="$GRAVITY_RPC" --bootstrap="$GRAVITY_BOOTSTRAP"
