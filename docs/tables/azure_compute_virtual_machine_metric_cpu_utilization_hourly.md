---
title: "Steampipe Table: azure_compute_virtual_machine_metric_cpu_utilization_hourly - Query Azure Compute Virtual Machine Metrics using SQL"
description: "Allows users to query Azure Compute Virtual Machine Metrics, specifically the hourly CPU utilization, providing insights into resource usage and potential performance bottlenecks."
folder: "Compute"
---

# Table: azure_compute_virtual_machine_metric_cpu_utilization_hourly - Query Azure Compute Virtual Machine Metrics using SQL

Azure Compute is a service within Microsoft Azure that provides scalable and secure virtual machines. It allows users to deploy and manage applications across a global network of Microsoft-managed data centers. Azure Compute provides a variety of virtual machine configurations to handle different workloads and performance requirements.

## Table Usage Guide

The `azure_compute_virtual_machine_metric_cpu_utilization_hourly` table provides insights into the CPU utilization of Azure Compute Virtual Machines on an hourly basis. As a system administrator or DevOps engineer, explore machine-specific details through this table, including CPU usage patterns, peak usage times, and potential performance bottlenecks. Utilize it to monitor and manage resource allocation, ensuring optimal performance and cost-effectiveness of your Azure Compute resources.

## Examples

### Basic info
Explore the utilization of virtual machine CPU over time to identify patterns or trends. This could help in efficient resource allocation and performance optimization.

```sql+postgres
select
  name,
  timestamp,
  minimum,
  maximum,
  average,
  sample_count
from
  azure_compute_virtual_machine_metric_cpu_utilization_hourly
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
  azure_compute_virtual_machine_metric_cpu_utilization_hourly
order by
  name,
  timestamp;
```

### CPU Over 80% average
Determine the areas in which the average CPU utilization exceeds 80% on Azure's virtual machines. This can be useful for identifying potential performance issues and ensuring efficient resource allocation.

```sql+postgres
select
  name,
  timestamp,
  round(minimum::numeric,2) as min_cpu,
  round(maximum::numeric,2) as max_cpu,
  round(average::numeric,2) as avg_cpu,
  sample_count
from
  azure_compute_virtual_machine_metric_cpu_utilization_hourly
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
  azure_compute_virtual_machine_metric_cpu_utilization_hourly
where
  average > 80
order by
  name,
  timestamp;
```