---
title: "Steampipe Table: azure_postgresql_flexible_server - Query Azure PostgreSQL Flexible Servers using SQL"
description: "Allows users to query Azure PostgreSQL Flexible Servers, providing insights into the configuration, status, and capabilities of these managed database instances."
---

# Table: azure_postgresql_flexible_server - Query Azure PostgreSQL Flexible Servers using SQL

Azure PostgreSQL Flexible Server is a fully managed relational database service, based on the open-source Postgres database engine. It provides capabilities for intelligent performance, high availability, and dynamic scalability, enabling you to focus on application development and business logic rather than database management tasks. This service helps you to securely manage, monitor, and scale your PostgreSQL databases in the cloud.

## Table Usage Guide

The `azure_postgresql_flexible_server` table provides insights into the configuration and status of Azure PostgreSQL Flexible Server instances. As a database administrator, you can leverage this table to explore server-specific details, including the server's state, version, location, and more. Utilize it to monitor and manage your PostgreSQL databases in Azure, ensuring optimal performance, security, and compliance.

## Examples

### Basic info
Uncover the details of your Azure PostgreSQL flexible servers including their names, IDs, and configurations. This information is essential for managing your cloud environment effectively and understanding where your servers are located.

```sql+postgres
select
  name,
  id,
  cloud_environment,
  flexible_server_configurations,
  location
from
  azure_postgresql_flexible_server;
```

```sql+sqlite
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
Identify the specific details of flexible servers, such as their unique identifiers and SKU details. This information can be beneficial in managing resources and understanding the tier level of each server.

```sql+postgres
select
  name,
  id,
  sku ->> 'name' as sku_name,
  sku ->> 'tier' as sku_tier
from
  azure_postgresql_flexible_server;
```

```sql+sqlite
select
  name,
  id,
  json_extract(sku, '$.name') as sku_name,
  json_extract(sku, '$.tier') as sku_tier
from
  azure_postgresql_flexible_server;
```

### List flexible servers that have geo-redundant backup enabled
Identify instances where flexible servers have geo-redundant backup enabled to ensure data redundancy and disaster recovery for your Azure PostgreSQL databases.

```sql+postgres
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

```sql+sqlite
select
  name,
  id,
  cloud_environment,
  flexible_server_configurations,
  json_extract(json_extract(server_properties, '$.backup'), '$.geoRedundantBackup') as geoRedundantBackup,
  location
from
  azure_postgresql_flexible_server
where
  json_extract(json_extract(server_properties, '$.backup'), '$.geoRedundantBackup') = 'Enabled';
```

### List flexible servers configured in more than one availability zones
Determine the areas in which flexible servers are configured across multiple availability zones. This is useful for understanding the distribution and redundancy of your servers, which can impact service availability and disaster recovery strategies.

```sql+postgres
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

```sql+sqlite
select
  name,
  id,
  cloud_environment,
  flexible_server_configurations,
  json_extract(server_properties, '$.availabilityZone'),
  location
from
  azure_postgresql_flexible_server
where
  CAST(json_extract(server_properties, '$.availabilityZone') AS INTEGER) > 1;
```