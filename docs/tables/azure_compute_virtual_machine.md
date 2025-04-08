---
title: "Steampipe Table: azure_compute_virtual_machine - Query Azure Compute Virtual Machines using SQL"
description: "Allows users to query Azure Compute Virtual Machines, providing detailed information about the configuration, status, and other operational aspects of each virtual machine."
folder: "Networking"
---

# Table: azure_compute_virtual_machine - Query Azure Compute Virtual Machines using SQL

Azure Compute is a service within Microsoft Azure that allows you to run applications on virtual machines in the cloud. It provides scalable, on-demand compute capacity in the cloud and lets you create and manage virtual machines to run applications. Azure Compute supports a range of operating systems, languages, tools, and frameworks.

## Table Usage Guide

The `azure_compute_virtual_machine` table provides insights into the virtual machines within Azure Compute. As a system administrator, you can explore detailed information about each virtual machine, including its configuration, status, and operational aspects. Utilize this table to manage and monitor your virtual machines effectively, ensuring optimal performance and resource usage.

## Examples

### Virtual machine configuration info
Explore the configuration of your virtual machines to gain insights into their power state, IP addresses, size, operating system, and image details. This can help in managing resources and ensuring optimal performance of your virtual machines.

```sql+postgres
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

```sql+sqlite
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
Analyze the distribution of virtual machines across different regions. This information can be useful for understanding your infrastructure's geographical spread and planning resource allocation.

```sql+postgres
select
  region,
  count(name)
from
  azure_compute_virtual_machine
group by
  region;
```

```sql+sqlite
select
  region,
  count(name)
from
  azure_compute_virtual_machine
group by
  region;
```

### List of VMs whose OS disk is not encrypted by customer managed key
Determine the areas in which virtual machines are potentially vulnerable due to their operating system disk not being encrypted by a customer-managed key. This query is useful in identifying security risks and enhancing data protection measures.

```sql+postgres
select
  vm.name,
  disk.encryption_type
from
  azure_compute_disk as disk
  join azure_compute_virtual_machine as vm on disk.name = vm.os_disk_name
where
  not disk.encryption_type = 'EncryptionAtRestWithCustomerKey';
```

```sql+sqlite
select
  vm.name,
  disk.encryption_type
from
  azure_compute_disk as disk
  join azure_compute_virtual_machine as vm on disk.name = vm.os_disk_name
where
  disk.encryption_type != 'EncryptionAtRestWithCustomerKey';
```

### List of VMs provisioned with undesired(for example Standard_D8s_v3 and Standard_DS3_v3 is desired) sizes.
Determine the areas in which virtual machines have been provisioned with non-standard sizes. This is useful for identifying potential inefficiencies or areas for optimization in your Azure Compute resources.

```sql+postgres
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

```sql+sqlite
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
Explore the relationship between virtual machines and their respective availability sets in Azure, including fault domain count, update domain count and SKU name. This can be beneficial for understanding the resilience and update strategy of your virtual machines.

```sql+postgres
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

```sql+sqlite
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
Explore which virtual machines are of the spot type and understand their eviction policies. This can be useful in managing costs and resource allocation in an Azure environment.

```sql+postgres
select
  name,
  vm_id,
  eviction_policy
from
  azure_compute_virtual_machine
where
  priority = 'Spot';
```

```sql+sqlite
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
This query is useful to gain insights into the disk storage usage across all virtual machines in an Azure environment. It helps in managing and optimizing storage resources by providing a summary of the number and total size of disks used by each virtual machine.

```sql+postgres
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

```sql+sqlite
select
  vm.name,
  count(d.disk_size_gb) as num_disks,
  sum(d.disk_size_gb) as total_disk_size_gb
from
  azure_compute_virtual_machine as vm
  left join azure_compute_disk as d on lower(vm.id) = lower(d.managed_by)
group by
  vm.name
order by
  vm.name;
```

### View Network Security Group Rules for a VM
Explore the security rules applied to a specific virtual machine within your network. This can be useful for auditing security configurations and identifying potential vulnerabilities.

```sql+postgres
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

```sql+sqlite
select
  vm.name,
  nsg.name,
  security_rules
from
  azure_compute_virtual_machine as vm,
  json_each(vm.network_interfaces) as vm_nic,
  azure_network_security_group as nsg,
  json_each(nsg.network_interfaces) as nsg_int
where
  lower(json_extract(vm_nic.value, '$.id')) = lower(json_extract(nsg_int.value, '$.id'))
  and vm.name = 'warehouse-01';
```

### List virtual machines with user assigned identities
This example helps you identify the virtual machines in your Azure environment that are configured with user-assigned identities. This is useful for understanding your identity management practices, specifically in scenarios where you want to delegate permissions to resources in your Azure environment.

```sql+postgres
select
  name,
  identity -> 'type' as identity_type,
  jsonb_pretty(identity -> 'userAssignedIdentities') as identity_user_assignedidentities
from
  azure_compute_virtual_machine
where
    exists (
      select
      from
        unnest(regexp_split_to_array(identity ->> 'type', ',')) elem
      where
        trim(elem) = 'UserAssigned'
  );
```

```sql+sqlite
select
  name,
  json_extract(identity, '$.type') as identity_type,
  identity_user_assignedidentities
from
  azure_compute_virtual_machine
where
  instr(identity_type, 'UserAssigned') > 0;
```

### List security profile details
Determine the areas in which encryption is being used at host level within Azure's virtual machines. This can be useful for assessing security measures and identifying potential vulnerabilities.

```sql+postgres
select
  name,
  vm_id,
  security_profile -> 'encryptionAtHost' as encryption_at_host
from
  azure_compute_virtual_machine;
```

```sql+sqlite
select
  name,
  vm_id,
  json_extract(security_profile, '$.encryptionAtHost') as encryption_at_host
from
  azure_compute_virtual_machine;
```