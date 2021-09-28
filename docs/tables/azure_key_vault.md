# Table: azure_key_vault

Azure Key Vault is a cloud service used to manage keys, secrets, and certificates.

## Examples

### List of key vaults where soft deletion is not enabled

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
