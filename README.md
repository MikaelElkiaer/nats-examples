# NATS Examples

## Leaf node reply across accounts

This works with nats protocol and direct interaction with the server.

It does not work when replying from a leaf node.

### Using nats

```bash
docker compose run requester-nats
docker compose down --volumes
```

### Using leaf node

```bash
docker compose run requester-leafnode
docker compose down --volumes
```
