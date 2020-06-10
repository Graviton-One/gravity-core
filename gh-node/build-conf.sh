#!/bin/bash


main () {
    local config_name='config.json'
    local nebula_address=''

    while [ -n "$1" ]
    do
        case "$1" in
            --config) config_name="$2.json" ;;
            --nebula) nebula_address=$2 ;;
        esac
        shift
    done

    JSON_FMT='{
    "GHNodeURL":"http://localhost:26657",
    "Timeout":0,
    "NebulaId":"0000000000000000000000000000000000000000000000000000000000000000",
    "NodeUrl":"http://localhost:8545",
    "NebulaContract":"%s"
    "ChainType": "Ethereum"\n}\n'
    printf "$JSON_FMT" "$nebula_address" > "$config_name"
}

main $@