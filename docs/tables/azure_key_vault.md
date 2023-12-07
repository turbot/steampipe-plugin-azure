---
title: "Steampipe Table: azure_key_vault - Query Azure Key Vaults using SQL"
description: "Allows users to query Azure Key Vaults, specifically the keys, secrets, and certificates stored within them, providing insights into security and access management."
---

# Table: azure_key_vault - Query Azure Key Vaults using SQL

Azure Key Vault is a service within Microsoft Azure that safeguards cryptographic keys and other secrets used by cloud apps and services. It allows you to securely store and tightly control access to tokens, passwords, certificates, API keys, and other secrets. Azure Key Vault simplifies the process of meeting the compliance needs for managing, protecting, and controlling access to sensitive business information.

## Table Usage Guide

The `azure_key_vault` table provides insights into the keys, secrets, and certificates stored within Azure Key Vaults. As a security administrator, explore vault-specific details through this table, including enabled status, recovery level, and associated metadata. Utilize it to uncover information about vaults, such as the access policies, tenant details, and the verification of enabled network rules.

## Examples

### List of key vaults where soft deletion is not enabled
Identify instances where key vaults in Azure are potentially vulnerable due to the lack of soft deletion feature. This can help in enhancing data security by pinpointing areas where improvements can be made.

```sql+postgres
select
  name,
  id,
  soft_delete_enabled,
  soft_delete_retention_in_days
from
  azure_key_vault
where
  not soft_delete_enabled;
```

```sql+sqlite
select
  name,
  id,
  soft_delete_enabled,
  soft_delete_retention_in_days
from
  azure_key_vault
where
  soft_delete_enabled = 0;
```


### List of key vaults where soft deletion retention period is less than 30 days
Determine the areas in which the soft deletion retention period of key vaults in Azure is less than 30 days. This query can be used to pinpoint specific locations where data retention policies may need to be strengthened for better security.

```sql+postgres
select
  name,
  id,
  soft_delete_enabled,
  soft_delete_retention_in_days
from
  azure_key_vault
where
  soft_delete_retention_in_days < 30;
```

```sql+sqlite
select
  name,
  id,
  soft_delete_enabled,
  soft_delete_retention_in_days
from
  azure_key_vault
where
  soft_delete_retention_in_days < 30;
```


### Key vaults access information
Determine the areas in which your Azure Key Vaults are being utilized by assessing whether they are enabled for deployment, disk encryption, or template deployment. This allows for a comprehensive understanding of your vault usage and can help optimize resource allocation.

```sql+postgres
select
  name,
  id,
  enabled_for_deployment,
  enabled_for_disk_encryption,
  enabled_for_template_deployment
from
  azure_key_vault;
```

```sql+sqlite
select
  name,
  id,
  enabled_for_deployment,
  enabled_for_disk_encryption,
  enabled_for_template_deployment
from
  azure_key_vault;
```


### List of premium category key vaults
Determine the areas in which premium category key vaults are being used within your Azure environment. This is useful for keeping track of high-security vaults and ensuring they are being used appropriately.

```sql+postgres
select
  name,
  id,
  sku_name,
  sku_family
from
  azure_key_vault
where
  sku_name = 'Premium';
```

```sql+sqlite
select
  name,
  id,
  sku_name,
  sku_family
from
  azure_key_vault
where
  sku_name = 'Premium';
```


### Key vaults access policies details for certificates, keys and secrets
Determine the access policies for certificates, keys, and secrets within Azure Key Vaults to enhance security and access management. This query is useful in understanding the permissions structure within your Key Vaults, which can aid in identifying potential security vulnerabilities.

```sql+postgres
select
  name,
  policy -> 'permissionsCertificates' as certificates_permissions,
  policy -> 'permissionsKeys' as keys_permissions,
  policy -> 'permissionsSecrets' as  secrets_permissions
from
  azure_key_vault,
  jsonb_array_elements(access_policies) as policy;
```

```sql+sqlite
select
  name,
  json_extract(policy.value, '$.permissionsCertificates') as certificates_permissions,
  json_extract(policy.value, '$.permissionsKeys') as keys_permissions,
  json_extract(policy.value, '$.permissionsSecrets') as secrets_permissions
from
  azure_key_vault,
  json_each(access_policies) as policy;
```


### List vaults with logging enabled
Determine the areas in which your Azure Key Vaults have logging enabled for auditing purposes. This can be useful to ensure compliance with security policies and regulations by identifying vaults that are actively recording and retaining audit events.

```sql+postgres
select
  name,
  setting -> 'properties' ->> 'storageAccountId' storage_account_id,
  log ->> 'category' category,
  log -> 'retentionPolicy' ->> 'days' log_retention_days
from
  azure_key_vault,
  jsonb_array_elements(diagnostic_settings) setting,
  jsonb_array_elements(setting -> 'properties' -> 'logs') log
where
  diagnostic_settings is not null
  and setting -> 'properties' ->> 'storageAccountId' <> ''
  and (log ->> 'enabled')::boolean
  and log ->> 'category' = 'AuditEvent'
  and (log -> 'retentionPolicy' ->> 'days')::integer > 0;
```

```sql+sqlite
select
  name,
  json_extract(setting.value, '$.properties.storageAccountId') storage_account_id,
  json_extract(log.value, '$.category') category,
  json_extract(log.value, '$.retentionPolicy.days') log_retention_days
from
  azure_key_vault,
  json_each(diagnostic_settings) as setting,
  json_each(json_extract(setting.value, '$.properties.logs')) as log
where
  diagnostic_settings is not null
  and json_extract(setting.value, '$.properties.storageAccountId') <> ''
  and json_extract(log.value, '$.enabled') = 1
  and json_extract(log.value, '$.category') = 'AuditEvent'
  and json_extract(log.value, '$.retentionPolicy.days') > 0;
```