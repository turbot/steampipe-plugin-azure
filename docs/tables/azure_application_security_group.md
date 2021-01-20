# Table: azure_application_security_groups

Application security groups enable you to configure network security as a natural extension of an application's structure, allowing you to group virtual machines and define network security policies based on those groups.

## Examples

### Basic info

```sql
select
  name,
  location,
  resource_group
from
  azure_application_security_group;
```


### List of application security group without application tag key

```sql
select
  name,
  tags
from
  azure_application_security_group
where
  not tags :: JSONB ? 'application';
```
