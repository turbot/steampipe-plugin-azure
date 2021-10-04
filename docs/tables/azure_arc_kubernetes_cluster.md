# Table: azure_arc_kubernetes_cluster

Hybrid Kubernetes Service allows you to manage your on-premise kubernetes clusters from azure by onboarding them to Azure Arc. The Hybrid Kubernetes API allows you to create, list, update and delete your Arc enabled kubernetes clusters.

## Examples

### Basic info

```sql
select
  name,
  id,
  connectivity_status,
  provisioning_state,
  region
from
  azure_arc_kubernetes_cluster;
```

### List expired clusters

```sql
select
  name,
  id,
  type,
  provisioning_state,
  connectivity_status,
  region
from
  azure_arc_kubernetes_cluster
where
  connectivity_status = 'Expired';
```
