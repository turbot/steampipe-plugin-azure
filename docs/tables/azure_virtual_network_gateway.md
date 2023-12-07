---
title: "Steampipe Table: azure_virtual_network_gateway - Query Azure Virtual Network Gateways using SQL"
description: "Allows users to query Azure Virtual Network Gateways, specifically providing details about the gateway's configuration, location, and associated resources."
---

# Table: azure_virtual_network_gateway - Query Azure Virtual Network Gateways using SQL

Azure Virtual Network Gateway is a component used to send network traffic between Azure virtual networks and on-premises locations. It works as a specific type of virtual network gateway, designed to send encrypted traffic between an Azure virtual network and an on-premises location over the public Internet. It also enables communication between Azure virtual networks.

## Table Usage Guide

The `azure_virtual_network_gateway` table provides insights into the configuration and status of Azure Virtual Network Gateways. As a network administrator, explore gateway-specific details through this table, including its IP configuration, SKU, and associated virtual network. Utilize it to manage and monitor your network gateways, ensuring secure and efficient communication between your Azure virtual networks and on-premises locations.

## Examples

### Basic info
Explore which Azure Virtual Network Gateways have Border Gateway Protocol (BGP) enabled. This can be useful for network administrators seeking to understand their network's configuration and routing protocols.

```sql+postgres
select
  name,
  id,
  enable_bgp,
  region,
  resource_group
from
  azure_virtual_network_gateway;
```

```sql+sqlite
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
Discover the segments of your Azure virtual network gateways that are not connected to any resources. This can help in identifying unused network gateways, potentially reducing infrastructure costs and improving network management.

```sql+postgres
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

```sql+sqlite
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