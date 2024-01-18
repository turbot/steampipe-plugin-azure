---
title: "Steampipe Table: azure_lb_rule - Query Azure Load Balancer Rules using SQL"
description: "Allows users to query Azure Load Balancer Rules, providing insights into the rules defined for load balancing traffic."
---

# Table: azure_lb_rule - Query Azure Load Balancer Rules using SQL

An Azure Load Balancer Rule is a rule within Microsoft Azure that determines how network traffic is distributed among service endpoints in a load balancer. It provides a way to direct traffic based on a variety of parameters, including port, protocol, and IP address. These rules play a crucial role in ensuring the smooth operation and scalability of applications hosted on Azure.

## Table Usage Guide

The `azure_lb_rule` table provides insights into the rules defined within Azure Load Balancer. As a network engineer or a system administrator, you can explore rule-specific details through this table, including the load balancing algorithm used, health probe settings, and associated metadata. Utilize it to uncover information about the rules, such as those directing traffic to specific ports or using certain protocols, and to verify their configurations.

## Examples

### Basic info
Explore the configuration and status of Azure load balancer rules. This aids in understanding the type and provisioning state of each rule, which can help in managing and troubleshooting your Azure load balancer setup.

```sql+postgres
select
  id,
  name,
  type,
  provisioning_state,
  etag
from
  azure_lb_rule;
```

```sql+sqlite
select
  id,
  name,
  type,
  provisioning_state,
  etag
from
  azure_lb_rule;
```

### List failed load balancer rules
Discover the segments that have unsuccessful load balancer rules, allowing you to focus on rectifying these specific areas to improve network traffic distribution.

```sql+postgres
select
  id,
  name,
  type,
  provisioning_state
from
  azure_lb_rule
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
  azure_lb_rule
where
  provisioning_state = 'Failed';
```

### List load balancer rules order by idle timeout
Assess the elements within your load balancer rules to understand which ones have been idle for the longest time. This could help optimize resource allocation and improve system performance.

```sql+postgres
select
  id,
  name,
  type,
  idle_timeout_in_minutes
from
  azure_lb_rule
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
  azure_lb_rule
order by 
  idle_timeout_in_minutes;
```