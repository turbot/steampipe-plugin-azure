---
title: "Steampipe Table: azure_virtual_network - Query Azure Virtual Networks using SQL"
description: "Allows users to query Azure Virtual Networks."
---

# Table: azure_virtual_network - Query Azure Virtual Networks using SQL

Azure Virtual Networks (VNet) is a fundamental building block for your private network in Azure. VNet enables many types of Azure resources, such as Azure Virtual Machines (VM), to securely communicate with each other, the internet, and on-premises networks. VNet is similar to a traditional network that you'd operate in your own data center, but brings with it additional benefits of Azure's infrastructure, such as scale, availability, and isolation.

## Table Usage Guide

The 'azure_virtual_network' table provides insights into Virtual Networks within Azure. As a DevOps engineer, explore network-specific details through this table, including address spaces, DNS servers, and associated subnets. Utilize it to uncover information about networks, such as those with certain security rules, the associated subnets, and the verification of DNS servers. The schema presents a range of attributes of the Virtual Network for your analysis, like the network ID, creation date, associated subnets, and associated tags.

## Examples

### List of virtual networks where DDoS(Distributed Denial of Service attacks) Protection is not enabled
Explore the virtual networks that lack protection against Distributed Denial of Service (DDoS) attacks. This allows for the identification of potential network vulnerabilities and aids in strengthening security measures.

```sql
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

### CIDR list for each virtual network
Explore which address blocks are associated with each virtual network in Azure. This can help you understand the network structure and manage IP address allocation efficiently.

```sql
select
  name,
  jsonb_array_elements_text(address_prefixes) as address_block
from
  azure_virtual_network;
```

### List VPCs with public CIDR blocks
Determine the areas in which Azure Virtual Networks are configured with public CIDR blocks, allowing you to assess potential exposure to the internet and take necessary security measures.

```sql
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


### Subnet details associated with the virtual network
Explore the configuration of your virtual network to understand the details of associated subnets. This can help in managing network policies, service endpoints, and routing tables efficiently.

```sql
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