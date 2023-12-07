---
title: "Steampipe Table: azure_key_vault_deleted_vault - Query Azure Key Vaults using SQL"
description: "Allows users to query deleted Azure Key Vaults, providing insights into the historical and current deletions of Azure Key Vaults."
---

# Table: azure_key_vault_deleted_vault - Query Azure Key Vaults using SQL

Azure Key Vault is a cloud service for securely storing and accessing secrets. A secret is anything that you want to tightly control access to, such as API keys, passwords, certificates, or cryptographic keys. Azure Key Vault handles requesting and renewing Transport Layer Security (TLS) certificates.

## Table Usage Guide

The `azure_key_vault_deleted_vault` table provides insights into the deleted vaults within Azure Key Vault. As a security analyst, explore vault-specific details through this table, including deletion dates, recovery levels, and associated metadata. Utilize it to uncover information about deleted vaults, such as their scheduled purge dates, recovery ids, and the geographical location of the vaults.

## Examples

### Basic info
Discover the segments that have been deleted from your Azure Key Vault, including when they were deleted and when they are scheduled for permanent removal. This can be useful for auditing purposes, ensuring data integrity, and managing your digital assets.

```sql+postgres
select
  name,
  id,
  type,
  deletion_date,
  scheduled_purge_date
from
  azure_key_vault_deleted_vault;
```

```sql+sqlite
select
  name,
  id,
  type,
  deletion_date,
  scheduled_purge_date
from
  azure_key_vault_deleted_vault;
```

### List deleted vaults with scheduled purge date more than 1 day
Explore which Azure Key Vaults have been deleted but are scheduled for purge after more than one day. This can be useful for reviewing and managing your data retention and recovery strategies.

```sql+postgres
select
  name,
  id,
  type,
  deletion_date,
  scheduled_purge_date
from
  azure_key_vault_deleted_vault
where
  scheduled_purge_date > (current_date - interval '1' day);
```

```sql+sqlite
select
  name,
  id,
  type,
  deletion_date,
  scheduled_purge_date
from
  azure_key_vault_deleted_vault
where
  scheduled_purge_date > date('now','-1 day');
```