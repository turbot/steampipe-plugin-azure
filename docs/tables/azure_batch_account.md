# Table: azure_batch_account

An Azure Batch account is a uniquely identified entity within the Batch service. Most Batch solutions use Azure Storage for storing resource files and output files, so each Batch account is usually associated with a corresponding storage account.

## Examples

### Basic info

```sql
select
  name,
  id,
  type,
  provisioning_state,
  dedicated_core_quota,
  region
from
  azure_batch_account;
```

### List failed batch accounts

```sql
select
  name,
  id,
  type,
  provisioning_state,
  dedicated_core_quota,
  region
from
  azure_batch_account
where
provisioning_state = 'Failed';
```
