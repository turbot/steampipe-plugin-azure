# Table: azure_data_factory_dataset

 Datasets identify data within different data stores, such as tables, files, folders, and documents.

## Examples

### Basic info

```sql
select
  name,
  id,
  description,
  etag,
  type
from
  azure_data_factory_dataset;
```


### List relational table type Datasets

```sql
select
  name,
  id,
  type,
  properties ->> 'type' as dataset_type
from
  azure_data_factory_dataset
where
  properties ->> 'type' = 'RelationalTable';
```