---
title: "Steampipe Table: azure_postgresql_flexible_server - Query Azure Database for PostgreSQL Flexible Servers using SQL"
description: "Allows users to query Azure Database for PostgreSQL Flexible Servers."
---

# Table: azure_postgresql_flexible_server - Query Azure Database for PostgreSQL Flexible Servers using SQL

Azure Database for PostgreSQL Flexible Server is a fully managed database service designed for developers. It offers built-in high availability with no additional cost. It also provides the capability to scale compute resources up or down based on your application's need.

## Table Usage Guide

The 'azure_postgresql_flexible_server' table provides insights into PostgreSQL Flexible Servers within Azure Database for PostgreSQL. As a DevOps engineer, explore server-specific details through this table, including server version, state, and associated metadata. Utilize it to uncover information about servers, such as those with public network access, the backup retention period, and the geo-redundant backup setting. The schema presents a range of attributes of the PostgreSQL Flexible Server for your analysis, like the server name, creation date, SKU name, and associated tags.

## Examples

### Basic info
Discover the details of your Azure PostgreSQL Flexible Server configurations, including their names, IDs, and cloud environments. This can be particularly useful for understanding the geographic distribution of your servers and assessing their various configurations.

```sql
select
  name,
  id,
  cloud_environment,
  flexible_server_configurations,
  location
from
  azure_postgresql_flexible_server;
```

### List SKU details of the flexible servers
Explore the specific details of your flexible servers, such as their SKU name and tier, to better understand and manage your resources within the Azure PostgreSQL environment. This can be particularly useful for resource allocation, cost management, and strategic planning.

```sql
select
  name,
  id,
  sku ->> 'name' as sku_name,
  sku ->> 'tier' as sku_tier
from
  azure_postgresql_flexible_server;
```

### List flexible servers that have geo-redundant backup enabled
Explore which flexible servers have geo-redundant backup enabled to ensure data security and continuity in case of a regional outage. This query is useful in identifying servers that have additional data protection measures in place.

```sql
select
  name,
  id,
  cloud_environment,
  flexible_server_configurations,
  server_properties -> 'backup' ->> 'geoRedundantBackup',
  location
from
  azure_postgresql_flexible_server
where
  server_properties -> 'backup' ->> 'geoRedundantBackup' = 'Enabled';
```

### List flexible servers configured in more than one availability zones
Explore which flexible servers are configured across multiple availability zones in Azure. This is particularly useful for ensuring high availability and disaster recovery, as it allows you to identify any servers that might be at risk due to being confined to a single zone.

```sql
select
  name,
  id,
  cloud_environment,
  flexible_server_configurations,
  server_properties ->> 'availabilityZone',
  location
from
  azure_postgresql_flexible_server
where
  (server_properties ->> 'availabilityZone')::int > 1;
```