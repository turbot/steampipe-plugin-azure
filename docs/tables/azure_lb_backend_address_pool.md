---
title: "Steampipe Table: azure_lb_backend_address_pool - Query Azure Load Balancer Backend Address Pools using SQL"
description: "Allows users to query Azure Load Balancer Backend Address Pools, providing insights into the configuration and status of backend address pools associated with Azure Load Balancers."
---

# Table: azure_lb_backend_address_pool - Query Azure Load Balancer Backend Address Pools using SQL

A Backend Address Pool is a part of Azure Load Balancer, which contains IP addresses served by the load balancer. It is a critical component for defining the resource allocation and traffic distribution in the Azure Load Balancer. It provides a centralized way to manage and distribute network traffic among multiple resources, such as virtual machines.

## Table Usage Guide

The `azure_lb_backend_address_pool` table provides insights into the backend address pools associated with Azure Load Balancers. As a Network Administrator, explore pool-specific details through this table, including associated load balancers, network interfaces, and IP configurations. Utilize it to understand the distribution of network traffic, manage resource allocation, and ensure optimal load balancing across your Azure resources.

## Examples

### Basic info
Explore the status and type of your Azure load balancer's backend address pool to understand its current operational state and configuration. This helps in managing and troubleshooting your network traffic operations effectively.

```sql+postgres
select
  id,
  name,
  load_balancer_name,
  provisioning_state,
  type
from
  azure_lb_backend_address_pool;
```

```sql+sqlite
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
Determine the areas in which Azure load balancer backend address pools have failed to provision. This can be useful for troubleshooting and identifying areas that require attention or modification.

```sql+postgres
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

```sql+sqlite
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