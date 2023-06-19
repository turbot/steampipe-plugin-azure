# Table: azure_storage_table_service

The properties of a storage account’s Table service endpoint, including properties for Storage Analytics and CORS (Cross-Origin Resource Sharing) rules.

## Examples

### Basic info

```sql
select
  name,
  storage_account_name,
  region,
  resource_group
from
  azure_storage_table_service;
```

### CORS rules info of each storage table service

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
  azure_storage_table_service,
  jsonb_array_elements(cors_rules) as cors;
```