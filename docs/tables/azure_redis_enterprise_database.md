---
title: "Steampipe Table: azure_redis_enterprise_database - Query Azure Redis Enterprise Databases using SQL"
description: "Allows users to query Azure Redis Enterprise Databases, providing details about port, client protocol, clustering policy, eviction policy, persistence, and geo-replication settings."
folder: "Redis"
---

# Table: azure_redis_enterprise_database - Query Azure Redis Enterprise Databases using SQL

Azure Cache for Redis Enterprise databases (`Microsoft.Cache/redisEnterprise/databases`) represent the individual databases hosted within a Redis Enterprise cluster. Each cluster can host one or more databases with configurable settings such as port, eviction policy, clustering policy, and geo-replication — none of which are exposed through the generic `azure_resource` table.

## Table Usage Guide

The `azure_redis_enterprise_database` table provides granular insights into each database within your Redis Enterprise clusters. As a database administrator or security engineer, you can use this table to audit eviction policies, verify clustering policies, check port configurations, and inspect geo-replication setup across all databases.

## Examples

### Basic info
List all Redis Enterprise databases with their key configuration attributes.

```sql+postgres
select
  cluster_name,
  name,
  redis_version,
  port,
  client_protocol,
  clustering_policy,
  eviction_policy,
  provisioning_state,
  region,
  resource_group
from
  azure_redis_enterprise_database;
```

```sql+sqlite
select
  cluster_name,
  name,
  redis_version,
  port,
  client_protocol,
  clustering_policy,
  eviction_policy,
  provisioning_state,
  region,
  resource_group
from
  azure_redis_enterprise_database;
```

### List databases by Redis version
Find databases running a specific Redis version — useful to identify Azure Managed Redis v2 clusters (version `7.4`).

```sql+postgres
select
  cluster_name,
  name,
  redis_version,
  access_keys_authentication,
  region,
  resource_group
from
  azure_redis_enterprise_database
where
  redis_version = '7.4';
```

```sql+sqlite
select
  cluster_name,
  name,
  redis_version,
  access_keys_authentication,
  region,
  resource_group
from
  azure_redis_enterprise_database
where
  redis_version = '7.4';
```

### List databases with access keys authentication enabled
Find databases where access key authentication is active — these may need additional audit controls.

```sql+postgres
select
  cluster_name,
  name,
  access_keys_authentication,
  region,
  resource_group
from
  azure_redis_enterprise_database
where
  access_keys_authentication = 'Enabled';
```

```sql+sqlite
select
  cluster_name,
  name,
  access_keys_authentication,
  region,
  resource_group
from
  azure_redis_enterprise_database
where
  access_keys_authentication = 'Enabled';
```

### List databases allowing plaintext (non-TLS) connections
Identify databases not enforcing TLS for client connections.

```sql+postgres
select
  cluster_name,
  name,
  client_protocol,
  region,
  resource_group
from
  azure_redis_enterprise_database
where
  client_protocol = 'Plaintext';
```

```sql+sqlite
select
  cluster_name,
  name,
  client_protocol,
  region,
  resource_group
from
  azure_redis_enterprise_database
where
  client_protocol = 'Plaintext';
```

### List databases with no eviction policy (NoEviction)
Find databases configured with NoEviction, which may cause write failures under memory pressure.

```sql+postgres
select
  cluster_name,
  name,
  eviction_policy,
  region,
  resource_group
from
  azure_redis_enterprise_database
where
  eviction_policy = 'NoEviction';
```

```sql+sqlite
select
  cluster_name,
  name,
  eviction_policy,
  region,
  resource_group
from
  azure_redis_enterprise_database
where
  eviction_policy = 'NoEviction';
```

### List databases with geo-replication configured
Discover databases that have active geo-replication links.

```sql+postgres
select
  cluster_name,
  name,
  geo_replication -> 'groupNickname' as geo_group,
  jsonb_array_length(geo_replication -> 'linkedDatabases') as linked_db_count,
  region,
  resource_group
from
  azure_redis_enterprise_database
where
  geo_replication is not null;
```

```sql+sqlite
select
  cluster_name,
  name,
  region,
  resource_group
from
  azure_redis_enterprise_database
where
  geo_replication is not null;
```

### List databases with persistence enabled
Check which databases have AOF or RDB persistence configured.

```sql+postgres
select
  cluster_name,
  name,
  persistence ->> 'aofEnabled' as aof_enabled,
  persistence ->> 'rdbEnabled' as rdb_enabled,
  region,
  resource_group
from
  azure_redis_enterprise_database
where
  (persistence ->> 'aofEnabled')::boolean = true
  or (persistence ->> 'rdbEnabled')::boolean = true;
```

```sql+sqlite
select
  cluster_name,
  name,
  json_extract(persistence, '$.aofEnabled') as aof_enabled,
  json_extract(persistence, '$.rdbEnabled') as rdb_enabled,
  region,
  resource_group
from
  azure_redis_enterprise_database
where
  json_extract(persistence, '$.aofEnabled') = 1
  or json_extract(persistence, '$.rdbEnabled') = 1;
```

### Get database details for a specific cluster
Retrieve all databases within a specific Redis Enterprise cluster.

```sql+postgres
select
  cluster_name,
  name,
  port,
  clustering_policy,
  eviction_policy,
  provisioning_state
from
  azure_redis_enterprise_database
where
  cluster_name = 'my-enterprise-cluster';
```

```sql+sqlite
select
  cluster_name,
  name,
  port,
  clustering_policy,
  eviction_policy,
  provisioning_state
from
  azure_redis_enterprise_database
where
  cluster_name = 'my-enterprise-cluster';
```
