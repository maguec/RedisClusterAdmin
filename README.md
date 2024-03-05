# Redis Cluster Admin

Run commands against all shards in a Redis Cluster.

## Usage

```sh
./rcadmin -s <host> -p <port> -v INFO KEYSPACE
```

## Options

```sh
./rcadmin -h
Usage: rcadmin [--server SERVER] [--port PORT] [--verbose] [COMMAND [COMMAND ...]]

Positional arguments:
  COMMAND                Command

Options:
  --server SERVER, -s SERVER
                         Cluster Server Host [default: localhost, env: CLUSTER_SERVER]
  --port PORT, -p PORT   Cluster Server Port [default: 6379, env: CLUSTER_PORT]
  --verbose, -v          Verbose
  --help, -h             display this help and exit
```
