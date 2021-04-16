# Table: azure_storage_account

An Azure storage account contains all of your Azure Storage data objects: blobs, files, queues, tables, and disks.

## Examples

### List of storage accounts where versioning is not enabled

```sql
select
  name,
  blob_versioning_enabled
from
  azure_storage_account
where
  not blob_versioning_enabled;
```


### List of storage accounts where Blob Soft delete is not enabled

```sql
select
  name,
  blob_soft_delete_enabled,
  blob_soft_delete_retention_days
from
  azure_storage_account
where
  not blob_soft_delete_enabled;
```


### List of storage accounts which allows blob public access

```sql
select
  name,
  allow_blob_public_access
from
  azure_storage_account
where
  allow_blob_public_access;
```


### List of storage accounts where encryption in transit is not enabled

```sql
select
  name,
  enable_https_traffic_only
from
  azure_storage_account
where
  not enable_https_traffic_only;
```


### Storage type info for storage accounts

```sql
select
  name,
  sku_name,
  sku_tier,
  primary_location,
  secondary_location
from
  azure_storage_account;
```


### Storage accounts which are not locked with `CanNotDelete` lock

```sql
select
  sg.name,
  ml.scope,
  ml.lock_level,
  ml.notes
from
  azure_storage_account as sg
  left join azure_management_lock as ml on lower(sg.id) = lower(ml.scope)
where
  (
    (ml.lock_level is null)
    or(ml.lock_level = 'ReadOnly')
  );
```


### List of storage accounts with logging enabled for Queue service for read, write, and delete requests

```sql
select
  name,
  queue_logging_delete,
  queue_logging_read,
  queue_logging_write
from
  azure_storage_account
where
  queue_logging_delete
  and queue_logging_read
  and queue_logging_write;
```
