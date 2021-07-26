# Table: azure_container_registry

The Azure container registry is Microsoft's own hosting platform for Docker images. It is a private registry where you can store and manage private docker container images and other related artifacts. These images can then be pulled and run locally or used for container-based deployments to hosting platforms.

## Examples

### Basic info

```sql
select
  name,
  id,
  provisioning_state,
  creation_date,
  sku_tier,
  region
from
  azure_container_registry;
```

### List registries not encrypted with a customer-managed key

```sql
select
  name,
  encryption ->> 'status' as encryption_status,
  region
from
  azure_container_registry;
```

### List registries not configured with virtual network service endpoint

```sql
select
  name,
  network_rule_set ->> 'defaultAction' as network_rule_default_action,
  network_rule_set ->> 'virtualNetworkRules' as virtual_network_rules
from
  azure_container_registry
where
  network_rule_set is not null
  and network_rule_set ->> 'defaultAction' = 'Allow';
```

### List registries with admin user account enabled

```sql
select
  name,
  admin_user_enabled,
  region
from
  azure_container_registry
where
  admin_user_enabled;
```
