---
title: "Steampipe Table: azure_compute_virtual_machine_metric_cpu_utilization_hourly - Query Azure Compute Virtual Machines using SQL"
description: "Allows users to query Azure Compute Virtual Machines' CPU Utilization metrics on an hourly basis."
---

# Table: azure_compute_virtual_machine_metric_cpu_utilization_hourly - Query Azure Compute Virtual Machines using SQL

Azure Compute is a service that provides on-demand, scalable compute resources in the cloud. It allows users to create and manage virtual machines (VMs) that run on Microsoft's data centers. The service is designed to support a wide range of workloads, including web applications, batch processing, and high-performance computing.

## Table Usage Guide

The 'azure_compute_virtual_machine_metric_cpu_utilization_hourly' table provides insights into the CPU utilization metrics of Azure Compute Virtual Machines on an hourly basis. As a system administrator, you can use this table to monitor and analyze the CPU usage of your virtual machines, helping you to optimize resource allocation and performance. The table provides detailed information such as the maximum and average CPU utilization, the time of the metric, and the resource group of the VM. Utilize it to uncover trends in CPU usage, identify potential performance bottlenecks, and make informed decisions about scaling and capacity planning. The schema presents a range of attributes of the VM's CPU utilization for your analysis, like the maximum and average utilization, the timestamp of the metric, and the resource group of the VM.

## Examples

### Basic info
Explore which virtual machines in your Azure Compute environment have the highest CPU utilization over the past hour. This can help you identify potential performance issues and optimize resource allocation.

```sql
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
Analyze the performance of Azure virtual machines by identifying instances where the average CPU utilization exceeds 80%. This can be useful for spotting potential bottlenecks or performance issues in your infrastructure.

```sql
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