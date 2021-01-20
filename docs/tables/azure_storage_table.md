# Table: azure_storage_table

Azure Table storage is a service that stores structured NoSQL data in the cloud, providing a key/attribute store with a schemaless design.

## Examples

### Storage table basic info

```sql
select
  name,
  storage_account_name,
  location,
  resource_group
from
  azure_storage_table;
```


### CORS rules info of each storage table

```sql
select
  name,
  storage_account_name,
  cors -> 'allowedHeaders' as allowed_headers,
  cors -> 'allowedMethods' as allowed_methods,
  cors -> 'allowedOrigins' as allowed_origins,
  cors -> 'exposedHeaders' as exposed_eaders,
  cors -> 'maxAgeInSeconds' as max_age_in_seconds
from
  azure_storage_table,
  jsonb_array_elements(cors_rules) as cors;
```