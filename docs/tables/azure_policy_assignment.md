# Table: azure_policy_assignment

Azure policy assignment retrieves the information of all policy assignments associated with the given subscription.

## Examples

### Basic info

```sql
select
  id,
  name,
  type
from
  azure_policy_assignment;
```
