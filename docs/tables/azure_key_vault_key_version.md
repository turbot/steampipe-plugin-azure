---
title: "Steampipe Table: azure_key_vault_key_version - Query Azure Key Vault Keys using SQL"
description: "Allows users to query versions of Azure Key Vault Keys."
---

# Table: azure_key_vault_key_version - Query Azure Key Vault Keys using SQL

Azure Key Vault is a cloud service for securely storing and accessing secrets. A secret is anything that you want to tightly control access to, such as API keys, passwords, certificates, or cryptographic keys. Key Vault service supports multiple key types and algorithms and enables the use of Hardware Security Modules (HSM) for high value keys.

## Table Usage Guide

The 'azure_key_vault_key_version' table provides insights into versions of keys within Azure Key Vault. As a security analyst, explore key-specific details through this table, including key type, key size, and key state. Utilize it to uncover information about keys, such as their creation date, update date, and the recovery level. The schema presents a range of attributes of the key for your analysis, like the key ID, enabled status, expiration date, and associated tags.


## Examples

### Basic info
Explore the settings of Azure Key Vault keys to understand their status and configuration. This is useful for assessing security measures and ensuring proper key management.

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
  azure_key_vault_key_version;
```

### List disabled key versions
Discover the segments that contain disabled key versions in Azure Key Vault. This is useful for assessing security configurations and maintaining proper access controls.

```sql
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

### List keys versions with no expiration time set
Explore which versions of keys in Azure Key Vault lack a set expiration time. This query is useful for identifying potential security risks, as keys without expiration times can be misused if they fall into the wrong hands.

```sql
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
Explore the distribution of different key versions within your Azure Key Vault. This is useful for assessing the overall version management and understanding if certain keys are being updated more frequently than others.

```sql
select
  key_name,
  count(name) as key_version_count
from
  azure_key_vault_key_version
group by
  key_name;
```