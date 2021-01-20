# Table: azure_storage_queue

Azure Queue Storage is a service for storing large numbers of messages which allows to access messages from anywhere in the world via authenticated calls using HTTP or HTTPS.

## Examples

### List of queues and their corresponding storage accounts

```sql
select
	name as queue_name,
	storage_account_name,
	location
from
	azure_storage_queue;
```


### List of storage queues without owner tag key

```sql
select
	name,
	tags
from
	azure_application_security_group
where
	not tags :: JSONB ? 'owner';
```