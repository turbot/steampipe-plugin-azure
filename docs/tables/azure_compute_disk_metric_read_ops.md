---
title: "Steampipe Table: azure_compute_disk_metric_read_ops - Query Azure Compute Disks using SQL"
description: "Allows users to query Azure Compute Disks read operations metrics."
---

# Table: azure_compute_disk_metric_read_ops - Query Azure Compute Disks using SQL

Azure Compute Disks are a type of Azure Storage that provide high-performance, durable block storage for Azure Virtual Machines. These disks are designed to support I/O-intensive workloads and offer seamless integration with Azure Virtual Machines. They provide consistent low-latency performance, deliver high IOPS/throughput, and ensure data durability and availability.

## Table Usage Guide

The 'azure_compute_disk_metric_read_ops' table provides insights into read operations metrics of Azure Compute Disks. As a DevOps engineer, explore specific details through this table, including the time grain, average, minimum, and maximum read operations. Utilize it to monitor and analyze the performance of your disks, such as those with high read operations, the average read operations over a period, and the peak read operations. The schema presents a range of attributes of the read operations metrics for your analysis, like the unit, timestamp, and total count.

## Examples

### Basic info
Explore the performance of Azure Compute Disks over time to identify potential bottlenecks or inefficiencies. This query provides a historical overview of disk operations, helping you pinpoint areas for optimization or resource allocation.

```sql
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
Determine the areas in which Azure compute disk read operations exceed an average of 10 bytes. This can be useful for identifying potential performance bottlenecks or areas where optimization may be beneficial.

```sql
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