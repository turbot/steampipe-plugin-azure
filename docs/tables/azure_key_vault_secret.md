---
title: "Steampipe Table: azure_key_vault_secret - Query Azure Key Vault Secrets using SQL"
description: "Allows users to query Azure Key Vault Secrets."
---

# Table: azure_key_vault_secret - Query Azure Key Vault Secrets using SQL

Azure Key Vault is a service in Microsoft Azure that allows you to securely store and tightly control access to tokens, passwords, certificates, API keys, and other secrets. It provides a centralized way to manage application secrets and control their distribution. Azure Key Vault helps you control your applications' secrets by keeping them off the code and allowing secure access to them.

## Table Usage Guide

The 'azure_key_vault_secret' table provides insights into secrets within Azure Key Vault. As a security engineer, explore secret-specific details through this table, including secret versions, enabled status, and associated metadata. Utilize it to uncover information about secrets, such as those with expirations, the recovery level of each secret, and the verification of content types. The schema presents a range of attributes of the Azure Key Vault secret for your analysis, like the secret name, vault name, enabled status, and creation date.

## Examples

### Basic info
Discover the secrets stored in your Azure Key Vault by examining the details such as name, ID, vault name, and status. This can help you manage and track your secrets, ensuring they are enabled and updated as needed.

```sql
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
Discover the segments that contain disabled secrets within your Azure Key Vault, allowing you to assess potential security vulnerabilities or areas requiring further management. This is particularly useful for maintaining the integrity of your system by identifying inactive or unused secrets.

```sql
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
Discover the segments that contain secrets in your Azure Key Vault that do not have an expiration date set. This can help in identifying potential security risks and ensuring that all secrets are managed according to best practices.

```sql
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
Discover the segments that consist of active secrets within your Azure Key Vault that have remained unchanged since their creation. This is beneficial for maintaining good security practices, as it allows you to identify and update stagnant secrets.

```sql
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

### Count the number of secrets by vault
Determine the quantity of secrets stored in each Azure Key Vault. This can help in managing and monitoring the distribution of secrets across your vaults.

```sql
select
  vault_name,
  count(vault_name) as count
from
  azure_key_vault_secret
group by
  vault_name;
```