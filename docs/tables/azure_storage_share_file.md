---
title: "Steampipe Table: azure_storage_share_file - Query Azure Storage Files using SQL"
description: "Allows users to query Azure Storage Files, specifically retrieving details such as file name, share name, content length, last modified time, and more."
folder: "Storage"
---

# Table: azure_storage_share_file - Query Azure Storage Files using SQL

Azure Storage Files is a service within Microsoft Azure that offers fully managed file shares in the cloud accessible via the industry-standard Server Message Block (SMB) protocol. Azure file shares can be mounted concurrently by cloud or on-premises deployments of Windows, Linux, and macOS. It provides a simple, secure, and scalable solution for sharing data between applications running in your virtual machines.

## Table Usage Guide

The `azure_storage_share_file` table provides insights into Azure Storage Files within Microsoft Azure. As a DevOps engineer, explore file-specific details through this table, including file name, share name, content length, last modified time, and more. Utilize it to uncover information about files, such as those with large content length, the shares they are associated with, and their last modified time.

## Examples

### Basic info
Explore the settings and configurations of your Azure storage shares to understand their storage capacity and accessibility. This can help in optimizing storage usage and ensuring the right protocols are enabled for secure and efficient data access.

```sql+postgres
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

```sql+sqlite
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
Determine the areas in which file shares are set to the default 'TransactionOptimized' access tier in Azure storage. This can help identify potential areas for optimization and cost savings.

```sql+postgres
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

```sql+sqlite
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
Explore which file share within your Azure storage has the maximum quota. This is useful for understanding your storage usage and managing resources effectively.

```sql+postgres
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
order by
  share_quota desc limit 1;
```

```sql+sqlite
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
order by
  share_quota desc limit 1;
```