---
title: "Steampipe Table: azure_virtual_network - Query Azure Virtual Networks using SQL"
description: "Allows users to query Azure Virtual Networks, specifically providing insights into the configuration and status of each virtual network within an Azure subscription."
folder: "Network"
---

# Table: azure_virtual_network - Query Azure Virtual Networks using SQL

Azure Virtual Networks is a fundamental building block for your private network in Azure. It enables many types of Azure resources, such as Azure Virtual Machines (VM), to securely communicate with each other, the internet, and on-premises networks. Azure virtual network is similar to a traditional network that you'd operate in your own data center but brings with it additional benefits of Azure's infrastructure such as scale, availability, and isolation.

## Table Usage Guide

The `azure_virtual_network` table provides detailed information about each virtual network within an Azure subscription. As a network administrator or cloud architect, you can use this table to gather data about the subnets, IP address ranges, and connected devices within each virtual network. This information can be used to monitor network usage, plan for capacity, and ensure the network is correctly configured for your applications' requirements.

## Examples

### List of virtual networks where DDoS(Distributed Denial of Service attacks) Protection is not enabled
Discover the segments of your virtual networks that are potentially vulnerable to Distributed Denial of Service (DDoS) attacks, as they do not have DDoS protection enabled. This information can help prioritize areas for security enhancement and risk mitigation.

```sql+postgres
select
  name,
  enable_ddos_protection,
  region,
  resource_group
from
  azure_virtual_network
where
  not enable_ddos_protection;
```

```sql+sqlite
select
  name,
  enable_ddos_protection,
  region,
  resource_group
from
  azure_virtual_network
where
  enable_ddos_protection is not 1;
```

### CIDR list for each virtual network
Determine the areas in which your Azure virtual networks operate by identifying their respective address blocks. This can help in network planning and management by providing a clear view of the network's structure and usage.

```sql+postgres
select
  name,
  jsonb_array_elements_text(address_prefixes) as address_block
from
  azure_virtual_network;
```

```sql+sqlite
select
  name,
  json_each.value as address_block
from
  azure_virtual_network,
  json_each(azure_virtual_network.address_prefixes);
```

### List VPCs with public CIDR blocks
Determine the areas in which Virtual Private Networks (VPCs) have public CIDR blocks, allowing you to assess network accessibility and security risks. This is particularly useful in identifying potential exposure of your Azure virtual networks to the public internet.

```sql+postgres
select
  name,
  cidr_block,
  region,
  resource_group
from
  azure_virtual_network
  cross join jsonb_array_elements_text(address_prefixes) as cidr_block
where
  not cidr_block :: cidr = '10.0.0.0/16'
  and not cidr_block :: cidr = '192.168.0.0/16'
  and not cidr_block :: cidr = '172.16.0.0/12';
```

```sql+sqlite
Error: SQLite does not support CIDR operations.
```


### Subnet details associated with the virtual network
Determine the areas in which subnets interact with your virtual network. This query helps to analyze the configuration of these subnets, providing insights into their address prefixes, network policies, service endpoints, and route tables, which can be useful for network management and troubleshooting.

```sql+postgres
select
  name,
  subnet ->> 'name' as subnet_name,
  subnet -> 'properties' ->> 'addressPrefix' as address_prefix,
  subnet -> 'properties' ->> 'privateEndpointNetworkPolicies' as private_endpoint_network_policies,
  subnet -> 'properties' ->> 'privateLinkServiceNetworkPolicies' as private_link_service_network_policies,
  subnet -> 'properties' ->> 'serviceEndpoints' as service_endpoints,
  split_part(subnet -> 'properties' ->> 'routeTable', '/', 9) as route_table
from
  azure_virtual_network
  cross join jsonb_array_elements(subnets) as subnet;
```

```sql+sqlite
Error: SQLite does not support split_part function.
```