---
title: "Steampipe Table: azure_compute_virtual_machine_metric_cpu_utilization_daily - Query Azure Compute Virtual Machines using SQL"
description: "Allows users to query Azure Compute Virtual Machine metrics, specifically the daily CPU utilization, providing insights into resource usage patterns and potential performance issues."
folder: "Compute"
---

# Table: azure_compute_virtual_machine_metric_cpu_utilization_daily - Query Azure Compute Virtual Machines using SQL

Azure Compute is a service within Microsoft Azure that offers scalable and secure virtual machines. These virtual machines provide the power to support large-scale, mission-critical applications. They allow users to deploy a wide range of computing solutions in an agile manner.

## Table Usage Guide

The `azure_compute_virtual_machine_metric_cpu_utilization_daily` table provides insights into the daily CPU utilization of Azure Compute Virtual Machines. As a system administrator or DevOps engineer, explore VM-specific CPU utilization details through this table to identify resource usage patterns and potential performance bottlenecks. Utilize it to monitor and optimize the performance of your Azure Compute resources effectively.

## Examples

### Basic info
Explore the performance metrics of Azure virtual machines on a daily basis to gain insights into CPU utilization trends. This can help identify instances of resource overload or inefficiency, assisting in better resource management and planning.

```sql+postgres
select
  name,
  timestamp,
  minimum,
  maximum,
  average,
  sample_count
from
  azure_compute_virtual_machine_metric_cpu_utilization_daily
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
  azure_compute_virtual_machine_metric_cpu_utilization_daily
order by
  name,
  timestamp;
```

### CPU Over 80% average
This example helps to pinpoint specific instances where the average CPU utilization of Azure virtual machines exceeds 80%. It's useful in identifying potential performance issues and ensuring efficient resource allocation.

```sql+postgres
select
  name,
  timestamp,
  round(minimum::numeric,2) as min_cpu,
  round(maximum::numeric,2) as max_cpu,
  round(average::numeric,2) as avg_cpu,
  sample_count
from
  azure_compute_virtual_machine_metric_cpu_utilization_daily
where
  average > 80
order by
  name,
  timestamp;
```

```sql+sqlite
select
  name,
  timestamp,
  round(minimum,2) as min_cpu,
  round(maximum,2) as max_cpu,
  round(average,2) as avg_cpu,
  sample_count
from
  azure_compute_virtual_machine_metric_cpu_utilization_daily
where
  average > 80
order by
  name,
  timestamp;
```