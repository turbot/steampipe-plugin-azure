# Table: azure_cognitive_account

Azure Cognitive Services are cloud-based services with REST APIs and client library SDKs available to help you build cognitive intelligence into your applications. You can add cognitive features to your applications without having artificial intelligence (AI) or data science skills. Azure Cognitive Services comprise various AI services that enable you to build cognitive solutions that can see, hear, speak, understand, and even make decisions.

## Examples

### Basic info

```sql
select
  name,
  id,
  kind,
  type,
  provisioning_state
from
  azure_cognitive_account;
```

### List accounts with enabled public network access

```sql
select
  name,
  id,
  kind,
  type,
  provisioning_state,
  public_network_access
from
  azure_cognitive_account
where
  public_network_access = 'Enabled';
```

### List private endpoint connection details for accounts

```sql
select
  name,
  id,
  connections ->> 'ID' as connection_id,
  connections ->> 'Name' as connection_name,
  jsonb_pretty(connections -> 'Properties') as connection_properties,
  connections ->> 'Type' as connection_type
from
  azure_cognitive_account,
  jsonb_array_elements(private_endpoint_connections) as connections;
```

### List restored accounts

```sql
select
  name,
  id,
  kind,
  type,
  provisioning_state,
  public_network_access
from
  azure_cognitive_account
where
  restore;
```

### List diagnostic setting details for accounts

```sql
select
  name,
  id,
  settings ->> 'id' as settings_id,
  settings ->> 'name' as settings_name,
  jsonb_pretty(settings -> 'properties' -> 'logs') as settings_properties_logs,
  jsonb_pretty(settings -> 'properties' -> 'metrics') as settings_properties_metrics,
  settings -> 'properties' ->> 'workspaceId' as settings_properties_workspaceId,
  settings ->> 'type' as settings_type
from
  azure_cognitive_account,
  jsonb_array_elements(diagnostic_settings) as settings;
```
