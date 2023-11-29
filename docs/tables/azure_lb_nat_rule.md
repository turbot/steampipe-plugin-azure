---
title: "Steampipe Table: azure_lb_nat_rule - Query Azure Load Balancer NAT Rules using SQL"
description: "Allows users to query Azure Load Balancer NAT Rules."
---

# Table: azure_lb_nat_rule - Query Azure Load Balancer NAT Rules using SQL

Azure Load Balancer is a high-performance, ultra low-latency Layer 4 load-balancing service (inbound and outbound) for all UDP and TCP protocols. Load Balancer NAT Rules are resources within Azure Load Balancer that allow you to control IP address translations. NAT rules use source network address translation (SNAT) and destination network address translation (DNAT) to translate IP addresses and ports.

## Table Usage Guide

The 'azure_lb_nat_rule' table provides insights into NAT rules within Azure Load Balancer. As a network administrator, explore NAT rule-specific details through this table, including inbound and outbound IP address translations, associated front-end IP configurations, and protocol types. Utilize it to uncover information about NAT rules, such as those with specific IP address translations, the associated load balancer, and the verification of protocol types. The schema presents a range of attributes of the NAT rule for your analysis, like the rule ID, resource group, subscription ID, and associated tags.

## Examples

### Basic info
Explore which Azure Load Balancer Network Address Translation (NAT) rules are currently in use. This can help in understanding the provisioning state and type of each rule for better management and optimization of network resources.

```sql
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
Explore which load balancer NAT rules in Azure have failed to provision, allowing you to identify potential issues and take corrective action.

```sql
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
Determine the areas in which load balancer NAT rules are prioritized based on their idle timeout duration. This helps in managing and optimizing network traffic flow by identifying rules that are inactive for longer periods.

```sql
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