# Table: azure_kusto_cluster

Azure Data Explorer aka ADX, is a fast, highly scalable and fully managed data analytics service for log, telemetry and streaming data. ... Previously known as 'codenamed Kusto', this tool uses SQL-like query language, Kusto query language (KQL) for analyzing fast-flowing data from IoT devices, applications, websites, etc.

## Examples

### Basic Info

```sql
select
  name,
  id,
  location,
  type,
  sku_name,
  uri
from
  azure_kusto_cluster;
```

### List kusto clusters with standard sku tier

```sql
select
  name,
  id,
  type,
  sku_name,
  sku_tier
from
  azure_kusto_cluster
where
  sku_tier = 'Standard';
```

### List running kusto clusters

```sql
select
  name,
  id,
  state
from
  azure_kusto_cluster
where
  state = 'Running';
```