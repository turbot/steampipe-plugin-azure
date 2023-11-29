---
title: "Steampipe Table: azure_eventhub_namespace - Query Azure Event Hubs Namespaces using SQL"
description: "Allows users to query Azure Event Hubs Namespaces."
---

# Table: azure_eventhub_namespace - Query Azure Event Hubs Namespaces using SQL

Azure Event Hubs is a big data streaming platform and event ingestion service, capable of receiving and processing millions of events per second. Event Hubs can process and analyze the data produced by connected devices and applications. A namespace is a container for all messaging components, multiple event hubs can reside within a single namespace, and namespaces are used as a way to isolate different sets of messaging components in separate environments.

## Table Usage Guide

The 'azure_eventhub_namespace' table provides insights into Azure Event Hubs Namespaces. As a DevOps engineer, explore namespace-specific details through this table, including the SKU name, capacity, tier, and associated metadata. Utilize it to uncover information about namespaces, such as their maximum throughput units, whether auto-inflate is enabled, and the network rules set. The schema presents a range of attributes of the namespace for your analysis, like the resource group, region, subscription ID, and associated tags.

## Examples

### Basic info
Explore the basic details of your Azure Eventhub namespaces, including their names, IDs, types, and provisioning states. This allows you to gain insights into their creation dates and current operational status for effective management and monitoring.

```sql
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
Determine the areas in which Azure EventHub namespaces are not making use of the virtual network service endpoint. This can be useful to identify potential network security gaps in your Azure environment.

```sql
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

### List unencrypted namespaces
Explore which Azure EventHub namespaces are unencrypted. This is useful for identifying potential security vulnerabilities within your Azure EventHub configuration.

```sql
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
Identify Azure EventHub namespaces where the auto-inflate feature is disabled. This can be useful for optimizing resource usage and managing costs.

```sql
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
Explore the details of private endpoint connections within your Azure EventHub Namespace. This can be useful in assessing the security and connectivity status of your system.

```sql
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