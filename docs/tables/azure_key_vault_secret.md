---
title: "Steampipe Table: azure_key_vault_secret - Query Azure Key Vault Secrets using SQL"
description: "Allows users to query Azure Key Vault Secrets, providing insights into the secrets stored in Azure Key Vaults, including their attributes, versions, and associated metadata."
---

# Table: azure_key_vault_secret - Query Azure Key Vault Secrets using SQL

Azure Key Vault Secret is a resource within Microsoft Azure that allows you to securely store and tightly control access to tokens, passwords, certificates, API keys, and other secrets. It provides a centralized way to manage application secrets and control their distribution. Azure Key Vault Secret helps maintain application secrets with a high level of security.

## Table Usage Guide

The `azure_key_vault_secret` table provides insights into the secrets stored in Azure Key Vaults. As a security engineer, explore secret-specific details through this table, including secret attributes, versions, and associated metadata. Utilize it to uncover information about secrets, such as their recovery level, enabled status, and expiration dates.

## Examples

### Basic info
Explore the status and details of your Azure Key Vault secrets. This query is useful to keep track of the secrets' status, enabling you to manage and monitor them effectively.

```sql+postgres
select
  name,
  id,
  vault_name,
  enabled,
  created_at,
  updated_at,
  value
from
  azure_key_vault_secret;
```

```sql+sqlite
select
  name,
  id,
  vault_name,
  enabled,
  created_at,
  updated_at,
  value
from
  azure_key_vault_secret;
```

### List disabled secrets
Explore which secrets within the Azure Key Vault are currently disabled. This can help in managing access and maintaining the security of your vault.

```sql+postgres
select
  name,
  vault_name,
  enabled
from
  azure_key_vault_secret
where
  not enabled;
```

```sql+sqlite
select
  name,
  vault_name,
  enabled
from
  azure_key_vault_secret
where
  not enabled;
```

### List secrets that do not expire
Discover the segments that consist of non-expiring secrets within Azure's key vault. This can be useful in managing and identifying potential security risks associated with indefinite secret keys.

```sql+postgres
select
  name,
  enabled,
  expires_at
from
  azure_key_vault_secret
where
  expires_at is null;
```

```sql+sqlite
select
  name,
  enabled,
  expires_at
from
  azure_key_vault_secret
where
  expires_at is null;
```

### List enabled secrets that have never been updated
Identify the enabled secrets within your Azure Key Vault that have remained unchanged since their creation. This is useful for security purposes and ensuring that secret keys are being regularly updated and managed properly.

```sql+postgres
select
  name,
  enabled,
  created_at,
  updated_at
from
  azure_key_vault_secret
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
  azure_key_vault_secret
where
  enabled
  and (julianday(updated_at) - julianday(created_at)) * 24 * 60 * 60 = 0;
```

### Count the number of secrets by vault
Assess the elements within your Azure Key Vault by counting the number of secrets each vault holds. This allows you to understand the distribution of secrets across your vaults, helping to manage and balance storage.

```sql+postgres
select
  vault_name,
  count(vault_name) as count
from
  azure_key_vault_secret
group by
  vault_name;
```

```sql+sqlite
select
  vault_name,
  count(vault_name) as count
from
  azure_key_vault_secret
group by
  vault_name;
```