# Table: azure_compute_virtual_machine_scale_set

Azure virtual machine scale sets let you create and manage a group of load balanced VMs. The number of VM instances can automatically increase or decrease in response to demand or a defined schedule. Scale sets provide high availability to your applications, and allow you to centrally manage, configure, and update a large number of VMs.

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
  azure_compute_virtual_machine_scale_set;
```

### List Standard tier virtual machine scale set

```sql
select
  name,
  id,
  sku_name,
  sku_tier
from
  azure_compute_virtual_machine_scale_set
where
   sku_tier = 'Standard';
```
