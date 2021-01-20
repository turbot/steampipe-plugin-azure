# Table: azure_compute_disk

Azure Managed Disks are the new and recommended disk storage offering for use with Azure virtual machines for persistent storage of data.

## Examples

### List of all premium tier compute disks

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

```sql
select
  name,
  encryption_type
from
  azure_compute_disk
where
  encryption_type <> 'EncryptionAtRestWithCustomerKey';
```
