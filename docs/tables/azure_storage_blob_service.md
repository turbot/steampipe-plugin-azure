---
title: "Steampipe Table: azure_storage_blob_service - Query Azure Storage Blob Services using SQL"
description: "Allows users to query Azure Blob Services."
---

# Table: azure_storage_blob_service - Query Azure Storage Blob Services using SQL

Azure Blob storage is a service for storing large amounts of unstructured object data, such as text or binary data, that can be accessed from anywhere in the world via HTTP or HTTPS. You can use Blob storage to expose data publicly to the world, or to store application data privately. Common uses of Blob storage include serving images or documents directly to a browser, storing files for distributed access, streaming video and audio, writing to log files, storing data for backup and restore, disaster recovery, and archiving.

## Table Usage Guide

The 'azure_storage_blob_service' table provides insights into Blob Services within Azure Storage. As a DevOps engineer, explore service-specific details through this table, including the status of blob services, the CORS rules in place, and associated metadata. Utilize it to uncover information about each blob service, such as its default service version, whether or not it supports HTTPS traffic only, and the last modified time. The schema presents a range of attributes of the blob service for your analysis, like the storage account name, resource group, and Azure region.

## Examples

### Basic info
Analyze the settings to understand the tier and region of your Azure storage accounts. This can help you manage resources and costs effectively.

```sql
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
Determine the areas in which the delete retention policy is not enabled for your Azure storage blob service. This query is useful for identifying potential vulnerabilities and maintaining data security within your storage services.

```sql
select
  name,
  storage_account_name,
  delete_retention_policy -> 'enabled' as delete_retention_policy_enabled
from
  azure_storage_blob_service
where
  delete_retention_policy -> 'enabled' = 'false';
```

### List of storage blob service where versioning is not enabled
Identify instances where Azure Blob Storage services do not have versioning enabled. This is useful for ensuring data recovery options are in place, as versioning allows restoration of previous versions of blobs in the event of accidental deletion or alteration.

```sql
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
Explore the Cross-Origin Resource Sharing (CORS) rules for your Azure Storage Blob Service to understand the permissions and restrictions in place. This can help ensure secure data transactions and identify potential areas for security optimization.

```sql
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