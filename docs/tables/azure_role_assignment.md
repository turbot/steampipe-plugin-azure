# Table: azure_role_assignment

Azure role assignments is the authorization system to manage access to Azure resources. To grant access, you assign roles to users, groups, service principals, or managed identities at a particular scope.

## Examples

### Role assignment basic info

```sql
select
  name,
  id,
  principal_id,
  principal_type
from
  azure_role_assignment;
```


### List of role assignments which has permission at root level

```sql
select
  name,
  id,
  scope
from
  azure_role_assignment
where
  scope = '/';
```


### List of role assignments which has subscription level permission and full access to the subscription

```sql
select
  ra.name as roll_assignment_name,
  rd.role_name
from
  azure_role_assignment ra
  join azure_role_definition rd on ra.role_definition_id = rd.id
  cross join jsonb_array_elements(rd.permissions) as perm
where
  ra.scope like '/subscriptions/%'
  and perm -> 'actions' = '["*"]';
```
