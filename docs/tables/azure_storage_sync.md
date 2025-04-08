---
title: "Steampipe Table: azure_storage_sync - Query Azure Storage Sync Services using SQL"
description: "Allows users to query Azure Storage Sync Services, specifically the synchronization details between Azure Blob storage and on-premises servers."
folder: "Storage"
---

# Table: azure_storage_sync - Query Azure Storage Sync Services using SQL

Azure Storage Sync Service is a feature within Microsoft Azure that enables synchronization of data across different Azure File shares. It allows for centralizing file services in Azure while maintaining local access to data. The service provides multi-site access, cloud tiering, integrated management, and change detection.

## Table Usage Guide

The `azure_storage_sync` table provides insights into Azure Storage Sync Services within Microsoft Azure. As a DevOps engineer, explore synchronization details through this table, including the sync group, registered servers, and associated metadata. Utilize it to uncover information about the synchronization status, such as those with pending synchronization, the relationships between servers, and the verification of synchronization activities.

## Examples

### Basic info
Determine the areas in which Azure's storage synchronization service is being utilized, along with its provisioning status. This can be useful for understanding the distribution and status of storage sync services across your Azure environment.

```sql+postgres
select
  name,
  id,
  type,
  provisioning_state
from
  azure_storage_sync;
```

```sql+sqlite
select
  name,
  id,
  type,
  provisioning_state
from
  azure_storage_sync;
```

### List storage sync which allows traffic through private endpoints only
Identify Azure storage syncs configured to accept incoming traffic solely through private network endpoints. This can be useful for maintaining security by ensuring data is only accessible within specific, controlled network environments.

```sql+postgres
select
  name,
  id,
  type,
  provisioning_state,
  incoming_traffic_policy
from
  azure_storage_sync
where
  incoming_traffic_policy = 'AllowVirtualNetworksOnly';
```

```sql+sqlite
select
  name,
  id,
  type,
  provisioning_state,
  incoming_traffic_policy
from
  azure_storage_sync
where
  incoming_traffic_policy = 'AllowVirtualNetworksOnly';
```

### List private endpoint connection details for accounts
This query allows you to explore the details of private endpoint connections associated with your accounts. It's particularly useful for gaining insights into the connection state and type, which can help assess the security and functionality of your data synchronization service.

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
  azure_storage_sync,
  jsonb_array_elements(private_endpoint_connections) as connections;
```

```sql+sqlite
select
  name,
  s.id,
  json_extract(connections.value, '$.ID') as connection_id,
  json_extract(connections.value, '$.Name') as connection_name,
  json_extract(connections.value, '$.PrivateEndpointPropertyID') as property_private_endpoint_id,
  connections.value as property_private_link_service_connection_state,
  json_extract(connections.value, '$.Type') as connection_type
from
  azure_storage_sync as s,
  json_each(private_endpoint_connections) as connections;
```