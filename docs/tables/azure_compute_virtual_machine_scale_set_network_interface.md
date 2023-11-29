---
title: "Steampipe Table: azure_compute_virtual_machine_scale_set_network_interface - Query Azure Compute Virtual Machine Scale Set Network Interfaces using SQL"
description: "Allows users to query Azure Compute Virtual Machine Scale Set Network Interfaces."
---

# Table: azure_compute_virtual_machine_scale_set_network_interface - Query Azure Compute Virtual Machine Scale Set Network Interfaces using SQL

A Virtual Machine Scale Set Network Interface in Azure is an interconnection between a Virtual Machine Scale Set and a Virtual Network. These network interfaces enable the virtual machines within the scale set to communicate with internet, Azure, and on-premises resources. Network security group rules and route tables can be applied directly to the network interfaces to filter network traffic.

## Table Usage Guide

The 'azure_compute_virtual_machine_scale_set_network_interface' table provides insights into Network Interfaces within Azure Compute Virtual Machine Scale Sets. As a DevOps engineer, explore network interface-specific details through this table, including the IP configuration, network security group association, and subnet details. Utilize it to uncover information about network interfaces, such as their private and public IP addresses, MAC address, and the states of IP forwarding and accelerated networking. The schema presents a range of attributes of the network interfaces for your analysis, like the network interface ID, IP configuration, network security group, and associated tags.

## Examples

### Basic info
Explore the configuration of your Azure virtual machine scale sets by understanding the provisioning state and location. This can be beneficial in managing resources and optimizing your cloud infrastructure.

```sql
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
Explore the configuration of network interfaces that have the IP forwarding rule enabled. This can be useful in identifying network instances that may allow for IP packet forwarding, which can be critical for understanding network traffic flow and potential security implications.

```sql
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

### List network interfaces with accelerated networking enabled
Explore which network interfaces have the accelerated networking feature enabled. This can be particularly useful for identifying areas where network performance can be improved.

```sql
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

### Get scale set virtual machine details for scale set network interface
Analyze the details of virtual machine scale sets to understand the associated network interfaces. This is beneficial in managing the configuration and performance of your network resources in a large-scale cloud environment.

```sql
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