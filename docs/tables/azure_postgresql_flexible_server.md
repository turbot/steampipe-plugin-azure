# Table: azure_postgresql_flexible_server

Azure Database for PostgreSQL - Flexible Server is a fully managed database service designed to provide more granular control and flexibility over database management functions and configuration settings. The service generally provides more flexibility and server configuration customizations based on user requirements. The flexible server architecture allows users to collocate the database engine with the client tier for lower latency and choose high availability within a single availability zone and across multiple availability zones. Flexible servers also provide better cost optimization controls with the ability to stop/start your server and a burstable compute tier ideal for workloads that don't need full compute capacity continuously.

## Examples

### Basic info

```sql
select
  name,
  id,
  cloud_environment,
  flexible_server_configurations,
  created_at,
  location
from
  azure_postgresql_flexible_server;
```

### List the flexible servers created in the last 7 days

```sql
select
  name,
  id,
  cloud_environment,
  flexible_server_configurations,
  created_at,
  location
from
  azure_postgresql_flexible_server
where
  created_at >= now() - interval '7' day;
```

### List the flexible servers that have never been modified

```sql
select
  name,
  id,
  cloud_environment,
  flexible_server_configurations,
  created_at,
  location
from
  azure_postgresql_flexible_server
where
  last_modified_at is null;
```

### List SKU details of the flexible servers

```sql
select
  name,
  id,
  sku ->> 'name' as sku_name,
  sku ->> 'tier' as sku_tier
from
  azure_postgresql_flexible_server
where
  last_modified_at is null;
```

### List flexible servers that have geo-redundant backup enabled

```sql
select
  name,
  id,
  cloud_environment,
  flexible_server_configurations,
  created_at,
  server_properties -> 'backup' ->> 'geoRedundantBackup',
  location
from
  azure_postgresql_flexible_server
where
  server_properties -> 'backup' ->> 'geoRedundantBackup' = 'Enabled';
```

### List flexible servers configured in more than one availability zones

```sql
select
  name,
  id,
  cloud_environment,
  flexible_server_configurations,
  server_properties ->> 'availabilityZone',
  created_at,
  location
from
  azure_postgresql_flexible_server
where
  server_properties ->> 'availabilityZone' > 1;
```