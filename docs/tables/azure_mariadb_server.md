# Table: azure_mariadb_server

Azure Database for MariaDB is a relational database service based on the open-source MariaDB Server engine. It's a fully managed database as a service offering that can handle mission-critical workloads with predictable performance and dynamic scalability.

## Examples

### Basic info

```sql
select
  name,
  version,
  sku_name,
  user_visible_state,
  region,
  resource_group
from
  azure_mariadb_server;
```

### List servers with Geo-redundant backup disabled

```sql
select
  name,
  version,
  region,
  geo_redundant_backup_enabled
from
  azure_mariadb_server
where
  geo_redundant_backup_enabled = 'Disabled';
```

### List servers with SSL enabled

```sql
select
  name,
  version,
  region,
  ssl_enforcement
from
  azure_mariadb_server
where
  ssl_enforcement = 'Enabled';
```

### List servers with backup retention days greater than 90 days

```sql
select
  name,
  version,
  region,
  backup_retention_days
from
  azure_mariadb_server
where
  backup_retention_days > 90;
```
