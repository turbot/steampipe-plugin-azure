---
title: "Steampipe Table: azure_postgresql_server - Query Azure PostgreSQL Servers using SQL"
description: "Allows users to query Azure PostgreSQL Servers, specifically providing access to configuration details, performance tiers, and resource usage."
---

# Table: azure_postgresql_server - Query Azure PostgreSQL Servers using SQL

Azure Database for PostgreSQL is a fully managed relational database service based on the open-source Postgres database engine. It's built to support the Postgres community edition, allowing users to leverage community-driven features and plugins. This service provides built-in high availability, security, and automated scaling to help businesses focus on application development rather than database management.

## Table Usage Guide

The `azure_postgresql_server` table provides insights into PostgreSQL servers within Azure Database for PostgreSQL. As a database administrator or developer, explore server-specific details through this table, including configuration settings, performance tiers, and resource usage. Utilize it to manage server settings, monitor resource consumption, and ensure optimal performance for your PostgreSQL databases within Azure.

## Examples

### Basic info
Explore which PostgreSQL servers are currently running on your Azure platform and where they are located. This information can assist in managing server distribution and planning for future resources.

```sql+postgres
select
  name,
  id,
  location
from
  azure_postgresql_server;
```

```sql+sqlite
select
  name,
  id,
  location
from
  azure_postgresql_server;
```

### List servers with encryption disabled
Discover the segments that contain servers with disabled encryption, enabling you to identify potential security vulnerabilities and take necessary action to enhance data protection.

```sql+postgres
select
  name,
  id,
  location,
  ssl_enforcement
from
  azure_postgresql_server
where
  ssl_enforcement = 'Disabled';
```

```sql+sqlite
select
  name,
  id,
  location,
  ssl_enforcement
from
  azure_postgresql_server
where
  ssl_enforcement = 'Disabled';
```

### List servers that allow access to Azure services
Explore which servers allow access to Azure services, a crucial element in managing security and controlling access. You can also pinpoint specific servers without an Active Directory admin, helping you identify potential vulnerabilities and areas that may require additional security measures.

```sql+postgres
select
  name,
  id,
  rule ->> 'Name' as rule_name,
  rule ->> 'Type' as rule_type,
  rule -> 'FirewallRuleProperties' ->> 'endIpAddress' as end_ip_address,
  rule -> 'FirewallRuleProperties' ->> 'startIpAddress' as start_ip_address
from
  azure_postgresql_server,
  jsonb_array_elements(firewall_rules) as rule
where
  rule ->> 'Name' = 'AllowAllWindowsAzureIps'
  and rule -> 'FirewallRuleProperties' ->> 'startIpAddress' = '0.0.0.0'
  and rule -> 'FirewallRuleProperties' ->> 'endIpAddress' = '0.0.0.0';
```

```sql+sqlite
select
  name,
  id,
  json_extract(rule.value, '$.Name') as rule_name,
  json_extract(rule.value, '$.Type') as rule_type,
  json_extract(rule.value, '$.FirewallRuleProperties.endIpAddress') as end_ip_address,
  json_extract(rule.value, '$.FirewallRuleProperties.startIpAddress') as start_ip_address
from
  azure_postgresql_server,
  json_each(firewall_rules) as rule
where
  json_extract(rule.value, '$.Name') = 'AllowAllWindowsAzureIps'
  and json_extract(rule.value, '$.FirewallRuleProperties.startIpAddress') = '0.0.0.0'
  and json_extract(rule.value, '$.FirewallRuleProperties.endIpAddress') = '0.0.0.0';
```

## List servers without an Active Directory admin

```sql+postgres
select
  name,
  id,
  location
from
  azure_postgresql_server
where
  server_administrators is null;
```

```sql+sqlite
select
  name,
  id,
  location
from
  azure_postgresql_server
where
  server_administrators is null;
```

### List servers with log checkpoints disabled
Determine the areas in which log checkpoints are disabled on your Azure PostgreSQL servers. This can help identify potential security vulnerabilities and improve your database management.

```sql+postgres
select
  name,
  configurations ->> 'Name' as configuration_name,
  configurations -> 'ConfigurationProperties' ->> 'value' as configuration_value
from
  azure_postgresql_server,
  jsonb_array_elements(server_configurations) as configurations
where
  configurations ->> 'Name' = 'log_checkpoints'
  and configurations -> 'ConfigurationProperties' ->> 'value' = 'OFF';
```

```sql+sqlite
select
  name,
  json_extract(configurations.value, '$.Name') as configuration_name,
  json_extract(configurations.value, '$.ConfigurationProperties.value') as configuration_value
from
  azure_postgresql_server,
  json_each(server_configurations) as configurations
where
  json_extract(configurations.value, '$.Name') = 'log_checkpoints'
  and json_extract(configurations.value, '$.ConfigurationProperties.value') = 'OFF';
```

### List servers with a logging retention period greater than 3 days
Determine the servers in your Azure PostgreSQL setup that have a logging retention period of more than 3 days. This is useful for ensuring your logging policies meet your organization's data retention requirements.

```sql+postgres
select
  name,
  configurations ->> 'Name' as configuration_name,
  configurations -> 'ConfigurationProperties' ->> 'value' as configuration_value
from
  azure_postgresql_server,
  jsonb_array_elements(server_configurations) as configurations
where
  configurations ->> 'Name' = 'log_retention_days'
  and (configurations -> 'ConfigurationProperties' ->> 'value')::INTEGER > 3;
```

```sql+sqlite
select
  name,
  json_extract(configurations.value, '$.Name') as configuration_name,
  json_extract(json_extract(configurations.value, '$.ConfigurationProperties'), '$.value') as configuration_value
from
  azure_postgresql_server,
  json_each(server_configurations) as configurations
where
  json_extract(configurations.value, '$.Name') = 'log_retention_days'
  and cast(json_extract(json_extract(configurations.value, '$.ConfigurationProperties'), '$.value') as INTEGER) > 3;
```

### List servers with geo-redundant backup storage disabled
Discover the segments where servers are running without geo-redundant backup storage. This is useful for identifying potential risk areas in your server infrastructure where data loss may occur in the event of a server failure.

```sql+postgres
select
  name,
  id,
  location,
  geo_redundant_backup
from
  azure_postgresql_server
where
  geo_redundant_backup = 'Disabled';
```

```sql+sqlite
select
  name,
  id,
  location,
  geo_redundant_backup
from
  azure_postgresql_server
where
  geo_redundant_backup = 'Disabled';
```

### List private endpoint connection details
Explore the status and details of private endpoint connections on your Azure PostgreSQL server. This can help in identifying any required actions or understanding the current provisioning state of these connections.

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
  azure_postgresql_server,
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
  azure_postgresql_server,
  json_each(private_endpoint_connections) as connections;
```