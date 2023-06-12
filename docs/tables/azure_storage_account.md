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

### List storage accounts with replication but unavailable secondary

```sql
select
  name,
  status_of_primary,
  status_of_secondary,
  sku_name
from
  azure_storage_account
where
  status_of_primary = 'available'
  and status_of_secondary != 'available'
  and sku_name in ('Standard_GRS', 'Standard_RAGRS')
```

### Get table properties of storage accounts

```sql
select
  name,
  table_properties -> 'Cors' as table_logging_cors,
  table_properties -> 'Logging' -> 'Read' as table_logging_read,
  table_properties -> 'Logging' -> 'Write' as table_logging_write,
  table_properties -> 'Logging' -> 'Delete' as table_logging_delete,
  table_properties -> 'Logging' ->> 'Version' as table_logging_version,
  table_properties -> 'Logging' -> 'RetentionPolicy' as table_logging_retention_policy,
  table_properties -> 'HourMetrics' -> 'Enabled' as table_hour_metrics_enabled,
  table_properties -> 'HourMetrics' -> 'IncludeAPIs' as table_hour_metrics_include_ap_is,
  table_properties -> 'HourMetrics' ->> 'Version' as table_hour_metrics_version,
  table_properties -> 'HourMetrics' -> 'RetentionPolicy' as table_hour_metrics_retention_policy,
  table_properties -> 'MinuteMetrics' -> 'Enabled' as table_minute_metrics_enabled,
  table_properties -> 'MinuteMetrics' -> 'IncludeAPIs' as table_minute_metrics_include_ap_is,
  table_properties -> 'MinuteMetrics' ->> 'Version' as table_minute_metrics_version,
  table_properties -> 'MinuteMetrics' -> 'RetentionPolicy' as table_minute_metrics_retention_policy
from
  azure_storage_account;
```