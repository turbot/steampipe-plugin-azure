---
title: "Steampipe Table: azure_compute_disk_metric_write_ops - Query Azure Compute Disks using SQL"
description: "Allows users to query Azure Compute Disks write operations metrics."
---

# Table: azure_compute_disk_metric_write_ops - Query Azure Compute Disks using SQL

Azure Compute Disks are a type of storage that can be attached to Azure Virtual Machines. They provide persistent, secured, and highly reliable storage capabilities, allowing you to read and write data. Azure Compute Disks come in different performance tiers to support a variety of workloads and applications.

## Table Usage Guide

The 'azure_compute_disk_metric_write_ops' table provides insights into the write operations metrics of Azure Compute Disks. As a DevOps engineer, explore disk-specific details through this table, including total write operations, average write operations, and maximum write operations. Utilize it to monitor and analyze the performance of your Azure Compute Disks, identify any unusual increase in write operations, and optimize disk usage. The schema presents a range of attributes of the Compute Disk write operations for your analysis, like the average, maximum, minimum, and total count of write operations.

## Examples

### Basic info
Explore which Azure Compute Disk has the most write operations over time. This can help in understanding disk usage patterns and planning for potential disk capacity upgrades.

```sql
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
Explore disk operations that have an average higher than 10 bytes. This can be useful to monitor and manage storage performance, ensuring efficient data handling and optimal system operation.

```sql
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