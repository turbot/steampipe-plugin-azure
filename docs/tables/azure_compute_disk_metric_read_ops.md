---
title: "Steampipe Table: azure_compute_disk_metric_read_ops - Query Azure Compute Disk Metrics using SQL"
description: "Allows users to query Azure Compute Disk Metrics, specifically read operations, providing insights into disk performance and potential bottlenecks."
folder: "Resource"
---

# Table: azure_compute_disk_metric_read_ops - Query Azure Compute Disk Metrics using SQL

Azure Compute Disk Metrics is a resource within Microsoft Azure that allows you to monitor and analyze the performance of your Azure managed disks. It provides detailed information about read and write operations, throughput, and latency for your disks. Azure Compute Disk Metrics helps you understand disk performance and identify potential bottlenecks or performance issues.

## Table Usage Guide

The `azure_compute_disk_metric_read_ops` table provides insights into read operations on Azure managed disks. As a system administrator or DevOps engineer, explore disk-specific details through this table, including the number of read operations, the time of the operations, and associated metadata. Utilize it to monitor and analyze disk performance, identify potential bottlenecks, and optimize disk usage.

## Examples

### Basic info
Explore the performance of Azure compute disks over time by assessing the minimum, maximum, and average read operations. This can help determine potential bottlenecks and optimize disk usage for better system performance.

```sql+postgres
select
  name,
  timestamp,
  minimum,
  maximum,
  average,
  sample_count
from
  azure_compute_disk_metric_read_ops
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
  azure_compute_disk_metric_read_ops
order by
  name,
  timestamp;
```

### Operations Over 10 Bytes average
Determine the performance of Azure Compute Disks by identifying instances where the average read operations exceed 10 bytes. This is useful to monitor and optimize disk usage for improved system performance.

```sql+postgres
select
  name,
  timestamp,
  round(minimum::numeric,2) as min_read_ops,
  round(maximum::numeric,2) as max_read_ops,
  round(average::numeric,2) as avg_read_ops,
  sample_count
from
  azure_compute_disk_metric_read_ops
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
  azure_compute_disk_metric_read_ops
where
  average > 10
order by
  name,
  timestamp;
```