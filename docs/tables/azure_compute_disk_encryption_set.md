---
title: "Steampipe Table: azure_compute_disk_encryption_set - Query Azure Compute Disk Encryption Sets using SQL"
description: "Allows users to query Azure Compute Disk Encryption Sets."
---

# Table: azure_compute_disk_encryption_set - Query Azure Compute Disk Encryption Sets using SQL

Azure Disk Encryption is a capability that helps you encrypt your Windows and Linux IaaS virtual machine disks. Disk Encryption Sets are a resource in Azure that contain and manage the key for server-side encryption of Azure managed disks and snapshots. It simplifies the key management for disk encryption and allows you to use Customer Managed Keys for managed disks instead of platform-managed keys.

## Table Usage Guide

The 'azure_compute_disk_encryption_set' table provides insights into Disk Encryption Sets within Azure Compute. As a security engineer, explore Disk Encryption Set-specific details through this table, including the encryption type, key URL, and source vault. Utilize it to uncover information about encryption sets, such as those with server-side encryption and customer-managed keys. The schema presents a range of attributes of the Disk Encryption Set for your analysis, like the id, name, type, location, and associated tags.

## Examples

### Key vault associated with each disk encryption set
Identify the specific key vault associated with each disk encryption set in your Azure Compute environment. This is useful for managing and auditing your encryption keys and their usage.

```sql
select
  name,
  split_part(active_key_source_vault_id, '/', 9) as vault_name,
  split_part(active_key_url, '/', 5) as key_name
from
  azure_compute_disk_encryption_set;
```


### List of encryption sets which are not using customer managed key
Explore which encryption sets in Azure's Compute Disk Encryption are not utilizing customer-managed keys, providing a way to identify potential areas for enhancing data security practices.

```sql
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


### Identity info of each disk encryption set
Explore which disk encryption sets in your Azure Compute resources have specific identities associated with them. This can help in assessing security configurations and managing access control within your environment.

```sql
select
  name,
  identity_type,
  identity_principal_id,
  identity_tenant_id
from
  azure_compute_disk_encryption_set;
```