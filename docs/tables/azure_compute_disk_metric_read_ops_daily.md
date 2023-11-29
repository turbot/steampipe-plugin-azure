---
title: "Steampipe Table: azure_compute_disk_metric_read_ops_daily - Query Azure Compute Disks using SQL"
description: "Allows users to query Azure Compute Disks daily read operations metrics."
---

# Table: azure_compute_disk_metric_read_ops_daily - Query Azure Compute Disks using SQL

Azure Compute Disks are a key component of Azure Infrastructure-as-a-Service (IaaS) based solutions, providing high-performance, reliable, and resilient block storage for Azure Virtual Machines. They support a wide range of workloads like relational databases, high-volume transactional systems, and big data applications. Azure Compute Disks offer a variety of disk storage options to meet varying workload requirements.

## Table Usage Guide

The 'azure_compute_disk_metric_read_ops_daily' table provides insights into the daily read operations of Azure Compute Disks. As a system administrator or a DevOps engineer, you can explore disk-specific details through this table, including the total number of read operations, maximum and average read operations, and the time at which the maximum read operations occurred. Utilize it to monitor the performance of your disks, identify potential bottlenecks, and plan capacity. The schema presents a range of attributes of the disk read operations for your analysis, like the resource group name, subscription ID, time grain, and unit type.

## Examples

### Basic info
Explore the daily read operations metrics for Azure compute disks to understand usage patterns and performance. This can help in identifying any unusual activity or potential areas for optimization.

```sql
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
Explore which operations have an average higher than 10 bytes. This is useful for identifying potential areas of heavy data usage or inefficiency in your Azure compute disk metrics.

```sql
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