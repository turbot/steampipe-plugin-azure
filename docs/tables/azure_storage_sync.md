# Table: azure_storage_sync

Azure File Sync is a service that allows you to cache several Azure file shares on an on-premises Windows Server or cloud VM.

## Examples

### Basic info

```sql
select
  name,
  id,
  type,
  provisioning_state
from
  azure_storage_sync;
```

### List storage sync which allows traffic only from private endpoints

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
