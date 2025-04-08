---
title: "Steampipe Table: azure_redis_cache - Query Azure Redis Cache using SQL"
description: "Allows users to query Azure Redis Cache, specifically details about the name, location, resource group, and subscription of each Redis Cache resource."
folder: "Redis"
---

# Table: azure_redis_cache - Query Azure Redis Cache using SQL

Azure Redis Cache is a fully managed, in-memory cache that enables high-performance and scalable architectures. It uses the popular open-source Redis data structure store, which supports a variety of data structures such as strings, hashes, lists, sets, sorted sets with range queries, bitmaps, and more. It's a part of Azure's suite of database services, providing a reliable and secure Redis cache environment.

## Table Usage Guide

The `azure_redis_cache` table provides insights into each Azure Redis Cache resource within your Azure environment. As a database administrator or developer, you can use this table to gain a comprehensive overview of your Redis Cache resources, including their names, locations, associated resource groups, and subscriptions. This information can be instrumental in optimizing your cache usage, managing resources, and planning capacity.

## Examples

### Basic info
Explore the details of your Azure Redis Cache instances to understand their current status, region, and version. This can help you manage your resources effectively and ensure they are correctly provisioned and operating in the expected regions.

```sql+postgres
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

```sql+sqlite
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
Identify instances where your cache servers are not utilizing the latest TLS protocol. This can be particularly useful for maintaining high security standards and ensuring data protection.

```sql+postgres
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

```sql+sqlite
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
Explore which cache servers are potentially vulnerable by identifying those with in-transit encryption disabled. This is crucial for enhancing your data security by ensuring all cache servers are encrypted.

```sql+postgres
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

```sql+sqlite
select
  name,
  region,
  resource_group,
  enable_non_ssl_port
from
  azure_redis_cache
where
  enable_non_ssl_port = 1;
```

### List premium cache servers
Discover the segments that utilize premium cache servers in Azure, enabling you to understand your resource distribution and manage costs effectively. This is particularly useful when assessing your premium services usage across different regions and resource groups.

```sql+postgres
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

```sql+sqlite
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