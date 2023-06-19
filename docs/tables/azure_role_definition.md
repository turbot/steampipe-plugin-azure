# Table: azure_role_definition

A role definition lists the operations that can be performed, such as read, write, and delete.

## Examples

### List the custom roles

```sql
select
  name,
  description,
  role_name,
  role_type,
  title
from
  azure_role_definition
where
  role_type = 'CustomRole';
```

### List of roles whose assignable scope is set to root('/') scope

```sql
select
  name,
  role_name,
  scope
from
  azure_role_definition,
  jsonb_array_elements_text(assignable_scopes) as scope
where
  scope = '/';
```

### Permissions of all custom roles

```sql
select
  name,
  role_name,
  role_type,
  permission -> 'actions' as action,
  permission -> 'dataActions' as data_action,
  permission -> 'notActions' as no_action,
  permission -> 'notDataActions' as not_data_actions
from
  azure_role_definition
  cross join jsonb_array_elements(permissions) as permission
where
  role_type = 'CustomRole';
```
