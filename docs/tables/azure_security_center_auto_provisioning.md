# Table: azure_security_center_auto_provisioning

Azure security center auto provisioning settings exposes the auto provisioning settings of the subscriptions.

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

### Ensure that Automatic provisioning of monitoring agent is set to On

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
