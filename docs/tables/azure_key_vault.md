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
  policy #> '{permissions, certificates}'  certificates_permissions,
  policy #> '{permissions, keys}'  keys_permissions,
  policy #> '{permissions, secrets}'  secrets_permissions
from
  azure_key_vault
  cross join jsonb_array_elements(access_policies) as policy;
```
