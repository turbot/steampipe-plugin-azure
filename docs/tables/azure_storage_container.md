# Table: azure_storage_container

A container organizes a set of blobs, similar to a directory in a file system. A storage account can include an unlimited number of containers, and a container can store an unlimited number of blobs.

## Examples

### Basic info

```sql
select
  name,
  id,
  type,
  account_name
from
  azure_storage_container;
```

### List containers which are not publicly accessible

```sql
select
  name,
  id,
  type,
  account_name
from
  azure_storage_container
where
  public_access = 'None';
```
