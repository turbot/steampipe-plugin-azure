---
title: "Steampipe Table: azure_key_vault - Query Azure Key Vault using SQL"
description: "Allows users to query Azure Key Vaults"
---

# Table: azure_key_vault - Query Azure Key Vault using SQL

Azure Key Vault is a cloud service for securely storing and accessing secrets. A secret is anything that you want to tightly control access to, such as API keys, passwords, certificates, or cryptographic keys. Azure Key Vault handles the storage and management of these secrets in a secure and scalable manner, reducing the chances of accidental secret leakage.

## Table Usage Guide

The 'azure_key_vault' table provides insights into Key Vaults within Azure Key Vault service. As a security engineer, explore details specific to each Key Vault through this table, including the vault's URI, resource group, subscription, and location. Utilize it to uncover information about Key Vaults' properties, such as enabled for deployment, disk encryption, template deployment, and soft delete. The schema presents a range of attributes of the Key Vault for your analysis, like the tenant ID, SKU name, family, vault URI, access policies, and associated tags.

## Examples

### List of key vaults where soft deletion is not enabled
Determine the areas in which soft deletion is not enabled within key vaults. This query can be useful for identifying potential security risks and ensuring data recovery options are in place.

```sql
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


### List of key vaults where soft deletion retention period is less than 30 days
Determine the areas in which your Azure Key Vaults have a soft deletion retention period of less than 30 days. This is useful to ensure that your data retention policies are in line with your organization's security standards.

```sql
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
Explore which Azure Key Vaults are enabled for deployment, disk encryption, and template deployment. This is useful for assessing your security configurations and identifying potential vulnerabilities.

```sql
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
Explore which key vaults fall under the premium category. This can be beneficial for understanding your usage and cost distribution in Azure.

```sql
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
Explore the access policies for certificates, keys, and secrets within your Azure Key Vaults. This can help you understand the permissions set up in your environment, ensuring the right access controls are in place.

```sql
select
  name,
  policy -> 'permissionsCertificates' as certificates_permissions,
  policy -> 'permissionsKeys' as keys_permissions,
  policy -> 'permissionsSecrets' as  secrets_permissions
from
  azure_key_vault,
  jsonb_array_elements(access_policies) as policy;
```


### List vaults with logging enabled
Discover the segments of your Azure Key Vaults where logging is enabled. This can be useful for auditing and compliance purposes, as it allows you to track and retain important security and access data.

```sql
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