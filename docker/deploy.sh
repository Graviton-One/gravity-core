#!/bin/bash

start_customnet() {
	GRAVITY_WAVES_CHAINID='S'
	GRAVITY_ETH_ADDRESS='0x605f2226b0451492Cdd72D776EF311926ceE0B92'
	GRAVITY_WAVES_ADDRESS='3MiFxwmcrkujBRsM9FzCxGAL6i1acYah1pJ'
	NETWORK=custom

	# NO ENDING "/"
	WAVES_NODE_URL=https://nodes-stagenet.wavesnodes.com
	ETH_NODE_URL=https://ropsten.infura.io/v3/55ce99b713ee4918896e979d172109cf

	docker run -itd \
	       -e ETH_NODE_URL=$ETH_NODE_URL \
	       -e GRAVITY_NETWORK=$NETWORK \
	       -e GRAVITY_ETH_ADDRESS=$GRAVITY_ETH_ADDRESS \
	       -e GRAVITY_WAVES_ADDRESS=$GRAVITY_WAVES_ADDRESS \
	       -e GRAVITY_WAVES_CHAINID=$GRAVITY_WAVES_CHAINID \
	       -e WAVES_NODE_URL=$WAVES_NODE_URL \
	       -p 26657:26657 -p 2500:2500 -v $GRAVITY_HOME:/etc/gravity gravity-ledger:latest
}

start_devnet() {
	docker run -itd -p 26657:26657 -p 2500:2500 -v $GRAVITY_HOME:/etc/gravity gravity-ledger:latest
}

while [ -n "$1" ]
do
	case "$1" in
                --pure) rm -rf $GRAVITY_HOME ;;
		--custom) start_customnet ;;
		--dev) start_devnet ;;
	esac
	shift
done

exit 0
#--------------------------------------------------------------------------------------

dirs=(/tmp/gravity/main /tmp/gravity/scnd)

index=0
for gravity_dir in dirs
do
  if [ ! -d "$gravity_dir" ]
  then
    mkdir -p "$gravity_dir"
  fi

  if [ $index -eq 0 ]
  then
    # start main ledger
    docker run -itd -e -v $gravity_dir:/etc/gravity gravity-ledger:latest
  else
    # start another ledger
    docker run -itd -e INIT_CONFIG=0 -v $gravity_dir:/etc/gravity gravity-ledger:latest
  fi
  
  index=$((index + 1)) 
done

