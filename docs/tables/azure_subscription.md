# Table: azure_subscription

An Azure subscription is a logical container used to provision resources in Azure.

## Examples

### Basic info

```sql
select
  id,
  subscription_id,
  display_name,
  tenant_id,
  state,
  authorization_source,
  subscription_policies
from
  azure_subscription;
```
