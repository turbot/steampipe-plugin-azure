# Table: azure_network_watcher

Network Watcher is a regional service that enables you to monitor and diagnose conditions at a network scenario level.

## Examples

### List of regions where network watcher is enabled

```sql
select
  name,
  region
from
  azure_network_watcher;
```

### List of Network watcher without application tag key

```sql
select
  name,
  tags
from
  azure_network_watcher
where
  not tags :: JSONB ? 'application';
```