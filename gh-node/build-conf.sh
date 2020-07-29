#!/bin/bash


main () {
    local config_name='config.json'
    local nebula_address=''
    local native_url='http://127.0.0.1:26657'
    local node_url='http://0.0.0.0:8545';
    local tcp_priv_key='';

    while [ -n "$1" ]
    do
        case "$1" in
            --config) config_name="$2.json" ;;
            --native-url) native_url=$2 ;;
            --node-url) node_url=$2 ;;
            --nebula) nebula_address=$2 ;;
            --tcp-priv-key) tcp_priv_key=$2 ;;
        esac
        shift
    done

    JSON_FMT='{
    "TCPrivateKey": "%s",   
    "GHNodeURL":"%s",
    "Timeout":0,
    "NebulaId":"0000000000000000000000000000000000000000000000000000000000000000",
    "NodeUrl":"%s",
    "NebulaContract":"%s",
    "ChainType": "Ethereum"\n}\n'
    printf "$JSON_FMT" "$tcp_priv_key" "$native_url" "$node_url" "$nebula_address" > "$config_name"
}

main $@
