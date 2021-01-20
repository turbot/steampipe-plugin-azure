# Table: azure_storage_blob

A Binary Large OBject (BLOB) is a collection of binary data stored as a single entity in a database management system. Blobs are typically images, audio or other multimedia objects, though sometimes binary executable code is stored as a blob.

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
  azure_storage_blob;
```


### List of storage blobs where delete retention policy is not enabled

```sql
select
  name,
  storage_account_name,
  delete_retention_policy -> 'enabled' as delete_retention_policy_enabled
from
  azure_storage_blob
where
  delete_retention_policy -> 'enabled' = 'false';
```


### List of storage blobs where versioning is not enabled

```sql
select
  name,
  storage_account_name,
  is_versioning_enabled
from
  azure_storage_blob
where
  not is_versioning_enabled;
```


### CORS rules info for storage blob

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
  azure_storage_blob
  cross join jsonb_array_elements(cors_rules) as cors;
```