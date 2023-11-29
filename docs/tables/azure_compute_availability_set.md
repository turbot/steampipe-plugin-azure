---
title: "Steampipe Table: azure_compute_availability_set - Query Azure Compute Availability Sets using SQL"
description: "Allows users to query Azure Compute Availability Sets."
---

# Table: azure_compute_availability_set - Query Azure Compute Availability Sets using SQL

An Azure Compute Availability Set is a logical grouping capability that you can use in Azure to ensure that the VM resources you place within it are isolated from each other when they are deployed within an Azure datacenter. Azure ensures that the VMs you place within an Availability Set run across multiple physical servers, compute racks, storage units, and network switches. This is particularly useful for building high availability applications and protecting your applications from planned or unplanned maintenance.

## Table Usage Guide

The 'azure_compute_availability_set' table provides insights into the Availability Sets within Azure Compute. As a DevOps engineer, explore Availability Set-specific details through this table, including fault domain count, update domain count, and associated metadata. Utilize it to uncover information about Availability Sets, such as those with specific virtual machine profiles, the virtual machines within an availability set, and the verification of fault and update domains. The schema presents a range of attributes of the Availability Set for your analysis, like the set name, resource group, subscription ID, and associated tags.

## Examples

### Basic info
Explore which Azure Compute Availability Sets are in a specific region and assess the count of fault and update domains within them. This can help in managing and planning resource distribution across various domains and regions.

```sql
select
  name,
  platform_fault_domain_count,
  platform_update_domain_count,
  region
from
  azure_compute_availability_set;
```


### List of availability sets which does not use managed disks configuration
Identify instances where availability sets in Azure are not utilizing the managed disks configuration. This is beneficial in pinpointing areas where you could optimize your resources for improved performance and management.

```sql
select
  name,
  sku_name
from
  azure_compute_availability_set
where
  sku_name = 'Classic';
```


### List of availability sets without application tag key
Discover the segments that lack specific application tag keys within the Azure compute availability sets. This query is useful for identifying potential areas of misconfiguration or missing data.

```sql
select
  name,
  tags
from
  azure_compute_availability_set
where
  not tags :: JSONB ? 'application';
```