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
  connections ->> 'privateEndpointPropertyId' as connection_private_endpoint_property_id,
  connections ->> 'privateLinkServiceConnectionStateActionsRequired' as connection_actions_required,
  connections ->> 'privateLinkServiceConnectionStateDescription' as connection_description,
  connections ->> 'privateLinkServiceConnectionStateStatus' as connection_status,
  connections ->> 'provisioningState' as connection_provisioning_state
from
  azure_mysql_server,
  jsonb_array_elements(private_endpoint_connections) as connections;
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

### List server configuration details

**Note:** `Server configurations` is the same as `Server parameters` as shown in Azure MySQL server console

```sql
select
  name as server_name,
  id as server_id,
  configurations ->> 'Name' as configuration_name,
  configurations -> 'ConfigurationProperties' ->> 'value' as value
from
  azure_mysql_server,
  jsonb_array_elements(server_configurations) as configurations;
```

### Current state of audit_log_enabled parameter for the servers

```sql
select
  name as server_name,
  id as server_id,
  configurations ->> 'Name' as configuration_name,
  configurations -> 'ConfigurationProperties' ->> 'value' as value
from
  azure_mysql_server,
  jsonb_array_elements(server_configurations) as configurations
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
  azure_mysql_server,
  jsonb_array_elements(server_configurations) as configurations
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
  azure_mysql_server,
  jsonb_array_elements(server_configurations) as configurations
where
   configurations ->'ConfigurationProperties' ->> 'value' = 'FILE'
   and configurations ->> 'Name' = 'log_output';
```

### Get VNET rules details of the server

```sql
select
  name as server_name,
  id as server_id,
  rules -> 'properties' ->> 'ignoreMissingVnetServiceEndpoint' as ignore_missing_vnet_service_endpoint,
  rules -> 'properties' ->> 'virtualNetworkSubnetId' as virtual_network_subnet_id
from
  azure_mysql_server,
  jsonb_array_elements(vnet_rules) as rules;
```