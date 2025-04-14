---
title: "Steampipe Table: azure_key_vault_key_version - Query Azure Key Vault Key Versions using SQL"
description: "Allows users to query Azure Key Vault Key Versions. This table provides detailed information about each version of a key in Azure Key Vault."
folder: "Key Vault"
---

# Table: azure_key_vault_key_version - Query Azure Key Vault Key Versions using SQL

Azure Key Vault is a service that provides a secure storage for secrets, keys, and certificates. It enables users to securely store and tightly control access to tokens, passwords, certificates, API keys, and other secrets. Azure Key Vault simplifies the process of meeting industry compliance and regulatory standards.

## Table Usage Guide

The `azure_key_vault_key_version` table provides insights into each version of a key stored in Azure Key Vault. As a security analyst, you can explore key-specific details through this table, including the key type, key state, and associated metadata. Use it to track the lifecycle of keys, verify the key state, and ensure compliance with security policies.

## Examples

### Basic info
Explore the status and details of various versions of keys in your Azure Key Vault. This will help you understand the lifecycle of your keys, their types, and their geographical locations, which can be crucial for managing security and compliance.

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
  azure_key_vault_key_version;
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
  azure_key_vault_key_version;
```

### List disabled key versions
Identify instances where key versions are disabled in Azure Key Vault, allowing you to review and manage your keys' security settings effectively.

```sql+postgres
select
  name,
  key_name,
  vault_name,
  enabled
from
  azure_key_vault_key_version
where
  not enabled;
```

```sql+sqlite
select
  name,
  key_name,
  vault_name,
  enabled
from
  azure_key_vault_key_version
where
  enabled = 0;
```

### List keys versions with no expiration time set
Explore which versions of keys in Azure Key Vault have not been assigned an expiration time. This is useful for identifying potential security risks and ensuring key management best practices are being followed.

```sql+postgres
select
  name,
  enabled,
  expires_at
from
  azure_key_vault_key_version
where
  expires_at is null;
```

```sql+sqlite
select
  name,
  enabled,
  expires_at
from
  azure_key_vault_key_version
where
  expires_at is null;
```

### Count the number of versions by key
Assess the elements within your Azure Key Vault by determining the quantity of versions for each key. This can be beneficial in managing key rotations and understanding the lifecycle of each key.

```sql+postgres
select
  key_name,
  count(name) as key_version_count
from
  azure_key_vault_key_version
group by
  key_name;
```

```sql+sqlite
select
  key_name,
  count(name) as key_version_count
from
  azure_key_vault_key_version
group by
  key_name;
```