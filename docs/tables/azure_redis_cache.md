# Table: azure_redis_cache

Azure Cache for Redis provides an in-memory data store based on the Redis software. Redis improves the performance and scalability of an application that uses backend data stores heavily. It's able to process large volumes of application requests by keeping frequently accessed data in the server memory, which can be written to and read from quickly. Redis brings a critical low-latency and high-throughput data storage solution to modern applications.

## Examples

### Basic info

```sql
select
  name,
  redis_version,
  provisioning_state,
  port,
  sku_name,
  region,
  subscription_id
from
  azure_redis_cache;
```

### List cache servers not using latest TLS protocol

```sql
select
  name,
  region,
  resource_group,
  minimum_tls_version
from
  azure_redis_cache
where
  minimum_tls_version is null
  or minimum_tls_version <> '1.2';
```

### List cache servers with in-transit encryption disabled

```sql
select
  name,
  region,
  resource_group,
  enable_non_ssl_port
from
  azure_redis_cache
where
  enable_non_ssl_port;
```

### List premium cache servers

```sql
select
  name,
  region,
  resource_group,
  sku_name
from
  azure_redis_cache
where
  sku_name = 'Premium';
```
