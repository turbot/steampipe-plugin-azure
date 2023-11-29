---
title: "Steampipe Table: azure_postgresql_server - Query Azure Database for PostgreSQL Servers using SQL"
description: "Allows users to query Azure Database for PostgreSQL Servers."
---

# Table: azure_postgresql_server - Query Azure Database for PostgreSQL Servers using SQL

Azure Database for PostgreSQL is a fully managed database service provided by Microsoft Azure. It is built on the open-source PostgreSQL database engine and offers compatibility with PostgreSQL, which allows users to use familiar PostgreSQL tools and scripts. This service provides built-in high availability, automatic backups, and scaling of resources in minutes without application downtime.

## Table Usage Guide

The 'azure_postgresql_server' table provides insights into PostgreSQL servers within Azure Database for PostgreSQL. As a database administrator or DevOps engineer, explore server-specific details through this table, including configurations, network settings, and associated metadata. Utilize it to uncover information about servers, such as those with specific configurations, the networking rules applied to servers, and the verification of server statuses. The schema presents a range of attributes of the PostgreSQL server for your analysis, like the server name, resource group, region, version, SSL enforcement, and storage capacity.

## Examples

### Basic info
Explore the details of your Azure PostgreSQL servers, such as their names, IDs, and locations. This can be useful for managing and organizing your servers across various locations.

```sql
select
  name,
  id,
  location
from
  azure_postgresql_server;
```

### List servers with encryption disabled
Discover the segments that have encryption disabled on their servers. This is crucial for identifying potential security risks and ensuring data protection standards are upheld.

```sql
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
1. Discover the segments that allow access to Azure services from any location, which could potentially indicate a security risk.
2. Identify instances where servers lack an assigned Active Directory admin, which could pose a management or security issue.

```sql
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

## List servers without an Active Directory admin

```sql
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
Explore which servers have disabled log checkpoints, which could potentially compromise data integrity and recovery. This can be useful for auditing server configurations and ensuring optimal data safety practices.

```sql
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

### List servers with a logging retention period greater than 3 days
This query is useful for identifying servers that maintain logs for more than three days, which can be beneficial for organizations that need to keep track of server activities for extended periods for auditing or troubleshooting purposes.

```sql
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

### List servers with geo-redundant backup storage disabled
Uncover the details of servers that have disabled geo-redundant backup storage, helping to highlight potential areas of risk in your Azure PostgreSQL Server setup. This is useful for ensuring data redundancy and disaster recovery planning.

```sql
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
Explore the status and details of private endpoint connections within a server. This can be useful to monitor and manage the connections' state and actions required for maintaining optimal server performance.

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
  azure_postgresql_server,
  jsonb_array_elements(private_endpoint_connections) as connections;
```