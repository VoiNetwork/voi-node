# Voi Node

This project contains a simple Docker image and associated golang tooling to
run a Voi node.

## Node types currently supported
- Relay node
- Archiver node

Node type will default to 'relay' if not specified.
Node configuration will use testnet defaults unless otherwise provided at image build time.

## Running with a pre-defined network
To run a Voi node with a pre-defined network, you can use the following command:

```bash
docker run -e VOINETWORK_NETWORK=testnet ghcr.io/voi-network/voi-node:latest
```

## Running with a custom network

```bash
docker run -e VOINETWORK_NETWORK=custom_name -e VOINETWORK_GENESIS=custom_url ghcr.io/voi-network/voi-node:latest
```

## Running with a specific profile

### Relay node
```bash
docker run -e VOINETWORK_NETWORK=custom_name -e VOINETWORK_GENESIS=custom_url -e VOINETWORK_PROFILE=relay ghcr.io/voi-network/voi-node:latest
```

### Archiver node
```bash
docker run -e VOINETWORK_NETWORK=custom_name -e VOINETWORK_GENESIS=custom_url -e VOINETWORK_PROFILE=archiver ghcr.io/voi-network/voi-node:latest
```