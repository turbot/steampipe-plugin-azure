---
title: "Steampipe Table: azure_compute_disk_metric_read_ops_daily - Query Azure Compute Disk Metrics using SQL"
description: "Allows users to query Azure Compute Disk Metrics, specifically the daily read operations, providing insights into disk read performance and potential bottlenecks."
---

# Table: azure_compute_disk_metric_read_ops_daily - Query Azure Compute Disk Metrics using SQL

Azure Compute Disk Metrics is a feature within Microsoft Azure that provides data about the performance of Azure managed disks. It provides detailed information about disk read operations, write operations, and other disk performance metrics. This feature helps users monitor and optimize the performance of their Azure managed disks.

## Table Usage Guide

The `azure_compute_disk_metric_read_ops_daily` table provides insights into the daily read operations of Azure managed disks. As a system administrator or DevOps engineer, use this table to monitor disk performance and identify potential bottlenecks or performance issues. This table can be particularly useful in optimizing disk usage and ensuring efficient operation of your Azure resources.

## Examples

### Basic info
Explore the daily read operations on Azure compute disks to gain insights into the average, minimum, and maximum operations, along with the sample count. This is useful for tracking disk usage patterns and identifying any unusual activity or potential bottlenecks.

```sql+postgres
select
  name,
  timestamp,
  minimum,
  maximum,
  average,
  sample_count
from
  azure_compute_disk_metric_read_ops_daily
order by
  name,
  timestamp;
```

```sql+sqlite
select
  name,
  timestamp,
  minimum,
  maximum,
  average,
  sample_count
from
  azure_compute_disk_metric_read_ops_daily
order by
  name,
  timestamp;
```

### Operations Over 10 Bytes average
This query is used to monitor the performance of Azure Compute Disk operations by identifying those with an average read operation count exceeding 10 per day. It allows for effective resource management by highlighting areas where usage may be higher than expected.

```sql+postgres
select
  name,
  timestamp,
  round(minimum::numeric,2) as min_read_ops,
  round(maximum::numeric,2) as max_read_ops,
  round(average::numeric,2) as avg_read_ops,
  sample_count
from
  azure_compute_disk_metric_read_ops_daily
where
  average > 10
order by
  name,
  timestamp;
```

```sql+sqlite
select
  name,
  timestamp,
  round(minimum,2) as min_read_ops,
  round(maximum,2) as max_read_ops,
  round(average,2) as avg_read_ops,
  sample_count
from
  azure_compute_disk_metric_read_ops_daily
where
  average > 10
order by
  name,
  timestamp;
```