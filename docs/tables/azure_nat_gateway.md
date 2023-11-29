---
title: "Steampipe Table: azure_nat_gateway - Query Azure NAT Gateways using SQL"
description: "Allows users to query Azure NAT Gateways."
---

# Table: azure_nat_gateway - Query Azure NAT Gateways using SQL

Azure NAT (Network Address Translation) Gateway is a resource that provides outbound internet connectivity for virtual networks. The NAT gateway sends outbound traffic from a virtual network to the internet. It also enables you to configure a static, outbound public IP address, which can be used for the services in your virtual network.

## Table Usage Guide

The 'azure_nat_gateway' table provides insights into NAT Gateways within Azure Networking. As a Network Engineer, explore NAT Gateway-specific details through this table, including subnet details, IP configuration, and associated metadata. Utilize it to uncover information about NAT Gateways, such as those with specific IP configurations, the subnet relationships, and the verification of IP addresses. The schema presents a range of attributes of the NAT Gateway for your analysis, like the NAT Gateway ID, creation date, subnet count, and associated tags.

## Examples

### Basic info
Explore the basic information of your Azure NAT Gateway to understand its provisioning state and type. This can be useful in managing resources and troubleshooting potential issues.

```sql
select
  name,
  id,
  provisioning_state,
  sku_name,
  type
from
  azure_nat_gateway;
```

### List public IP address details for each nat gateway
This query aids in identifying the details of public IP addresses associated with each NAT gateway. It's useful for managing network traffic and ensuring secure and efficient data routing.

```sql
select
  n.name,
  i.ip_address as ip_address,
  i.ip_configuration_id as ip_configuration_id,
  i.public_ip_address_version as public_ip_address_version,
  i.public_ip_allocation_method as public_ip_allocation_method
from
  azure_nat_gateway as n,
  azure_public_ip as i,
  jsonb_array_elements(n.public_ip_addresses) as ip
where
  ip ->> 'id' = i.id;
```

### List subnet details associated with each nat gateway
Analyze the settings to understand the association between each NAT gateway and the related subnet details in your Azure environment. This can be beneficial in managing network topology and ensuring correct routing configurations.

```sql
select
  n.name as name,
  s.name as subnet_name,
  s.virtual_network_name as virtual_network_name
from
  azure_nat_gateway as n,
  azure_subnet as s,
  jsonb_array_elements(n.subnets) as sb
where
  sb ->> 'id' = s.id;
```