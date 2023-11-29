---
title: "Steampipe Table: azure_lb_rule - Query Azure Load Balancer Rules using SQL"
description: "Allows users to query Azure Load Balancer Rules."
---

# Table: azure_lb_rule - Query Azure Load Balancer Rules using SQL

Azure Load Balancer is a highly available network performance utility that distributes incoming network traffic across many servers. It ensures the delivery of network traffic to various services in the Microsoft Azure public cloud, virtual machines (VMs) and other operations. Azure Load Balancer supports inbound and outbound scenarios, provides low latency and high throughput, and scales up to millions of flows for all Transmission Control Protocol (TCP) and User Datagram Protocol (UDP) applications.

## Table Usage Guide

The 'azure_lb_rule' table provides insights into Load Balancer Rules within Azure Load Balancer. As a network administrator, explore rule-specific details through this table, including protocol type, frontend and backend port, and associated metadata. Utilize it to uncover information about rules, such as the load distribution method, whether direct server return is enabled, and the idle timeout in minutes. The schema presents a range of attributes of the Load Balancer Rule for your analysis, like the rule ID, provisioning state, and associated tags.

## Examples

### Basic info
Analyze the settings of your Azure load balancer rules to understand their current state and type. This can be beneficial for assessing your network traffic management and ensuring it aligns with your intended configuration.

```sql
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
Analyze the settings to understand which load balancer rules have failed in their setup process, providing insights to troubleshoot and rectify the issues.

```sql
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
Pinpoint the specific load balancer rules based on their idle timeout duration. This can help in optimizing system performance and managing resource allocation effectively.

```sql
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