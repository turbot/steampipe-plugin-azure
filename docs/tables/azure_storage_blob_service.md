# Table: azure_storage_blob_service

The properties of a storage account's Blob service endpoint, including properties for Storage Analytics, CORS (Cross-Origin Resource Sharing) rules and soft delete settings.

## Examples

### Basic info

```sql
select
  name,
  storage_account_name,
  location,
  sku_name,
  sku_tier
from
  azure_storage_blob_service;
```


### List of storage blob service where delete retention policy is not enabled

```sql
select
  name,
  storage_account_name,
  delete_retention_policy -> 'enabled' as delete_retention_policy_enabled
from
  azure_storage_blob_service
where
  delete_retention_policy -> 'enabled' = 'false';
```


### List of storage blob service where versioning is not enabled

```sql
select
  name,
  storage_account_name,
  is_versioning_enabled
from
  azure_storage_blob_service
where
  not is_versioning_enabled;
```


### CORS rules info for storage blob service

```sql
select
  name,
  storage_account_name,
  cors -> 'allowedHeaders' as allowed_headers,
  cors -> 'allowedMethods' as allowed_methods,
  cors -> 'allowedMethods' as allowed_methods,
  cors -> 'exposedHeaders' as exposed_headers,
  cors -> 'maxAgeInSeconds' as max_age_in_seconds
from
  azure_storage_blob_service
  cross join jsonb_array_elements(cors_rules) as cors;
```