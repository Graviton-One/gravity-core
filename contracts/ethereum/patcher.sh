#!/bin/bash

eth_address=''

eth_port='8545'
eth_network=''

replace_address_in_migration () {
    local migration_name='2_initial_contracts.js'
    local route_to_file="migrations/$migration_name"
    
    local updated_file=$(cat "$route_to_file" | sed "s/ADDRESS/$eth_address/")

    > "$route_to_file"

    echo "$updated_file" >> "$route_to_file"
}

update_truffle_config () {

    local template='
    module.exports = {
        external: {
        host: "%s",     // Localhost (default: none)
        port: 8545,            // Standard Ethereum port (default: none)
        network_id: "*",       // Any network (default: none)
        gas: 0,
        skipDryRun: false
        }
    }
    '

    printf "$template" $eth_network > "truffle-config.js"
}

main () {
    while [ -n "$1" ]
    do
        case "$1" in
            --eth-address) eth_address=$2; replace_address_in_migration ;;
            --eth-network) eth_network=$2; update_truffle_config ;;
        esac 
        shift
    done
}

main $@