---
title: "Steampipe Table: azure_compute_disk_metric_read_ops_hourly - Query Azure Compute Disks using SQL"
description: "Allows users to query Azure Compute Disks' hourly read operations metrics."
---

# Table: azure_compute_disk_metric_read_ops_hourly - Query Azure Compute Disks using SQL

Azure Compute Disks are data storage units available in Microsoft Azure, used to manage and store data persistently. These disks are designed to provide secure, scalable storage for virtual machines. They offer high-performance, durable storage for Azure Virtual Machines instances.

## Table Usage Guide

The 'azure_compute_disk_metric_read_ops_hourly' table provides insights into the read operations metrics of Azure Compute Disks on an hourly basis. As a DevOps engineer, you can use this table to explore disk-specific details such as the number of read operations, their time duration, and other related metadata. This can be particularly useful for monitoring disk performance, identifying potential bottlenecks, and ensuring optimal data management. The schema presents a range of attributes for your analysis, such as the disk name, resource group, subscription ID, and the count of read operations.

## Examples

### Basic info
Explore the performance of your Azure Compute Disks by analyzing the hourly read operations. This allows you to identify periods of high or low activity, assisting in capacity planning and troubleshooting performance issues.

```sql
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
This query is useful to analyze disk operations that exceed an average of 10 bytes in Azure's Compute Disk service. It can help optimize system performance by identifying potential bottlenecks in disk operations.

```sql
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