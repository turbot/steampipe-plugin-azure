---
title: "Steampipe Table: azure_compute_disk_metric_write_ops - Query Azure Compute Disk Metrics using SQL"
description: "Allows users to query Azure Compute Disk Metrics, specifically the write operations, providing insights into disk usage and potential performance issues."
---

# Table: azure_compute_disk_metric_write_ops - Query Azure Compute Disk Metrics using SQL

Azure Compute Disk Metrics is a feature within Microsoft Azure that allows monitoring and analysis of disk performance and usage. It offers detailed information on various metrics, such as write operations, enabling users to understand disk behavior and identify potential performance issues. This feature is crucial for maintaining optimal disk performance and managing storage resources efficiently.

## Table Usage Guide

The `azure_compute_disk_metric_write_ops` table provides insights into write operations on Azure Compute Disks. As a system administrator or DevOps engineer, you can explore disk-specific details through this table, including the number of write operations, to understand disk usage patterns and potential performance bottlenecks. Utilize it to monitor and optimize disk performance, and ensure efficient resource management in your Azure environment.

## Examples

### Basic info
Explore which Azure compute disk has the most write operations over time. This can help optimize disk usage by identifying high-usage periods and potentially underutilized resources.

```sql+postgres
select
  name,
  timestamp,
  minimum,
  maximum,
  average,
  sample_count
from
  azure_compute_disk_metric_write_ops
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
  azure_compute_disk_metric_write_ops
order by
  name,
  timestamp;
```

### Operations Over 10 Bytes average
Determine the areas in which the average write operations on Azure compute disks exceed 10 bytes. This can be useful in identifying potential bottlenecks or high usage periods, enabling proactive management and optimization of disk resources.

```sql+postgres
select
  name,
  timestamp,
  round(minimum::numeric,2) as min_write_ops,
  round(maximum::numeric,2) as max_write_ops,
  round(average::numeric,2) as avg_write_ops,
  sample_count
from
  azure_compute_disk_metric_write_ops
where average > 10
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
  azure_compute_disk_metric_write_ops
where average > 10
order by
  name,
  timestamp;
```