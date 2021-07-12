# Table: azure_data_factory_pipeline

A Data Factory pipeline is a logical grouping of activities that together perform a task. The activities in a pipeline define actions to perform on data.

## Examples

### Basic info

```sql
select
  name,
  id,
  factory_name,
  type,
  etag
from
  azure_data_factory_pipeline;
```
