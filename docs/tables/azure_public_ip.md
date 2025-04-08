---
title: "Steampipe Table: azure_public_ip - Query Azure Public IPs using SQL"
description: "Allows users to query Public IPs in Azure, specifically the allocated IP address, providing insights into IP configurations and potential network anomalies."
folder: "Networking"
---

# Table: azure_public_ip - Query Azure Public IPs using SQL

Azure Public IP is a service in Microsoft Azure that allows you to allocate a public IP address to Azure resources such as virtual machines, load balancers, and VPN gateways. It provides a way to communicate with the internet, a private network, or both. Azure Public IP helps you manage network connectivity and access, ensuring that your Azure resources are reachable and responsive.

## Table Usage Guide

The `azure_public_ip` table provides insights into Public IPs within Microsoft Azure. As a network administrator, explore IP-specific details through this table, including IP address, allocation method, and associated metadata. Utilize it to uncover information about network configurations, such as those with static or dynamic allocation, the IP version (IPv4 or IPv6), and the verification of IP tags.

## Examples

### List of unassociated elastic IPs
Discover the segments that consist of unassigned public IPs in your Azure infrastructure. This is useful in identifying potential cost savings, as you may be billed for these unassociated resources.

```sql+postgres
select
  name,
	ip_configuration_id
from
  azure_public_ip
where
  ip_configuration_id is null;
```

```sql+sqlite
select
  name,
  ip_configuration_id
from
  azure_public_ip
where
  ip_configuration_id is null;
```

### List of IP addresses with corresponding associations
Explore the relationships between various IP addresses and their corresponding resources in your Azure environment. This can aid in managing network configurations and identifying potential issues.

```sql+postgres
select
  name,
  ip_address,
  split_part(ip_configuration_id, '/', 8) as resource,
  split_part(ip_configuration_id, '/', 9) as resource_name
from
  azure_public_ip;
```

```sql+sqlite
Error: SQLite does not support split or string_to_array functions.
```

### List of dynamic IP addresses
Discover the segments that utilize dynamic IP allocation in your Azure environment. This helps in understanding the networking configuration and managing resources effectively.

```sql+postgres
select
  name,
  public_ip_allocation_method
from
  azure_public_ip
where
  public_ip_allocation_method = 'Dynamic';
```

```sql+sqlite
select
  name,
  public_ip_allocation_method
from
  azure_public_ip
where
  public_ip_allocation_method = 'Dynamic';
```