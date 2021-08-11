# Table: azure_search_service

Azure Cognitive Search is the only cloud search service with built-in AI capabilities that enrich all types of information to help you identify and explore relevant content at scale. Use cognitive skills for vision, language and speech or use custom machine learning models to uncover insights from all types of content.

## Examples

### Basic info

```sql
select
  name,
  id,
  type,
  provisioning_state,
  status,
  sku_name,
  replica_count
from
  azure_search_service;
```

### List virtual machine type search service

```sql
select
  name,
  id,
  type
from
  azure_search_service
where
  type = 'Microsoft.Compute/virtualMachines';
```
