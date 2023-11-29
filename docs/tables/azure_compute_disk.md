---
title: "Steampipe Table: azure_compute_disk - Query Azure Compute Disks using SQL"
description: "Allows users to query Azure Compute Disks."
---

# Table: azure_compute_disk - Query Azure Compute Disks using SQL

Azure Compute Disks are durable, high-performance, secure disk storage for Azure Virtual Machines. They provide persistent, secured disk storage and support for industry-leading data protection capabilities. Azure Compute Disks can be used with Azure Virtual Machines to deliver high-performance and highly durable disk storage.

## Table Usage Guide

The 'azure_compute_disk' table provides insights into Azure Compute Disks within Azure Compute service. As a DevOps engineer, explore disk-specific details through this table, including disk size, creation time, encryption settings, and associated metadata. Utilize it to uncover information about disks, such as the ones with specific encryption settings, the type of disks, and their provisioning state. The schema presents a range of attributes of the Azure Compute Disk for your analysis, like the disk ID, creation time, disk state, and associated tags.

## Examples

### List of all premium tier compute disks
Determine the areas in which premium tier compute disks are being utilized within the Azure environment. This can be beneficial for cost management and resource optimization.

```sql
select
  name,
  sku_name,
  sku_tier
from
  azure_compute_disk
where
  sku_tier = 'Premium';
```


### List of unattached disks
Determine the areas in which there are unattached disks within your Azure Compute service. This can help you identify unused resources and potential cost savings.

```sql
select
  name,
  disk_state
from
  azure_compute_disk
where
  disk_state = 'Unattached';
```


### Size and performance info of each disk
Explore the performance and capacity of each disk in your Azure Compute environment. This information can be crucial for optimizing resource allocation and ensuring efficient data operations.

```sql
select
  name,
  disk_size_gb as disk_size,
  disk_iops_read_only as disk_iops_read_only,
  disk_iops_read_write as provision_iops,
  disk_iops_mbps_read_write as bandwidth,
  disk_iops_mbps_read_only as disk_mbps_read_write
from
  azure_compute_disk;
```


### List of compute disks which are not available in multiple az
Determine the areas in which certain compute disks are not available across multiple Azure availability zones. This is useful in identifying potential risks to data redundancy and disaster recovery plans.

```sql
select
  name,
  az,
  region
from
  azure_compute_disk
  cross join jsonb_array_elements(zones) az
where
  zones is not null;
```


### List of compute disks which are not encrypted with customer key
Discover the segments of your Azure compute disks that are not utilizing customer key encryption. This is beneficial in identifying potential security vulnerabilities and ensuring data protection standards are met.

```sql
select
  name,
  encryption_type
from
  azure_compute_disk
where
  encryption_type <> 'EncryptionAtRestWithCustomerKey';
```