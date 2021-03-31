# Table: azure_key_vault_key

Azure Key Vault Keys are 'Cryptographic keys' used to encrypt information without releasing the private key to the consumer. It acts like a black box to encrypt and decrypt content using the RSA algorithm. The RSA algorithm, involves a public key and private key.

## Examples

### Basic info

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
  azure_key_vault_key;
```

### List keys which are not enabled

```sql
select
  name,
  vault_name,
  enabled
from
  azure_key_vault_key
where
  not enabled;
```

### List keys for which expiration time is not set

```sql
select
  name,
  enabled,
  expired_at
from
  azure_key_vault_key
where
  expired_at is null;
```

### List keys which have never updated

```sql
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

### Count of keys by Key Vault

```sql
select
  vault_name,
  count(vault_name) as count
from
  azure_key_vault_key
group by
  vault_name;
```
