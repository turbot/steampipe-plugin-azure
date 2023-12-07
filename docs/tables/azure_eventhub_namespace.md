---
title: "Steampipe Table: azure_eventhub_namespace - Query Azure Event Hubs Namespaces using SQL"
description: "Allows users to query Azure Event Hubs Namespaces, providing insights into the details of each namespace, including its name, region, resource group, SKU, and more."
---

# Table: azure_eventhub_namespace - Query Azure Event Hubs Namespaces using SQL

Azure Event Hubs is a big data streaming platform and event ingestion service. It can receive and process millions of events per second. A namespace is a scoping container for Event Hubs under an Azure subscription.

## Table Usage Guide

The `azure_eventhub_namespace` table provides insights into Azure Event Hubs Namespaces. As a Data Engineer, you can explore namespace-specific details through this table, including its name, region, resource group, SKU, and more. Utilize it to manage and monitor the health and performance of your Azure Event Hubs Namespaces.

## Examples

### Basic info
Discover the segments that provide you with a comprehensive overview of your Azure EventHub namespaces. This includes details like the provisioning status and creation date, which can help you track and manage your resources more effectively.

```sql+postgres
select
  name,
  id,
  type,
  provisioning_state,
  created_at
from
  azure_eventhub_namespace;
```

```sql+sqlite
select
  name,
  id,
  type,
  provisioning_state,
  created_at
from
  azure_eventhub_namespace;
```

### List namespaces not configured to use virtual network service endpoint
Determine the areas in which Azure EventHub namespaces are not utilizing the virtual network service endpoint. This query is beneficial in identifying potential security loopholes, as these namespaces might be exposed to risks without the added protection of a virtual network.

```sql+postgres
select
  name,
  id,
  type,
  network_rule_set -> 'properties' -> 'virtualNetworkRules' as virtual_network_rules
from
  azure_eventhub_namespace
where
  network_rule_set -> 'properties' -> 'virtualNetworkRules' = '[]';
```

```sql+sqlite
select
  name,
  id,
  type,
  json_extract(network_rule_set, '$.properties.virtualNetworkRules') as virtual_network_rules
from
  azure_eventhub_namespace
where
  json_extract(network_rule_set, '$.properties.virtualNetworkRules') = '[]';
```

### List unencrypted namespaces
Discover the segments that are unencrypted within the Azure EventHub namespace. This is useful for identifying potential security vulnerabilities where sensitive data might not be adequately protected.

```sql+postgres
select
  name,
  id,
  type,
  encryption
from
  azure_eventhub_namespace
where
  encryption is null;
```

```sql+sqlite
select
  name,
  id,
  type,
  encryption
from
  azure_eventhub_namespace
where
  encryption is null;
```

### List namespaces with auto-inflate disabled
Identify the Azure EventHub namespaces where the auto-inflate feature is turned off. This can be useful to pinpoint potential resource limitations or throttling issues in your Azure EventHub service.

```sql+postgres
select
  name,
  id,
  type,
  is_auto_inflate_enabled
from
  azure_eventhub_namespace
where
  not is_auto_inflate_enabled;
```

```sql+sqlite
select
  name,
  id,
  type,
  is_auto_inflate_enabled
from
  azure_eventhub_namespace
where
  not is_auto_inflate_enabled;
```

### List private endpoint connection details
Determine the details of private endpoint connections within your Azure EventHub Namespace. This can help understand the state and type of connections, which is useful for managing and troubleshooting your network connectivity.

```sql+postgres
select
  name,
  id,
  connections ->> 'id' as connection_id,
  connections ->> 'name' as connection_name,
  connections ->> 'privateEndpointPropertyID' as property_private_endpoint_id,
  connections ->> 'provisioningState' as property_provisioning_state,
  jsonb_pretty(connections -> 'privateLinkServiceConnectionState') as property_private_link_service_connection_state,
  connections ->> 'type' as connection_type
from
  azure_eventhub_namespace,
  jsonb_array_elements(private_endpoint_connections) as connections;
```

```sql+sqlite
select
  name,
  id,
  json_extract(connections.value, '$.id') as connection_id,
  json_extract(connections.value, '$.name') as connection_name,
  json_extract(connections.value, '$.privateEndpointPropertyID') as property_private_endpoint_id,
  json_extract(connections.value, '$.provisioningState') as property_provisioning_state,
  connections.value as property_private_link_service_connection_state,
  json_extract(connections.value, '$.type') as connection_type
from
  azure_eventhub_namespace,
  json_each(private_endpoint_connections) as connections;
```