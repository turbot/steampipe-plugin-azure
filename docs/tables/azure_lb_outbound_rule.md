---
title: "Steampipe Table: azure_lb_outbound_rule - Query Azure Load Balancer Outbound Rules using SQL"
description: "Allows users to query Azure Load Balancer Outbound Rules, specifically providing details about outbound rules configured for Azure Load Balancers."
---

# Table: azure_lb_outbound_rule - Query Azure Load Balancer Outbound Rules using SQL

Azure Load Balancer Outbound Rules are part of the Azure Load Balancer service, which allows you to manage and distribute network traffic. Outbound Rules in Azure Load Balancer provide you with the ability to control the outbound network traffic from a virtual network to the internet. They offer you the flexibility to scale and tune your outbound connectivity.

## Table Usage Guide

The `azure_lb_outbound_rule` table provides insights into the outbound rules of Azure Load Balancers. As a network administrator, you can use this table to get detailed information about each outbound rule, including the associated Load Balancer, the allocated outbound ports, and the protocol used. This can help you better understand your network traffic and potentially identify any configuration issues or areas for optimization.

## Examples

### Basic info
Determine the status and type of your Azure Load Balancer outbound rules to understand their current provisioning state and manage your network traffic more effectively.

```sql+postgres
select
  id,
  name,
  type,
  provisioning_state,
  etag
from
  azure_lb_outbound_rule;
```

```sql+sqlite
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
Discover the segments that have failed in the provisioning process of outbound rules in Azure's load balancer. This can be useful to identify and troubleshoot problematic areas in your network infrastructure.

```sql+postgres
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

```sql+sqlite
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
Analyze the settings to understand the sequence of outbound rules for a load balancer based on their idle timeout duration. This can help in effective management and optimization of network resources.

```sql+postgres
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

```sql+sqlite
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