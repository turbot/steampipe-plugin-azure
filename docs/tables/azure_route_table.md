---
title: "Steampipe Table: azure_route_table - Query Azure Route Tables using SQL"
description: "Allows users to query Route Tables in Azure, specifically the information related to routes and subnets, providing insights into the network flow within an Azure Virtual Network."
folder: "Network"
---

# Table: azure_route_table - Query Azure Route Tables using SQL

A Route Table contains a set of rules, called routes, that are used to determine where network traffic is directed. Each subnet in an Azure virtual network is configured with a route table, which can be associated to one or more virtual network subnets. These tables enable you to control the flow of traffic for a subnet.

## Table Usage Guide

The `azure_route_table` table provides insights into Route Tables within Microsoft Azure. As a network administrator, explore route-specific details through this table, including associated subnets, address prefixes, and next hop types. Utilize it to uncover information about network traffic flow, such as the routing of packets, the direction of traffic, and the configuration of subnets.

## Examples

### List of subnets associated with route table
Discover the segments of your network by identifying the subnets associated with a specific route table in your Azure environment. This can help in network management and security by providing insights into the organization of your network infrastructure.

```sql+postgres
select
  name,
  split_part(subnet ->> 'id', '/', 11) subnet,
  region
from
  azure_route_table
  cross join jsonb_array_elements(subnets) as subnet;
```

```sql+sqlite
Error: SQLite does not support split or string_to_array functions.
```

### List of route tables where route propagation is enabled
Determine the areas in which route propagation is active in your Azure Route Table. This is beneficial for understanding network traffic flow and ensuring optimal routing configurations.

```sql+postgres
select
  name,
  disable_bgp_route_propagation,
  region
from
  azure_route_table
where
  not disable_bgp_route_propagation;
```

```sql+sqlite
select
  name,
  disable_bgp_route_propagation,
  region
from
  azure_route_table
where
  disable_bgp_route_propagation = 0;
```

### Route info of each routes table
This query helps users gain insights into the routing information of each route in their Azure network. The practical application of this query is to understand the network flow and the next hop type for each route, which is crucial for network troubleshooting and optimization.

```sql+postgres
select
  name,
  route ->> 'name' route_name,
  route -> 'properties' ->> 'addressPrefix' address_prefix,
  route -> 'properties' ->> 'nextHopType' next_hop_type
from
  azure_route_table
  cross join jsonb_array_elements(routes) as route;
```

```sql+sqlite
select
  name,
  json_extract(route.value, '$.name') as route_name,
  json_extract(route.value, '$.properties.addressPrefix') as address_prefix,
  json_extract(route.value, '$.properties.nextHopType') as next_hop_type
from
  azure_route_table,
  json_each(routes) as route;
```