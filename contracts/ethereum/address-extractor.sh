#!/bin/bash

entered=0
while read -r in; 
do
    nebula_line=$(echo "$in" | sed -n -e "s/^.*\(Deploying 'Nebula'\)/\1/p")

    if [ -n "$nebula_line" ]
    then
        entered=1
    fi

    if [ $entered -eq 1 ]
    then
        nebula_address=$(echo "$in" | grep 'contract address:' | head -n1 | awk '{ print $4 }')

        if [ -n "$nebula_address" ]
        then
            echo "$nebula_address"
            break
        fi
    fi
done