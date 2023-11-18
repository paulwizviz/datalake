#!/bin/bash

export GRPC_SERVER_IMAGE=datalake/server:current
export GRPC_CLIENT_IMAGE=datalake/client:current
export NETWORK=datalake_network

COMMAND="$1"
SUBCOMMAND="$2"

function image(){
    local cmd="$1"
    case $cmd in
        "build")
            docker-compose -f ./build/builder.yml build
            ;;
        "clean")
            docker rmi -f ${GRPC_SERVER_IMAGE}
            docker rmi -f ${GRPC_CLIENT_IMAGE}
            docker rmi -f $(docker images --filter "dangling=true" -q)
            ;;
        *)
            echo "image [ build | clean]"
            ;;
    esac
}

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

function client(){
    docker run -it --rm \
           -e GRPC_SERVER="grpcserver" \
           -e DB_URL="postgres://postgres:postgres@defaultserver:5432/postgres" \
           --network=$NETWORK \
           $GRPC_CLIENT_IMAGE /bin/bash
}

case $COMMAND in
    "image")
        image $SUBCOMMAND
        ;;
    "network")
        network $SUBCOMMAND
        ;;
    "client")
        client
        ;;
    *)
        echo "$0 <command>
commands:
    image     build or clean
    network   clean, start and stop
"
        ;;
esac