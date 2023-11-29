---
title: "Steampipe Table: azure_compute_disk_metric_write_ops_daily - Query Azure Compute Disks using SQL"
description: "Allows users to query daily write operations metrics of Azure Compute Disks."
---

# Table: azure_compute_disk_metric_write_ops_daily - Query Azure Compute Disks using SQL

Azure Compute Disks are a key component of the Azure Infrastructure-as-a-Service (IaaS) offering. These disks provide durable, secure, and scalable storage for the data that drives your applications and services. Azure Compute Disks support a variety of workloads, like relational databases, large-scale NoSQL databases, and enterprise applications, with the flexibility and security required for Azure-based virtual machines.

## Table Usage Guide

The 'azure_compute_disk_metric_write_ops_daily' table provides insights into the daily write operations metrics of Azure Compute Disks. As a system administrator or a DevOps engineer, you can explore disk-specific details through this table, including the time grain, average, minimum, and maximum write operations. Utilize it to uncover information about disk performance, such as spikes in write operations, periods of low activity, and overall write operation trends. The schema presents a range of attributes of the disk's write operations for your analysis, like the unit, timestamp, and the total count of write operations.

## Examples

### Basic info
Analyze the daily write operations on Azure Compute Disks to understand performance trends and identify potential areas of concern. This can help in proactive resource management and ensure optimal application performance.

```sql
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
Uncover the details of disk operations in your Azure Compute instances that exceed an average of 10 bytes. This allows you to monitor and manage disk usage effectively, ensuring optimal performance.

```sql
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