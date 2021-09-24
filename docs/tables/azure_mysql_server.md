# Table: azure_mysql_server

Azure Database for MySQL Server is a fully managed database service designed to provide more granular control and flexibility over database management functions and configuration settings.

## Examples

### Basic info

```sql
select
  name,
  id,
  location,
  ssl_enforcement,
  minimal_tls_version
from
  azure_mysql_server;
```

### List servers with SSL enabled

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

### List servers with public network access disabled

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

### List servers with storage profile auto growth disabled

```sql
select
  name,
  id,
  storage_auto_grow
from
  azure_mysql_server
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
  azure_mysql_server
where
  backup_retention_days > 90;
```

### List servers with minimum TLS version lower than 1.2

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

### List private endpoint connection details

```sql
select
  name as server_name,
  id as server_id,
  connections ->> 'id' as connection_id,
  jsonb_pretty(connections -> 'properties') as connection_property,
  connections ->> 'provisioningState' as connection_provisioning_state
from
  azure_mysql_server,
  jsonb_array_elements(private_endpoint_connections) as connections;
```
