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
