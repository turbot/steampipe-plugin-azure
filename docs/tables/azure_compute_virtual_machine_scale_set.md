---
title: "Steampipe Table: azure_compute_virtual_machine_scale_set - Query Azure Compute Virtual Machine Scale Sets using SQL"
description: "Allows users to query Azure Compute Virtual Machine Scale Sets, specifically providing details about the scale set configuration, capacity, and status."
folder: "Resource"
---

# Table: azure_compute_virtual_machine_scale_set - Query Azure Compute Virtual Machine Scale Sets using SQL

Azure Compute Virtual Machine Scale Sets are a service within Microsoft Azure that allows you to create and manage a group of identical, load balanced VMs. They enable you to centrally manage, configure, and update a large number of VMs in minutes to provide highly available applications. The scale set adjusts the number of VMs in response to demand or a defined schedule.

## Table Usage Guide

The `azure_compute_virtual_machine_scale_set` table provides insights into Azure Compute Virtual Machine Scale Sets within Microsoft Azure. As a system administrator or DevOps engineer, explore scale set-specific details through this table, including configuration, capacity, and status. Utilize it to uncover information about scale sets, such as their current capacity, configuration details, and overall status, aiding in efficient management and monitoring of your virtual machine resources.

## Examples

### Basic info
Explore the configuration of your virtual machine scale sets in Azure to identify their associated regions and resource groups. This can help you manage and organize your resources more efficiently.

```sql+postgres
select
  name,
  id,
  identity,
  region,
  resource_group
from
  azure_compute_virtual_machine_scale_set;
```

```sql+sqlite
select
  name,
  id,
  identity,
  region,
  resource_group
from
  azure_compute_virtual_machine_scale_set;
```

### List Standard tier virtual machine scale set
Determine the areas in which standard-tier virtual machine scale sets are being used within your Azure Compute environment. This query helps to understand resource allocation and cost management.

```sql+postgres
select
  name,
  id,
  sku_name,
  sku_tier
from
  azure_compute_virtual_machine_scale_set
where
  sku_tier = 'Standard';
```

```sql+sqlite
select
  name,
  id,
  sku_name,
  sku_tier
from
  azure_compute_virtual_machine_scale_set
where
  sku_tier = 'Standard';
```