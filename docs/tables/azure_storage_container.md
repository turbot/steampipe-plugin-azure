---
title: "Steampipe Table: azure_storage_container - Query Azure Storage Containers using SQL"
description: "Allows users to query Azure Storage Containers. The table provides details about each container in the Azure Storage Account, including metadata, public access level, and more."
---

# Table: azure_storage_container - Query Azure Storage Containers using SQL

Azure Storage Containers are a part of Azure Blob Storage service. They are used to organize blobs in a similar way as a directory in a file system. Containers provide a grouping of a set of blobs, and all blobs must be in a container.

## Table Usage Guide

The `azure_storage_container` table provides insights into Azure Storage Containers within Azure Blob Storage service. As a data engineer, explore container-specific details through this table, including metadata, public access level, and more. Utilize it to uncover information about containers, such as those with public access, the metadata associated with containers, and the verification of access policies.

## Examples

### Basic info
Explore which Azure storage containers are linked to your account. This can help in managing resources and identifying potential areas for optimization or restructuring.

```sql+postgres
select
  name,
  id,
  type,
  account_name
from
  azure_storage_container;
```

```sql+sqlite
select
  name,
  id,
  type,
  account_name
from
  azure_storage_container;
```

### List containers which are publicly accessible
Explore which Azure storage containers are set to public access, allowing you to identify potential security risks and rectify them to prevent unauthorized access to sensitive data.

```sql+postgres
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

```sql+sqlite
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
Discover the segments that have legal hold enabled in their Azure storage containers. This is beneficial for understanding which areas have additional data preservation measures in place for legal or compliance reasons.

```sql+postgres
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

```sql+sqlite
select
  name,
  id,
  type,
  account_name,
  has_legal_hold
from
  azure_storage_container
where
  has_legal_hold = 1;
```

### List containers which are either leased or have a broken lease state
Determine the areas in which Azure storage containers are either currently leased or have a broken lease state. This is useful for managing resources and identifying potential issues with container leases.

```sql+postgres
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

```sql+sqlite
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
Discover the segments that have an unlimited lease duration in Azure Storage, helping you identify potential areas for cost optimization and better resource management.

```sql+postgres
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

```sql+sqlite
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
Determine the areas in which Azure storage containers are nearing the end of their retention period. This is useful for proactive management of storage resources, allowing you to take timely action before the containers expire.

```sql+postgres
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

```sql+sqlite
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
Analyze the settings to understand the immutability policies of your Azure storage containers. This can help you manage data retention and protect your data from being modified or deleted.

```sql+postgres
select
  name,
  account_name,
  jsonb_pretty(immutability_policy) as immutability_policy
from
  azure_storage_container;
```

```sql+sqlite
select
  name,
  account_name,
  immutability_policy
from
  azure_storage_container;
```