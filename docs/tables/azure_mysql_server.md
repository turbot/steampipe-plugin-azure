# Table: azure_mysql_server

Azure Database for MySQL Server is a fully managed database service designed to provide more granular control and flexibility over database management functions and configuration settings.

## Examples

### List servers for which Enforce SSL connection is enabled

```sql
select
  name,
  id,
  location,
  ssl_enforcement
from
  azure_mysql_server
where
  ssl_enforcement = 'Enabled';
```

### List servers for which public network access is disabled

```sql
select
  name,
  id,
  public_network_access
from
  azure_mysql_server
where
  public_network_access = 'Disabled';
```

### List servers for which storage profile auto grow feature is disabled

```sql
select
  name,
  id,
  storage_profile_storage_auto_grow
from
  azure_mysql_server
where
  storage_profile_storage_auto_grow = 'Disabled';
```

### List servers for which 'backup_retention_days' is greater than 3 days

```sql
select
  name,
  id,
  backup_retention_days
from
  azure_mysql_server
where
  backup_retention_days > 3;
```

### List servers with minimum TLS version is lower than 1.2

```sql
select
  name,
  id,
  minimal_tls_version
from
  azure_mysql_server
where
  minimal_tls_version = 'TLS1_0'
  or minimal_tls_version = 'TLS1_1';
```