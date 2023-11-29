---
title: "Steampipe Table: azure_compute_virtual_machine_scale_set_vm - Query Azure Compute Virtual Machine Scale Sets using SQL"
description: "Allows users to query Azure Compute Virtual Machine Scale Sets"
---

# Table: azure_compute_virtual_machine_scale_set_vm - Query Azure Compute Virtual Machine Scale Sets using SQL

Azure Compute is a cloud computing service that provides on-demand, high-scale compute capacity for applications and workloads. One of its resources, Virtual Machine Scale Sets, allows for the creation, management, and scaling of a set of identical, load-balanced VMs. This service is ideal for building large-scale services, such as big data, containerized applications, and distributed systems.

## Table Usage Guide

The 'azure_compute_virtual_machine_scale_set_vm' table provides insights into Virtual Machine Scale Sets within Azure Compute. As a DevOps engineer, explore specific details through this table, including the status, location, and configuration of each VM in the scale set. Utilize it to uncover information about VMs, such as their operating system, network profile, and associated metadata. The schema presents a range of attributes of the VM for your analysis, like the VM ID, instance ID, virtual network, and associated tags.

## Examples

### Basic info
Explore which virtual machines are part of your Azure scale set to manage resources effectively. This can help in identifying instances where resources are underutilized or overprovisioned, ensuring optimal resource allocation and cost management.

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

### List Standard tier scale set virtual machine
Explore which scale set virtual machines operate on the 'Standard' tier. This query is useful for understanding the distribution and usage of different tiered resources within your Azure environment.

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
Explore which virtual machines are part of a specific set to understand the scale and tier of your Azure computing resources. This aids in resource management and capacity planning.

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

### View Network Security Group Rules for a virtual machine
Determine the security rules applied to a specific virtual machine within your network. This is useful for assessing the security measures in place and identifying any potential vulnerabilities.

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