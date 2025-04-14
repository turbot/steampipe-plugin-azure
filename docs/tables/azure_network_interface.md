---
title: "Steampipe Table: azure_network_interface - Query Azure Network Interfaces using SQL"
description: "Allows users to query Network Interfaces in Azure, providing detailed information about each network interface, including its associated network security group, IP configuration, and subnet."
folder: "Network"
---

# Table: azure_network_interface - Query Azure Network Interfaces using SQL

A Network Interface in Azure is the interconnection between a Virtual Machine (VM) and the underlying Azure VNet. This interface enables an Azure VM to communicate with internet, Azure, and on-premises resources. Network interfaces consist of one or more IP configurations and a network security group.

## Table Usage Guide

The `azure_network_interface` table provides insights into Network Interfaces within Azure. As an Infrastructure Engineer, explore detailed information about each network interface through this table, including its IP configuration, associated network security group, and subnet. Use this table to manage and optimize your network interface configurations, ensuring seamless communication between your Azure VMs and other resources.

## Examples

### Basic IP address info
Explore the configuration of your Azure network interface to gain insights into your private IP address details. This can help you understand your IP allocation methods and versions, which is useful for managing your network resources effectively.

```sql+postgres
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

```sql+sqlite
select
  name,
  json_extract(ip.value, '$.name') as config_name,
  json_extract(ip.value, '$.properties.privateIPAddress') as private_ip_address,
  json_extract(ip.value, '$.properties.privateIPAddressVersion') as private_ip_address_version,
  json_extract(ip.value, '$.properties.privateIPAllocationMethod') as private_ip_address_allocation_method
from
  azure_network_interface,
  json_each(ip_configurations) as ip;
```

### Find all network interfaces with private IPs that are in a given subnet (10.66.0.0/16)
Determine the areas in which your Azure network interfaces have private IPs within a specific subnet. This is useful for understanding how your network resources are distributed and identifying potential areas of congestion or security vulnerabilities.

```sql+postgres
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

```sql+sqlite
select
  name,
  json_extract(ip.value, '$.name') as config_name,
  json_extract(ip.value, '$.properties.privateIPAddress') as private_ip_address
from
  azure_network_interface,
  json_each(ip_configurations) as ip
where
  json_extract(ip.value, '$.properties.privateIPAddress') = '10.66.0.0/16';
```

### Security groups attached to each network interface
Explore which security groups are linked to each network interface in your Azure environment. This can help in managing and improving the security posture of your network.

```sql+postgres
select
  name,
  split_part(network_security_group_id, '/', 8) as security_groups
from
  azure_network_interface;
```

```sql+sqlite
Error: SQLite does not support split functions.
```