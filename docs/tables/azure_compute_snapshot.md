---
title: "Steampipe Table: azure_compute_snapshot - Query Azure Compute Snapshots using SQL"
description: "Allows users to query Azure Compute Snapshots"
---

# Table: azure_compute_snapshot - Query Azure Compute Snapshots using SQL

Azure Compute Snapshots are a point-in-time copy of data. They are used to back up data and can be used to restore a virtual machine to the state at the time of the snapshot. Snapshots are incremental, capturing only the changes since the last snapshot, and are thus space-efficient.

## Table Usage Guide

The 'azure_compute_snapshot' table provides insights into snapshots within Azure Compute. As a DevOps engineer, explore snapshot-specific details through this table, including snapshot state, creation time, and associated metadata. Utilize it to uncover information about snapshots, such as those that are incremental, the disk size, and the source disk. The schema presents a range of attributes of the snapshot for your analysis, like the snapshot ID, resource group, and associated tags.

## Examples

### Disk info of each snapshot
Analyze the settings to understand the disk information for each snapshot in Azure, including its size and encryption set ID, which can help in managing storage and security aspects. This is particularly useful in assessing the storage consumption and encryption status of each snapshot.

```sql
select
  name,
  split_part(disk_access_id, '/', 8) as disk_name,
  disk_encryption_set_id,
  disk_size_gb,
  region
from
  azure_compute_snapshot;
```


### List of snapshots which are publicly accessible
Discover the segments that contain snapshots which are publicly accessible. This query is useful to identify potential security risks by pinpointing areas where data might be exposed.

```sql
select
  name,
  network_access_policy
from
  azure_compute_snapshot
where
  network_access_policy = 'AllowAll';
```


### List of all incremental type snapshots
Gain insights into all snapshots that are incremental in nature within the Azure compute service. This is useful for managing storage and tracking changes made over time.

```sql
select
  name,
  incremental
from
  azure_compute_snapshot
where
  incremental;
```