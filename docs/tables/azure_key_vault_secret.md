# Table: azure_key_vault_secret

Azure Key Vault is a cloud service for securely storing and accessing secrets.
A secret is anything that you want to tightly control access to, such as API keys, passwords, certificates, or cryptographic keys.

## Examples

### Basic info

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

### List secrets which are not enabled

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

### List secrets for which expiration time is not set

```sql
select
  name,
  enabled,
  expired_at
from
  azure_key_vault_secret
where
  expired_at is null;
```

### List secrets which have never updated

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

### Count of secrets by Key Vault

```sql
select
  vault_name,
  count(vault_name) as count
from
  azure_key_vault_secret
group by
  vault_name;
```
