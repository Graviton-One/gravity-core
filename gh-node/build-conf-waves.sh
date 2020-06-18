#!/bin/bash


main () {
    local config_name='config-waves.json'
    local native_url='http://127.0.0.1:26657'
    local node_url='http://0.0.0.0:6869'

    while [ -n "$1" ]
    do
        case "$1" in
            --config) config_name="$2.json" ;;
            --node-url) node_url=$2 ;;
            --native-url) native_url=$2 ;;
        esac
        shift
    done

    JSON_FMT='{
    "GHNodeURL":"%s",
    "Timeout": 0,
    "NebulaId":"0000000000000000000000000000000000000000000000000000000000000001",
    "NodeUrl":"%s",
    "NebulaContract":"3MHQtdw8myT6b42RD4a9zc346MxgeKeKTkE",
    "ChainType": "Waves"\n}\n'
    printf "$JSON_FMT" "$native_url" "$node_url" > "$config_name"
}

main $@
