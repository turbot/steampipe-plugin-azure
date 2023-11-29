---
title: "Steampipe Table: azure_storage_account - Query Azure Storage Accounts using SQL"
description: "Allows users to query Azure Storage Accounts."
---

# Table: azure_storage_account - Query Azure Storage Accounts using SQL

Azure Storage Account is a service provided by Microsoft Azure that offers highly scalable and secure data storage. It allows you to store and retrieve large amounts of unstructured data, such as documents and media files, and structured data, such as databases. Azure Storage Account supports different data types including blobs, files, queues, tables, and disks.

## Table Usage Guide

The 'azure_storage_account' table provides insights into Storage Accounts within Microsoft Azure. As a DevOps engineer, explore account-specific details through this table, including creation time, primary location, and associated metadata. Utilize it to uncover information about accounts, such as those with public access, the replication type, and the status of primary and secondary locations. The schema presents a range of attributes of the Storage Account for your analysis, like the account name, resource group, subscription ID, and associated tags.

## Examples

### Basic info
Explore which Azure storage accounts are in use, their associated SKU names and tiers, and their primary and secondary locations. This can help in understanding the distribution and classification of storage resources within your Azure environment.

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
Discover the segments that have disabled versioning within their storage accounts, enabling you to identify potential risks and ensure data recovery options.

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
Explore which Azure storage accounts have the blob soft delete feature disabled. This is useful in identifying potential data loss risks, as these accounts do not have a recovery option for accidentally deleted blobs.

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
Explore which Azure storage accounts permit public access to blob data. This is useful for assessing potential security risks and ensuring appropriate access controls are in place.

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
Explore which Azure storage accounts lack encryption in transit, a feature crucial for maintaining data security during transmission. This query is useful for identifying potential security vulnerabilities within your cloud storage infrastructure.

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
Analyze the settings to understand which storage accounts lack a 'cannot-delete' lock, therefore potentially posing a risk of accidental deletion. This query is useful in identifying areas that need improved security measures.

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
Explore which Azure storage accounts have queue logging enabled for all actions, such as delete, read, and write. This is useful in monitoring activity and maintaining security within your storage accounts.

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
Discover the storage accounts that lack a lifecycle management policy. This is useful for identifying areas where data retention and deletion policies may not be properly enforced, potentially leading to unnecessary storage costs or compliance issues.

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
Explore the diagnostic settings of your Azure storage accounts. This can help you better understand and manage the logging and monitoring capabilities of your storage resources.

```sql
select
  name,
  jsonb_pretty(diagnostic_settings) as diagnostic_settings
from
  azure_storage_account;
```

### List storage accounts with replication but unavailable secondary
Determine the areas in which storage accounts have replication enabled but the secondary is unavailable. This is useful to identify potential risks and ensure data redundancy in case of primary failure.

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
Explore the properties of your storage accounts to understand their configurations, such as logging settings and metrics, which can help in optimizing storage usage and improving data management practices.

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