---
title: "Steampipe Table: azure_compute_disk_metric_write_ops_daily - Query Azure Compute Disk Metrics using SQL"
description: "Allows users to query Azure Compute Disk Metrics, specifically focusing on daily write operations. This provides valuable insights into disk usage and performance."
folder: "Resource"
---

# Table: azure_compute_disk_metric_write_ops_daily - Query Azure Compute Disk Metrics using SQL

Azure Compute Disk is a resource within Microsoft Azure that provides scalable and secure disk storage for Azure Virtual Machines. It offers high-performance, highly durable block storage for your mission-critical workloads. You can use it to persist data by writing it to the disk, or to read data from the disk.

## Table Usage Guide

The `azure_compute_disk_metric_write_ops_daily` table provides insights into daily write operations on Azure Compute Disks. As a system administrator or a DevOps engineer, you can use this table to monitor disk performance and usage, enabling you to proactively address any potential issues. This can help you ensure optimal performance and availability of your Azure resources.

## Examples

### Basic info
Analyze the daily write operations on Azure compute disks to understand their usage patterns and performance metrics. This query can be useful in identifying potential bottlenecks or areas for optimization in your storage infrastructure.

```sql+postgres
select
  name,
  timestamp,
  minimum,
  maximum,
  average,
  sample_count
from
  azure_compute_disk_metric_write_ops_daily
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
  azure_compute_disk_metric_write_ops_daily
order by
  name,
  timestamp;
```

### Operations Over 10 Bytes average
Determine the areas in which the average daily write operations on Azure Compute Disk exceed 10 bytes. This can help optimize disk usage by identifying potential inefficiencies or areas of high activity.

```sql+postgres
select
  name,
  timestamp,
  round(minimum::numeric,2) as min_write_ops,
  round(maximum::numeric,2) as max_write_ops,
  round(average::numeric,2) as avg_write_ops,
  sample_count
from
  azure_compute_disk_metric_write_ops_daily
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
  azure_compute_disk_metric_write_ops_daily
where
  average > 10
order by
  name,
  timestamp;
```