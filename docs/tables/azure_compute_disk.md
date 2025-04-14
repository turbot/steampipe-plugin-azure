---
title: "Steampipe Table: azure_compute_disk - Query Azure Compute Disks using SQL"
description: "Allows users to query Azure Compute Disks, specifically providing details about each disk's properties including its type, size, location, and encryption settings."
folder: "Compute"
---

# Table: azure_compute_disk - Query Azure Compute Disks using SQL

Azure Compute Disk is a resource within Microsoft Azure that allows you to create and manage disks for your virtual machines. These disks can be used as system disks or data disks and come in different types, including standard HDD, standard SSD, and premium SSD. Azure Compute Disk also supports disk encryption for enhanced security.

## Table Usage Guide

The `azure_compute_disk` table provides insights into the disks used in Azure Compute. As a system administrator or developer, you can explore disk-specific details through this table, including the type, size, location, and encryption settings of each disk. Utilize it to manage disk resources effectively, ensuring optimal allocation and enhanced security.

## Examples

### List of all premium tier compute disks
Determine the areas in which premium tier compute disks are being utilized across your Azure environment. This can help in resource management and cost optimization by identifying areas of high-end usage.

```sql+postgres
select
  name,
  sku_name,
  sku_tier
from
  azure_compute_disk
where
  sku_tier = 'Premium';
```

```sql+sqlite
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
Discover the segments that consist of unused storage resources within your Azure infrastructure. This can aid in optimizing resource allocation and reducing unnecessary costs.

```sql+postgres
select
  name,
  disk_state
from
  azure_compute_disk
where
  disk_state = 'Unattached';
```

```sql+sqlite
select
  name,
  disk_state
from
  azure_compute_disk
where
  disk_state = 'Unattached';
```

### Size and performance info of each disk
Gain insights into the performance and size of each disk in your Azure Compute service. This helps in optimizing resource allocation and identifying potential performance bottlenecks.

```sql+postgres
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

```sql+sqlite
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
Determine the areas in which Azure Compute Disks are not available across multiple availability zones. This is useful for identifying potential vulnerabilities in your system's redundancy and disaster recovery capabilities.

```sql+postgres
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

```sql+sqlite
select
  name,
  az.value as az,
  region
from
  azure_compute_disk,
  json_each(zones) az
where
  zones is not null;
```

### List of compute disks which are not encrypted with customer key
Discover the segments that utilize compute disks not encrypted with a customer key, enabling you to identify potential security risks and take necessary actions to enhance data protection.

```sql+postgres
select
  name,
  encryption_type
from
  azure_compute_disk
where
  encryption_type <> 'EncryptionAtRestWithCustomerKey';
```

```sql+sqlite
select
  name,
  encryption_type
from
  azure_compute_disk
where
  encryption_type != 'EncryptionAtRestWithCustomerKey';
```