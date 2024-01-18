---
title: "Steampipe Table: azure_key_vault_key - Query Azure Key Vault Keys using SQL"
description: "Allows users to query Azure Key Vault Keys, providing access to key details, including key type, key state, and key attributes."
---

# Table: azure_key_vault_key - Query Azure Key Vault Keys using SQL

Azure Key Vault is a service within Microsoft Azure that provides a secure store for secrets, keys, and certificates. It provides a centralized way to manage cryptographic keys and secrets in cloud applications, without having to maintain an in-house key management infrastructure. Azure Key Vault helps users safeguard cryptographic keys and secrets used by cloud apps and services.

## Table Usage Guide

The `azure_key_vault_key` table provides insights into keys within Azure Key Vault. As a security engineer, explore key-specific details through this table, including key type, key state, and key attributes. Utilize it to uncover information about keys, such as those with specific attributes, the state of the keys, and the verification of key properties.

## Examples

### Basic info
Explore the status and details of your Azure Key Vault keys to understand their configurations and keep track of their activity. This is useful for maintaining security and ensuring that keys are up-to-date and correctly enabled.

```sql+postgres
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

```sql+sqlite
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
Identify instances where Azure Key Vault keys are disabled to ensure proper security measures are in place and access control is effectively managed.

```sql+postgres
select
  name,
  vault_name,
  enabled
from
  azure_key_vault_key
where
  not enabled;
```

```sql+sqlite
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
Identify instances where certain keys within Azure's Key Vault service have not been assigned an expiration time. This could be useful in managing security practices, as keys without set expiration times could potentially pose a risk.

```sql+postgres
select
  name,
  enabled,
  expires_at
from
  azure_key_vault_key
where
  expires_at is null;
```

```sql+sqlite
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
Discover the keys in your Azure Key Vault that have remained unmodified since their creation. This can be useful to identify any keys that may have been overlooked or forgotten, ensuring all keys are up-to-date and secure.

```sql+postgres
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

```sql+sqlite
select
  name,
  enabled,
  created_at,
  updated_at
from
  azure_key_vault_key
where
  enabled
  and (strftime('%s', updated_at) - strftime('%s', created_at)) = 0;
```

### Count the number of keys by key vault
Determine the distribution of keys across various vaults to understand your security setup better. This can help identify any potential vaults that may be overloaded or underutilized.

```sql+postgres
select
  vault_name,
  count(vault_name) as count
from
  azure_key_vault_key
group by
  vault_name;
```

```sql+sqlite
select
  vault_name,
  count(vault_name) as count
from
  azure_key_vault_key
group by
  vault_name;
```