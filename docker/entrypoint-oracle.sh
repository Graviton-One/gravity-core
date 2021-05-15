#!/bin/bash


init_oracle() {
  /bin/gravity oracle --home=/etc/gravity init \
        $NEBULA_ADDRESS \
        $CHAIN_TYPE \ 
        $GRAVITY_PUBLIC_LEDGER_RPC \
        $GRAVITY_TARGET_CHAIN_NODE_URL \
        $GRAVITY_EXTRACTOR_ENDPOINT
}

start_oracle() {
  /bin/gravity oracle --home=/etc/gravity \
        start \
        $NEBULA_ADDRESS
}

if [ "$INIT_CONFIG" -eq 1 ]
then
  if [ -z "$(ls -A /etc/gravity/)" ]; then
		init_oracle
  fi
fi

start_oracle
