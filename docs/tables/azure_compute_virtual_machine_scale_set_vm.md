# Table: azure_compute_virtual_machine_scale_set_vm

You can scale the number of VMs in the scale set manually, or define rules to autoscale based on resource usage like CPU, memory demand, or network traffic. An Azure load balancer then distributes traffic to the VM instances in the scale set.

## Examples

### Basic info

```sql
select
  name,
  scale_set_name,
  instance_id,
  id,
  vm_id,
  region,
  resource_group
from
  azure_compute_virtual_machine_scale_set_vm;
```

### List Standard tier scale set vms

```sql
select
  name,
  scale_set_name,
  id,
  sku_name,
  sku_tier
from
  azure_compute_virtual_machine_scale_set_vm
where
   sku_tier = 'Standard';
```

### List all virtual machines under a specific scale set

```sql
select
  name,
  scale_set_name,
  id,
  sku_name,
  sku_tier
from
  azure_compute_virtual_machine_scale_set_vm
where 
  scale_set_name = 'my_vm_scale';
```

### View Network Security Group Rules for a VM

```sql
select
  vm.name,
  nsg.name,
  jsonb_pretty(security_rules)
from
  azure.azure_compute_virtual_machine_scale_set_vm as vm,
  jsonb_array_elements(vm.virtual_machine_network_profile) as vm_nic,
  azure_network_security_group as nsg,
  jsonb_array_elements(nsg.network_interfaces) as nsg_int
where
  lower(vm_nic ->> 'id') = lower(nsg_int ->> 'id')
  and vm.name = 'warehouse-01';
```