# Table: azure_hybrid_compute_machine

Azure Arc enables you to manage servers running outside of Azure using Azure Resource Manager. Each server is represented in Azure as a hybrid compute machine resource. Once a server is managed with Azure Arc, you can deploy agents, scripts, or configurations to the machine using extensions. The Hybrid Compute API allows you to create, list, update and delete your Azure Arc enabled servers and any extensions associated with them.

## Examples

### Basic info

```sql
select
  name,
  id,
  status,
  provisioning_state,
  region
from
  azure_hybrid_compute_machine;
```

### List disconnected machines

```sql
select
  name,
  id,
  type,
  provisioning_state,
  status,
  region
from
  azure_hybrid_compute_machine
where
  status = 'Disconnected';
```