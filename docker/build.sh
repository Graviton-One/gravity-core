#!/bin/bash

cd ..
docker build -t gravity-ledger -f docker/ledger.dockerfile .
