# Table: azure_compute_disk_encryption_set

Disk Encryption Set simplifies the key management for managed disks. When a disk encryption set is created, a system-assigned managed identity is created in Azure Active Directory (AD) and associated with the disk encryption set.

## Examples

### Key vault associated with each disk encryption set

```sql
select
  name,
  split_part(active_key_source_vault_id, '/', 9) as vault_name,
  split_part(active_key_url, '/', 5) as key_name
from
  azure_compute_disk_encryption_set;
```


### List of encryption sets which are not using customer managed key

```sql
select
  name,
  encryption_type
from
  azure_compute_disk_encryption_set
where
  (
    encryption_type <> 'EncryptionAtRestWithPlatformAndCustomerKeys'
    and encryption_type <> 'EncryptionAtRestWithCustomerKey'
  );
```


### Identity info of each disk encryption set

```sql
select
  name,
  identity_type,
  identity_principal_id,
  identity_tenant_id
from
  azure_compute_disk_encryption_set;
```