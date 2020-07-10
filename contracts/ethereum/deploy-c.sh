#!/bin/bash


# cd ./contracts/ethereum && \
    # bash patcher.sh --eth-network $ETH_NETWORK --eth-address $ETH_ADDRESS && \
    # cat truffle-config.js && sleep 1 && \
    # truffle migrate --network external | tee migration.txt

# RUN echo "Migration file: \n" && cat ./contracts/ethereum/migration.txt

# RUN cd ./contracts/ethereum && cat migration.txt | bash address-extractor.sh >> nebula-address.txt

main () {
  nebula_addr=$1

  if [ $nebula_addr = '0' ]
  then
     bash patcher.sh --eth-network $ETH_NETWORK --eth-address $ETH_ADDRESS && \
     cat truffle-config.js && sleep 1 && \
     truffle migrate --network external | tee migration.txt

     cat migration.txt | bash address-extractor.sh >> nebula-address.txt
  else
     echo "$nebula_addr" > nebula-address.txt     
  fi
}

main $@
