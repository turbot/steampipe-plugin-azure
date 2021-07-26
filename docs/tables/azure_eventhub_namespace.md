# Table: azure_eventhub_namespace

An Event Hubs namespace provides a unique scoping container, in which you create one or more event hubs.

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


### List Namespaces with no virtual network service endpoint configured

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


### List Namespaces with encryption disabled

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


### List Namespaces with auto inflate disabled

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