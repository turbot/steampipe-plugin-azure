# Table: azure_virtual_machine_scale_set

A virtual network gateway is used to establish secure, cross-premises connectivity.

## Examples

### Basic info

```sql
select
  name,
  id,
  identity,
  region,
  resource_group
from
  azure_virtual_machine_scale_set;
```

### List Standard tier virtual machine scale set

```sql
select
  name,
  id,
  sku_name,
  sku_tier
from
  azure_virtual_machine_scale_set
where
   sku_tier = 'Standard';
```
