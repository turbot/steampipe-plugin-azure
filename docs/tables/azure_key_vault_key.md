---
title: "Steampipe Table: azure_key_vault_key - Query Azure Key Vault Keys using SQL"
description: "Allows users to query Azure Key Vault Keys."
---

# Table: azure_key_vault_key - Query Azure Key Vault Keys using SQL

Azure Key Vault is a service that safeguards cryptographic keys and secrets used by cloud applications and services. It provides secure key management, ensures that keys are available when needed, and prevents unauthorized access. Azure Key Vault Keys are the keys that are stored in the Azure Key Vault for use in applications and services.

## Table Usage Guide

The 'azure_key_vault_key' table provides insights into keys within Azure Key Vault. As a security engineer, explore key-specific details through this table, including the key type, key state, and associated metadata. Utilize it to uncover information about keys, such as those that are disabled, the verification of key attributes, and the creation and expiry dates. The schema presents a range of attributes of the Key Vault key for your analysis, like the key ID, creation date, updated date, and vault details.

## Examples

### Basic info
This query allows you to review the details of your Azure Key Vault keys. It is particularly useful in auditing and managing these keys by providing information such as their status, creation and modification dates, and location.

```sql
select
  name,
  vault_name,
  enabled,
  created_at,
  updated_at,
  key_type,
  location
from
  azure_key_vault_key;
```

### List disabled keys
Explore which keys in your Azure Key Vault are currently disabled. This can help in maintaining security by identifying and managing inactive keys.

```sql
select
  name,
  vault_name,
  enabled
from
  azure_key_vault_key
where
  not enabled;
```

### List keys with no expiration time set
Explore which keys in the Azure Key Vault have no expiration time set. This can help in identifying potential security risks, as keys without an expiration can be misused if they fall into the wrong hands.

```sql
select
  name,
  enabled,
  expires_at
from
  azure_key_vault_key
where
  expires_at is null;
```

### List keys which have never been updated
Explore which keys in Azure Key Vault are active but have never been modified since their creation. This helps in identifying unused or potentially obsolete keys, aiding in better security management.

```sql
select
  name,
  enabled,
  created_at,
  updated_at
from
  azure_key_vault_key
where
  enabled
  and age(updated_at, created_at) = '00:00:00';
```

### Count the number of keys by key vault
Explore which Azure Key Vault has the most keys, providing a useful overview of your key distribution and aiding in the management and organization of your security assets.

```sql
select
  vault_name,
  count(vault_name) as count
from
  azure_key_vault_key
group by
  vault_name;
```