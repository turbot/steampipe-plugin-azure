# Table: azure_deleted_key_vault

Azure Key Vault's soft-delete feature allows recovery of the deleted vaults and deleted key vault objects.

## Examples

### Basic info

```sql
select
  name,
  id,
  type,
  deletion_date,
  scheduled_purge_date
from
  azure_deleted_key_vault;
```

### List vaults with scheduled purge date more than 1 day

```sql
select
  name,
  id,
  type,
  deletion_date,
  scheduled_purge_date
from
  azure_deleted_key_vault
where
scheduled_purge_date > (current_date - interval '1' day);
```
