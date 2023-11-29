---
title: "Steampipe Table: azure_compute_virtual_machine_metric_cpu_utilization - Query Azure Compute Virtual Machines using SQL"
description: "Allows users to query Azure Compute Virtual Machine CPU Utilization metrics."
---

# Table: azure_compute_virtual_machine_metric_cpu_utilization - Query Azure Compute Virtual Machines using SQL

Azure Compute is a service within Microsoft Azure that provides on-demand, high-scale, secure, virtualized infrastructure using Microsoft's advanced data centers. With Azure Compute, users can deploy a wide range of computing solutions, including virtual machines (VMs). This service is particularly useful for workloads that require high-performance computing, analytics, AI, real-time applications, and low-latency applications.

## Table Usage Guide

The 'azure_compute_virtual_machine_metric_cpu_utilization' table provides insights into CPU utilization metrics of Azure Compute Virtual Machines. As a systems administrator, you can explore VM-specific details through this table, including the average, minimum, and maximum CPU utilization, and the timestamps for these metrics. Utilize it to uncover information about VM performance, such as identifying VMs with high CPU utilization, understanding the CPU usage pattern over time, and taking necessary actions to optimize resource usage. The schema presents a range of attributes of the VM CPU utilization metrics for your analysis, like the average CPU utilization, minimum CPU utilization, maximum CPU utilization, and the timestamps for these metrics.

## Examples

### Basic info
Explore the use patterns of virtual machines in your Azure environment by analyzing CPU utilization metrics. This can help identify periods of high demand or underutilization, allowing for better resource management and cost optimization.

```sql
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
Determine the areas in which virtual machine CPU utilization exceeds 80% on average. This can help in identifying potential performance issues and ensuring efficient resource management.

```sql
select
  name,
  timestamp,
  round(minimum::numeric,2) as min_cpu,
  round(maximum::numeric,2) as max_cpu,
  round(average::numeric,2) as avg_cpu,
  sample_count
from
  azure_compute_virtual_machine_metric_cpu_utilization
where average > 80
order by
  name,
  timestamp;
```