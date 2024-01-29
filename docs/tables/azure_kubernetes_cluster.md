---
title: "Steampipe Table: azure_kubernetes_cluster - Query Azure Kubernetes Services using SQL"
description: "Allows users to query Azure Kubernetes Services, specifically providing information about the configuration, health and performance of each Kubernetes cluster deployed in Azure."
---

# Table: azure_kubernetes_cluster - Query Azure Kubernetes Services using SQL

Azure Kubernetes Service (AKS) is a managed container orchestration service provided by Microsoft Azure. AKS simplifies the deployment, scaling, and operations of Kubernetes, an open-source system for automating the deployment, scaling, and management of containerized applications. It provides developers with a scalable and highly available infrastructure that's ideal for deploying microservice apps.

## Table Usage Guide

The `azure_kubernetes_cluster` table provides insights into each Kubernetes cluster within Azure Kubernetes Service (AKS). As a DevOps engineer, you can use this table to explore details about each cluster, including its configuration, health status, and performance metrics. This information can be useful for monitoring the state of your clusters, troubleshooting issues, and optimizing resource usage.

## Examples

### Basic Info
Analyze the settings to understand the fundamental details of your Azure Kubernetes clusters. This information can help you monitor and manage your clusters more effectively by providing insights into aspects such as their location, type, and SKU.

```sql+postgres
select
  name,
  id,
  location,
  type,
  sku
from
  azure_kubernetes_cluster;
```

```sql+sqlite
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
Identify instances where your Azure Kubernetes clusters are using a system assigned identity. This is useful in managing and securing cluster resources, as system assigned identities allow Azure to automatically manage the credentials.

```sql+postgres
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

```sql+sqlite
select
  name,
  id,
  location,
  type,
  json_extract(identity, '$.type') as identity_type,
  sku
from
  azure_kubernetes_cluster
where
  json_extract(identity, '$.type') = 'SystemAssigned';
```

### List clusters that have role-based access control (RBAC) disabled
Determine the areas in your Azure Kubernetes clusters where role-based access control (RBAC) is disabled. This can help enhance your security measures by identifying potential vulnerabilities and ensuring appropriate access controls are in place.

```sql+postgres
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

```sql+sqlite
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
Identify instances where your clusters are running on an outdated version (older than 1.20.5) in Azure Kubernetes. This is beneficial for maintaining system security and performance by ensuring your clusters are up-to-date.

```sql+postgres
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

```sql+sqlite
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