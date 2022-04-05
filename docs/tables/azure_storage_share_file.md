# Table: azure_storage_share_file

Azure Files is Microsoft's easy-to-use cloud file system. Azure file shares can be mounted in Windows, Linux, and macOS.

## Examples

### Basic info

```sql
select
  name,
  storage_account_name,
  type,
  access_tier,
  share_quota,
  enabled_protocols
from
  azure_storage_share_file;
```

### List file shares with default access tier

```sql
select
  name,
  storage_account_name,
  type,
  access_tier,
  access_tier_change_time,
  share_quota,
  enabled_protocols
from
  azure_storage_share_file
where
  access_tier = 'TransactionOptimized';
```

### Get file share with maximum share quota

```sql
select
  name,
  storage_account_name,
  type,
  access_tier,
  access_tier_change_time,
  share_quota,
  enabled_protocols
from
  azure_storage_share_file
order by share_quota desc limit 1;
```
