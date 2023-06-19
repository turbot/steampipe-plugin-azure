# Table: azure_storage_table

Azure Table storage is a service that stores structured NoSQL data in the cloud, providing a key/attribute store with a schema less design.

## Examples

### Basic info

```sql
select
  name,
  id,
  storage_account_name,
  resource_group,
  region,
  subscription_id
from
  azure_storage_table;
```
