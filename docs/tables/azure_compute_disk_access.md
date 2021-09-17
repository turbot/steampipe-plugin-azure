# Table: azure_compute_disk_access

Disk access resource is required to use private link to export and import managed disks, first you create a disk access resource and link it to a virtual network in the same subscription by creating a private endpoint. Then, associate a disk or a snapshot with a disk access instance.

## Examples

### Basic info

```sql
select
  name,
  id,
  type,
  provisioning_state,
  resource_group
from
  azure_compute_disk_access;
```

### List failed disk accesses

```sql
select
  name,
  id,
  type,
  provisioning_state,
  resource_group
from
  azure_compute_disk_access
where
  provisioning_state = 'Failed';
```
