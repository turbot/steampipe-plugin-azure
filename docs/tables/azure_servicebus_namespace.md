---
title: "Steampipe Table: azure_servicebus_namespace - Query Azure Service Bus Namespaces using SQL"
description: "Allows users to query Azure Service Bus Namespaces, providing critical insights into the properties, status, and configurations of each namespace."
---

# Table: azure_servicebus_namespace - Query Azure Service Bus Namespaces using SQL

Azure Service Bus is a fully managed enterprise integration message broker. Service Bus can decouple applications and services. Service Bus offers secure and reliable message delivery.

## Table Usage Guide

The `azure_servicebus_namespace` table provides insights into namespaces within Azure Service Bus. As a DevOps engineer, explore namespace-specific details through this table, including active message count, scheduled message count, and transfer message count. Utilize it to uncover information about namespaces, such as their status, SKU, and properties.

## Examples

### Basic info
Explore the status and tier level of your Azure Service Bus namespaces to assess their setup and monitor their creation time. This helps in managing resources and understanding their distribution across different tiers.

```sql+postgres
select
  name,
  id,
  sku_tier,
  provisioning_state,
  created_at
from
  azure_servicebus_namespace;
```

```sql+sqlite
select
  name,
  id,
  sku_tier,
  provisioning_state,
  created_at
from
  azure_servicebus_namespace;
```

### List premium namespaces
Analyze the settings to understand the distribution of premium-tier service bus namespaces across different regions in your Azure environment. This can help optimize resource allocation and cost management.

```sql+postgres
select
  name,
  sku_tier,
  region
from
  azure_servicebus_namespace
where
  sku_tier = 'Premium';
```

```sql+sqlite
select
  name,
  sku_tier,
  region
from
  azure_servicebus_namespace
where
  sku_tier = 'Premium';
```

### List unencrypted namespaces
Analyze the settings to understand the premium Azure Service Bus namespaces that lack encryption. This can be useful for identifying potential security risks and ensuring data protection standards are met.

```sql+postgres
select
  name,
  sku_tier,
  encryption
from
  azure_servicebus_namespace
where
  sku_tier = 'Premium'
  and encryption is null;
```

```sql+sqlite
select
  name,
  sku_tier,
  encryption
from
  azure_servicebus_namespace
where
  sku_tier = 'Premium'
  and encryption is null;
```

### List namespaces not using a virtual network service endpoint
Determine the areas in your Azure service bus setup where premium tier namespaces are not utilizing a virtual network service endpoint. This can be useful to improve security by ensuring all communication within your service bus happens over your virtual network.

```sql+postgres
select
  name,
  region,
  network_rule_set -> 'properties' -> 'virtualNetworkRules' as virtual_network_rules
from
  azure_servicebus_namespace
where
  sku_tier = 'Premium'
  and (
    jsonb_array_length(network_rule_set -> 'properties' -> 'virtualNetworkRules') = 0
    or exists (
      select
        *
      from
        jsonb_array_elements(network_rule_set -> 'properties' -> 'virtualNetworkRules') as t
      where
        t -> 'subnet' ->> 'id' is null
    )
  );
```

```sql+sqlite
select
  name,
  region,
  json_extract(network_rule_set, '$.properties.virtualNetworkRules') as virtual_network_rules
from
  azure_servicebus_namespace
where
  sku_tier = 'Premium'
  and (
    json_array_length(json_extract(network_rule_set, '$.properties.virtualNetworkRules')) = 0
    or exists (
      select
        1
      from
        json_each(json_extract(network_rule_set, '$.properties.virtualNetworkRules')) as t
      where
        json_extract(t.value, '$.subnet.id') is null
    )
  );
```

### List private endpoint connection details
Explore the details of private endpoint connections in Azure Service Bus Namespace to understand their provisioning state and connection types. This is useful for assessing the security and configuration of your cloud resources.

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
  azure_servicebus_namespace,
  jsonb_array_elements(private_endpoint_connections) as connections;
```

```sql+sqlite
select
  name,
  n.id,
  json_extract(connections.value, '$.id') as connection_id,
  json_extract(connections.value, '$.name') as connection_name,
  json_extract(connections.value, '$.privateEndpointPropertyID') as property_private_endpoint_id,
  json_extract(connections.value, '$.provisioningState') as property_provisioning_state,
  connections.value as property_private_link_service_connection_state,
  json_extract(connections.value, '$.type') as connection_type
from
  azure_servicebus_namespace as n,
  json_each(private_endpoint_connections) as connections;
```

### List encryption details
Determine the encryption specifications of your Azure Service Bus namespaces. This can provide insights into your security configurations, helping you understand if your data is properly secured and whether infrastructure encryption is required.

```sql+postgres
select
  name,
  id,
  encryption ->> 'keySource' as key_source,
  jsonb_pretty(encryption -> 'keyVaultProperties') as key_vault_properties,
  encryption -> 'requireInfrastructureEncryption' as require_infrastructure_encryption
from
  azure_servicebus_namespace;
```

```sql+sqlite
select
  name,
  id,
  json_extract(encryption, '$.keySource') as key_source,
  encryption as key_vault_properties,
  json_extract(encryption, '$.requireInfrastructureEncryption') as require_infrastructure_encryption
from
  azure_servicebus_namespace;
```