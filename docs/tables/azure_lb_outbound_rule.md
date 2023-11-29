---
title: "Steampipe Table: azure_lb_outbound_rule - Query Azure Load Balancer Outbound Rules using SQL"
description: "Allows users to query Azure Load Balancer Outbound Rules"
---

# Table: azure_lb_outbound_rule - Query Azure Load Balancer Outbound Rules using SQL

An Azure Load Balancer is a network performance utility within Microsoft Azure that enables you to manage network traffic to your applications. It operates at layer four of the Open Systems Interconnection (OSI) model and provides high availability by distributing incoming traffic among healthy service instances in cloud services or virtual machines in a load balancer set. Outbound Rules in Azure Load Balancer are used to control outbound connectivity for virtual machines (VMs) in your virtual network.

## Table Usage Guide

The 'azure_lb_outbound_rule' table provides insights into Outbound Rules within Azure Load Balancer. As a network administrator, explore rule-specific details through this table, including protocol type, backend pool, frontend IP configuration, and associated metadata. Utilize it to uncover information about outbound rules, such as those with specific protocols, the associated backend pool, and the configured frontend IP. The schema presents a range of attributes of the Outbound Rule for your analysis, like the rule id, provisioning state, protocol type, and associated tags.

## Examples

### Basic info
Explore which outbound rules are currently being provisioned within your Azure load balancer. This query allows you to keep track of the state of your rules, ensuring that your network traffic is being managed effectively.

```sql
select
  id,
  name,
  type,
  provisioning_state,
  etag
from
  azure_lb_outbound_rule;
```

### List failed load balancer outbound rules
Identify instances where load balancer outbound rules have failed in Azure. This can be beneficial in troubleshooting and understanding the network issues that might be affecting your services.

```sql
select
  id,
  name,
  type,
  provisioning_state
from
  azure_lb_outbound_rule
where
  provisioning_state = 'Failed';
```

### List load balancer outbound rules order by idle timeout
Explore which load balancer outbound rules have the longest idle timeouts to optimize resource allocation and efficiency. This could help in identifying areas where resources might be underutilized and could be better deployed elsewhere.

```sql
select
  id,
  name,
  type,
  idle_timeout_in_minutes
from
  azure_lb_outbound_rule
order by 
  idle_timeout_in_minutes;
```