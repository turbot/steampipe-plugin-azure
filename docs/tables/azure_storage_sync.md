---
title: "Steampipe Table: azure_storage_sync - Query Azure Storage Sync Services using SQL"
description: "Allows users to query Azure Storage Sync Services, which are used to synchronize files across multiple Azure File shares."
---

# Table: azure_storage_sync - Query Azure Storage Sync Services using SQL

Azure Storage Sync Services is a feature within Microsoft Azure that allows you to synchronize files across multiple Azure File shares. It provides a centralized way to manage and synchronize files across different regions and offices. Azure Storage Sync Services helps you to keep your data close to where it is being used, irrespective of whether it's being used on-premises or in the cloud.

## Table Usage Guide

The 'azure_storage_sync' table provides insights into Azure Storage Sync Services within Microsoft Azure. As a DevOps engineer, explore service-specific details through this table, including the synchronization status, last synchronization time, and associated metadata. Utilize it to uncover information about storage sync services, such as those with synchronization issues, the relationships between different services, and the verification of synchronization health. The schema presents a range of attributes of the Azure Storage Sync Services for your analysis, like the service name, id, type, and associated tags.

## Examples

### Basic info
Explore the status and types of your Azure storage synchronization services. This can help in managing and monitoring your storage resources effectively.

```sql
select
  name,
  id,
  type,
  provisioning_state
from
  azure_storage_sync;
```

### List storage sync which allows traffic through private endpoints only
Determine the areas in which your Azure storage sync is configured to allow traffic through private endpoints only. This is particularly useful for enhancing security by ensuring that network traffic is restricted to virtual networks only.

```sql
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
This query is useful for gaining insights into the details of private endpoint connections associated with your Azure storage sync accounts. It helps in analyzing the connection settings to understand the status and type of each connection, which can be critical for auditing and compliance purposes.

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
  azure_storage_sync,
  jsonb_array_elements(private_endpoint_connections) as connections;
```