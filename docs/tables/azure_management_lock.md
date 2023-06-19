# Table: azure_management_lock

Management locks help you prevent accidental deletion or modification of your Azure resources.

## Examples

### List of resources where the management locks are applied

```sql
select
  name,
  split_part(id, '/', 8) as resource_type,
  split_part(id, '/', 9) as resource_name
from
  azure_management_lock;
```


### Resources and lock levels

```sql
select
  name,
  split_part(id, '/', 8) as resource_type,
  split_part(id, '/', 9) as resource_name,
  lock_level
from
  azure_management_lock;
```
