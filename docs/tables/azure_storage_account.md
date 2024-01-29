---
title: "Steampipe Table: azure_storage_account - Query Azure Storage Accounts using SQL"
description: "Allows users to query Azure Storage Accounts, providing detailed information about each storage account within the Azure subscription."
---

# Table: azure_storage_account - Query Azure Storage Accounts using SQL

Azure Storage Account is a service within Microsoft Azure that provides scalable and secure data storage. It offers services like Blob Storage, File Storage, Queue Storage, and Table Storage. Azure Storage Account supports both Standard and Premium storage account types, allowing users to store large amounts of unstructured and structured data.

## Table Usage Guide

The `azure_storage_account` table provides insights into Storage Accounts within Microsoft Azure. As a Cloud Architect or DevOps engineer, explore account-specific details through this table, including the storage account type, creation date, access tier, and associated metadata. Utilize it to uncover information about storage accounts, such as their replication strategy, the network rules set, and the status of secure transfer.

## Examples

### Basic info
Explore the different tiers and locations of your Azure storage accounts. This can help you understand your storage distribution and make informed decisions about resource allocation.

```sql+postgres
select
  name,
  sku_name,
  sku_tier,
  primary_location,
  secondary_location
from
  azure_storage_account;
```

```sql+sqlite
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
Explore which Azure storage accounts have not enabled blob versioning. This is useful for identifying potential vulnerabilities in data backup and recovery systems.

```sql+postgres
select
  name,
  blob_versioning_enabled
from
  azure_storage_account
where
  not blob_versioning_enabled;
```

```sql+sqlite
select
  name,
  blob_versioning_enabled
from
  azure_storage_account
where
  blob_versioning_enabled is not 1;
```

### List storage accounts with blob soft delete disabled
Determine the areas in which storage accounts have the blob soft delete feature disabled. This is useful for identifying potential risk points where data might be permanently lost if accidentally deleted.

```sql+postgres
select
  name,
  blob_soft_delete_enabled,
  blob_soft_delete_retention_days
from
  azure_storage_account
where
  not blob_soft_delete_enabled;
```

```sql+sqlite
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
Determine the areas in which your Azure storage accounts are configured to allow public access to blobs. This can be used to identify potential security risks and ensure appropriate access controls are in place.

```sql+postgres
select
  name,
  allow_blob_public_access
from
  azure_storage_account
where
  allow_blob_public_access;
```

```sql+sqlite
select
  name,
  allow_blob_public_access
from
  azure_storage_account
where
  allow_blob_public_access;
```

### List storage accounts with encryption in transit disabled
Determine the areas in which data security may be compromised due to the lack of encryption during data transit in your Azure storage accounts. This query is useful to identify potential vulnerabilities and enhance your security measures.

```sql+postgres
select
  name,
  enable_https_traffic_only
from
  azure_storage_account
where
  not enable_https_traffic_only;
```

```sql+sqlite
select
  name,
  enable_https_traffic_only
from
  azure_storage_account
where
  enable_https_traffic_only = 0;
```

### List storage accounts that do not have a cannot-delete lock
Determine the areas in which storage accounts in Azure lack a 'cannot-delete' lock, which could potentially leave them vulnerable to unintentional deletion or modification. This query is useful for identifying and rectifying potential security risks within your storage management system.

```sql+postgres
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

```sql+sqlite
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
Discover the segments that have all types of queue logging enabled in their Azure storage accounts. This is useful to assess the storage accounts that are actively tracking and recording all queue activities for auditing or troubleshooting purposes.

```sql+postgres
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

```sql+sqlite
select
  name,
  queue_logging_delete,
  queue_logging_read,
  queue_logging_write
from
  azure_storage_account
where
  queue_logging_delete = 1
  and queue_logging_read = 1
  and queue_logging_write = 1;
```

### List storage accounts without lifecycle
Determine the storage accounts that lack a lifecycle management policy. This is useful for identifying potential risks or inefficiencies related to data retention and storage management.

```sql+postgres
select
  name,
  lifecycle_management_policy -> 'properties' -> 'policy' -> 'rules' as lifecycle_rules
from
  azure_storage_account
where
  lifecycle_management_policy -> 'properties' -> 'policy' -> 'rules' is null;
```

```sql+sqlite
select
  name,
  json_extract(lifecycle_management_policy, '$.properties.policy.rules') as lifecycle_rules
from
  azure_storage_account
where
  json_extract(lifecycle_management_policy, '$.properties.policy.rules') is null;
```

### List diagnostic settings details
Explore the diagnostic settings of your Azure storage accounts to gain insights into their configurations. This is beneficial to ensure optimal settings are in use for efficient data storage and management.

```sql+postgres
select
  name,
  jsonb_pretty(diagnostic_settings) as diagnostic_settings
from
  azure_storage_account;
```

```sql+sqlite
select
  name,
  diagnostic_settings
from
  azure_storage_account;
```

### List storage accounts with replication but unavailable secondary
Determine the areas in which Azure storage accounts have available primary status but unavailable secondary status, specifically within the 'Standard_GRS' and 'Standard_RAGRS' SKU categories. This is useful for identifying potential risk areas in your storage infrastructure where data replication might not be functioning as expected.

```sql+postgres
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
  and sku_name in ('Standard_GRS', 'Standard_RAGRS');
```

```sql+sqlite
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
  and sku_name in ('Standard_GRS', 'Standard_RAGRS');
```

### Get table properties of storage accounts
Explore the properties of your storage accounts to gain insights into their configuration. This can help you understand and manage your access and retention policies, as well as monitor their usage metrics.

```sql+postgres
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

```sql+sqlite
select
  name,
  json_extract(table_properties, '$.Cors') as table_logging_cors,
  json_extract(table_properties, '$.Logging.Read') as table_logging_read,
  json_extract(table_properties, '$.Logging.Write') as table_logging_write,
  json_extract(table_properties, '$.Logging.Delete') as table_logging_delete,
  json_extract(table_properties, '$.Logging.Version') as table_logging_version,
  json_extract(table_properties, '$.Logging.RetentionPolicy') as table_logging_retention_policy,
  json_extract(table_properties, '$.HourMetrics.Enabled') as table_hour_metrics_enabled,
  json_extract(table_properties, '$.HourMetrics.IncludeAPIs') as table_hour_metrics_include_ap_is,
  json_extract(table_properties, '$.HourMetrics.Version') as table_hour_metrics_version,
  json_extract(table_properties, '$.HourMetrics.RetentionPolicy') as table_hour_metrics_retention_policy,
  json_extract(table_properties, '$.MinuteMetrics.Enabled') as table_minute_metrics_enabled,
  json_extract(table_properties, '$.MinuteMetrics.IncludeAPIs') as table_minute_metrics_include_ap_is,
  json_extract(table_properties, '$.MinuteMetrics.Version') as table_minute_metrics_version,
  json_extract(table_properties, '$.MinuteMetrics.RetentionPolicy') as table_minute_metrics_retention_policy
from
  azure_storage_account;
```