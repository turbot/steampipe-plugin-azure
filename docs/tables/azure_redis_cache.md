---
title: "Steampipe Table: azure_redis_cache - Query Azure Cache for Redis instances using SQL"
description: "Allows users to query Azure Cache for Redis instances."
---

# Table: azure_redis_cache - Query Azure Cache for Redis instances using SQL

Azure Cache for Redis is an in-memory data store that is used to power fast, scalable applications. It provides secure and dedicated Redis server instances and additional features like Azure Virtual Network, full Redis command-set support, and premium tier features like clustering, persistence, and virtual network support.

## Table Usage Guide

The 'azure_redis_cache' table provides insights into Azure Cache for Redis instances. As a DevOps engineer, explore instance-specific details through this table, including configuration, access keys, and associated metadata. Utilize it to uncover information about instances, such as configuration settings, the number of clients connected, and the memory usage. The schema presents a range of attributes of the Redis instances for your analysis, like the Redis version, creation date, SKU name, and associated tags.

## Examples

### Basic info
Explore the configuration and status of your Azure Redis Cache instances. This is useful for understanding the versions in use, their locations, and the state of provisioning to ensure optimal performance and resource allocation.

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
Explore which cache servers in your Azure Redis Cache are not utilizing the latest TLS protocol. This helps ensure optimal security by identifying areas where updates may be needed.

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
Explore which cache servers in Azure have in-transit encryption disabled. This is useful to identify potential security risks and ensure that all your data is securely transmitted.

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
Explore which cache servers are of premium type in your Azure Redis Cache setup. This can help in managing resources and costs more effectively.

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