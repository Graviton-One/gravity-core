#!/bin/bash

build_oracle() {
  cd ..
  docker build -t gravity-oracle -f docker/oracle.dockerfile .
}

build_ledger() {
  cd ..
  docker build -t gravity-ledger -f docker/ledger.dockerfile .
}

main() {

  while [ -n "$1" ]
  do
    case "$1" in
      --ledger) build_ledger ;;
      --oracle) build_oracle ;;
      *) echo "Enter --oracle or --ledger in order to run operation"; exit 1; ;; 
    esac

    exit 0;
  done
}

main $@
