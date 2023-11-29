---
title: "Steampipe Table: azure_lb - Query Azure Load Balancers using SQL"
description: "Allows users to query Azure Load Balancers."
---

# Table: azure_lb - Query Azure Load Balancers using SQL

Azure Load Balancers support the distribution of network traffic across Azure resources in a manner that is scalable and highly available. They provide low latency and high throughput, making applications highly responsive and robust. Load Balancers can be configured to provide public or private network access, and support both inbound and outbound scenarios.

## Table Usage Guide

The 'azure_lb' table provides insights into Load Balancers within Azure. As a DevOps engineer, explore Load Balancer-specific details through this table, including SKU, type, and associated metadata. Utilize it to uncover information about Load Balancers, such as those with specific provisioning states, the IP configurations, and the verification of backend address pools. The schema presents a range of attributes of the Load Balancer for your analysis, like the resource group name, subscription ID, and associated tags.

## Examples

### Basic info
Explore the basic details of your Azure Load Balancer to understand its operational state and location. This could be useful for assessing the load distribution and performance optimization in your network.

```sql
select
  id,
  name,
  type,
  provisioning_state,
  etag,
  region
from
  azure_lb;
```

### List failed load balancers
Identify instances where Azure load balancers have failed to provision correctly. This can help in diagnosing issues and ensuring optimal system performance.

```sql
select
  id,
  name,
  type,
  provisioning_state
from
  azure_lb
where
  provisioning_state = 'Failed';
```