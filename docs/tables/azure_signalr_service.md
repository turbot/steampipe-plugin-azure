---
title: "Steampipe Table: azure_signalr_service - Query Azure SignalR Services using SQL"
description: "Allows users to query Azure SignalR Services."
---

# Table: azure_signalr_service - Query Azure SignalR Services using SQL

Azure SignalR Service is a fully-managed service that allows developers to focus on building real-time web experiences without worrying about capacity provisioning, reliable connections, scaling, encryption, or authentication. It is an Azure cloud-based service that is designed to support real-time web technologies like WebSockets and has built-in support for scaling your applications instantly. It also provides robust client SDKs for .NET, JavaScript, and Java, making it easier to build web applications with real-time features.

## Table Usage Guide

The 'azure_signalr_service' table provides insights into SignalR services within Azure. As a DevOps engineer, explore service-specific details through this table, including the service mode, primary and secondary connection strings, and associated metadata. Utilize it to uncover information about services, such as those with specific features, the connections between services, and the verification of connection strings. The schema presents a range of attributes of the SignalR service for your analysis, like the service tier, unit count, host name, and associated tags.

## Examples

### Basic info
Explore the status and types of Azure SignalR services to gain insights into their provisioning status, which can help in managing and troubleshooting these services efficiently.

```sql
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
Explore the network access control lists (ACLs) for SignalR service to understand their configuration and status. This can help you assess security measures and pinpoint areas for potential improvement.

```sql
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

### List private endpoint connection details for SignalR service
This example helps you identify the details of private endpoint connections for the SignalR service. It's useful for understanding the state and type of your connections, providing insights that can aid in service configuration and management.

```sql
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