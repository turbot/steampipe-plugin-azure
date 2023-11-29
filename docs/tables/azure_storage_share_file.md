---
title: "Steampipe Table: azure_storage_share_file - Query Azure Storage File Shares using SQL"
description: "Allows users to query Azure Storage File Shares, providing details about each file stored within these resources."
---

# Table: azure_storage_share_file - Query Azure Storage File Shares using SQL

Azure Storage File Shares service is a feature of Azure Storage that allows you to create a serverless file share in the cloud. It provides highly available network file shares that can be accessed by using the standard Server Message Block (SMB) protocol. Azure File Shares can be used to replace or supplement on-premises file servers or NAS devices.

## Table Usage Guide

The 'azure_storage_share_file' table provides insights into the files stored within Azure Storage File Shares. As a DevOps engineer, explore file-specific details through this table, including the file's URL, content type, last modification time, and associated metadata. Utilize it to uncover information about each file, such as its size, type, and any lease status. The schema presents a range of attributes of the file for your analysis, like the file's Etag, content MD5, and whether it is a directory or not.

## Examples

### Basic info
Explore which storage shares are available in your Azure account, focusing on their types and capabilities. This can help you understand your storage utilization and optimize resource allocation.

```sql
select
  name,
  storage_account_name,
  type,
  access_tier,
  share_quota,
  enabled_protocols
from
  azure_storage_share_file;
```

### List file shares with default access tier
Explore which file shares are set with a default access tier of 'TransactionOptimized'. This is useful for understanding how your storage is optimized for transactional workloads.

```sql
select
  name,
  storage_account_name,
  type,
  access_tier,
  access_tier_change_time,
  share_quota,
  enabled_protocols
from
  azure_storage_share_file
where
  access_tier = 'TransactionOptimized';
```

### Get file share with maximum share quota
Analyze the settings to understand which file share has the largest quota within your Azure storage. This can be useful to determine where the majority of your storage resources are allocated.

```sql
select
  name,
  storage_account_name,
  type,
  access_tier,
  access_tier_change_time,
  share_quota,
  enabled_protocols
from
  azure_storage_share_file
order by share_quota desc limit 1;
```