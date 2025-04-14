---
title: "Steampipe Table: azure_compute_disk_access - Query Azure Compute Disk Accesses using SQL"
description: "Allows users to query Azure Compute Disk Accesses, providing detailed information about access configurations and their related resources."
folder: "Compute"
---

# Table: azure_compute_disk_access - Query Azure Compute Disk Accesses using SQL

Azure Compute Disk Access is a feature within Microsoft Azure that enables granular access control to managed disks. It provides a secure way to authorize specific virtual machines to access specific managed disks. Azure Compute Disk Access enhances the security and management of your Azure resources by controlling access at the disk level.

## Table Usage Guide

The `azure_compute_disk_access` table provides insights into disk access configurations within Azure Compute. As a Security Analyst, explore disk access-specific details through this table, including access locations, permissions, and associated virtual machines. Utilize it to uncover information about disk accesses, such as those with specific permissions, the relationships between disk accesses and virtual machines, and the verification of access policies.

## Examples

### Basic info
Explore the fundamental details of your Azure disk access resources to understand their status and organization. This can help in managing resources and ensuring optimal utilization.

```sql+postgres
select
  name,
  id,
  type,
  provisioning_state,
  resource_group
from
  azure_compute_disk_access;
```

```sql+sqlite
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
Explore which disk accesses in your Azure Compute resource have failed. This is beneficial for identifying potential issues with your resources and taking necessary corrective actions.

```sql+postgres
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

```sql+sqlite
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