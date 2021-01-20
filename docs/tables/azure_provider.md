# Table: azure_provider

A Azure Resource Provider (RP, for short) is simply an HTTPS RESTful API contract that Add-on owners will implement so a trusted Azure endpoint can provision, delete, and manage services on a user's behalf.

## Examples

### Basic info

```sql
select
	id,
	namespace,
	registration_state
from
	azure_provider;
```


### List of azure providers which are not registered for use

```sql
select
	namespace,
	registration_state
from
	azure_provider
where
	registration_state = 'NotRegistered';
```
