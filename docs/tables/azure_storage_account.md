# Table: azure_storage_account

An Azure storage account contains all of your Azure Storage data objects: blobs, files, queues, tables, and disks.

## Examples

### Basic info

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


### List storage accounts with versioning disabled

```sql
select
  name,
  blob_versioning_enabled
from
  azure_storage_account
where
  not blob_versioning_enabled;
```


### List storage accounts with blob soft delete disabled

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


### List storage accounts that allow blob public access

```sql
select
  name,
  allow_blob_public_access
from
  azure_storage_account
where
  allow_blob_public_access;
```


### List storage accounts with encryption in transit disabled

```sql
select
  name,
  enable_https_traffic_only
from
  azure_storage_account
where
  not enable_https_traffic_only;
```


### List storage accounts that do not have a cannot-delete lock

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


### List storage accounts with queue logging enabled

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


### List storage accounts without lifecycle

```sql
select
  name,
  lifecycle_management_policy -> 'properties' -> 'policy' -> 'rules' as lifecycle_rules
from
  azure_storage_account
where
  lifecycle_management_policy -> 'properties' -> 'policy' -> 'rules' is null;
```

### List diagnostic settings details

```sql
select
  name,
  jsonb_pretty(diagnostic_settings) as diagnostic_settings
from
  azure_storage_account;
```