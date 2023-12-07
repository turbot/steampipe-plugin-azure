---
title: "Steampipe Table: azure_mysql_server - Query Azure MySQL Servers using SQL"
description: "Allows users to query Azure MySQL Servers, providing detailed information about the configuration, status, and capabilities of each server instance."
---

# Table: azure_mysql_server - Query Azure MySQL Servers using SQL

Azure Database for MySQL is a managed service that you use to run, manage, and scale highly available MySQL databases in the cloud. This service offers built-in high availability, security at every level of the application stack, and scaling in seconds with Azure. Azure Database for MySQL integrates with popular open-source frameworks and languages, and it's built on the trusted foundation of MySQL community edition.

## Table Usage Guide

The `azure_mysql_server` table provides insights into MySQL servers within Azure. As a database administrator, explore server-specific details through this table, including server version, storage capacity, and location. Utilize it to uncover information about servers, such as those with specific configurations, the status of each server, and the backup retention period.

## Examples

### Basic info
Explore the configuration of your Azure MySQL servers to understand their geographical locations and security settings, such as SSL enforcement and the minimal TLS version. This can help ensure your servers are optimally configured for both performance and security.

```sql+postgres
select
  name,
  id,
  location,
  ssl_enforcement,
  minimal_tls_version
from
  azure_mysql_server;
```

```sql+sqlite
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
Identify instances where servers have SSL enabled to ensure secure data transmission and safeguard against potential security risks.

```sql+postgres
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

```sql+sqlite
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
Explore which servers have enhanced security by having public network access disabled. This is useful for assessing potential vulnerabilities and ensuring that your servers are not exposed to unnecessary risks.

```sql+postgres
select
  name,
  id,
  public_network_access
from
  azure_mysql_server
where
  public_network_access = 'Disabled';
```

```sql+sqlite
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
Identify instances where Azure MySQL servers have the storage profile auto growth feature disabled. This can be useful for optimizing storage management and preventing unexpected storage limitations.

```sql+postgres
select
  name,
  id,
  storage_auto_grow
from
  azure_mysql_server
where
  storage_auto_grow = 'Disabled';
```

```sql+sqlite
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
Identify instances where your Azure MySQL servers are set to retain backups for more than 90 days. This can help in assessing your data retention strategy and ensuring compliance with your organization's data policies.

```sql+postgres
select
  name,
  id,
  backup_retention_days
from
  azure_mysql_server
where
  backup_retention_days > 90;
```

```sql+sqlite
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
Determine the areas in which your Azure MySQL servers may have security vulnerabilities by identifying those running with a minimum TLS version lower than 1.2. This can be used to enhance your server's security by upgrading to a higher TLS version.

```sql+postgres
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

```sql+sqlite
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
Explore the details of private endpoint connections in your Azure MySQL server. This query is useful in identifying the status and actions required for each connection, which can help in managing and troubleshooting your private endpoint connections.

```sql+postgres
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

```sql+sqlite
select
  name as server_name,
  id as server_id,
  json_extract(connections.value, '$.id') as connection_id,
  json_extract(connections.value, '$.privateEndpointPropertyId') as connection_private_endpoint_property_id,
  json_extract(connections.value, '$.privateLinkServiceConnectionStateActionsRequired') as connection_actions_required,
  json_extract(connections.value, '$.privateLinkServiceConnectionStateDescription') as connection_description,
  json_extract(connections.value, '$.privateLinkServiceConnectionStateStatus') as connection_status,
  json_extract(connections.value, '$.provisioningState') as connection_provisioning_state
from
  azure_mysql_server,
  json_each(private_endpoint_connections) as connections;
```

### List server keys
Explore the creation and configuration details of server keys in Azure MySQL servers. This can be useful to manage and track key usage and ensure security compliance across servers.

```sql+postgres
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

```sql+sqlite
select
  name as server_name,
  id as server_id,
  json_extract(keys.value, '$.creationDate') as keys_creation_date,
  json_extract(keys.value, '$.id') as keys_id,
  json_extract(keys.value, '$.kind') as keys_kind,
  json_extract(keys.value, '$.name') as keys_name,
  json_extract(keys.value, '$.serverKeyType') as keys_server_key_type,
  json_extract(keys.value, '$.type') as keys_type,
  json_extract(keys.value, '$.uri') as keys_uri
from
  azure_mysql_server,
  json_each(server_keys) as keys;
```

### List server configuration details
This query can be used to analyze and understand the configuration details of your servers on Azure MySQL. It's particularly useful when you need to assess the current settings of your servers for optimization or troubleshooting purposes.
**Note:** `Server configurations` is the same as `Server parameters` as shown in Azure MySQL server console


```sql+postgres
select
  name as server_name,
  id as server_id,
  configurations ->> 'Name' as configuration_name,
  configurations -> 'ConfigurationProperties' ->> 'value' as value
from
  azure_mysql_server,
  jsonb_array_elements(server_configurations) as configurations;
```

```sql+sqlite
select
  name as server_name,
  id as server_id,
  json_extract(configurations.value, '$.Name') as configuration_name,
  json_extract(json_extract(configurations.value, '$.ConfigurationProperties'), '$.value') as value
from
  azure_mysql_server,
  json_each(server_configurations) as configurations;
```

### Current state of audit_log_enabled parameter for the servers
This query is used to analyze the status of the audit log feature across various servers in Azure's MySQL service. It provides valuable insights into which servers have the audit log enabled, which is crucial for maintaining security and compliance within the system.

```sql+postgres
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

```sql+sqlite
select
  name as server_name,
  id as server_id,
  json_extract(configurations.value, '$.Name') as configuration_name,
  json_extract(json_extract(configurations.value, '$.ConfigurationProperties'), '$.value') as value
from
  azure_mysql_server,
  json_each(server_configurations) as configurations
where
  json_extract(configurations.value, '$.Name') = 'audit_log_enabled';
```

### List servers with slow_query_log parameter enabled
Determine the areas in which the slow query log parameter is enabled on Azure MySQL servers. This is useful for identifying potential performance issues and optimizing database operations.

```sql+postgres
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

```sql+sqlite
select
  name as server_name,
  id as server_id,
  json_extract(configurations.value, '$.Name') as configuration_name,
  json_extract(json_extract(configurations.value, '$.ConfigurationProperties'), '$.value') as value
from
  azure_mysql_server,
  json_each(server_configurations) as configurations
where
  json_extract(json_extract(configurations.value, '$.ConfigurationProperties'), '$.value') = 'ON'
  and json_extract(configurations.value, '$.Name') = 'slow_query_log';
```

### List servers with log_output parameter set to file
This example helps identify Azure MySQL servers that have their log output parameter configured to file. This can be useful for administrators who want to ensure that their server logs are being written to a file for easier access and review.

```sql+postgres
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

```sql+sqlite
select
  name as server_name,
  id as server_id,
  json_extract(configurations.value, '$.Name') as configuration_name,
  json_extract(json_extract(configurations.value, '$.ConfigurationProperties'), '$.value') as value
from
  azure_mysql_server,
  json_each(server_configurations) as configurations
where
   json_extract(json_extract(configurations.value, '$.ConfigurationProperties'), '$.value') = 'FILE'
   and json_extract(configurations.value, '$.Name') = 'log_output';
```

### Get VNET rules details of the server
Explore the configuration of your server to identify whether it is set to ignore missing Virtual Network Service Endpoints. This allows you to assess the security of your server by understanding its network connectivity settings.

```sql+postgres
select
  name as server_name,
  id as server_id,
  rules -> 'properties' ->> 'ignoreMissingVnetServiceEndpoint' as ignore_missing_vnet_service_endpoint,
  rules -> 'properties' ->> 'virtualNetworkSubnetId' as virtual_network_subnet_id
from
  azure_mysql_server,
  jsonb_array_elements(vnet_rules) as rules;
```

```sql+sqlite
select
  name as server_name,
  id as server_id,
  json_extract(rules.value, '$.properties.ignoreMissingVnetServiceEndpoint') as ignore_missing_vnet_service_endpoint,
  json_extract(rules.value, '$.properties.virtualNetworkSubnetId') as virtual_network_subnet_id
from
  azure_mysql_server,
  json_each(vnet_rules) as rules;
```

### Get the security alert policy for a particular server
Determine the security alert policy for a specific server within a given resource group. This is useful for assessing the security measures in place for that server.

```sql+postgres
select
  name,
  id,
  type,
  server_security_alert_policy
from
  azure_mysql_server
where
  resource_group = 'demo'
  and name = 'server-test-for-pr';
```

```sql+sqlite
select
  name,
  id,
  type,
  server_security_alert_policy
from
  azure_mysql_server
where
  resource_group = 'demo'
  and name = 'server-test-for-pr';
```