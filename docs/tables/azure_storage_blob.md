# Table: azure_storage_blob

A Binary Large OBject (BLOB) is a collection of binary data stored as a single entity in a database management system. Blobs are typically images, audio or other multimedia objects, though sometimes binary executable code is stored as a blob.

## Examples

### Basic info

```sql
select
  name,
  storage_account_name,
  container_name,
  region,
  type,
  is_snapshot
from
  azure_storage_blob;
```
