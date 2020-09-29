#!/bin/bash


set -e

ETH_NODE_URL=https://ropsten.infura.io/v3/55ce99b713ee4918896e979d172109cf

docker run -itd -e ETH_NODE_URL=$ETH_NODE_URL -p 26657:26657 -p 2500:2500 -v $GRAVITY_HOME:/etc/gravity gravity-ledger:latest

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

