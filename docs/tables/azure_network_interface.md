---
title: "Steampipe Table: azure_network_interface - Query Azure Network Interfaces using SQL"
description: "Allows users to query Azure Network Interfaces."
---

# Table: azure_network_interface - Query Azure Network Interfaces using SQL

An Azure Network Interface is a virtual network interface card (NIC) in Azure that is attached to a virtual machine (VM). It enables Azure VMs to communicate with internet, Azure, and on-premises resources. Network Interfaces can include IP addresses, subnets, and network security groups.

## Table Usage Guide

The 'azure_network_interface' table offers insights into Network Interfaces within Azure. As a DevOps engineer, you can delve into interface-specific details via this table, including private and public IP addresses, network security group associations, and subnet information. Utilize it to uncover information about interfaces, such as their IP configurations, DNS settings, and associated subnets. The schema presents a range of attributes of the Network Interface for your analysis, like the interface ID, IP configurations, associated network security groups, and subnet details.

## Examples

### Basic IP address info
Explore network configurations by identifying the private IP addresses, their versions, and allocation methods in Azure. This can be beneficial in understanding the distribution and management of IP addresses within your Azure network interface.

```sql
select
  name,
  ip ->> 'name' as config_name,
  ip -> 'properties' ->> 'privateIPAddress' as private_ip_address,
  ip -> 'properties' ->> 'privateIPAddressVersion' as private_ip_address_version,
  ip -> 'properties' ->> 'privateIPAllocationMethod' as private_ip_address_allocation_method
from
  azure_network_interface
  cross join jsonb_array_elements(ip_configurations) as ip;
```

### Find all network interfaces with private IPs that are in a given subnet (10.66.0.0/16)
This query is useful for pinpointing specific network interfaces within a designated subnet that are utilizing private IP addresses. This can aid in network management and security by identifying potential areas of vulnerability or inefficiency.

```sql
select
  name,
  ip ->> 'name' as config_name,
  ip -> 'properties' ->> 'privateIPAddress' as private_ip_address
from
  azure_network_interface
  cross join jsonb_array_elements(ip_configurations) as ip
where
  ip -> 'properties' ->> 'privateIPAddress' = '10.66.0.0/16';
```

### Security groups attached to each network interface
Analyze the settings to understand the security groups linked with each network interface in your Azure network. This can be useful for assessing your network's security configuration and identifying potential vulnerabilities.

```sql
select
  name,
  split_part(network_security_group_id, '/', 8) as security_groups
from
  azure_network_interface;
```