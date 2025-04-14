---
title: "Steampipe Table: azure_compute_snapshot - Query Azure Compute Snapshots using SQL"
description: "Allows users to query Azure Compute Snapshots, specifically the snapshot details including status, creation time, and disk size, providing insights into the state and usage of virtual machine disk snapshots."
folder: "Compute"
---

# Table: azure_compute_snapshot - Query Azure Compute Snapshots using SQL

Azure Compute Snapshots are a resource within Microsoft Azure that allows you to create point-in-time backups of Azure managed disks, native blobs, or other data. These snapshots are read-only and can be used for data backup, disaster recovery, or migrating data across different regions or subscriptions. Azure Compute Snapshots help ensure data durability and accessibility, and are crucial for maintaining data integrity and system resilience in Azure.

## Table Usage Guide

The `azure_compute_snapshot` table provides insights into the snapshots within Azure Compute. As a system administrator or DevOps engineer, explore snapshot-specific details through this table, including creation time, disk size, and status. Utilize it to uncover information about snapshots, such as those associated with specific virtual machines, the state of these snapshots, and their usage for data backup or disaster recovery.

## Examples

### Disk info of each snapshot
Discover the specifics of each snapshot in your Azure Compute service, such as disk size and region, to better manage your storage resources and understand where your data is physically located.

```sql+postgres
select
  name,
  split_part(disk_access_id, '/', 8) as disk_name,
  disk_encryption_set_id,
  disk_size_gb,
  region
from
  azure_compute_snapshot;
```

```sql+sqlite
Error: SQLite does not support split_part function.
```

### List of snapshots which are publicly accessible
Determine the areas in which snapshots are set to be publicly accessible. This is useful for identifying potential security risks and ensuring appropriate access controls are in place.

```sql+postgres
select
  name,
  network_access_policy
from
  azure_compute_snapshot
where
  network_access_policy = 'AllowAll';
```

```sql+sqlite
select
  name,
  network_access_policy
from
  azure_compute_snapshot
where
  network_access_policy = 'AllowAll';
```

### List of all incremental type snapshots
Explore which snapshots in your Azure Compute resources are of the incremental type. This can help manage storage efficiently and reduce costs by identifying and focusing on snapshots that only capture changes since the last snapshot.

```sql+postgres
select
  name,
  incremental
from
  azure_compute_snapshot
where
  incremental;
```

```sql+sqlite
select
  name,
  incremental
from
  azure_compute_snapshot
where
  incremental = 1;
```