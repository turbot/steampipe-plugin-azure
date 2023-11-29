---
title: "Steampipe Table: azure_storage_container - Query Azure Storage Containers using SQL"
description: "Allows users to query Azure Storage Containers."
---

# Table: azure_storage_container - Query Azure Storage Containers using SQL

Azure Storage Containers are a part of Azure Blob Storage, which provides scalable, secure, performance-efficient storage services in the cloud. The containers organize blobs in a similar way that directories organize files in a file system. They are useful in storing and managing data objects, such as text or binary data, which can be accessed from anywhere in the world via HTTP or HTTPS.

## Table Usage Guide

The 'azure_storage_container' table provides insights into Azure Storage Containers within Azure Blob Storage. As a DevOps engineer, explore container-specific details through this table, including metadata, properties, and associated storage account information. Utilize it to uncover information about containers, such as public access level, last modified time, and the lease status. The schema presents a range of attributes of the Azure Storage Container for your analysis, like the storage account name, resource group name, and associated tags.

## Examples

### Basic info
Explore the basic details of your Azure storage containers to identify their types and associated accounts. This could be beneficial for managing resources and ensuring correct account allocation.

```sql
select
  name,
  id,
  type,
  account_name
from
  azure_storage_container;
```

### List containers which are publicly accessible
Discover the segments that are publicly accessible within your Azure storage containers to ensure data privacy and security. This query is useful for identifying potential vulnerabilities and implementing necessary access control measures.

```sql
select
  name,
  id,
  type,
  account_name,
  public_access
from
  azure_storage_container
where
  public_access <> 'None';
```

### List containers with legal hold enabled
Explore which Azure storage containers have the legal hold feature enabled. This is useful for identifying instances where data preservation is enforced for compliance or litigation purposes.

```sql
select
  name,
  id,
  type,
  account_name,
  has_legal_hold
from
  azure_storage_container
where
  has_legal_hold;
```

### List containers which are either leased or have a broken lease state
Explore which Azure storage containers are currently leased or have a broken lease state. This query is useful for managing resources and troubleshooting issues related to container leases.

```sql
select
  name,
  id,
  type,
  account_name,
  lease_state
from
  azure_storage_container
where
  lease_state = 'Leased'
  or lease_state = 'Broken';
```

### List containers with infinite lease duration
Explore which Azure storage containers have been set with an unlimited lease duration. This can help in managing storage resources effectively and identifying areas that may require attention to prevent unnecessary storage consumption.

```sql
select
  name,
  id,
  type,
  account_name,
  lease_duration
from
  azure_storage_container
where
  lease_duration = 'Infinite';
```

### List containers with a remaining retention period of 7 days
Explore which Azure storage containers have a remaining retention period of exactly 7 days. This is useful for managing resources and planning ahead for storage needs or potential data loss.

```sql
select
  name,
  id,
  type,
  account_name,
  remaining_retention_days
from
  azure_storage_container
where
  remaining_retention_days = 7;
```

### List containers ImmutabilityPolicy details
Explore the immutability policy details of your Azure storage containers to understand their data preservation settings. This can help in maintaining data integrity and ensuring compliance with data retention policies.

```sql
select
  name,
  account_name,
  jsonb_pretty(immutability_policy) as immutability_policy
from
  azure_storage_container;
```