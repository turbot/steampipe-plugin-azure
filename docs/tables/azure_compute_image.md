---
title: "Steampipe Table: azure_compute_image - Query Azure Compute Images using SQL"
description: "Allows users to query Azure Compute Images."
---

# Table: azure_compute_image - Query Azure Compute Images using SQL

Azure Compute Images are resources within Microsoft Azure that represent a virtual machine's operating system, applications, and configuration settings. These images can be used to create multiple identical virtual machines within Azure. They provide an efficient way to package, provision, and manage VMs in your cloud environment.

## Table Usage Guide

The 'azure_compute_image' table provides insights into Azure Compute Images. As a DevOps engineer, explore image-specific details through this table, including publisher details, offer information, and associated metadata. Utilize it to uncover information about images, such as those used in multiple VM deployments, the publishers of these images, and the verification of image configurations. The schema presents a range of attributes of the Azure Compute Image for your analysis, like the image name, resource group, publisher, offer, SKU, and version.


## Examples

### Basic compute image info
Explore the types and regional distribution of virtual machine images in your Azure environment. This can help in understanding the configuration and usage patterns of virtual machines, thereby aiding in resource management and optimization.

```sql
select
  name,
  type,
  region,
  hyper_v_generation,
  split_part(source_virtual_machine_id, '/', 9) as source_virtual_machine
from
  azure_compute_image;
```


### Storage profile's OS disk info of each compute image
Determine the storage characteristics of each compute image in your Azure environment. This could help optimize storage utilization and cost by revealing details such as disk size, snapshot ID, storage account type, state, and disk type.

```sql
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
Determine the areas in which your compute images are using premium disk storage type. This query can be useful for understanding your storage usage and optimizing costs.

```sql
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


### List of compute images which do not have owner or app_id tag key
Discover the segments that lack either an 'owner' or 'app_id' tag key within your Azure compute images. This query can be used to identify potential gaps in your image tagging strategy, which can help improve resource tracking and management.

```sql
select
  id,
  name
from
  azure_compute_image
where
  tags -> 'owner' is null
  or tags -> 'app_id' is null;
```