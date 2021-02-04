# Table: azure_compute_virtual_machine

Azure Virtual Machines (VM) is one of several types of on-demand, scalable computing resources that Azure offers.

## Examples

### Virtual machine configuration info

```sql
select
  name,
  power_state,
  private_ips,
  public_ips,
  vm_id,
  size,
  os_type,
  image_offer,
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

### Disk Storage Summary, by VM

```sql
select
  vm.name,
  count(d) as num_disks,
  sum(d.disk_size_gb) as total_disk_size_gb
from
  azure.azure_compute_virtual_machine as vm
  left join azure_compute_disk as d on lower(vm.id) = lower(d.managed_by)
group by
  vm.name
order by
  vm.name;
```

### View Network Security Group Rules for a VM

```sql
select
  vm.name,
  nsg.name,
  jsonb_pretty(security_rules)
from
  azure.azure_compute_virtual_machine as vm,
  jsonb_array_elements(vm.network_interfaces) as vm_nic,
  azure_network_security_group as nsg,
  jsonb_array_elements(nsg.network_interfaces) as nsg_int
where
  lower(vm_nic ->> 'id') = lower(nsg_int ->> 'id')
  and vm.name = 'warehouse-01';
```
