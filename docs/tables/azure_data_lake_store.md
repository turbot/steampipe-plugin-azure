# Table: azure_data_lake_store

Azure Data Lake Storage is an enterprise-wide hyper-scale repository for big data analytic workloads. Azure Data Lake enables you to capture data of any size, type, and ingestion speed in one single place for operational and exploratory analytics.

## Examples

### Basic info

```sql
select
  name,
  id,
  type,
  provisioning_state
from
  azure_data_lake_store;
```

### List data lake stores with encryption disabled

```sql
select
  name,
  id,
  type,
  provisioning_state
from
  azure_data_lake_store
where
  encryption_state = 'Disabled';
```

### List data lake stores with firewall disabled

```sql
\
```
