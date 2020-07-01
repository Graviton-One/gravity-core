#!/bin/bash

should_install=1
should_install_waves_deps=1

install_deps () {
    # Deps
    curl -sL https://deb.nodesource.com/setup_14.x | bash -
    apt-get install -y nodejs && curl -L https://npmjs.org/install.sh | sh
    apt-get update && \
    apt-get -y install gcc mono-mcs && \
    rm -rf /var/lib/apt/lists/*
    npm i -g --unsafe-perm=true --allow-root truffle
} 

install_deps_waves_gh_node () {
    curl -sL https://deb.nodesource.com/setup_14.x | bash -
    apt-get install -y nodejs && curl -L https://npmjs.org/install.sh | sh
    apt-get update && \
   	 apt-get -y install gcc mono-mcs && \
    	 rm -rf /var/lib/apt/lists/*
    npm i -g --unsafe-perm=true --allow-root @waves/surfboard
    echo "Building Surfboard config..."
    cd ./contracts/waves && bash patch-surfboard.sh --api $NODE_URL --out surfboard.config.json 
    cat surfboard.config.json
    surfboard test deploy.js
}


main () {
    while [ -n "$1" ]
    do
	case "$1" in
        -i) should_install=$2 ;;
        -wi) should_install_waves_deps=$2 ;;
        esac
        shift
    done

    if [ $should_install -eq 0 ]; then
    	install_deps
    fi 

    if [ $should_install_waves_deps -eq 0 ]; then
    	install_deps_waves_gh_node
    fi 

}

main $@
