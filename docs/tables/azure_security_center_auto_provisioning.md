# Table: azure_security_center_auto_provisioning

Azure security center auto provisioning settings expose the auto provisioning settings of the subscriptions.

## Examples

### Basic info

```sql
select
  id,
  name,
  type,
  auto_provision
from
  azure_security_center_auto_provisioning;
```

### List subscriptions that have automatic provisioning of VM monitoring agent enabled

```sql
select
  id,
  name,
  type,
  auto_provision
from
  azure_security_center_auto_provisioning
where
  auto_provision = 'On';
```
