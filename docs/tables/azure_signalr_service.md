# Table: azure_signalr_service

Azure SignalR Service is a fully-managed service which allows developers to focus on building real-time web experiences without worrying about capacity provisioning, reliable connections, scaling, encryption or authentication.

## Examples

### Basic info

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
