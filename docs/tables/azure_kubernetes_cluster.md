# Table: azure_kubernetes_cluster

Azure Kubernetes orchestrates clusters of virtual machines and schedules containers to run on those virtual machines based on their available compute resources and the resource requirements of each container.

## Examples

### Basic Info

```sql
select
  name,
  id,
  location,
  type,
  sku
from
  azure_kubernetes_cluster;
```


### List clusters with a system assigned identity

```sql
select
  name,
  id,
  location,
  type,
  identity ->> 'type' as identity_type,
  sku
from
  azure_kubernetes_cluster
where
  identity ->> 'type' = 'SystemAssigned';
```


### List clusters that have role-based access control (RBAC) disabled

```sql
select
  name,
  id,
  location,
  type,
  identity,
  enable_rbac,
  sku
from
  azure_kubernetes_cluster
where
  not enable_rbac;
```


### List clusters with an undesirable version (older than 1.20.5)

```sql
select
  name,
  id,
  location,
  type,
  kubernetes_version
from
  azure_kubernetes_cluster
where
  kubernetes_version < '1.20.5';
```
