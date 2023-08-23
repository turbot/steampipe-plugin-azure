# Table: azure_postgresql_flexible_server

Azure Database for PostgreSQL - Flexible Server is a fully managed database service designed to provide more granular control and flexibility over database management functions and configuration settings. The service generally provides more flexibility and server configuration customizations based on user requirements. The flexible server architecture allows users to collocate the database engine with the client tier for lower latency and choose high availability within a single availability zone and across multiple availability zones. Flexible servers also provide better cost optimization controls with the ability to stop/start your server and a burstable compute tier ideal for workloads that don't need full compute capacity continuously.

## Examples

### Basic info

```sql
select
  name,
  id,
  cloud_environment,
  properties,
  system_data ->> 'createdAt' as creation_time,
  location
from
  azure_postgresql_flexible_server;
```

### List SKU details of the flexible servers

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

```sql
select
  name,
  id,
  cloud_environment,
  system_data ->> 'createdAt',
  properties -> 'backup' ->> 'geoRedundantBackup' as geo_redundant_backup,
  location
from
  azure_postgresql_flexible_server
where
  properties -> 'backup' ->> 'geoRedundantBackup' = 'Enabled';
```

### List flexible servers configured in more than one availability zones

```sql
select
  name,
  id,
  cloud_environment,
  properties ->> 'availabilityZone',
  system_data ->> 'createdAt',
  location
from
  azure_postgresql_flexible_server
where
  (properties ->> 'availabilityZone')::int > 1;
```

### List flexible servers that have high availability mode enabled

```sql
select
  name,
  id,
  cloud_environment,
  properties -> 'highAvailability' ->> 'mode' as high_availability_mode,
  system_data ->> 'createdAt',
  location
from
  azure_postgresql_flexible_server
where
  properties -> 'highAvailability' ->> 'mode' = 'Disabled';
```

### List flexible servers that have password authentication enabled

```sql
select
  name,
  id,
  cloud_environment,
  properties -> 'authConfig' ->> 'passwordAuth' as password_authentication,
  system_data ->> 'createdAt',
  location
from
  azure_postgresql_flexible_server
where
  properties -> 'authConfig' ->> 'passwordAuth' = 'Enabled';
```

### List flexible servers that have AD authentication enabled

```sql
select
  name,
  id,
  cloud_environment,
  properties -> 'authConfig' ->> 'activeDirectoryAuth' as active_directory_authentication,
  system_data ->> 'createdAt',
  location
from
  azure_postgresql_flexible_server
where
  properties -> 'authConfig' ->> 'activeDirectoryAuth' = 'Enabled';
```

### List the servers that are publicly accessible

```sql
select
  name,
  id,
  cloud_environment,
  properties -> 'network' ->> 'publicNetworkAccess' as public_network_access,
  system_data ->> 'createdAt',
  location
from
  azure_postgresql_flexible_server
where
  properties -> 'network' ->> 'publicNetworkAccess' = 'Enabled';
```