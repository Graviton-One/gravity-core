#!/bin/bash

while read -r in; 
do 
    echo "$in" | sed -ne '/Deploying 'Nebula'/,$ p' | grep 'contract address:' | head -n1 | awk '{ print $4 }'
done