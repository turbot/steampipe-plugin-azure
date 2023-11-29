---
title: "Steampipe Table: azure_lb_probe - Query Azure Load Balancer Probes using SQL"
description: "Allows users to query Azure Load Balancer Probes."
---

# Table: azure_lb_probe - Query Azure Load Balancer Probes using SQL

Azure Load Balancer is a high-performance, ultra low-latency Layer 4 load-balancing service for all UDP and TCP protocols. Probes in Azure Load Balancer monitor the health of the resources in your load balancer's backend pool. They automatically detect failures and take steps to ensure that traffic only goes to healthy resources.

## Table Usage Guide

The 'azure_lb_probe' table provides insights into the probes within Azure Load Balancer. As a DevOps engineer, explore probe-specific details through this table, including protocol, port, request path, and associated metadata. Utilize it to uncover information about probes, such as their interval and timeout settings, the number of unhealthy responses before marking a resource as "unhealthy", and the load balancer that each probe is associated with. The schema presents a range of attributes of the probe for your analysis, like the probe's ID, name, and type, as well as the resource group and subscription it belongs to.

## Examples

### Basic info
Explore which Azure load balancer probes are currently active. This can help in determining the operational status and managing the load balancing configuration effectively.

```sql
select
  id,
  name,
  type,
  provisioning_state,
  load_balancer_name,
  port
from
  azure_lb_probe;
```

### List failed load balancer probes
Discover the segments that have failed load balancer probes to identify potential issues with your Azure load balancer setup. This could help in troubleshooting and enhancing the overall performance and reliability of your network infrastructure.

```sql
select
  id,
  name,
  type,
  provisioning_state
from
  azure_lb_probe
where
  provisioning_state = 'Failed';
```

### List load balancer probes order by interval
Analyze the settings to understand the frequency of load balancer probes within your Azure environment. This can help optimize network performance by identifying probes with unusually high or low intervals.

```sql
select
  id,
  name,
  type,
  interval_in_seconds
from
  azure_lb_probe
order by 
  interval_in_seconds;
```