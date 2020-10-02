#!/bin/bash


init_oracle() {
  gravity oracle --home="$GRAVITY_HOME" init \
        $GRAVITY_NEBULA_ADDRESS \
        $CHAIN_TYPE \ 
        $GRAVITY_PUBLIC_LEDGER_RPC \
        $GRAVITY_TARGET_CHAIN_NODE_URL \
        $GRAVITY_EXTRACTOR_ENDPOINT
}

start_oracle() {
  gravity oracle --home="$GRAVITY_HOME" \
        start \
        $GRAVITY_NEBULA_ADDRESS
}

validate_args() {
  required_arg_keys=('Gravity Home' 'Nebula address' 'Chain type' 'Gravity Public Ledger RPC' 'Gravity target chain node url' 'Gravity extractor endpoint')
  required_args=($GRAVITY_HOME $NEBULA_ADDRESS $CHAIN_TYPE $GRAVITY_HOME $GRAVITY_PUBLIC_LEDGER_RPC $GRAVITY_TARGET_CHAIN_NODE_URL $GRAVITY_EXTRACTOR_ENDPOINT)
  args_len="${#required_args[@]}"


	for ((i = 0; i < args_len; i++))
	do
	  arg_key="${required_arg_keys[i]}"
	  arg_value="${required_args[i]}"

	  if [ -z $arg_value ]
	  then
	    echo "$arg_key is not set. Quitting..."
	    exit 1
	  fi
	done
}


run() {
  init_oracle

  start_oracle
}

while [ -n "$1" ]
do
  case "$1" in
    --validate) validate_args ;;
    --run) run ;;
  esac
  shift
done
