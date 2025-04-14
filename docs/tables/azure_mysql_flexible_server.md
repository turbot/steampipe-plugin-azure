---
title: "Steampipe Table: azure_mysql_flexible_server - Query Azure MySQL Flexible Servers using SQL"
description: "Allows users to query Azure MySQL Flexible Servers, providing detailed information on server configurations, geographical location, and other server-related aspects."
folder: "MySQL"
---

# Table: azure_mysql_flexible_server - Query Azure MySQL Flexible Servers using SQL

Azure MySQL Flexible Server is a fully managed database service with built-in high availability and flexible scaling. It allows you to run your MySQL server workloads on Azure and includes features like stop/start, burstable compute, and adjustable storage. This service makes it easy to build cloud-native applications or modernize existing applications using a managed platform.

## Table Usage Guide

The `azure_mysql_flexible_server` table provides insights into Azure MySQL Flexible Servers within Azure Database for MySQL. As a database administrator, you can explore server-specific details through this table, including server configurations, geographical location, and other server-related aspects. Utilize it to uncover information about servers, such as their current state, performance tier, and the associated resource group.

## Examples

### Basic info
Explore the key details of your Azure MySQL flexible servers such as location, backup retention days, storage IOPS, and public network access. This can help in understanding the configuration and performance of your servers.

```sql+postgres
select
  name,
  id,
  location,
  backup_retention_days,
  storage_iops,
  public_network_access
from
  azure_mysql_flexible_server;
```

```sql+sqlite
select
  name,
  id,
  location,
  backup_retention_days,
  storage_iops,
  public_network_access
from
  azure_mysql_flexible_server;
```

### List servers with public network access disabled
Determine the areas in which servers have public network access turned off, enabling you to assess potential security risks and ensure compliance with your organization's policies.

```sql+postgres
select
  name,
  id,
  public_network_access
from
  azure_mysql_flexible_server
where
  public_network_access = 'Disabled';
```

```sql+sqlite
select
  name,
  id,
  public_network_access
from
  azure_mysql_flexible_server
where
  public_network_access = 'Disabled';
```

### List servers with storage auto grow disabled
Determine the areas in which servers have the automatic storage growth feature disabled. This can be useful to identify potential risks of running out of storage space unexpectedly.

```sql+postgres
select
  name,
  id,
  storage_auto_grow
from
  azure_mysql_flexible_server
where
  storage_auto_grow = 'Disabled';
```

```sql+sqlite
select
  name,
  id,
  storage_auto_grow
from
  azure_mysql_flexible_server
where
  storage_auto_grow = 'Disabled';
```

### List servers with backup retention days greater than 90 days
Determine the areas in which server backup retention exceeds a 90-day period, which could assist in identifying potential resource optimization and cost-saving opportunities.

```sql+postgres
select
  name,
  id,
  backup_retention_days
from
  azure_mysql_flexible_server
where
  backup_retention_days > 90;
```

```sql+sqlite
select
  name,
  id,
  backup_retention_days
from
  azure_mysql_flexible_server
where
  backup_retention_days > 90;
```

### List server configuration details
Assess the elements within your Azure MySQL flexible server by understanding the specific server configurations in use. This allows you to identify potential areas for optimization and ensure your server is set up according to your organization's requirements.
**Note:** `Flexible Server configurations` is the same as `Server parameters` as shown in Azure MySQL Flexible Server console

```sql+postgres
select
  name as server_name,
  id as server_id,
  configurations ->> 'Name' as configuration_name,
  configurations -> 'ConfigurationProperties' ->> 'value' as value
from
  azure_mysql_flexible_server,
  jsonb_array_elements(flexible_server_configurations) as configurations;
```

```sql+sqlite
select
  name as server_name,
  s.id as server_id,
  json_extract(configurations.value, '$.Name') as configuration_name,
  json_extract(configurations.value, '$.ConfigurationProperties.value') as value
from
  azure_mysql_flexible_server as s,
  json_each(flexible_server_configurations) as configurations;
```

### Current state of audit_log_enabled parameter for the servers
This query is used to assess the status of the audit log feature on your Azure MySQL flexible servers. It helps in maintaining security and compliance by identifying servers where this feature is not enabled.

```sql+postgres
select
  name as server_name,
  id as server_id,
  configurations ->> 'Name' as configuration_name,
  configurations -> 'ConfigurationProperties' ->> 'value' as value
from
  azure_mysql_flexible_server,
  jsonb_array_elements(flexible_server_configurations) as configurations
where
   configurations ->> 'Name' = 'audit_log_enabled';
```

```sql+sqlite
select
  name as server_name,
  s.id as server_id,
  json_extract(configurations.value, '$.Name') as configuration_name,
  json_extract(json_extract(configurations.value, '$.ConfigurationProperties'), '$.value') as value
from
  azure_mysql_flexible_server as s,
  json_each(flexible_server_configurations) as configurations
where
  json_extract(configurations.value, '$.Name') = 'audit_log_enabled';
```

### List servers with slow_query_log parameter enabled
Explore which servers have the slow_query_log parameter enabled, allowing you to identify potential performance issues and optimize your database operations. This is particularly useful for monitoring and improving the efficiency of your Azure MySQL flexible servers.

```sql+postgres
select
  name as server_name,
  id as server_id,
  configurations ->> 'Name' as configuration_name,
  configurations -> 'ConfigurationProperties' ->> 'value' as value
from
  azure_mysql_flexible_server,
  jsonb_array_elements(flexible_server_configurations) as configurations
where
  configurations ->'ConfigurationProperties' ->> 'value' = 'ON'
  and configurations ->> 'Name' = 'slow_query_log';
```

```sql+sqlite
select
  name as server_name,
  s.id as server_id,
  json_extract(configurations.value, '$.Name') as configuration_name,
  json_extract(json_extract(configurations.value, '$.ConfigurationProperties'), '$.value') as value
from
  azure_mysql_flexible_server as s,
  json_each(flexible_server_configurations) as configurations
where
  json_extract(json_extract(configurations.value, '$.ConfigurationProperties'), '$.value') = 'ON'
  and json_extract(configurations.value, '$.Name') = 'slow_query_log';
```

### List servers with log_output parameter set to file
Determine the areas in which servers have their log output parameter set to a file. This is useful for identifying servers that are configured to log activity directly to a file, which could be a requirement for certain security or auditing purposes.

```sql+postgres
select
  name as server_name,
  id as server_id,
  configurations ->> 'Name' as configuration_name,
  configurations -> 'ConfigurationProperties' ->> 'value' as value
from
  azure_mysql_flexible_server,
  jsonb_array_elements(flexible_server_configurations) as configurations
where
  configurations ->'ConfigurationProperties' ->> 'value' = 'FILE'
  and configurations ->> 'Name' = 'log_output';
```

```sql+sqlite
select
  name as server_name,
  s.id as server_id,
  json_extract(configurations.value, '$.Name') as configuration_name,
  json_extract(json_extract(configurations.value, '$.ConfigurationProperties'), '$.value') as value
from
  azure_mysql_flexible_server as s,
  json_each(flexible_server_configurations) as configurations
where
  json_extract(json_extract(configurations.value, '$.ConfigurationProperties'), '$.value') = 'FILE'
  and json_extract(configurations.value, '$.Name') = 'log_output';
```