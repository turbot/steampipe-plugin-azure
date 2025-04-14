---
title: "Steampipe Table: azure_compute_disk_metric_write_ops_hourly - Query Azure Compute Disk Metrics using SQL"
description: "Allows users to query Azure Compute Disk Metrics, specifically the hourly write operations, providing insights into disk usage patterns and potential anomalies."
folder: "Compute"
---

# Table: azure_compute_disk_metric_write_ops_hourly - Query Azure Compute Disk Metrics using SQL

Azure Compute Disk Metrics is a service within Microsoft Azure that allows users to monitor and track the performance of their Azure disks. It provides detailed data about the number of read and write operations, the amount of data transferred, and the latency of these operations. This service is crucial for understanding disk usage patterns, identifying potential bottlenecks, and optimizing performance.

## Table Usage Guide

The `azure_compute_disk_metric_write_ops_hourly` table provides insights into the hourly write operations of Azure Compute Disks. As a system administrator or DevOps engineer, explore disk-specific details through this table, including the number of write operations and the time of these operations. Utilize it to understand disk usage patterns, identify potential performance bottlenecks, and optimize your Azure disk configurations.

## Examples

### Basic info
Explore the performance of Azure compute disks over time by tracking the minimum, maximum, and average write operations per hour. This can help in identifying usage patterns, planning capacity, and troubleshooting performance issues.

```sql+postgres
select
  name,
  timestamp,
  minimum,
  maximum,
  average,
  sample_count
from
  azure_compute_disk_metric_write_ops_hourly
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
  azure_compute_disk_metric_write_ops_hourly
order by
  name,
  timestamp;
```

### Operations Over 10 Bytes average
This query is used to track the performance of Azure compute disks, specifically focusing on those with an average of more than 10 write operations per hour. By doing so, it helps in identifying potential bottlenecks and ensuring optimal disk performance.

```sql+postgres
select
  name,
  timestamp,
  round(minimum::numeric,2) as min_write_ops,
  round(maximum::numeric,2) as max_write_ops,
  round(average::numeric,2) as avg_write_ops,
  sample_count
from
  azure_compute_disk_metric_write_ops_hourly
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
  round(minimum,2) as min_write_ops,
  round(maximum,2) as max_write_ops,
  round(average,2) as avg_write_ops,
  sample_count
from
  azure_compute_disk_metric_write_ops_hourly
where
  average > 10
order by
  name,
  timestamp;
```