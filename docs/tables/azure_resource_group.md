# Table: azure_resource_group

A resource group is a container that holds related resources for an Azure solution.

## Examples

### List of resource groups with their locations

```sql
select
	name,
	location
from
	azure_resource_group;
```


### List of resource groups without owner tag key

```sql
select
	name,
	tags
from
	azure_resource_group
where
	not tags :: JSONB ? 'owner';
```