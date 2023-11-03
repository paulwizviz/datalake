# Overview

This is a project to simulate operations to manipulate data from an Ethereum node using a combination of GRPC and Postgres technologies.

## Architecture

<u>Source</u>

The source architecture follows the hexagonal pattern.

The project layout are:

* `build` - This contains build scripts to create docker images
* `cmd` - This contain main packages for client and server executable
* `internal` - This contains shared code for consumption by packages under cmd
* `protos` - This contains protobufs and grpc service specification
* `scripts` - This contains bash scripts to trigger operations to generate codes from protos and start/stop network
* `testdata` - This contains a series of pre-loaded raw protobufs simulating cache after it is download from S3 containers

To generate code from protos run the command: `./scripts/protoc.sh`. This will generate Go codes to the location `./internal/block`.

Testdata folder is attached to `grpc_server` container and mapped to `/var/blocks`. The latter simulate a cache following protobufs downloaded from [S3 container](https://s3.us-east-1.amazonaws.com/public.blocks.datalake). The testdata is preloaded by running the command `go run cmd/s3reader/main.go`.  The command pulls from S3 and store protobuf in `*.pb` files.

<u>Deployment</u>

For local testing, we use docker compose network. Please refer to [docker-compose file](./deployment/docker-compose.yml). The docker compose network runs three containers:

* GRPC server serving raw protobuf from S3 containers held in a file-base cache.
* Postgres server
* PGAdmin portal to enable use to view postgres DB

There is a fourth container acting a client. You activate this container via the command `./scripts/ops.sh client`. This opens up a shell for you to execute client operations.

## Running the network

Pre-requiste:

* Install docker
* Install Go version 1.21
* Install GRPC tooling

**NOTE**: This project has only been verified in to run on macOS. It should run on Linux. There is no support for Windows.

To run the application, following the following steps:

1. Run the script `./scripts/ops.sh image build` (macOS) or `sudo ./scripts/ops.sh image build` (Linux)
1. Start the network `./scripts/ops.sh network start` (macOS) and add sudo for Linux
1. After the network has started the client by running this `./scripts/ops.sh client`

When you start the client shell execute the client command as follows:

* Fetch by block number
```bash
$ client number {block number}
```

* Fetch by hash
```bash
$ client hash {block hash}
```
