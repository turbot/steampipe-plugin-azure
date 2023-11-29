---
title: "Steampipe Table: azure_storage_blob - Query Azure Storage Blobs using SQL"
description: "Allows users to query Azure Storage Blobs."
---

# Table: azure_storage_blob - Query Azure Storage Blobs using SQL

Azure Storage Blobs are scalable, object storage for unstructured data. They are ideal for serving images or documents directly to a browser, storing files for distributed access, streaming video and audio, writing to log files, storing data for backup and restore, disaster recovery, and archiving. Azure Storage Blobs are accessible from anywhere in the world via HTTP or HTTPS.

## Table Usage Guide

The 'azure_storage_blob' table provides insights into the storage blobs within Azure Storage. As a DevOps engineer, explore blob-specific details through this table, including type, content settings, and associated metadata. Utilize it to uncover information about blobs, such as their lease status, server encrypted status, and the last modified date. The schema presents a range of attributes of the Azure Storage Blob for your analysis, like the blob name, creation date, content type, and associated tags.

## Examples

### Basic info
Explore the details of specific Azure storage blobs within a designated resource group, storage account, and region. This is useful for managing and organizing your data storage in Azure, particularly when dealing with large numbers of blobs.

```sql
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
Discover the segments that contain snapshot type blobs with import data in a specific Azure storage account located in a certain region. This could be useful to assess the elements within a particular resource group for better data management and security compliance.

```sql
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