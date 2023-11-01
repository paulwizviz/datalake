#!/bin/bash

export NETWORK=datalake_network

COMMAND="$1"
SUBCOMMAND="$2"

function network(){
    local cmd=$1
    case $cmd in
        "clean")
            docker-compose -f ./deployment/docker-compose.yml down
            rm -rf ./tmp
            ;;
        "start")
            docker-compose -f ./deployment/docker-compose.yml up
            ;;
        "stop")
            docker-compose -f ./deployment/docker-compose.yml down
            ;;
        *)
            echo "network [ clean | start | stop ]"
            ;;
    esac
}

case $COMMAND in
    "network")
        network $SUBCOMMAND
        ;;
    *)
        echo "$0 <command>
commands:
    network   clean, start and stop
"
        ;;
esac