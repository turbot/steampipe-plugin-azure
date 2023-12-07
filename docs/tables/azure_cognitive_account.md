---
title: "Steampipe Table: azure_cognitive_account - Query Azure Cognitive Services Accounts using SQL"
description: "Allows users to query Azure Cognitive Services Accounts, providing insights into various cognitive services such as AI, speech analysis, language understanding, and search capabilities."
---

# Table: azure_cognitive_account - Query Azure Cognitive Services Accounts using SQL

Azure Cognitive Services is a collection of AI services and cognitive APIs to help you build intelligent apps. These services enable you to easily add cognitive features into your applications. The features include vision, speech, language, knowledge, and search capabilities.

## Table Usage Guide

The `azure_cognitive_account` table offers insights into the Azure Cognitive Services Accounts. As a developer or AI engineer, you can explore details about these accounts, such as the types of cognitive services being used, their configurations, and associated metadata. This information can be crucial for understanding the cognitive capabilities integrated into your applications and for optimizing their performance and usage.

## Examples

### Basic info
Determine the areas in which your Azure Cognitive Service accounts are provisioned, to better understand your resource usage and management. This is particularly useful for identifying any inconsistencies in provisioning and for gaining insights into your overall Azure resource allocation.

```sql+postgres
select
  name,
  id,
  kind,
  type,
  provisioning_state
from
  azure_cognitive_account;
```

```sql+sqlite
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
Discover the segments that have public network access enabled on their accounts. This is beneficial for identifying potential security risks and ensuring appropriate network access controls are in place.

```sql+postgres
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

```sql+sqlite
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
Determine the details of private endpoint connections for Azure cognitive accounts. This can help in managing and monitoring the security and access control of your cognitive services in Azure.

```sql+postgres
select
  name,
  id,
  connections ->> 'ID' as connection_id,
  connections ->> 'Name' as connection_name,
  connections ->> 'PrivateEndpointID' as property_private_endpoint_id,
  jsonb_pretty(connections -> 'PrivateLinkServiceConnectionState') as property_private_link_service_connection_state,
  connections ->> 'Type' as connection_type
from
  azure_cognitive_account,
  jsonb_array_elements(private_endpoint_connections) as connections;
```

```sql+sqlite
select
  name,
  id,
  json_extract(connections.value, '$.ID') as connection_id,
  json_extract(connections.value, '$.Name') as connection_name,
  json_extract(connections.value, '$.PrivateEndpointID') as property_private_endpoint_id,
  connections.value as property_private_link_service_connection_state,
  json_extract(connections.value, '$.Type') as connection_type
from
  azure_cognitive_account,
  json_each(private_endpoint_connections) as connections;
```

### List diagnostic setting details for accounts
Determine the diagnostic settings of Azure cognitive accounts to understand how they're configured. This is useful for auditing and managing account settings for optimal performance and security.

```sql+postgres
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

```sql+sqlite
select
  name,
  id,
  json_extract(settings.value, '$.id') as settings_id,
  json_extract(settings.value, '$.name') as settings_name,
  settings.value -> 'properties' -> 'logs' as settings_properties_logs,
  settings.value -> 'properties' -> 'metrics' as settings_properties_metrics,
  json_extract(settings.value, '$.properties.workspaceId') as settings_properties_workspaceId,
  json_extract(settings.value, '$.type') as settings_type
from
  azure_cognitive_account,
  json_each(diagnostic_settings) as settings;
```