# Table: azure_batch_account

An Azure Batch account is a uniquely identified entity within the Batch service. Most Batch solutions use Azure Storage for storing resource files and output files, so each Batch account is usually associated with a corresponding storage account.

## Examples

### Basic info

```sql
select
  display_name,
  user_principal_name,
  given_name,
  mail,
  account_enabled,
  object_id
from
  azure_batch_account;
```
