# Table: azure_management_group

Management groups provide a governance scope above subscriptions. You organize subscriptions into management groups in the governance conditions you apply cascade by inheritance to all associated subscriptions. Management groups give you enterprise-grade management at a scale no matter what type of subscriptions you might have. However, all subscriptions within a single management group must trust the same Azure Active Directory (Azure AD) tenant.

Note: To query this table, you need to have read access to the specific management group.

## Examples

### Basic info

```sql
select
  id,
  name,
  type,
  tenant_id,
  updated_by
from
  azure_management_group;
```

### List children for a specific management group

```sql
select
  name,
  updated_by,
  jsonb_pretty(children) as children
from
  azure_management_group;
```

### Get parent detail for a specific management group

```sql
select
  name,
  updated_by,
  jsonb_pretty(parent) as parent
from
  azure_management_group;
```
