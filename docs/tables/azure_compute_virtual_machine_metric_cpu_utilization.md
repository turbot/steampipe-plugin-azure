---
title: "Steampipe Table: azure_compute_virtual_machine_metric_cpu_utilization - Query Azure Compute Virtual Machine Metrics using SQL"
description: "Allows users to query Azure Compute Virtual Machine CPU Utilization Metrics, providing insights into the CPU usage of virtual machines."
folder: "Resource"
---

# Table: azure_compute_virtual_machine_metric_cpu_utilization - Query Azure Compute Virtual Machine Metrics using SQL

Azure Compute is a service within Microsoft Azure that allows you to deploy and manage virtual machines. These virtual machines can be used to run applications, host databases, and perform other computing tasks. The CPU utilization metric provides information on the percentage of total CPU resources that are being used by a virtual machine.

## Table Usage Guide

The `azure_compute_virtual_machine_metric_cpu_utilization` table provides insights into the CPU utilization of virtual machines within Azure Compute. As a system administrator or DevOps engineer, explore CPU-specific details through this table, including the percentage of total CPU resources that are being used. Utilize it to monitor the performance of your virtual machines, identify those that are under heavy load, and make informed decisions about resource allocation and scaling.

## Examples

### Basic info
Determine the areas in which your Azure virtual machines' CPU utilization varies over time. This query helps you analyze performance trends and optimize resource allocation for improved efficiency.

```sql+postgres
select
  name,
  timestamp,
  minimum,
  maximum,
  average,
  sample_count
from
  azure_compute_virtual_machine_metric_cpu_utilization
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
  azure_compute_virtual_machine_metric_cpu_utilization
order by
  name,
  timestamp;
```

### CPU Over 80% average
Determine the areas in which the average CPU usage of Azure virtual machines exceeds 80%. This can be useful to identify potential performance issues and optimize resource allocation.

```sql+postgres
select
  name,
  timestamp,
  round(minimum::numeric,2) as min_cpu,
  round(maximum::numeric,2) as max_cpu,
  round(average::numeric,2) as avg_cpu,
  sample_count
from
  azure_compute_virtual_machine_metric_cpu_utilization
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
  azure_compute_virtual_machine_metric_cpu_utilization
where
  average > 80
order by
  name,
  timestamp;
```