# Table: azure_network_watcher_flow_log

Network security group (NSG) flow logs is a feature of Azure Network Watcher that allows user to log information about IP traffic flowing through an NSG. Flow data is sent to Azure Storage accounts from where the user can access it.

## Examples

### Basic info

```sql
select
  name,
  enabled,
  network_watcher_name,
  target_resource_id
from
  azure_network_watcher_flow_log;
```

### List disabled flow logs

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

### List flow logs with a retention period less than 90 days

```sql
select
  name,
  region,
  enabled,
  retention_policy_days
from
  azure_network_watcher_flow_log
where
  enabled and retention_policy_days < 90;
```

### Get storage account details for each flow log

```sql
select
  name,
  file_type,
  storage_id
from
  azure_network_watcher_flow_log;
```
