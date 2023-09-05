# Leaf node reply using allow responses

This works with nats protocol and direct interaction with the server.

It does not work when replying from a leaf node.

## Using nats

```bash
export COMPOSE_FILE=docker-compose.yaml:docker-compose.nats.yaml
export ALLOW_RESPONSES=true
docker compose run --rm requester
docker compose down
```

## Using leaf node

```bash
export COMPOSE_FILE=docker-compose.yaml:docker-compose.leafnode.yaml
export ALLOW_RESPONSES=true
docker compose run --rm requester
docker compose down
```
