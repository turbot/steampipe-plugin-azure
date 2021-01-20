# Table: azure_compute_snapshot

A snapshot is a full, read-only copy of a virtual hard drive (VHD).

## Examples

### Disk info of each snapshot

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

```sql
select
  name,
  incremental
from
  azure_compute_snapshot
where
  incremental;
```