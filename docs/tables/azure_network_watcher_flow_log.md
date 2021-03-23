# Table: azure_network_watcher_flow_log

Network security group (NSG) flow logs is a feature of Azure Network Watcher that allows user to log information about IP traffic flowing through an NSG. Flow data is sent to Azure Storage accounts from where the user can access it.

## Examples

### List flow logs with their corresponding NSG and Network Watcher details

```sql
select
  name,
  enabled,
  network_watcher_name,
  target_resource_id
from
  azure_network_watcher_flow_log;
```

### List of flow logs which are not enabled

```sql
select
  name,
  id,
  region,
  enabled
from
  azure_network_watcher_flow_log
where
  not enabled;
```

### List of flow logs for which retention period is less than 90 days

```sql
select
  name,
  region,
  enabled,
  flowlog_retention_days
from
  azure_network_watcher_flow_log
where
  enabled and flowlog_retention_days < 90;
```

### Storage account details used to store the flow log

```sql
select
  name,
  flowlog_format_type,
  storage_id
from
  azure_network_watcher_flow_log;
```
