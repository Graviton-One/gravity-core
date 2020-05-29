#!/bin/bash

eth_address=''

replace_address_in_migration () {
    local migration_name='2_initial_contracts.js'
    local route_to_file="migrations/$migration_name"
    
    local updated_file=$(cat "$route_to_file" | sed "s/ADDRESS/$eth_address/")

    > "$route_to_file"

    echo "$updated_file" >> "$route_to_file"
}

main () {

    while [ -n "$1" ]
    do
        case "$1" in
            --eth-address) eth_address=$2 ;;
        esac 
        shift
    done

    replace_address_in_migration
}

main $@