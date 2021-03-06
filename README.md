# YugabyteDB Client API

YugabyteDB Client API provides high level API operations and command line interface for YugabyteDB RPC.

Command line interface is very roughly modelled after the `yb-admin` tool.

## Usage

```
go run ./main.go [command] [flags]
```

## Common flags

- `--master`: string, repeated, host port of the master to query, default `127.0.0.1:7100, 127.0.0.1:7101, 127.0.0.1:7102`
- `--operation-timeout`: RPC operation timeout, duration string (`5s`, `1m`, ...), default `60s`
- `--tls-ca-cert-file-path`: full path to the CA certificate file, default `empty string`
- `--tls-cert-file-path`: full path to the certificate file, default `empty string`
- `--tls-key-file-path`: full path to the key file, default `empty string`

Logging flags:

- `--log-level`: log level, default `info`
- `--log-as-json`: log entries as JSON, default `false`
- `--log-color`: log colored output, default `false`
- `--log-force-color`: force colored output, default `false`

## Commands

- [Universe and cluster commands](https://github.com/radekg/yugabyte-db-go-client#universe-and-cluster-commands)
- [Table commands](https://github.com/radekg/yugabyte-db-go-client#table-commands)
- [Backup and snapshot commands](https://github.com/radekg/yugabyte-db-go-client#backup-and-snapshot-commands)
- [Multi-zone and multi-region deployment commands](https://github.com/radekg/yugabyte-db-go-client#multi-zone-and-multi-region-deployment-commands)
- [Change data capture (CDC) commands](https://github.com/radekg/yugabyte-db-go-client#change-data-capture-cdc-commands)
- [Decommissioning commands](https://github.com/radekg/yugabyte-db-go-client#decommissioning-commands)
- [Rebalancing commands](https://github.com/radekg/yugabyte-db-go-client#rebalancing-commands)
- [Other commands](https://github.com/radekg/yugabyte-db-go-client#other-commands)

## Universe and cluster commands

### get-universe-config

Get the placement info and blacklist info of the universe.

### change-config

TODO: [track](https://github.com/radekg/yugabyte-db-go-client/issues/2).

### change-master-config

TODO: [track](https://github.com/radekg/yugabyte-db-go-client/issues/35).

### list-tablet-servers

List all the tablet servers in this database.

- `--primary-only`: boolean, list primary tablet servers only, default `false`

### list-masters

List all the masters in this database.

### list-replica-type-counts

TODO: [track](https://github.com/radekg/yugabyte-db-go-client/issues/30).

### dump-masters-state

TODO: [track](https://github.com/radekg/yugabyte-db-go-client/issues/36).

### list-tablets-for-table-server

TODO: [track](https://github.com/radekg/yugabyte-db-go-client/issues/37).

### split-tablet

TODO: [track](https://github.com/radekg/yugabyte-db-go-client/issues/52).

### master-leader-stepdown

Try to force the current master leader to step down.

### ysql-catalog-version

Fetch current YSQL catalog version.

## Table commands

### list-tables

List all tables in this database.

- `--name-filter`: string, When used, only returns tables that satisfy a substring match on `name_filter`, default `empty string`
- `--keyspace`: string, the namespace name to fetch info, default `empty string`
- `--exclude-system-tables`: boolean, exclude system tables, default `false`
- `--include-not-running`: boolean, include not running, default `false`
- `--relation-type`: list of strings, filter tables based on RelationType - supported values: `system_table`, `user_table`, `index_table`, default: all values

Examples:

- list all PostgreSQL `system_platform` relations: `cli list-tables --keyspace ysql.system_platform`
- list all PostgreSQL `postgres` relations: `cli list-tables --keyspace ysql.postgres`
- list all PostgreSQL `yugabyte` relations: `cli list-tables --keyspace ysql.yugabyte`
- list all PostgreSQL `template0` relations: `cli list-tables --keyspace ysql.template0`
- list all CQL `system_schema` relations: `cli list-tables --keyspace ycql.system_schema`
- list all Redis `system_redis` relations: `cli list-tables --keyspace yedis.system_redis`

### compact-table

TODO: requires an issue.

### modify-table-placement-info

TODO: requires an issue.

## Backup and snapshot commands

### create-snapshot

Creates a snapshot of an entire keyspace or selected tables in a keyspace.

- `--keyspace`: string, keyspace name to create snapshot of, default `<empty string>`
- `--name`: repeated string, table name to create snapshot of, default `empty list`
- `--uuid`: repeated string, table ID to create snapshot of, default `empty list`
- `--schedule-id`: base64 encoded, create snapshot to this schedule, other fields are ignored, default `empty`
- `--base64-encoded`: boolean, base64 decode given schedule ID before handling over to the API, default `false`

Remarks:

- Multiple `--name` and `--uuid` values can be combined together.
- YSQL keyspace snapshots do not support explicit `--name` and `--uuid` selection.
- To create a snapshot of an entire keyspace, do not specify any `--name` or `--uuid`. YCQL only.
- `yedis.*` keyspaces are not supported.

Examples:

- create a snapshot of an entire YSQL `yugabyte` database: `cli create-snapshot --keyspace ysql.yugabyte`
- create a snapshot of selected YCQL tables in the `example` database: `cli create-snapshot --keyspace ycql.example --name table`

### delete-snapshot

Delete a snapshot.

- `--snapshot-id`: string, snapshot identifier - literal ID or Base64 encoded value from YugabyteDB RPC API, required, default `empty string` (not defined)

### list-snapshots

List snapshots.

- `--snapshot-id`: string, snapshot identifier - literal ID or Base64 encoded value from YugabyteDB RPC API, optional, default `empty string` (not defined)
- `--list-deleted-snapshots`: boolean, list deleted snapshots, default `false`
- `--prepare-for-backup`: boolean, prepare for backup, default `false`

### list-snapshot-restorations

List snapshot restorations.

- `--schedule-id`: string, snapshot identifier - literal ID or Base64 encoded value from YugabyteDB RPC API, optional, default `empty string` (not defined)
- `--restoration-id`: string, restoration identifier - literal ID or Base64 encoded value from YugabyteDB RPC API, optional, default `empty string` (not defined)

### export-snapshot

Export a snapshot.

- `--snapshot-id`: string, snapshot identifier- literal ID or Base64 encoded value from YugabyteDB RPC API, required, default `empty string` (not defined)
- `--file-path`: string, full path to the export file, parent directories must exist, default `empty`

### import-snapshot

Import a snapshot.

- `--file-path`: string, full path to the exported snapshot file
- `--keyspace`: string, fully qualified keyspace name, for example `ycql.system_namespace`, no effect for YSQL imports, default `empty`
- `--table-name`: string, repeated, table name to import, no effect for YSQL snapshots, default `empty list`

### restore-snapshot

Restore a snapshot.

- `--schedule-id`: string, schedule identifier - literal ID or Base64 encoded value from YugabyteDB RPC API, required, default `empty string` (not defined)
- `--restore-target`: exact past HT (`16 digit literal`) or duration expression (`1h`, `5h15m`, ...), absolute Timing Option: Max HybridTime, or relative past interval, default `empty` (undefined)

### create-snapshot-schedule

Create a new snapshot schedule.

- `--keyspace`: string, keyspace name to create snapshot of, default `<empty string>`
- `--interval`: duration expression (`1h`, `1d`, ...), interval for taking snapshot in seconds, default `0` (undefined)
- `--retention-duration`: duration expression (`1h`, `1d`, ...), how long store snapshots in seconds, default `0` (undefined)
- `--delete-after`: exact future HT (`16 digit literal`) or duration expression (`1h`, `5h15m`, ...), how long until schedule is removed in seconds, hybrid time will be calculated by fetching server hybrid time and adding this value, default `0` (undefined)

Examples:

- create a snapshot schedule of an entire YSQL `yugabyte` database: `cli create-snapshot-schedule --keyspace ysql.yugabyte --interval 1h --retention-duration 2h --delete-after 1h`
- create a snapshot schedule of selected YSQL tables in the `yugabyte` database: `cli create-snapshot-schedule --keyspace ysql.yugabyte --name table --name another-table`

### delete-snapshot-schedule

Delete a snapshot schedule.

- `--schedule-id`: string, snapshot identifier - literal ID or Base64 encoded value from YugabyteDB RPC API, required, default `empty string` (not defined)

### list-snapshot-schedules

List snapshot schedules.

- `--schedule-id`: string, snapshot identifier - literal ID or Base64 encoded value from YugabyteDB RPC API, optional, default `empty string` (not defined)

### restore-snapshot-schedule

Restore a snapshot schedule.

- `--snapshot-id`: string, snapshot identifier - literal ID or Base64 encoded value from YugabyteDB RPC API, required, default `empty string` (not defined)
- `--restore-target`: exact past HT (`16 digit literal`) or duration expression (`1h`, `5h15m`, ...), absolute Timing Option: Max HybridTime, or relative past interval, default `empty` (undefined)

## Multi-zone and multi-region deployment commands

### modify-placement-info

Modifies the placement information (cloud, region, and zone) for a deployment.

- `--placement-info`: string, repeated, placement for cloud.region.zone, default cluster value is `cloud1.datacenter1.rack1`, default `empty`, at least one required
- `--replication-factor`: uint32, the number of replicas for each tablet, default `0` (must be explicitly specified)
- `--placement-uuid`: string, the identifier of the primary cluster, which can be any unique string, optional, if not set, a randomly-generated ID will be used, default `not set`

### set-preferred-zones

Sets the preferred availability zones (AZs) and regions.

- `--zone-info`: string, repeated, specifies the cloud, region, and zone, default `empty`

## Change data capture (CDC) commands

### create-cdc-stream

TODO: [track](https://github.com/radekg/yugabyte-db-go-client/issues/44).

### delete-cdc-stream

TODO: [track](https://github.com/radekg/yugabyte-db-go-client/issues/46).

### list-cdc-streams

TODO: [track](https://github.com/radekg/yugabyte-db-go-client/issues/45).

## Decommissioning commands

### get-leader-blacklist-completion

TODO: [track](https://github.com/radekg/yugabyte-db-go-client/issues/41).

### change-blacklist

TODO: [track](https://github.com/radekg/yugabyte-db-go-client/issues/1).

### change-leader-blacklist

TODO: [track](https://github.com/radekg/yugabyte-db-go-client/issues/42).

### leader-stepdown

- `--destination-uuid`: UUID of server this request is addressed to, default `empty` - not specified
- `--disable-graceful-transition`: boolean, if `new-leader-uuid` is not specified, the current leader will attempt to gracefully transfer leadership to another peer; setting this flag disables that behavior, default `false`
- `--new-leader-uuid`: UUID of the server that should run the election to become the new leader, default `empty` - not specified
- `--tablet-id`: the id of the tablet, default `empty` - not specified

## Rebalancing commands

### set-load-balancer-state

Options are mutually exclusive, exactly one has to be set:

- `--enabled`: boolean, default `false`, new desired state: enabled
- `--disabled`: boolean, default `false`, new desired state: disabled

### get-load-balancer-idle

Finds out if the load balancer is idle.

### get-load-move-completion

Get the completion percentage of tablet load move from blacklisted servers.

### is-load-balanced

Check if master leader thinks that the load is balanced across TServers.

- `--expected-num-servers`: int32, how many servers to include in this check, default `-1` (`undefined`)

## Other commands

### check-exists

Check that a table exists.

- `--keyspace`: string, keyspace name to check in, default `<empty string>`
- `--name`: string, table name to check for, default `<empty string>`
- `--uuid`: string, table identified (uuid) to check for, default `<empty string>`

### describe-table

Info on a table in this database.

- `--keyspace`: string, keyspace name to check in, default `<empty string>`, ignored when using `--uuid`
- `--name`: string, table name to check for, default `<empty string>`
- `--uuid`: string, table identified (uuid) to check for, default `<empty string>`

Examples:

- describe table `test` in the `yugabyte` database: `cli describe-table --keyspace yugabyte --name test`
- describe table with ID `000033c0000030008000000000004000`: `cli describe-table --uuid 000033c0000030008000000000004000`

### get-master-registration

Get master registration info.

### get-tablets-for-table

Fetch tablet information for a given table.

- `--keyspace`: string, keyspace to describe the table in, default `empty string`
- `--name`: string, table name to check for, default `empty string`
- `--uuid`: string, table identifier to check for, default `empty string`
- `--partition-key-start`: base64 encoded, partition key range start, default `empty`
- `--partition-key-end`: base64 encoded, partition key range end, default `empty`
- `--max-returned-locations`: uint32, maximum number of returned locations, default `10`
- `--require-tablet-running`: boolean, require tablet running, default `false`

### is-server-ready

Check if server is ready to serve IO requests.

- `--host`: string, host to check, default `<empty string>`
- `--port`: int, port to check, default `0`, must be higher than `0`
- `--is-tserver`: boolean, when `true` - indicated a TServer, default `false`

### ping

Ping a certain YB server.

- `--host`: string, host to ping, default `<empty string>`
- `--port`: int, port to ping, default `0`, must be higher than `0`

## Minimal YugabyteDB cluster in Docker compose

This repository contains a minimal YugabyteDB Docker compose setup which can be used for client testing or validation.

To start the cluster:

```sh
cd .compose/
docker compose -f yugabytedb-minimal.yml up
```

To restart:

```sh
docker compose -f yugabytedb-minimal.yml rm
docker compose -f yugabytedb-minimal.yml up
```

## Docker image

Build the Docker image:

```sh
make docker-image
```

Run against the provided minimal YugabyteDB cluster:

```sh
docker run --rm \
    --net yb-client-minimal \
    -ti local/ybdb-go-cli:0.0.1 \
    list-masters --master yb-master-1:7100 \
                 --master yb-master-2:7100 \
                 --master yb-master-3:7100
```
