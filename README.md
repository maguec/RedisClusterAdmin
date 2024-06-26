# Redis Cluster Admin

Run commands against all shards in a Redis Cluster.

## Building

```sh
make
```

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

The verbose option will include the shard information as a comment line before the output

## Examples

### Get Keyspace on all shards 

```sh
./rcadmin -s localhost -p 30001  INFO KEYSPACE | egrep "^db0:"
db0:keys=331,expires=0,avg_ttl=0
db0:keys=335,expires=0,avg_ttl=0
db0:keys=334,expires=0,avg_ttl=0
```
### Find all matching keys across a cluster

*WARNING* Running this against production may have severe performance implications - be careful

```sh
./rcadmin -s localhost -p 30001  -v KEYS "MYKEY5*"
# 127.0.0.1:30001
[MYKEY516 MYKEY523 MYKEY505 MYKEY538 MYKEY553 MYKEY527 MYKEY593 MYKEY545 MYKEY584 MYKEY562 MYKEY592 MYKEY530 MYKEY579 MYKEY580 MYKEY512 MYKEY552 MYKEY588 MYKEY53 MYKEY571 MYKEY556 MYKEY574 MYKEY597 MYKEY509 MYKEY581 MYKEY541 MYKEY578 MYKEY585 MYKEY575 MYKEY549 MYKEY534 MYKEY566 MYKEY567 MYKEY57 MYKEY570 MYKEY563 MYKEY596 MYKEY589 MYKEY501]
# 127.0.0.1:30002
[MYKEY514 MYKEY561 MYKEY558 MYKEY540 MYKEY531 MYKEY518 MYKEY535 MYKEY572 MYKEY504 MYKEY544 MYKEY510 MYKEY517 MYKEY598 MYKEY56 MYKEY587 MYKEY548 MYKEY550 MYKEY594 MYKEY543 MYKEY500 MYKEY507 MYKEY583 MYKEY513 MYKEY526 MYKEY590 MYKEY557 MYKEY522 MYKEY565 MYKEY554 MYKEY547 MYKEY576 MYKEY508 MYKEY52 MYKEY539 MYKEY569]
# 127.0.0.1:30003
[MYKEY536 MYKEY599 MYKEY560 MYKEY595 MYKEY51 MYKEY533 MYKEY502 MYKEY577 MYKEY529 MYKEY555 MYKEY532 MYKEY528 MYKEY564 MYKEY524 MYKEY586 MYKEY503 MYKEY58 MYKEY519 MYKEY559 MYKEY542 MYKEY50 MYKEY59 MYKEY551 MYKEY573 MYKEY54 MYKEY537 MYKEY546 MYKEY506 MYKEY5 MYKEY525 MYKEY568 MYKEY55 MYKEY511 MYKEY515 MYKEY520 MYKEY591 MYKEY521 MYKEY582]
```
