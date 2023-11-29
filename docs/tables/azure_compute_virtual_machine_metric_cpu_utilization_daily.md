---
title: "Steampipe Table: azure_compute_virtual_machine_metric_cpu_utilization_daily - Query Azure Compute Virtual Machines using SQL"
description: "Allows users to query Azure Compute Virtual Machine daily CPU utilization metrics."
---

# Table: azure_compute_virtual_machine_metric_cpu_utilization_daily - Query Azure Compute Virtual Machines using SQL

Azure Compute is a service that provides on-demand, scalable compute resources in Microsoft Azure. It allows you to deploy and manage virtual machines and containers, and supports a range of operating systems, tools, and frameworks. Virtual machines are a core part of Azure Compute, providing the ability to quickly scale up or down with demand, and offering a range of options for CPU, memory, storage, and networking capacity.

## Table Usage Guide

The 'azure_compute_virtual_machine_metric_cpu_utilization_daily' table provides insights into daily CPU utilization metrics of Azure Compute Virtual Machines. As a DevOps engineer, you can use this table to monitor and analyze the daily CPU usage of your virtual machines, helping you to understand the performance and resource demands of your applications and services. The schema presents a range of attributes for your analysis, such as the maximum, minimum, and average CPU utilization, the time of the metric, and the resource group and subscription ID of the virtual machine. Utilize this table to identify trends in resource usage, detect potential issues, and optimize your Azure Compute resources.

## Examples

### Basic info
Explore which Azure virtual machines have high CPU utilization over time. This can help in managing resources efficiently by identifying machines that may need upgrades or load balancing.

```sql
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
Determine the areas in which the average CPU utilization of Azure virtual machines exceeds 80%. This query can help identify potential performance issues and optimize resource allocation.

```sql
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