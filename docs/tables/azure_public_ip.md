---
title: "Steampipe Table: azure_public_ip - Query Azure Public IP Addresses using SQL"
description: "Allows users to query Azure Public IP Addresses."
---

# Table: azure_public_ip - Query Azure Public IP Addresses using SQL

Azure Public IP Address is a resource within Microsoft Azure that allows you to assign public IP addresses to Azure resources such as virtual machines, Azure Load Balancers, and Azure VPN Gateways. These public IP addresses are used to communicate with internet resources, and can be either dynamic or static. Azure Public IP Addresses provide a reliable and secure connection to the internet for your Azure resources.

## Table Usage Guide

The 'azure_public_ip' table provides insights into Public IP Addresses within Microsoft Azure. As a Network Administrator, explore IP-specific details through this table, including the IP version, IP configuration, and associated metadata. Utilize it to uncover information about IP addresses, such as their allocation method, their assigned resource, and their location. The schema presents a range of attributes of the Public IP Address for your analysis, like the IP address, the SKU name, the domain name label, and the reverse FQDN.

## Examples

### List of unassociated elastic IPs
Explore which Azure public IP addresses are not associated with any IP configuration. This is useful to identify any unused resources that could potentially be costing you money.

```sql
select
  name,
	ip_configuration_id
from
  azure_public_ip
where
  ip_configuration_id is null;
```

### List of IP addresses with corresponding associations
Explore which IP addresses are associated with specific resources in your Azure environment. This can help you manage your network configuration and identify potential issues or inefficiencies.

```sql
select
  name,
  ip_address,
  split_part(ip_configuration_id, '/', 8) as resource,
  split_part(ip_configuration_id, '/', 9) as resource_name
from
  azure_public_ip;
```

### List of dynamic IP addresses
Determine the areas in which Azure's public IP addresses are dynamically allocated to gain insights into the flexibility and scalability of your network resources.

```sql
select
  name,
  public_ip_allocation_method
from
  azure_public_ip
where
  public_ip_allocation_method = 'Dynamic';
```