---
title: "Steampipe Table: azure_storage_blob_service - Query Azure Storage Blob Services using SQL"
description: "Allows users to query Azure Storage Blob Services, providing insights into storage accounts and their blob service properties."
folder: "Storage"
---

# Table: azure_storage_blob_service - Query Azure Storage Blob Services using SQL

Azure Storage Blob Service is a feature within Microsoft Azure that provides scalable, secure, performance-efficient storage for unstructured data. It is optimized for storing massive amounts of unstructured data, such as text or binary data, that can be accessed globally via HTTP or HTTPS. The service includes features to process data and build sophisticated analytics solutions, recover from disaster, and archive data.

## Table Usage Guide

The `azure_storage_blob_service` table provides insights into Azure Storage Blob Services within Microsoft Azure. As a data analyst or storage administrator, explore blob service-specific details through this table, including storage account name, resource group, and associated metadata. Utilize it to uncover information about blob services, such as default service version, change feed enabled status, and delete retention policy details.

## Examples

### Basic info
Analyze the settings to understand the distribution of your Azure storage blob services across different regions, their associated storage accounts, and their respective pricing tiers. This can help in optimizing resource allocation and cost management.

```sql+postgres
select
  name,
  storage_account_name,
  region,
  sku_name,
  sku_tier
from
  azure_storage_blob_service;
```

```sql+sqlite
select
  name,
  storage_account_name,
  region,
  sku_name,
  sku_tier
from
  azure_storage_blob_service;
```

### List of storage blob service where delete retention policy is not enabled
Identify Azure storage blob services that have not enabled the delete retention policy. This query is useful for pinpointing potential areas of risk where deleted data cannot be recovered.

```sql+postgres
select
  name,
  storage_account_name,
  delete_retention_policy -> 'enabled' as delete_retention_policy_enabled
from
  azure_storage_blob_service
where
  delete_retention_policy -> 'enabled' = 'false';
```

```sql+sqlite
select
  name,
  storage_account_name,
  json_extract(delete_retention_policy, '$.enabled') as delete_retention_policy_enabled
from
  azure_storage_blob_service
where
  json_extract(delete_retention_policy, '$.enabled') = 'false';
```

### List of storage blob service where versioning is not enabled
Explore which Azure storage blob services do not have versioning enabled. This is useful in identifying potential data loss risks due to accidental deletion or overwriting.

```sql+postgres
select
  name,
  storage_account_name,
  is_versioning_enabled
from
  azure_storage_blob_service
where
  not is_versioning_enabled;
```

```sql+sqlite
select
  name,
  storage_account_name,
  is_versioning_enabled
from
  azure_storage_blob_service
where
  not is_versioning_enabled;
```

### CORS rules info for storage blob service
This query is useful for gaining insights into the Cross-Origin Resource Sharing (CORS) rules set up for Azure's storage blob service. It's a practical tool for understanding what headers and methods are permitted, which headers are exposed, and the maximum age for these settings, thereby aiding in ensuring secure and efficient data transfers.

```sql+postgres
select
  name,
  storage_account_name,
  cors -> 'allowedHeaders' as allowed_headers,
  cors -> 'allowedMethods' as allowed_methods,
  cors -> 'allowedMethods' as allowed_methods,
  cors -> 'exposedHeaders' as exposed_headers,
  cors -> 'maxAgeInSeconds' as max_age_in_seconds
from
  azure_storage_blob_service
  cross join jsonb_array_elements(cors_rules) as cors;
```

```sql+sqlite
select
  name,
  storage_account_name,
  json_extract(cors.value, '$.allowedHeaders') as allowed_headers,
  json_extract(cors.value, '$.allowedMethods') as allowed_methods,
  json_extract(cors.value, '$.allowedMethods') as allowed_methods,
  json_extract(cors.value, '$.exposedHeaders') as exposed_headers,
  json_extract(cors.value, '$.maxAgeInSeconds') as max_age_in_seconds
from
  azure_storage_blob_service,
  json_each(cors_rules) as cors;
```