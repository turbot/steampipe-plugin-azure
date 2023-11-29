---
title: "Steampipe Table: azure_route_table - Query Azure Networking Route Tables using SQL"
description: "Allows users to query Azure Networking Route Tables."
---

# Table: azure_route_table - Query Azure Networking Route Tables using SQL

Azure Networking Route Tables are a resource within Microsoft Azure that allow for control over the routing paths of network traffic. They provide a way to direct network traffic based on source, destination, and other key factors. This enables more granular control over network communication within and across Azure services.

## Table Usage Guide

The 'azure_route_table' table provides insights into Route Tables within Azure Networking. As a Network Administrator, explore route-specific details through this table, including associated routes, subnets, and related metadata. Utilize it to uncover information about the routing paths, such as those with specific next hops, the association between subnets and routes, and the verification of route properties. The schema presents a range of attributes of the Route Table for your analysis, like the route table ID, creation date, attached subnets, and associated tags.

## Examples

### List of subnets associated with route table
Explore the association between subnets and route tables within a specific region in Azure. This can help in understanding the network infrastructure and identifying potential issues related to network routing.

```sql
select
  name,
  split_part(subnet ->> 'id', '/', 11) subnet,
  region
from
  azure_route_table
  cross join jsonb_array_elements(subnets) as subnet;
```

### List of route tables where route propagation is enabled
Explore the route tables in your Azure network where route propagation is enabled. This can be useful in understanding how your network traffic is being directed and managed.

```sql
select
  name,
  disable_bgp_route_propagation,
  region
from
  azure_route_table
where
  not disable_bgp_route_propagation;
```

### Route info of each routes table
Explore the details of each route within your Azure network to understand the direction of traffic flow. This can help in optimizing network performance and managing traffic effectively.

```sql
select
  name,
  route ->> 'name' route_name,
  route -> 'properties' ->> 'addressPrefix' address_prefix,
  route -> 'properties' ->> 'nextHopType' next_hop_type
from
  azure_route_table
  cross join jsonb_array_elements(routes) as route;
```