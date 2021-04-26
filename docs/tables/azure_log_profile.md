# Table: azure_log_profile

Log profiles are the legacy method for sending the Activity log to Azure storage or event hubs. Use the following procedure to continue working with a log profile or to disable it in preparation for migrating to a diagnostic setting.

## Examples

### Basic info

```sql
select
  name,
  id,
  storage_account_id,
  service_bus_rule_id
from
  azure_log_profile;
```