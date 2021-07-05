# Table: azure_data_factory_pipeline

A pipeline is a logical grouping of activities that together perform a task. The activities in a pipeline define actions to perform on data.

## Examples

### Basic info

```sql
select
  name,
  id,
  type,
  etag,
  factory_name
from
  azure_data_factory_pipeline;
```