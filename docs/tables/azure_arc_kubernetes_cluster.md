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

### List failed batch accounts

```sql
select
  name,
  id,
  type,
  provisioning_state,
  dedicated_core_quota,
  region
from
  azure_batch_account
where
  provisioning_state = 'Failed';
```
