#!/bin/bash

start_mainnet() {
	NETWORK=mainnet

	docker run -itd \
				 -e INIT_CONFIG=1 \
	       -p 26657:26657 -p 2500:2500 -v "$GRAVITY_HOME":/etc/gravity gravityhuborg/gravity-ledger:master
}

start_devnet() {
	docker run -itd -p 26657:26657 -p 2500:2500 -v "$GRAVITY_HOME":/etc/gravity gravityhuborg/gravity-ledger:master
}

while [ -n "$1" ]
do
	case "$1" in
    --pure) rm -rf "$GRAVITY_HOME" ;;
		# --custom) start_customnet ;;
		--mainnet) start_mainnet ;;
		# --dev) start_devnet ;;
	esac
	shift
done
