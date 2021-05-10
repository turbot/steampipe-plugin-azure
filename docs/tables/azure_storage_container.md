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

### List containers which are publicly accessible

```sql
select
  name,
  id,
  type,
  account_name,
  public_access
from
  azure_storage_container
where
  public_access <> 'None';
```

### List containers with legal hold enabled

```sql
select
  name,
  id,
  type,
  account_name,
  has_legal_hold
from
  azure_storage_container
where
  has_legal_hold;
```

### List containers which are either leased or have a broken lease state

```sql
select
  name,
  id,
  type,
  account_name,
  lease_state
from
  azure_storage_container
where
  lease_state = 'Leased'
  or lease_state = 'Broken';
```

### List containers with infinite lease duration

```sql
select
  name,
  id,
  type,
  account_name,
  lease_duration
from
  azure_storage_container
where
  lease_duration = 'Infinite';
```

### List containers with a remaining retention period of 7 days

```sql
select
  name,
  id,
  type,
  account_name,
  remaining_retention_days
from
  azure_storage_container
where
  remaining_retention_days = 7;
```
