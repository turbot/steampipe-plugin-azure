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

### List containers for which the legal_hold is true

```sql
select
  name,
  id,
  type,
  account_name
from
  azure_storage_container
where
  has_legal_hold;
```

### List containers which are either leased or lease broken state

```sql
select
  name,
  id,
  type,
  account_name
from
  azure_storage_container
where
  lease_state = 'Leased'
  or lease_state = 'Broken';
```

### List containers where lease duration is infinite.

```sql
select
  name,
  id,
  type,
  account_name
from
  azure_storage_container
where
  lease_duration = 'Infinite';
```

### List containers where retention days ending in next 7 days

```sql
select
  name,
  id,
  type,
  account_name
from
  azure_storage_container
where
  remaining_retention_days = 7;
```
