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

### List disabled keys

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

### List keys with no expiration time set

```sql
select
  name,
  enabled,
  expires_at
from
  azure_key_vault_key
where
  expires_at is null;
```

### List keys which have never been updated

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

### Count the number of keys by key vault

```sql
select
  vault_name,
  count(vault_name) as count
from
  azure_key_vault_key
group by
  vault_name;
```
