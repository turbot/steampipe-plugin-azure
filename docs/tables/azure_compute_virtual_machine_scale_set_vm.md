---
title: "Steampipe Table: azure_compute_virtual_machine_scale_set_vm - Query Azure Compute Virtual Machine Scale Set VMs using SQL"
description: "Allows users to query Azure Compute Virtual Machine Scale Set VMs, providing insights into the configuration, state, and associated metadata of each virtual machine in a scale set."
folder: "Compute"
---

# Table: azure_compute_virtual_machine_scale_set_vm - Query Azure Compute Virtual Machine Scale Set VMs using SQL

Azure Compute Virtual Machine Scale Sets are a group of identical, load-balanced VMs. They are designed to support true auto-scale, no pre-provisioning of VMs is required, and they let you centrally manage, configure, and update a large number of VMs. With auto-scale, VMs get automatically created and added to a load balancer and get removed when not in use.

## Table Usage Guide

The `azure_compute_virtual_machine_scale_set_vm` table provides insights into the individual virtual machines within an Azure Virtual Machine Scale Set. As a system administrator, you can explore VM-specific details through this table, including the current state, configuration, and associated metadata. This table is useful for monitoring the status and performance of each VM in a scale set, enabling efficient resource management and troubleshooting.

## Examples

### Basic info
Analyze the settings to understand the distribution and organization of your virtual machine scale sets in Azure. This can be helpful to manage resources and monitor regional deployment of your virtual machines effectively.

```sql+postgres
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

```sql+sqlite
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

### List Standard tier scale set virtual machine
Determine the areas in which standard tier virtual machine scale sets are being used within your Azure environment. This allows for efficient resource allocation and cost management.

```sql+postgres
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

```sql+sqlite
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
Explore which virtual machines are part of a particular scale set in Azure. This is useful for managing resources and understanding the distribution of your virtual machines within specific scale sets.

```sql+postgres
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

```sql+sqlite
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

### View Network Security Group Rules for a virtual machine
Determine the security rules applied to a specific virtual machine within a network. This is useful for assessing the security measures in place and identifying any potential vulnerabilities or areas for improvement.

```sql+postgres
select
  vm.name,
  nsg.name,
  jsonb_pretty(security_rules)
from
  azure_compute_virtual_machine_scale_set_vm as vm,
  jsonb_array_elements(vm.virtual_machine_network_profile) as vm_nic,
  azure_network_security_group as nsg,
  jsonb_array_elements(nsg.network_interfaces) as nsg_int
where
  lower(vm_nic ->> 'id') = lower(nsg_int ->> 'id')
  and vm.name = 'warehouse-01';
```

```sql+sqlite
select
  vm.name,
  nsg.name,
  security_rules
from
  azure_compute_virtual_machine_scale_set_vm as vm,
  json_each(vm.virtual_machine_network_profile) as vm_nic,
  azure_network_security_group as nsg,
  json_each(nsg.network_interfaces) as nsg_int
where
  lower(json_extract(vm_nic.value, '$.id')) = lower(json_extract(nsg_int.value, '$.id'))
  and vm.name = 'warehouse-01';
```

### View power state of virtual machines
Determine the power state of virtual machines in all scale sets.

```sql+postgres
select
  vm.name,
  vm.power_state
from
  azure_compute_virtual_machine_scale_set_vm as vm;
```

```sql+sqlite
select
  vm.name,
  vm.power_state
from
  azure_compute_virtual_machine_scale_set_vm as vm;
```
