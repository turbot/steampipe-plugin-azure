# Table: azure_eventhub_namespace

An Event Hubs namespace provides DNS integrated network endpoints and a range of access control and network integration management features such as IP filtering, virtual network service endpoint, and Private Link and is the management container for one of multiple Event Hub instances (or topics, in Kafka parlance).

## Examples

### Basic info

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
