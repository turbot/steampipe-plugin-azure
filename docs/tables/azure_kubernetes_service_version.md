---
title: "Steampipe Table: azure_kubernetes_service_version - Query Azure Kubernetes Services using SQL"
description: "Allows users to query Azure Kubernetes Service Versions."
---

# Table: azure_kubernetes_service_version - Query Azure Kubernetes Services using SQL

Azure Kubernetes Service (AKS) is a managed container orchestration service provided by Microsoft Azure. AKS simplifies the deployment, scaling, and operations of Kubernetes by hosting the Kubernetes environment on Azure. With AKS, you can easily manage and scale your applications using Kubernetes, without the complexities of handling the underlying infrastructure.

## Table Usage Guide

The 'azure_kubernetes_service_version' table provides insights into the versions of Azure Kubernetes Services (AKS). As a DevOps engineer, explore version-specific details through this table, including the release date, Kubernetes version, and whether it's a preview version. Utilize it to uncover information about the availability of different versions, their status, and the upgrade paths. The schema presents a range of attributes of the AKS version for your analysis, like the version name, release date, and whether it's a default version.

## Examples

### Basic info
Explore which versions of Azure Kubernetes Service are available in the 'eastus2' location. This could be useful when planning deployments or upgrades in that specific region.

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
  location = 'eastus2';
```

### List major kubernetes versions
Explore major versions of Kubernetes services deployed in the 'eastus2' region of Azure. This can help you understand the types of Kubernetes orchestrators used and their versions for better management and updates.

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
  location = 'eastus2';
```

### List kubernetes orchestrator type
Determine the areas in which Kubernetes is used as the orchestrator type within the Azure Kubernetes service in the East US 2 region to understand its prevalence and preview status.

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
  location = 'eastus2';
```

### List kubernetes versions that are not in preview
Explore the various versions of Kubernetes that are fully released and available for use in the East US 2 region. This can be useful for planning and implementing your Kubernetes deployments in that specific region.

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
  location = 'eastus2';
```

### Get upgrade details of each kubernetes version
Explore the details of each Kubernetes upgrade, including the orchestrator type and version, and understand whether it is a preview version. This is particularly useful for managing and planning upgrades in the 'eastus2' location.

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
  location = 'eastus2';
```