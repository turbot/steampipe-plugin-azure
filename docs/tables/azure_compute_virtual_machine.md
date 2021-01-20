# Table: azure_compute_virtual_machine

Azure Virtual Machines (VM) is one of several types of on-demand, scalable computing resources that Azure offers.

## Examples


### Virtual machine configuration info

```sql
select
  name,
  vm_id,
  size,
  os_type,
  os_disk_size_gb image_offer,
  image_sku
from
  azure_compute_virtual_machine;
```


### Virtual machine count in each region

```sql
select
  region,
  count(name)
from
  azure_compute_virtual_machine
group by
  region;
```


### List of VMs whose OS disk is not encrypted by customer managed key

```sql
select
  vm.name,
  disk.encryption_type
from
  azure_compute_disk as disk
  join azure_compute_virtual_machine as vm on disk.name = vm.os_disk_name
where
  not disk.encryption_type = 'EncryptionAtRestWithCustomerKey';
```


### List of VMs provisioned with undesired(for example Standard_D8s_v3 and Standard_DS3_v3 is desired) sizes.

```sql
select
  size,
  count(*) as count
from
  azure_compute_virtual_machine
where
  size not in ('Standard_D8s_v3', 'Standard_DS3_v3')
group by
  size;
```


### Availability set info of VMs

```sql
select
  vm.name vm_name,
  aset.name availability_set_name,
  aset.platform_fault_domain_count,
  aset.platform_update_domain_count,
  aset.sku_name
from
  azure_compute_availability_set as aset
  join azure_compute_virtual_machine as vm on lower(aset.id) = lower(vm.availability_set_id);
```


### List of all spot type VM and their eviction policy

```sql
select
  name,
  vm_id,
  eviction_policy
from
  azure_compute_virtual_machine
where
  priority = 'Spot';
```
