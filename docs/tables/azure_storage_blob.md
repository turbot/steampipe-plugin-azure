---
title: "Steampipe Table: azure_storage_blob - Query Azure Storage Blobs using SQL"
description: "Allows users to query Azure Storage Blobs, specifically providing information about blob properties, blob metadata, and blob service properties."
folder: "Storage"
---

# Table: azure_storage_blob - Query Azure Storage Blobs using SQL

Azure Storage Blobs are objects in Azure Storage which can hold large amounts of text or binary data, ranging from hundreds of gigabytes to a petabyte. They are ideal for storing documents, videos, pictures, backups, and other unstructured text or binary data. Azure Storage Blobs are part of the Azure Storage service, which provides scalable, durable, and highly available storage for data.

## Table Usage Guide

The `azure_storage_blob` table provides insights into the blobs within Azure Storage. As a data analyst or a data engineer, you can explore blob-specific details through this table, including blob properties, blob metadata, and blob service properties. Utilize it to uncover information about blobs, such as those with public access, the types of blobs, and the verification of service properties.

## Examples

### Basic info
Explore which storage blobs within a specific resource group, storage account, and region in Azure. This is particularly useful to gain insights into your Azure storage configuration and identify instances where snapshots are being used.

```sql+postgres
select
  name,
  container_name,
  storage_account_name,
  region,
  type,
  is_snapshot
from
  azure_storage_blob
where
  resource_group = 'turbot'
  and storage_account_name = 'mystorageaccount'
  and region = 'eastus';
```

```sql+sqlite
select
  name,
  container_name,
  storage_account_name,
  region,
  type,
  is_snapshot
from
  azure_storage_blob
where
  resource_group = 'turbot'
  and storage_account_name = 'mystorageaccount'
  and region = 'eastus';
```

### List snapshot type blobs with import data
Explore the snapshot type blobs that have imported data in a specific Azure storage account and resource group. This can be useful for auditing purposes, such as ensuring that sensitive data is properly encrypted and stored in the correct region.

```sql+postgres
select
  name,
  type,
  access_tier,
  server_encrypted,
  metadata,
  creation_time,
  container_name,
  storage_account_name,
  resource_group,
  region
from
  azure_storage_blob
where
  resource_group = 'turbot'
  and storage_account_name = 'mystorageaccount'
  and region = 'eastus'
  and is_snapshot;
```

```sql+sqlite
select
  name,
  type,
  access_tier,
  server_encrypted,
  metadata,
  creation_time,
  container_name,
  storage_account_name,
  resource_group,
  region
from
  azure_storage_blob
where
  resource_group = 'turbot'
  and storage_account_name = 'mystorageaccount'
  and region = 'eastus'
  and is_snapshot = 1;
```