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


### List Clusters for a specific location

```sql
select
  name,
  id,
  location,
  type,
  identity,
  sku
from
  azure_kubernetes_cluster
where
  location = 'westus';
```


### List Clusters having System Assigned identity

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


### List Clusters having role-based access control (RBAC) disabled

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