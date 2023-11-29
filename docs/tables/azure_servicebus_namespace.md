---
title: "Steampipe Table: azure_servicebus_namespace - Query Azure Service Bus Namespaces using SQL"
description: "Allows users to query Azure Service Bus Namespaces, providing insights into their properties, statuses, and configurations."
---

# Table: azure_servicebus_namespace - Query Azure Service Bus Namespaces using SQL

Azure Service Bus is a fully managed enterprise integration message broker. It can decouple applications and services, enabling them to communicate independently and reliably through messages. A namespace is a scoping container for all messaging components, providing a unique environment within the Service Bus where the queues, topics, and subscriptions reside.

## Table Usage Guide

The 'azure_servicebus_namespace' table provides insights into Azure Service Bus Namespaces, allowing you to explore details such as their properties, statuses, and configurations. As a DevOps engineer, leverage this table to understand the setup and management of your Service Bus Namespaces, including their SKU details, provisioning states, and associated tags. The schema presents a range of attributes of the Service Bus Namespace for your analysis, such as the name, region, resource group, subscription ID, and more. Utilize it to monitor the health and performance of your Azure Service Bus Namespaces, ensuring they meet predefined conditions and standards.

## Examples

### Basic info
Explore which Azure Service Bus namespaces are currently in use, to understand their provisioning status and when they were created. This can help in managing resources and planning for future capacity needs.

```sql
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
Explore which service bus namespaces in your Azure environment are operating on a premium tier, allowing you to assess your resource allocation and optimize cost management.

```sql
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
Explore the premium tier of your Azure Service Bus to identify namespaces that lack encryption. This is useful for improving your security measures and ensuring data protection.

```sql
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
Identify premium Azure Service Bus namespaces that are not utilizing a virtual network service endpoint. This can be used to enhance network security by ensuring all namespaces are connected to a secure network.

```sql
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

### List private endpoint connection details
Explore the details of private endpoint connections in your Azure Service Bus Namespace. This can be useful to understand the state and type of each connection, which can assist in managing and optimizing your network's performance.

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
  azure_servicebus_namespace,
  jsonb_array_elements(private_endpoint_connections) as connections;
```

### List encryption details
Explore the encryption details of your Azure Service Bus namespaces to understand their security configurations and ensure that they meet your organization's requirements. This query is particularly useful for auditing and compliance purposes.

```sql
select
  name,
  id,
  encryption ->> 'keySource' as key_source,
  jsonb_pretty(encryption -> 'keyVaultProperties') as key_vault_properties,
  encryption -> 'requireInfrastructureEncryption' as require_infrastructure_encryption
from
  azure_servicebus_namespace;
```