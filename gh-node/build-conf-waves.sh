#!/bin/bash


main () {
    local config_name='config-waves.json'

    while [ -n "$1" ]
    do
        case "$1" in
            --config) config_name="$2.json" ;;
        esac
        shift
    done

    JSON_FMT='{
    "GHNodeURL":"http://localhost:26657",
    "Timeout": 0,
    "NebulaId":"0000000000000000000000000000000000000000000000000000000000000001",
    "NodeUrl":"http://127.0.0.1:6869",
    "NebulaContract":"3MHQtdw8myT6b42RD4a9zc346MxgeKeKTkE",
    "ChainType": "Waves"\n}\n'
    printf "$JSON_FMT" > "$config_name"
}

main $@