---
title: "Steampipe Table: azure_compute_virtual_machine_scale_set - Query Azure Compute Virtual Machine Scale Sets using SQL"
description: "Allows users to query Azure Compute Virtual Machine Scale Sets."
---

# Table: azure_compute_virtual_machine_scale_set - Query Azure Compute Virtual Machine Scale Sets using SQL

Azure Compute Virtual Machine Scale Sets are a service that allows you to deploy and manage a set of identical, auto-scaling virtual machines. You can scale the number of VMs in the scale set manually, or define rules to auto-scale based on resource usage like CPU, memory demand, or network traffic. An Azure load balancer then distributes network traffic to the VM instances in the scale set.

## Table Usage Guide

The 'azure_compute_virtual_machine_scale_set' table provides insights into Virtual Machine Scale Sets within Azure Compute. As a DevOps engineer, explore scale set-specific details through this table, including scaling configurations, virtual machine profiles, and associated metadata. Utilize it to uncover information about scale sets, such as those with specific scaling policies, the network configurations of the scale sets, and the verification of virtual machine profiles. The schema presents a range of attributes of the Virtual Machine Scale Set for your analysis, like the scale set name, resource group, location, and associated tags.

## Examples

### Basic info
Explore which virtual machine scale sets are located in specific regions and resource groups within your Azure Compute environment. This enables effective management and allocation of resources.

```sql
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
Explore the standard tier virtual machine scale sets within your Azure environment. This is useful for understanding your resource allocation and managing your cloud infrastructure more efficiently.

```sql
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