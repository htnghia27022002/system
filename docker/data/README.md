# Persistent data

Docker bind mounts store service data here so `docker compose down` does not delete it.

| Path | Service | Container path |
|------|---------|----------------|
| `postgres/` | PostgreSQL | `/var/lib/postgresql/data` |
| `redis/` | Redis | `/data` |

These directories are created on first `make up` and are ignored by git.

To wipe data intentionally, stop the stack and remove the subfolder contents:

```bash
make down
rm -rf docker/data/postgres docker/data/redis
```

Do not delete this folder while containers are running.
