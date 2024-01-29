---
title: "Steampipe Table: azure_lb_probe - Query Azure Load Balancer Probes using SQL"
description: "Allows users to query Azure Load Balancer Probes, providing valuable insights into the health and performance of the load balancer."
---

# Table: azure_lb_probe - Query Azure Load Balancer Probes using SQL

Azure Load Balancer Probes are a feature of Azure Load Balancer that enables monitoring of the health and performance of the load balancer. Probes are used to detect the health of the backend resources and make decisions about sending network traffic. They provide essential information to ensure the efficient and reliable operation of the load balancer.

## Table Usage Guide

The `azure_lb_probe` table provides insights into Azure Load Balancer Probes. Network administrators and DevOps engineers can use this table to monitor the health and performance of the load balancer, making it a valuable resource for maintaining optimal network performance. Furthermore, it can be utilized to detect anomalies and troubleshoot potential issues, ensuring the reliability and efficiency of the load balancer.

## Examples

### Basic info
Explore which Azure load balancer probes are currently in use to understand their configuration and state. This can help in managing network traffic and ensuring optimal performance.

```sql+postgres
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

```sql+sqlite
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
Discover the segments that have failed load balancer probes in your Azure environment. This information can help you identify potential issues and improve your resource management.

```sql+postgres
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

```sql+sqlite
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
Assess the elements within your Azure load balancer probes to prioritize them based on frequency of checks, allowing you to understand and manage the performance and availability of your services.

```sql+postgres
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

```sql+sqlite
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