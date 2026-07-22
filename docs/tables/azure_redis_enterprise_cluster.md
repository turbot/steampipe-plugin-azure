---
title: "Steampipe Table: azure_redis_enterprise_cluster - Query Azure Redis Enterprise Clusters using SQL"
description: "Allows users to query Azure Redis Enterprise Clusters, providing details about cluster name, SKU, host name, Redis version, provisioning state, and more."
folder: "Redis"
---

# Table: azure_redis_enterprise_cluster - Query Azure Redis Enterprise Clusters using SQL

Azure Cache for Redis Enterprise is the highest-tier, enterprise-grade Redis offering on Azure. It maps to the `Microsoft.Cache/redisEnterprise` resource type and exposes dedicated cluster metadata such as hostname, SKU, Redis version, and TLS configuration that is not available through the generic `azure_resource` table.

## Table Usage Guide

The `azure_redis_enterprise_cluster` table provides insights into each Azure Redis Enterprise cluster within your Azure environment. As a cloud architect or security engineer, you can use this table to inventory clusters, audit TLS settings, and understand SKU/capacity configurations across subscriptions and resource groups.

## Examples

### Basic info
Retrieve a summary of all Redis Enterprise clusters including their SKU, host name, Redis version, and kind.

```sql+postgres
select
  name,
  kind,
  host_name,
  redis_version,
  sku_name,
  sku_capacity,
  provisioning_state,
  region,
  resource_group
from
  azure_redis_enterprise_cluster;
```

```sql+sqlite
select
  name,
  kind,
  host_name,
  redis_version,
  sku_name,
  sku_capacity,
  provisioning_state,
  region,
  resource_group
from
  azure_redis_enterprise_cluster;
```

### List Azure Managed Redis clusters
Identify Azure Managed Redis (kind `v2`) clusters and inspect their high availability, redundancy, and public network access settings.

```sql+postgres
select
  name,
  kind,
  high_availability,
  redundancy_mode,
  public_network_access,
  region,
  resource_group
from
  azure_redis_enterprise_cluster
where
  kind = 'v2';
```

```sql+sqlite
select
  name,
  kind,
  high_availability,
  redundancy_mode,
  public_network_access,
  region,
  resource_group
from
  azure_redis_enterprise_cluster
where
  kind = 'v2';
```

### List clusters with public network access enabled
Find clusters that are reachable from public networks, which may require additional network controls.

```sql+postgres
select
  name,
  public_network_access,
  region,
  resource_group
from
  azure_redis_enterprise_cluster
where
  public_network_access = 'Enabled';
```

```sql+sqlite
select
  name,
  public_network_access,
  region,
  resource_group
from
  azure_redis_enterprise_cluster
where
  public_network_access = 'Enabled';
```

### List clusters not enforcing TLS 1.2
Identify clusters that may be accepting older TLS versions, which can be a security risk.

```sql+postgres
select
  name,
  region,
  resource_group,
  minimum_tls_version
from
  azure_redis_enterprise_cluster
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
  azure_redis_enterprise_cluster
where
  minimum_tls_version is null
  or minimum_tls_version <> '1.2';
```

### List clusters that are not in a succeeded provisioning state
Find clusters that may be in a degraded or transitional state.

```sql+postgres
select
  name,
  provisioning_state,
  resource_state,
  region,
  resource_group
from
  azure_redis_enterprise_cluster
where
  provisioning_state <> 'Succeeded';
```

```sql+sqlite
select
  name,
  provisioning_state,
  resource_state,
  region,
  resource_group
from
  azure_redis_enterprise_cluster
where
  provisioning_state <> 'Succeeded';
```

### List clusters with private endpoint connections
Discover which clusters have private endpoints configured.

```sql+postgres
select
  name,
  region,
  resource_group,
  jsonb_array_length(private_endpoint_connections) as private_endpoint_count
from
  azure_redis_enterprise_cluster
where
  private_endpoint_connections is not null
  and jsonb_array_length(private_endpoint_connections) > 0;
```

```sql+sqlite
select
  name,
  region,
  resource_group
from
  azure_redis_enterprise_cluster
where
  private_endpoint_connections is not null;
```

### Count clusters per region
Get a distribution of Redis Enterprise clusters across Azure regions.

```sql+postgres
select
  region,
  count(*) as cluster_count
from
  azure_redis_enterprise_cluster
group by
  region
order by
  cluster_count desc;
```

```sql+sqlite
select
  region,
  count(*) as cluster_count
from
  azure_redis_enterprise_cluster
group by
  region
order by
  cluster_count desc;
```
