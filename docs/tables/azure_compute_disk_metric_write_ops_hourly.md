---
title: "Steampipe Table: azure_compute_disk_metric_write_ops_hourly - Query Azure Compute Disks using SQL"
description: "Allows users to query Azure Compute Disks metrics on an hourly basis."
---

# Table: azure_compute_disk_metric_write_ops_hourly - Query Azure Compute Disks using SQL

Azure Compute Disks is a service that allows you to create and manage disks for your virtual machines. These disks can be used as system disks or data disks, and are available in different performance tiers to meet the needs of various applications and workloads. Azure Compute Disks also provide capabilities such as disk snapshots and disk backups for data protection and recovery.

## Table Usage Guide

The 'azure_compute_disk_metric_write_ops_hourly' table provides insights into the write operations metrics of Azure Compute Disks on an hourly basis. As a system administrator, you can use this table to explore the write operations performance of your disks, including the frequency and volume of data written to the disks. The table offers detailed metrics such as the timestamp of the data, minimum, maximum, and average write operations, and total count of write operations. Utilize it to monitor your disk performance, identify potential bottlenecks, and optimize your disk utilization for improved application performance. The schema presents a range of attributes of the disk write operations for your analysis, like the disk name, resource group, subscription ID, and more.

## Examples

### Basic info
Explore the performance of Azure compute disks by examining hourly write operations. This information can help identify potential bottlenecks or performance issues, allowing you to optimize your disk usage.

```sql
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
Explore which operations have an above-average rate, allowing you to assess potential areas of high activity or strain on your system. This can be useful in managing resources and identifying potential bottlenecks or areas for optimization.

```sql
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