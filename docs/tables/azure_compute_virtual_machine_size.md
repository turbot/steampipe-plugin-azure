---
title: "Steampipe Table: azure_compute_virtual_machine_size - Query Azure Virtual Machine Sizes using SQL"
description: "Allows users to query available virtual machine sizes in Azure, providing details on cores, memory, disk sizes, and more."
folder: "Compute"
---

# Table: azure_compute_virtual_machine_size - Query Azure Virtual Machine Sizes using SQL

Azure virtual machine sizes determine the computing resources, such as memory, CPU cores, and disk capacity, available to your virtual machines (VMs). The `azure_compute_virtual_machine_size` table in Steampipe allows you to query and explore the different virtual machine sizes available in Azure, enabling you to compare memory, disk sizes, and other attributes to choose the right VM for your workload.

## Table Usage Guide

The `azure_compute_virtual_machine_size` table is useful for cloud architects, administrators, and DevOps engineers who need insights into the available VM sizes in Azure. You can query attributes like the number of cores, memory, and disk sizes to make informed decisions about the resources required for your virtual machines.

## Examples

### Basic VM size information
Retrieve basic information about the virtual machine sizes, including name, region, and number of cores.

```sql+postgres
select
  name,
  region,
  number_of_cores,
  memory_in_mb
from
  azure_compute_virtual_machine_size;
```

```sql+sqlite
select
  name,
  region,
  number_of_cores,
  memory_in_mb
from
  azure_compute_virtual_machine_size;
```

### List VM sizes with high memory
Identify VM sizes that have more than a specified amount of memory (e.g., 32 GB).

```sql+postgres
select
  name,
  memory_in_mb,
  number_of_cores
from
  azure_compute_virtual_machine_size
where
  memory_in_mb > 32768;
```

```sql+sqlite
select
  name,
  memory_in_mb,
  number_of_cores
from
  azure_compute_virtual_machine_size
where
  memory_in_mb > 32768;
```

### List VM sizes with large OS disk size
Fetch VM sizes that support large OS disk sizes (e.g., over 500 GB).

```sql+postgres
select
  name,
  os_disk_size_in_mb,
  max_data_disk_count
from
  azure_compute_virtual_machine_size
where
  os_disk_size_in_mb > 512000;
```

```sql+sqlite
select
  name,
  os_disk_size_in_mb,
  max_data_disk_count
from
  azure_compute_virtual_machine_size
where
  os_disk_size_in_mb > 512000;
```

### List VM sizes by core count
Retrieve a list of VM sizes based on the number of cores available.

```sql+postgres
select
  name,
  number_of_cores,
  region
from
  azure_compute_virtual_machine_size
where
  number_of_cores >= 8;
```

```sql+sqlite
select
  name,
  number_of_cores,
  region
from
  azure_compute_virtual_machine_size
where
  number_of_cores >= 8;
```

### List VM sizes with a specific data disk count
Identify VM sizes that support a certain number of data disks.

```sql+postgres
select
  name,
  max_data_disk_count,
  resource_disk_size_in_mb
from
  azure_compute_virtual_machine_size
where
  max_data_disk_count >= 8;
```

```sql+sqlite
select
  name,
  max_data_disk_count,
  resource_disk_size_in_mb
from
  azure_compute_virtual_machine_size
where
  max_data_disk_count >= 8;
```