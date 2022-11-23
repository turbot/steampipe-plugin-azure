# Table: azure_key_vault_key_version

Azure Key Vault Keys are 'Cryptographic keys' used to encrypt information without releasing the private key to the consumer. It acts like a black box to encrypt and decrypt content using the RSA algorithm. The RSA algorithm, involves a public key and private key. They can roll to a new version of the key, back it up, and do related tasks.

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
  azure_key_vault_key_version;
```

### List disabled key versions

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

```sql
select
  key_name,
  count(key_name) as count
from
  azure_key_vault_key_version
group by
  key_name;
```
