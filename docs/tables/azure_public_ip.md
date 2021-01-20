# Table: azure_public_ip

Public IP addresses enable Azure resources to communicate to Internet and public-facing Azure services.

## Examples

### List of unassociated elastic IPs

```sql
select
	name,
	ip_configuration_id
from
	azure_public_ip
where
	ip_configuration_id is null;
```


### List of IP addresses with corresponding associations

```sql
select
	name,
	ip_address,
	split_part(ip_configuration_id, '/', 8) as resource,
	split_part(ip_configuration_id, '/', 9) as resource_name
from
	azure_public_ip;
```


### List of dynamic IP addresses

```sql
select
	name,
	public_ip_allocation_method
from
	azure_public_ip
where
	public_ip_allocation_method = 'Dynamic';
```
