---
title: "Steampipe Table: azure_signalr_service - Query Azure SignalR Services using SQL"
description: "Allows users to query Azure SignalR Services, providing insights into real-time web functionality for applications."
folder: "SignalR Service"
---

# Table: azure_signalr_service - Query Azure SignalR Services using SQL

Azure SignalR Service is a fully-managed service that allows developers to focus on building real-time web experiences without worrying about capacity provisioning, reliable connections, scaling, encryption, or authentication. It is an Azure service that helps to simplify the process of adding real-time web functionality to applications over HTTP. This functionality enables applications to stream content updates to connected clients instantly.

## Table Usage Guide

The `azure_signalr_service` table offers insights into Azure SignalR Services within Microsoft Azure. As a developer or system administrator, you can explore service-specific details through this table, including service tiers, client negotiation endpoints, and associated metadata. Use this table to monitor and manage your SignalR Services, identify their capacity and usage, and ensure optimal performance and security for your real-time web applications.

## Examples

### Basic info
Explore the status and type of your Azure SignalR services to understand their current operational state and categorization. This is beneficial for managing and monitoring your application's real-time messaging services.

```sql+postgres
select
  name,
  id,
  type,
  kind,
  provisioning_state
from
  azure_signalr_service;
```

```sql+sqlite
select
  name,
  id,
  type,
  kind,
  provisioning_state
from
  azure_signalr_service;
```

### List network ACL details for SignalR service
This query helps you explore the network access control list (ACL) details for your SignalR service. It's useful for understanding the default actions and the configuration of private and public networks, which in turn can aid in managing access control and enhancing security.

```sql+postgres
select
  name,
  id,
  type,
  provisioning_state,
  network_acls ->> 'defaultAction' as default_action,
  jsonb_pretty(network_acls -> 'privateEndpoints') as private_endpoints,
  jsonb_pretty(network_acls -> 'publicNetwork') as public_network
from
  azure_signalr_service;
```

```sql+sqlite
select
  name,
  id,
  type,
  provisioning_state,
  json_extract(network_acls, '$.defaultAction') as default_action,
  network_acls as private_endpoints,
  network_acls as public_network
from
  azure_signalr_service;
```

### List private endpoint connection details for SignalR service
Determine the areas in which private endpoint connections are established for the SignalR service. This is useful for understanding and managing the security and access of your SignalR services.

```sql+postgres
select
  name,
  id,
  connections ->> 'ID' as connection_id,
  connections ->> 'Name' as connection_name,
  connections ->> 'PrivateEndpointPropertyID' as property_private_endpoint_id,
  jsonb_pretty(connections -> 'PrivateLinkServiceConnectionState') as property_private_link_service_connection_state,
  connections ->> 'Type' as connection_type
from
  azure_signalr_service,
  jsonb_array_elements(private_endpoint_connections) as connections;
```

```sql+sqlite
select
  name,
  id,
  json_extract(connections.value, '$.ID') as connection_id,
  json_extract(connections.value, '$.Name') as connection_name,
  json_extract(connections.value, '$.PrivateEndpointPropertyID') as property_private_endpoint_id,
  connections.value as property_private_link_service_connection_state,
  json_extract(connections.value, '$.Type') as connection_type
from
  azure_signalr_service,
  json_each(private_endpoint_connections) as connections;
```