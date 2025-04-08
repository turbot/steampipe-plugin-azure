---
title: "Steampipe Table: azure_compute_image - Query Azure Compute Images using SQL"
description: "Allows users to query Azure Compute Images, providing detailed information about the virtual machine images available in Azure."
folder: "SQL"
---

# Table: azure_compute_image - Query Azure Compute Images using SQL

Azure Compute Images are pre-configured operating system images used to create virtual machines within the Azure platform. These images include a set of pre-installed applications and configurations, which can be used to quickly deploy new virtual machines. Azure Compute Images provide a convenient way to manage and maintain consistent configurations across multiple virtual machines.

## Table Usage Guide

The `azure_compute_image` table provides insights into Azure Compute Images within Azure. As a DevOps engineer, explore image-specific details through this table, including image versions, operating system types, and associated metadata. Utilize it to uncover information about images, such as those with specific versions, the operating system types, and the verification of configurations.

## Examples

### Basic compute image info
This query allows you to gain insights into the basic information of your Azure compute images, including their name, type, and region. It's particularly useful when you need to understand the specifics of the source virtual machine associated with each image.

```sql+postgres
select
  name,
  type,
  region,
  hyper_v_generation,
  split_part(source_virtual_machine_id, '/', 9) as source_virtual_machine
from
  azure_compute_image;
```

```sql+sqlite
Error: SQLite does not support split_part function.
```

### Storage profile's OS disk info of each compute image
Determine the size, type, and status of your operating system disk within each compute image in Azure. This query can help you manage your storage resources more effectively by identifying potential areas for optimization.

```sql+postgres
select
  name,
  storage_profile_os_disk_size_gb,
  storage_profile_os_disk_snapshot_id,
  storage_profile_os_disk_storage_account_type,
  storage_profile_os_disk_state,
  storage_profile_os_disk_type
from
  azure_compute_image;
```

```sql+sqlite
select
  name,
  storage_profile_os_disk_size_gb,
  storage_profile_os_disk_snapshot_id,
  storage_profile_os_disk_storage_account_type,
  storage_profile_os_disk_state,
  storage_profile_os_disk_type
from
  azure_compute_image;
```

### List of compute images where disk storage type is Premium_LRS
This example helps you identify the compute images that use Premium_LRS as their disk storage type. Understanding the storage type of your compute images can assist in optimizing performance and cost in your Azure environment.

```sql+postgres
select
  name,
  split_part(disk -> 'managedDisk' ->> 'id', '/', 9) as disk_name,
  disk ->> 'storageAccountType' as storage_account_type,
  disk ->> 'diskSizeGB' as disk_size_gb,
  disk ->> 'caching' as caching
from
  azure_compute_image
  cross join jsonb_array_elements(storage_profile_data_disks) as disk
where
  disk ->> 'storageAccountType' = 'Premium_LRS';
```

```sql+sqlite
Error: SQLite does not support split or string_to_array functions.
```

### List of compute images which do not have owner or app_id tag key
Explore which Azure compute images lack either an owner or app_id tag, helping to identify potential issues with image management and organization. This can be useful for maintaining a clean and efficient cloud environment.

```sql+postgres
select
  id,
  name
from
  azure_compute_image
where
  tags -> 'owner' is null
  or tags -> 'app_id' is null;
```

```sql+sqlite
select
  id,
  name
from
  azure_compute_image
where
  json_extract(tags, '$.owner') is null
  or json_extract(tags, '$.app_id') is null;
```