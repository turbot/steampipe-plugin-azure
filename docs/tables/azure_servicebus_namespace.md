# Table: azure_servicebus_namespace

A ServiceBus namespace is a container for all messaging components (queues and topics). Multiple queues and topics can be in a single namespace, and namespaces often serve as application containers. A Service Bus namespace is your own capacity slice of a large cluster made up of dozens of all-active virtual machines.

## Examples

### Basic info

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
