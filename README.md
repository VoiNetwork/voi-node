# Voi Node

This project contains a simple Docker image and associated golang tooling to
run a Voi node.

## Node types currently supported

- Relay node
- Archiver node
- Developer node
- Participation node
- Conduit node

Node type will default to 'relay' if not specified.
Node configuration will use testnet defaults unless otherwise provided at image build time.

## Running a Voi node

### Running with default settings

```bash
docker run ghcr.io/voinetwork/voi-node:latest
```

This will run a relay node on the testnet network.

### Running with a configuration value override

```bash
docker run -e VOINETWORK_INCOMING_CONNECTIONS_LIMIT=30 ghcr.io/voinetwork/voi-node:latest 
```

### Running with default mainnet

```bash
docker run -e VOINETWORK_NETWORK=testnet ghcr.io/voinetwork/voi-node:latest
```

### Running with a pre-defined network

To run a Voi node with a pre-defined network, you can use the following command:

```bash
docker run -e VOINETWORK_NETWORK=testnet ghcr.io/voinetwork/voi-node:latest
```

### Running with a custom network

```bash
docker run -e VOINETWORK_NETWORK=betanet -e VOINETWORK_GENESIS=custom_url ghcr.io/voinetwork/voi-node:latest
```

### Running with a specific profile

#### Relay node

```bash
docker run -e VOINETWORK_PROFILE=relay -p 5011:8080 ghcr.io/voinetwork/voi-node:latest
```

This maps local port 5011 to the blockchain service running on port 8080 within the container.

#### Archiver node

```bash
docker run -e VOINETWORK_PROFILE=archiver -v <my_local_path>:/algod/data ghcr.io/voinetwork/voi-node:latest
```

#### Developer node

```bash
docker run -e VOINETWORK_PROFILE=relay -p 5011:8080 ghcr.io/voinetwork/voi-node:latest
```

This maps local port 5011 to the blockchain service running on port 8080 within the container.

#### Participation node

```bash
docker run -e VOINETWORK_PROFILE=participation -p 5011:8080 ghcr.io/voinetwork/voi-node:latest
```

This maps local port 5011 to the blockchain service running on port 8080 within the container.

#### Conduit node

```bash
docker run -e VOINETWORK_PROFILE=conduit -p 5011:8080 -v <my_local_path>:/algod/data ghcr.io/voinetwork/voi-node:latest
```

This maps local port 5011 to the blockchain service running on port 8080 within the container.
This mode is used if you want to connect an indexer and the conduit framework to the network, 
in order to query network such, such as looking up transaction history.

Guidance on how to install Conduit and the Indexer can be found on the [Conduit GitHub page](https://github.com/algorand/conduit),
as well as the [Indexer GitHub page](https://github.com/algorand/indexer)
The API token required can be found in the folder that you have mounted via the `-v` flag, after
startup. 

The API token is persistent and generated randomly on first start. If you need to ensure it 
stays constant pre-create the folder and create a `algod.token` and `algod.admin.token` file with the wanted sha256 tokens inside.

Example `algod.token`:
```bash
c3818ac81a91c6e58df2b151a43caf7245e806ec02e4851bca01a1be2ba72da7
```