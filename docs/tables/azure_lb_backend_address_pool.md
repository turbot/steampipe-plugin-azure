---
title: "Steampipe Table: azure_lb_backend_address_pool - Query Azure Load Balancer Backend Address Pools using SQL"
description: "Allows users to query Azure Load Balancer Backend Address Pools"
---

# Table: azure_lb_backend_address_pool - Query Azure Load Balancer Backend Address Pools using SQL

A Backend Address Pool is a part of Azure Load Balancer, which contains IP addresses for the backend servers. Azure Load Balancer distributes inbound flows that arrive at the load balancer's front end to backend pool instances. These flows are according to configured load balancing rules and health probes.

## Table Usage Guide

The 'azure_lb_backend_address_pool' table provides insights into Backend Address Pools within Azure Load Balancer. As a DevOps engineer, explore details specific to Backend Address Pools through this table, including the backend IP configurations, load balancing rules, and associated metadata. Utilize it to uncover information about Backend Address Pools, such as their health probe settings, load balancing rules, and the verification of backend IP configurations. The schema presents a range of attributes of the Backend Address Pool for your analysis, like the name, ID, type, region, and associated tags.

## Examples

### Basic info
Discover the segments that are part of your Azure load balancer's backend address pool. This query can help you assess the elements within your infrastructure, particularly useful in understanding the provisioning state and types of your resources for better resource management.

```sql
select
  id,
  name,
  load_balancer_name,
  provisioning_state,
  type
from
  azure_lb_backend_address_pool;
```

### List failed load balancer backend address pools
Identify instances where load balancer backend address pools in Azure have failed to provision. This can help in troubleshooting and ensuring optimal resource allocation.

```sql
select
  id,
  name,
  type,
  provisioning_state
from
  azure_lb_backend_address_pool
where
  provisioning_state = 'Failed';
```