---
title: "Steampipe Table: azure_compute_virtual_machine - Query Azure Compute Virtual Machines using SQL"
description: "Allows users to query Azure Compute Virtual Machines."
---

# Table: azure_compute_virtual_machine - Query Azure Compute Virtual Machines using SQL

Azure Compute is a service within Microsoft Azure that allows you to deploy and manage virtual machines (VMs). It provides the flexibility of virtualization for a wide range of computing solutions—development and testing, running applications, and extending your datacenter. Azure Virtual Machines provide on-demand, high-scale, secure, virtualized infrastructure using Windows servers or Linux servers.

## Table Usage Guide

The 'azure_compute_virtual_machine' table provides insights into Virtual Machines within Azure Compute. As a DevOps engineer, explore VM-specific details through this table, including VM sizes, operating systems, network interfaces, and associated metadata. Utilize it to uncover information about VMs, such as their power states, the virtual networks they are associated with, and the disks they use. The schema presents a range of attributes of the VM for your analysis, like the VM ID, creation date, location, and associated tags.

## Examples

### Virtual machine configuration info
Analyze the settings to understand the configuration and status of your virtual machines in Azure. This can assist in managing machine resources, tracking machine states, and ensuring optimal utilization of your Azure cloud resources.

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
Gain insights into the distribution of virtual machines across different regions. This helps in understanding resource allocation and planning for capacity management.

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
Discover the segments that include virtual machines (VMs) where the operating system disk is not encrypted using a customer-managed key. This can be useful for identifying potential security risks and ensuring compliance with data protection policies.

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
Explore which virtual machines have been provisioned with sizes other than the desired ones. This is useful for identifying potential inefficiencies or mismatches in resource allocation.

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
Explore which virtual machines are part of a specific availability set in Azure. This can help you understand how your VMs are distributed across fault and update domains, allowing for better management of redundancy and availability.

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
Explore the comprehensive list of all spot type Virtual Machines and their corresponding eviction policies. This information can be used to understand and manage resource allocation and cost-efficiency in your Azure cloud environment.

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
Explore the disk storage usage across different virtual machines in your Azure environment. This helps in managing resources and planning for storage needs more effectively.

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
Discover the security rules applied to a specific virtual machine in your Azure network. This query is useful for understanding the security parameters and restrictions currently in place for a given machine.

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

### List virtual machines with user assigned identities
Explore which virtual machines have user assigned identities. This can be beneficial for managing access control and ensuring secure operations in your Azure environment.

```sql
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

### List security profile details
Explore the security profiles of your virtual machines in Azure to understand if the 'encryption at host' setting is enabled. This can aid in assessing your data security and compliance.

```sql
select
  name,
  vm_id,
  security_profile -> 'encryptionAtHost' as encryption_at_host
from
  azure_compute_virtual_machine;
```