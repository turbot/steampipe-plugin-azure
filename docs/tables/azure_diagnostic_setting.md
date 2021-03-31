# Table: azure_diagnostic_setting

Each Azure resource requires its own diagnostic setting, which defines the following criteria:

Categories of logs and metric data sent to the destinations defined in the setting. The available categories will vary for different resource types.
One or more destinations to send the logs. Current destinations include Log Analytics workspace, Event Hubs, and Azure Storage.

## Examples

### Basic info

```sql
select
  name,
  id,
  type
from
  azure_diagnostic_setting;
```


### Ensure Diagnostic Setting captures appropriate categories

```sql
select
  jsonb_pretty(diagnostic_settings) as diagnostic_settings
from
  azure_diagnostic_setting;
```