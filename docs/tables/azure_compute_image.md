# Table: azure_compute_image

Compute Engine offers many preconfigured public images that have compatible Linux or Windows operating systems. Compute Engine uses selected image to create a persistent boot disk for each instance.

## Examples

### Basic compute image info

```sql
select
  name,
  type,
  location,
  hyper_v_generation,
  split_part(source_virtual_machine_id, '/', 9) as source_virtual_machine
from
  azure_compute_image;
```


### Storage profile's OS disk info of each compute image

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
