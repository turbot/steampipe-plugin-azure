---
title: "Steampipe Table: azure_kubernetes_service_version - Query Azure Kubernetes Service Versions using SQL"
description: "Allows users to query Azure Kubernetes Service Versions, providing detailed information about the different versions of the Kubernetes service available in Azure."
folder: "AKS"
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
  version,
  cloud_environment,
  is_preview,
  title
from
  azure_kubernetes_service_version
where
  location = 'eastus2';
```

```sql+sqlite
select
  version,
  cloud_environment,
  is_preview,
  title
from
  azure_kubernetes_service_version
where
  location = 'eastus2';
```

### List patch kubernetes versions

Determine the patch versions of Kubernetes orchestration service in the East US 2 region within Azure. This is useful for understanding the available Kubernetes versions in a specific location for planning deployments or upgrades.

```sql+postgres
select
  version,
  is_preview,
  jsonb_pretty(patch_versions) as patch_versions
from
  azure_kubernetes_service_version
where
  location = 'eastus2';
```

```sql+sqlite
select
  version,
  is_preview,
  jsonb(patch_versions) as patch_versions
from
  azure_kubernetes_service_version
where
  location = 'eastus2';
```

### List kubernetes versions that are not in preview

Determine the versions of Kubernetes in the 'eastus2' location that are fully released and not in a preview stage. This could be useful for organizations planning to use stable versions of Kubernetes for their operations in the specified location.

```sql+postgres
select
  version,
  cloud_environment,
  is_preview,
  title
from
  azure_kubernetes_service_version
where
  not is_preview
and
  location = 'eastus2';
```

```sql+sqlite
select
  version,
  cloud_environment,
  is_preview,
  title
from
  azure_kubernetes_service_version
where
  not is_preview
and
  location = 'eastus2';
```

### Get capabilities of each kubernetes version

Determine the capabilities for each version of Kubernetes within a specific location. This can be useful for planning and managing version capabilities, especially in identifying whether the version is still in preview or fully released.

```sql+postgres
select
  version,
  is_preview,
  jsonb_pretty(capabilities) as capabilities
from
  azure_kubernetes_service_version
where
  location = 'eastus2';
```

```sql+sqlite
select
  version,
  is_preview,
  json(capabilities) as capabilities
from
  azure_kubernetes_service_version
where
  location = 'eastus2';
```
