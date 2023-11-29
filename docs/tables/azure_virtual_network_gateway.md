---
title: "Steampipe Table: azure_virtual_network_gateway - Query Azure Virtual Network Gateways using SQL"
description: "Allows users to query Azure Virtual Network Gateways"
---

# Table: azure_virtual_network_gateway - Query Azure Virtual Network Gateways using SQL

Azure Virtual Network Gateway is a component that provides a point-to-point network connection from an Azure virtual network to an on-premises location over the public internet. It can also be used to send encrypted traffic between an Azure virtual network and an on-premises location over a VPN tunnel, or to route traffic between virtual networks.

## Table Usage Guide

The 'azure_virtual_network_gateway' table provides insights into Virtual Network Gateways within Azure. As a DevOps engineer, explore gateway-specific details through this table, including gateway type, VPN type, and associated metadata. Utilize it to uncover information about gateways, such as their active-active status, the private IP allocated to the gateway, and the verification of gateway SKU. The schema presents a range of attributes of the Virtual Network Gateway for your analysis, like the gateway name, resource group, subscription ID, and associated tags.

## Examples

### Basic info
Explore the configuration of your Azure Virtual Network Gateway to gain insights into settings such as BGP status and regional distribution. This can be useful in assessing network performance and identifying areas for optimization.

```sql
select
  name,
  id,
  enable_bgp,
  region,
  resource_group
from
  azure_virtual_network_gateway;
```

### List network gateways with no connections
Determine the areas in which network gateways are not connected in your Azure virtual network. This can help identify potential network vulnerabilities or inefficiencies.

```sql
select
  name,
  id,
  enable_bgp,
  region,
  resource_group
from
  azure_virtual_network_gateway
where
   gateway_connections is null;
```