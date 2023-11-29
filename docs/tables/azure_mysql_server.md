---
title: "Steampipe Table: azure_mysql_server - Query Azure MySQL Servers using SQL"
description: "Allows users to query Azure MySQL Servers, fetching detailed information about the configuration and status of these database servers."
---

# Table: azure_mysql_server - Query Azure MySQL Servers using SQL

Azure MySQL Server is a fully managed database service that makes it easy to set up, maintain, manage, and administer your MySQL relational databases on the cloud. It provides built-in high availability with no additional cost and you can scale up or down quickly to meet your workload needs. Azure MySQL Server also supports connecting your MySQL databases to popular analytics tools for comprehensive insights and business intelligence.

## Table Usage Guide

The 'azure_mysql_server' table provides insights into MySQL servers within Azure. As a DevOps engineer, explore server-specific details through this table, including server name, location, resource group, SKU name, and associated metadata. Utilize it to uncover information about servers, such as the version of MySQL running, SSL enforcement status, and storage auto-grow settings. The schema presents a range of attributes of the MySQL server for your analysis, like the server ID, creation date, administrator login name, and more.

## Examples

### Basic info
Explore the settings of your Azure MySQL server to understand its location and security enforcement policies, such as SSL enforcement and the minimum TLS version. This is useful for ensuring your server is properly configured for secure data transmission.

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
Explore which servers in your Azure MySQL Server have SSL enforcement enabled. This is useful for ensuring that your servers are secure and adhering to best practices for data protection.

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
Explore which servers within your Azure MySQL setup have public network access disabled. This can help enhance security by identifying servers that are not exposed to potential external threats.

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
Identify instances where Azure MySQL servers have their storage profile auto growth feature disabled. This is useful to manage storage and avoid unexpected database growth.

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
Analyze the settings to understand which servers have their backup retention period set for more than 90 days. This is useful for ensuring data retention compliance and managing storage costs.

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
Discover the servers that are potentially vulnerable due to lower than recommended TLS versions. This is useful in identifying and addressing security risks in your Azure MySQL server configurations.

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
Determine the areas in which private endpoint connections are needed for your Azure MySQL server. This query helps you understand the state of these connections, including any actions required, providing valuable insights for managing and optimizing your server's security.

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
This query is used to examine the keys associated with each server in your Azure MySQL database. It's useful for understanding the types and creation dates of these keys, which can aid in managing security and access controls.

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
Explore the configuration details of your servers to gain insights into their set-up and manage them effectively. This query is particularly useful for understanding and managing the settings of your Azure MySQL servers.
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
Analyze the settings to understand the current status of the audit log feature for your servers. This can be useful for ensuring compliance with security protocols and maintaining a record of server activity.

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
Determine the areas in which certain servers have the 'slow_query_log' parameter enabled. This can be useful to identify potential performance issues and optimize server configurations accordingly.

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
Explore which Azure MySQL servers have their log output parameter set to a file. This is useful to identify servers that are storing their logs as files, which could potentially take up a lot of storage space.

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
Explore the configuration of a server to understand its Virtual Network (VNET) rules. This is useful for assessing network security and connectivity settings for your Azure MySQL server.

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

### Get the security alert policy for a particular server
Analyze the settings to understand the security alert policy associated with a specific server in a given resource group. This is particularly useful when you need to assess the security configurations of your servers for compliance or auditing purposes.

```sql
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