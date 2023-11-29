---
title: "Steampipe Table: azure_compute_disk_access - Query Azure Compute Disk Accesses using SQL"
description: "Allows users to query Azure Compute Disk Accesses"
---

# Table: azure_compute_disk_access - Query Azure Compute Disk Accesses using SQL

Azure Compute Disk Access is a feature within Microsoft Azure that enables and controls access to managed disks, snapshots, and images. It provides a secure way to grant permissions to read or write data from these resources. Disk Access resources are Azure Resource Manager resources that can be created and managed just like other Azure resources.

## Table Usage Guide

The 'azure_compute_disk_access' table provides insights into Disk Accesses within Azure Compute. As a DevOps engineer, explore specific details through this table, including the network access policy, disk encryption set ID, and associated metadata. Utilize it to uncover information about disk accesses, such as those with unrestricted network access, the associated disk encryption sets, and the verification of network access policies. The schema presents a range of attributes of the Disk Access for your analysis, like the resource ID, creation date, provisioning state, and associated tags.

## Examples

### Basic info
Explore the basic details of your Azure compute disk access to understand its state and group allocation. This can help you manage and optimize your resources effectively.

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
Explore which disk accesses in your Azure Compute resource have failed. This is useful for diagnosing system issues and ensuring optimal performance of your resources.

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