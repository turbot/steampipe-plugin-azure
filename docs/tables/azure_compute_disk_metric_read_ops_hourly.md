---
title: "Steampipe Table: azure_compute_disk_metric_read_ops_hourly - Query Azure Compute Disk Metrics using SQL"
description: "Allows users to query Azure Compute Disk Metrics, specifically the hourly read operations count, providing insights into disk usage patterns and potential performance issues."
folder: "Compute"
---

# Table: azure_compute_disk_metric_read_ops_hourly - Query Azure Compute Disk Metrics using SQL

Azure Compute Disks are data storage units available in Microsoft Azure. They are used to store data for Azure Virtual Machines and other services. Azure Compute Disks provide high-performance, durable storage for I/O-intensive workloads.

## Table Usage Guide

The `azure_compute_disk_metric_read_ops_hourly` table provides insights into read operations of Azure Compute Disks on an hourly basis. As a system administrator or a DevOps engineer, explore disk-specific details through this table, including the number of read operations, the time of operations, and associated metadata. Utilize it to monitor disk performance, identify usage patterns, and detect potential performance issues.

## Examples

### Basic info
Assess the elements within the Azure compute disk's read operations on an hourly basis. This can help in identifying patterns, understanding usage trends, and planning for capacity or performance optimization.

```sql+postgres
select
  name,
  timestamp,
  minimum,
  maximum,
  average,
  sample_count
from
  azure_compute_disk_metric_read_ops_hourly
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
  azure_compute_disk_metric_read_ops_hourly
order by
  name,
  timestamp;
```

### Operations Over 10 Bytes average
This query is used to monitor disk read operations on Azure, specifically focusing on instances where the average read operation exceeds 10 bytes. This is useful for identifying potential performance issues or bottlenecks in the system, allowing for proactive management and optimization of resources.

```sql+postgres
select
  name,
  timestamp,
  round(minimum::numeric,2) as min_read_ops,
  round(maximum::numeric,2) as max_read_ops,
  round(average::numeric,2) as avg_read_ops,
  sample_count
from
  azure_compute_disk_metric_read_ops_hourly
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
  azure_compute_disk_metric_read_ops_hourly
where
  average > 10
order by
  name,
  timestamp;
```