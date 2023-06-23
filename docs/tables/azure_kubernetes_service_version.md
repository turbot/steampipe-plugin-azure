# Table: azure_kubernetes_service_version

Azure AKS (Azure Kubernetes Service) orchestrator is a managed container orchestration service provided by Microsoft Azure. It simplifies the deployment, management, and scaling of containerized applications using Kubernetes. AKS allows you to deploy and manage containerized applications without the need to manage the underlying infrastructure. It provides automated Kubernetes upgrades, built-in monitoring and diagnostics, and seamless integration with other Azure services. AKS enables developers and DevOps teams to focus on application development and deployment, while Azure takes care of the underlying Kubernetes infrastructure.

**Note:** You must need to pass the `region` in where clause to query this table.

## Examples

### Basic info

```sql
select
  name,
  id,
  type,
  orchestrator_type,
  orchestrator_version
from
  azure_kubernetes_service_version
where
  region = 'eastus2';
```

### List major kubernetes versions

```sql
select
  name,
  id,
  orchestrator_type,
  orchestrator_version
from
  azure_kubernetes_service_version
where
  orchestrator_version = 'major'
and
  region = 'eastus2';
```

### List kubernetes orchestrator type

```sql
select
  name,
  id,
  type,
  orchestrator_type,
  is_preview
from
  azure_kubernetes_service_version
where
  orchestrator_type = 'Kubernetes'
and
  region = 'eastus2';
```

### List kubernetes versions that are not in preview

```sql
select
  name,
  id,
  orchestrator_type,
  orchestrator_version,
  is_preview
from
  azure_kubernetes_service_version
where
  not is_preview
and
  region = 'eastus2';
```

### Get upgrade details of each kubernetes version

```sql
select
  name,
  u ->> 'orchestratorType' as orchestrator_type,
  u ->> 'orchestratorVersion' as orchestrator_version,
  u ->> 'isPreview' as is_preview
from
  azure_kubernetes_service_version,
  jsonb_array_elements(upgrades) as u
where
  region = 'eastus2';
```
