---
title: "Steampipe Table: azure_key_vault_deleted_vault - Query Azure Key Vault Deleted Vaults using SQL"
description: "Allows users to query Azure Key Vault Deleted Vaults for detailed information."
---

# Table: azure_key_vault_deleted_vault - Query Azure Key Vault Deleted Vaults using SQL

Azure Key Vault is a cloud service for securely storing and accessing secrets. A secret is anything that you want to tightly control access to, such as API keys, passwords, certificates, or cryptographic keys. Azure Key Vault Deleted Vaults are vaults that have been deleted but are still recoverable for a certain period of time.

## Table Usage Guide

The 'azure_key_vault_deleted_vault' table provides insights into deleted vaults within Azure Key Vault. As a security analyst or DevOps engineer, explore deleted vault-specific details through this table, including deletion date, recovery level, and scheduled purge date. Utilize it to uncover information about deleted vaults, such as those scheduled for permanent deletion or those still recoverable. The schema presents a range of attributes of the deleted vault for your analysis, like the vault name, location, deletion date, and scheduled purge date.

## Examples

### Basic info
Explore which Azure Key Vault resources have been deleted and when they are scheduled for permanent removal. This can be useful for auditing purposes or to recover resources before they are permanently purged.

```sql
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
Identify instances where Azure Key Vaults have been deleted and are scheduled for purging in more than a day. This can be useful in assessing data cleanup strategies and preventing accidental loss of important keys.

```sql
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