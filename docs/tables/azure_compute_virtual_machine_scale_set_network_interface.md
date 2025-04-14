---
title: "Steampipe Table: azure_compute_virtual_machine_scale_set_network_interface - Query Azure Compute Virtual Machine Scale Sets Network Interfaces using SQL"
description: "Allows users to query Azure Compute Virtual Machine Scale Sets Network Interfaces, providing detailed information about the network interfaces of each scale set."
folder: "Compute"
---

# Table: azure_compute_virtual_machine_scale_set_network_interface - Query Azure Compute Virtual Machine Scale Sets Network Interfaces using SQL

A Network Interface within Azure Compute Virtual Machine Scale Sets is a virtual network interface card (NIC) attached to a Virtual Machine Scale Set in Azure. It provides the interconnection between a Virtual Machine Scale Set and the underlying Azure virtual network. Each Network Interface can have one or more IP configurations associated with it.

## Table Usage Guide

The `azure_compute_virtual_machine_scale_set_network_interface` table provides insights into the network interfaces associated with Azure Compute Virtual Machine Scale Sets. As a network administrator, you can use this table to explore details about each network interface, including its IP configurations, subnet information, and associated scale set. This can be particularly useful for managing network connectivity and troubleshooting network-related issues within your Azure Compute Virtual Machine Scale Sets.

## Examples

### Basic info
Explore the status and location of your Azure virtual machine scale sets to gain insights into their deployment and management. This is useful for assessing the distribution and provisioning of your resources across different regions and groups.

```sql+postgres
select
  name,
  id,
  scale_set_name,
  provisioning_state,
  region,
  resource_group
from
  azure_compute_virtual_machine_scale_set_network_interface;
```

```sql+sqlite
select
  name,
  id,
  scale_set_name,
  provisioning_state,
  region,
  resource_group
from
  azure_compute_virtual_machine_scale_set_network_interface;
```

### List network interfaces with IP forwarding rule enabled
Determine the areas in which IP forwarding is enabled on network interfaces for better control and management of network traffic. This is particularly useful in scenarios where data packets need to be redirected or rerouted for specific purposes.

```sql+postgres
select
  name,
  id,
  enable_ip_forwarding,
  scale_set_name
from
  azure_compute_virtual_machine_scale_set_network_interface
where
  enable_ip_forwarding;
```

```sql+sqlite
select
  name,
  id,
  enable_ip_forwarding,
  scale_set_name
from
  azure_compute_virtual_machine_scale_set_network_interface
where
  enable_ip_forwarding = 1;
```

### List network interfaces with accelerated networking enabled
Explore which network interfaces are utilizing accelerated networking within your Azure virtual machine scale sets. This information can be useful for optimizing network performance and troubleshooting connectivity issues.

```sql+postgres
select
  name,
  id,
  enable_accelerated_networking,
  scale_set_name
from
  azure_compute_virtual_machine_scale_set_network_interface
where
  enable_accelerated_networking;
```

```sql+sqlite
select
  name,
  id,
  enable_accelerated_networking,
  scale_set_name
from
  azure_compute_virtual_machine_scale_set_network_interface
where
  enable_accelerated_networking = 1;
```

### Get scale set virtual machine details for scale set network interface
This query helps to map network interfaces to their corresponding virtual machines within a specified scale set. It's particularly useful for managing and monitoring network traffic and performance across multiple instances within a scale set.

```sql+postgres
select
  i.name as name,
  i.id as id,
  v.instance_id as instance_id,
  v.scale_set_name as scale_set_name,
  v.sku_name as vm_sku_name
from
  azure_compute_virtual_machine_scale_set_network_interface as i,
  azure_compute_virtual_machine_scale_set_vm as v
where
  i.virtual_machine ->> 'id' = v.id;
```

```sql+sqlite
select
  i.name as name,
  i.id as id,
  v.instance_id as instance_id,
  v.scale_set_name as scale_set_name,
  v.sku_name as vm_sku_name
from
  azure_compute_virtual_machine_scale_set_network_interface as i,
  azure_compute_virtual_machine_scale_set_vm as v
where
  json_extract(i.virtual_machine, '$.id') = v.id;
```