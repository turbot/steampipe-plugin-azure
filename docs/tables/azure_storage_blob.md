# Table: azure_storage_blob

Azure Blob Storage helps you create data lakes for your analytics needs and provides storage to build powerful cloud-native and mobile apps. Optimise costs with tiered storage for your long-term data and flexibly scale up for high-performance computing and machine learning workloads.

## Examples

### Basic info

```sql
select
  name,
  container_name,
  storage_account_name,
  region,
  type,
  is_snapshot
from
  azure_storage_blob
where
  resource_group = 'turbot'
  and storage_account_name = 'mystorageaccount'
  and region = 'eastus';
```

### List snapshot type blobs with import data

```sql
select
  name,
  type,
  access_tier,
  server_encrypted,
  metadata,
  creation_time,
  container_name,
  storage_account_name,
  resource_group,
  region
from
  azure_storage_blob
where
  resource_group = 'turbot'
  and storage_account_name = 'mystorageaccount'
  and region = 'eastus'
  and is_snapshot;
```
