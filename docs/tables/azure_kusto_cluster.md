# Table: azure_kusto_cluster

An Azure Data Explorer cluster (Previously known as Kusto) is a pair of engine and data management clusters which uses several Azure resources such as Azure Linux VM’s and Storage. The applicable VMs, Azure Storage, Azure Networking and Azure Load balancer costs are billed directly to the customer subscriptions, applications, websites, etc.

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