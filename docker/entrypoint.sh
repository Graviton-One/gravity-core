#!/bin/bash

if [ $INIT_CONFIG -eq 1 ]
then
  gravity ledger --home=/etc/gravity/ init --network="$GRAVITY_NETWORK"
fi

if [ ! -z $ETH_NODE_URL ]
then
  temp_config=/etc/gravity/config_tmp.json
  cat /etc/gravity/config.json | jq ".Adapters.ethereum.NodeUrl = \"$ETH_NODE_URL\"" > "$temp_config"
  cat $temp_config > /etc/gravity/config.json 
  rm $temp_config
fi

gravity ledger --home=/etc/gravity start --rpc="$GRAVITY_RPC" --bootstrap="$GRAVITY_BOOTSTRAP"
