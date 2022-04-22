# Table: azure_mysql_flexible_server

Azure Database for MySQL Flexible Server is a fully managed MySQL database as a service offering that can handle mission-critical workloads with predictable performance and dynamic scalability.

## Examples

### Basic info

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

### List servers with 'storage_auto_grow' disabled

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

### List servers with 'backup_retention_days' greater than 90 days

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
