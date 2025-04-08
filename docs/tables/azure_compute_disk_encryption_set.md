---
title: "Steampipe Table: azure_compute_disk_encryption_set - Query Azure Compute Disk Encryption Sets using SQL"
description: "Allows users to query Azure Compute Disk Encryption Sets, specifically the encryption settings and associated metadata, providing insights into data security and compliance."
folder: "Resource"
---

# Table: azure_compute_disk_encryption_set - Query Azure Compute Disk Encryption Sets using SQL

Azure Compute Disk Encryption Sets is a resource within Microsoft Azure that manages the encryption of Azure Disk Storage. It provides a centralized way to manage and enforce encryption policies for data at rest. Azure Compute Disk Encryption Sets help you meet organizational security and compliance commitments.

## Table Usage Guide

The `azure_compute_disk_encryption_set` table provides insights into encryption sets within Azure Compute Disk. As a security analyst, explore encryption set-specific details through this table, including encryption settings, associated keys, and metadata. Utilize it to uncover information about encryption sets, such as those with outdated keys, the associations between encryption sets and disks, and the verification of encryption policies.

## Examples

### Key vault associated with each disk encryption set
Determine the areas in which a specific key vault is associated with each disk encryption set. This can be useful for understanding the security configuration of your Azure resources and identifying potential vulnerabilities.

```sql+postgres
select
  name,
  split_part(active_key_source_vault_id, '/', 9) as vault_name,
  split_part(active_key_url, '/', 5) as key_name
from
  azure_compute_disk_encryption_set;
```

```sql+sqlite
Error: SQLite does not support split_part function.
```


### List of encryption sets which are not using customer managed key
Determine the areas in which disk encryption sets in Azure are not utilizing customer-managed keys. This is useful for identifying potential security vulnerabilities where data is not being encrypted using customer's keys.

```sql+postgres
select
  name,
  encryption_type
from
  azure_compute_disk_encryption_set
where
  (
    encryption_type <> 'EncryptionAtRestWithPlatformAndCustomerKeys'
    and encryption_type <> 'EncryptionAtRestWithCustomerKey'
  );
```

```sql+sqlite
select
  name,
  encryption_type
from
  azure_compute_disk_encryption_set
where
  (
    encryption_type != 'EncryptionAtRestWithPlatformAndCustomerKeys'
    and encryption_type != 'EncryptionAtRestWithCustomerKey'
  );
```


### Identity info of each disk encryption set
Assess the elements within each disk encryption set to gain insights into their identity information. This can help in managing and tracking the sets effectively across your network.

```sql+postgres
select
  name,
  identity_type,
  identity_principal_id,
  identity_tenant_id
from
  azure_compute_disk_encryption_set;
```

```sql+sqlite
select
  name,
  identity_type,
  identity_principal_id,
  identity_tenant_id
from
  azure_compute_disk_encryption_set;
```