# Table: azure_postgresql_server

Azure Database for PostgreSQL is a relational database service based on the open-source Postgres database engine. It's a fully managed database-as-a-service that can handle mission-critical workloads with predictable performance, security, high availability, and dynamic scalability.

## Examples

### List servers for which SSL is enabled

```sql
select
  name,
  id,
  location,
  ssl_enforcement
from
  azure_postgresql_server
where
  ssl_enforcement = 'Enabled';
```

### List servers for which 'Allow access to Azure services' feature is enabled

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

## List servers for which Active Directory Admin is not configured

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

### List servers for which 'log_checkpoints' is set to 'OFF'

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

### List servers for which 'log_retention_days' is greater than 3 days

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
