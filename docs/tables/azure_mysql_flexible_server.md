---
title: "Steampipe Table: azure_mysql_flexible_server - Query Azure MySQL Flexible Servers using SQL"
description: "Allows users to query Azure MySQL Flexible Servers"
---

# Table: azure_mysql_flexible_server - Query Azure MySQL Flexible Servers using SQL

Azure MySQL Flexible Server is a fully managed database service with built-in high availability and the flexibility to adjust compute and storage resources on demand. It supports the diverse needs of your workloads requiring MySQL and allows you to choose the right compute and storage resources for your server. Azure MySQL Flexible Server also provides cost-effectiveness with stop/start capabilities and burstable compute tier.

## Table Usage Guide

The 'azure_mysql_flexible_server' table provides insights into MySQL Flexible Servers within Azure. As a DevOps engineer, explore server-specific details through this table, including server state, version, storage capacity, and associated metadata. Utilize it to uncover information about servers, such as those with high storage capacity, the administrator login name, and the verification of SSL enforcement. The schema presents a range of attributes of the MySQL Flexible Server for your analysis, like the server name, creation date, SKU name, and associated tags.

## Examples

### Basic info
Explore the settings of your Azure MySQL flexible servers to understand their locations, backup retention periods, storage IOPS, and public network access status. This helps in managing resources efficiently and ensuring optimal server configuration.

```sql
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
Explore which servers have disabled public network access to ensure a higher level of security and prevent unauthorized access. This can be beneficial in maintaining data privacy and safeguarding sensitive information.

```sql
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
Identify servers where the automatic storage growth feature is turned off. This is useful for understanding which servers might run out of storage unexpectedly, potentially disrupting operations.

```sql
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
Explore which servers in your Azure MySQL Flexible Server have a backup retention period exceeding 90 days. This is beneficial in understanding your organization's data retention practices and ensuring compliance with internal or regulatory data backup policies.

```sql
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
Explore the configuration details of your servers on Azure's MySQL Flexible Server. This can be useful to understand and manage the properties of your servers, such as identifying any unusual settings that may impact your server's performance.
**Note:** `Flexible Server configurations` is the same as `Server parameters` as shown in Azure MySQL Flexible Server console


```sql
select
  name as server_name,
  id as server_id,
  configurations ->> 'Name' as configuration_name,
  configurations -> 'ConfigurationProperties' ->> 'value' as value
from
  azure_mysql_flexible_server,
  jsonb_array_elements(flexible_server_configurations) as configurations;
```

### Current state of audit_log_enabled parameter for the servers
Analyze the settings to understand the status of the audit log enablement feature across your Azure MySQL flexible servers. This can help ensure that audit logs are active for security and compliance monitoring.

```sql
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

### List servers with slow_query_log parameter enabled
Discover the segments that have the 'slow_query_log' parameter enabled on Azure MySQL Flexible servers. This can be useful for identifying servers that may be experiencing performance issues due to slow queries.

```sql
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

### List servers with log_output parameter set to file
Discover the segments that have the 'log_output' parameter set to 'FILE' within Azure's MySQL Flexible Server. This is particularly useful when you need to identify servers that are logging output to files for auditing or troubleshooting purposes.

```sql
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