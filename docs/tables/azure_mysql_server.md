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

### List server keys

```sql
select
  name as server_name,
  id as server_id,
  keys ->> 'creationDate' as keys_creation_date,
  keys ->> 'id' as keys_id,
  keys ->> 'kind' as keys_kind,
  keys ->> 'name' as keys_name,
  keys ->> 'serverKeyType' as keys_server_key_type,
  keys ->> 'type' as keys_type,
  keys ->> 'uri' as keys_uri
from
  azure_mysql_server,
  jsonb_array_elements(server_keys) as keys;
```
