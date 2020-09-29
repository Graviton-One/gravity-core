#!/bin/bash

export GRAVITY_CMD=/var/www/gravity-core/cmd/gravity

cd $GRAVITY_CMD

./gravity ledger --home="$GRAVITY_HOME" start --rpc="127.0.0.1:2500" --bootstrap=""
