---
title: "Steampipe Table: azure_lb_nat_rule - Query Azure Load Balancer NAT Rules using SQL"
description: "Allows users to query Azure Load Balancer NAT Rules, providing insights into the network traffic routing configurations."
folder: "Load Balancer"
---

# Table: azure_lb_nat_rule - Query Azure Load Balancer NAT Rules using SQL

Azure Load Balancer NAT Rules are part of the Azure Load Balancer service, which ensures high availability and network performance to your applications. NAT Rules are responsible for translating the public IP address and port of a packet to a private IP address and port. They play a crucial role in managing network traffic and routing.

## Table Usage Guide

The `azure_lb_nat_rule` table provides insights into the NAT rules within Azure Load Balancer. As a network engineer, explore NAT rule-specific details through this table, including the associated load balancer, protocol, and ports. Utilize it to uncover information about NAT rules, such as their configuration, associated resources, and the effectiveness of the network routing.

## Examples

### Basic info
Explore which Azure Load Balancer NAT rules are currently in use and assess their provisioning states to ensure optimal performance and resource allocation. This query is particularly useful in managing and troubleshooting network traffic within your Azure environment.

```sql+postgres
select
  id,
  name,
  type,
  provisioning_state,
  etag
from
  azure_lb_nat_rule;
```

```sql+sqlite
select
  id,
  name,
  type,
  provisioning_state,
  etag
from
  azure_lb_nat_rule;
```

### List failed load balancer nat rules
Explore instances where load balancer NAT rules have failed in Azure. This helps in pinpointing areas of concern and aids in troubleshooting the issues for smooth operation.

```sql+postgres
select
  id,
  name,
  type,
  provisioning_state
from
  azure_lb_nat_rule
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
  azure_lb_nat_rule
where
  provisioning_state = 'Failed';
```

### List load balancer nat rules order by idle timeout
Analyze the settings to understand the order of NAT rules based on their idle timeout within a load balancer. This can be useful in optimizing system performance and managing network traffic more efficiently.

```sql+postgres
select
  id,
  name,
  type,
  idle_timeout_in_minutes
from
  azure_lb_nat_rule
order by 
  idle_timeout_in_minutes;
```

```sql+sqlite
select
  id,
  name,
  type,
  idle_timeout_in_minutes
from
  azure_lb_nat_rule
order by 
  idle_timeout_in_minutes;
```