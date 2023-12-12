---
title: "Steampipe Table: azure_kubernetes_service_version - Query Azure Kubernetes Service Versions using SQL"
description: "Allows users to query Azure Kubernetes Service Versions, providing detailed information about the different versions of the Kubernetes service available in Azure."
---

# Table: azure_kubernetes_service_version - Query Azure Kubernetes Service Versions using SQL

Azure Kubernetes Service (AKS) is a managed container orchestration service provided by Azure. It simplifies the deployment, scaling, and operations of containerized applications using Kubernetes, an open-source platform for automating deployment, scaling, and management of containerized applications. The service versions table provides information about the different versions of the Kubernetes service available in Azure.

## Table Usage Guide

The `azure_kubernetes_service_version` table provides insights into the different versions of Azure Kubernetes Service available. As a DevOps engineer or system administrator, you can use this table to understand the features, improvements, and fixes associated with each version of the service. This can help in making informed decisions when planning for version upgrades or when troubleshooting issues related to specific versions.

**Important notes:**
- You must specify the `location` in the `where` clause to query this table.

## Examples

### Basic info
Discover the segments of Azure's Kubernetes service located in the 'eastus2' region to understand their orchestration types and versions. This can be useful to identify and manage services based on their orchestration details.

```sql+postgres
select
  name,
  id,
  type,
  orchestrator_type,
  orchestrator_version
from
  azure_kubernetes_service_version
where
  location = 'eastus2';
```

```sql+sqlite
select
  name,
  id,
  type,
  orchestrator_type,
  orchestrator_version
from
  azure_kubernetes_service_version
where
  location = 'eastus2';
```

### List major kubernetes versions
Determine the major versions of Kubernetes orchestration service in the East US 2 region within Azure. This is useful for understanding the available Kubernetes versions in a specific location for planning deployments or upgrades.

```sql+postgres
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
  location = 'eastus2';
```

```sql+sqlite
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
  location = 'eastus2';
```

### List kubernetes orchestrator type
Determine the areas in which Kubernetes is used as the orchestrator type within the Azure Kubernetes Service in the East US 2 region. This can be useful for organizations to assess their use of Kubernetes in specific geographical locations.

```sql+postgres
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
  location = 'eastus2';
```

```sql+sqlite
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
  location = 'eastus2';
```

### List kubernetes versions that are not in preview
Determine the versions of Kubernetes in the 'eastus2' location that are fully released and not in a preview stage. This could be useful for organizations planning to use stable versions of Kubernetes for their operations in the specified location.

```sql+postgres
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
  location = 'eastus2';
```

```sql+sqlite
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
  location = 'eastus2';
```

### Get upgrade details of each kubernetes version
Determine the upgrade details for each version of Kubernetes within a specific location. This can be useful for planning and managing version upgrades, especially in identifying whether the version is still in preview or fully released.

```sql+postgres
select
  name,
  u ->> 'orchestratorType' as orchestrator_type,
  u ->> 'orchestratorVersion' as orchestrator_version,
  u ->> 'isPreview' as is_preview
from
  azure_kubernetes_service_version,
  jsonb_array_elements(upgrades) as u
where
  location = 'eastus2';
```

```sql+sqlite
select
  name,
  json_extract(u.value, '$.orchestratorType') as orchestrator_type,
  json_extract(u.value, '$.orchestratorVersion') as orchestrator_version,
  json_extract(u.value, '$.isPreview') as is_preview
from
  azure_kubernetes_service_version,
  json_each(upgrades) as u
where
  location = 'eastus2';
```