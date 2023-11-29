---
title: "Steampipe Table: azure_kubernetes_cluster - Query Azure Kubernetes Services using SQL"
description: "Allows users to query Azure Kubernetes Services"
---

# Table: azure_kubernetes_cluster - Query Azure Kubernetes Services using SQL

Azure Kubernetes Service (AKS) is a managed container orchestration service provided by Microsoft Azure. AKS simplifies the deployment, scaling, and operations of Kubernetes. It provides an integrated developer experience for building, deploying, and scaling containerized applications.

## Table Usage Guide

The 'azure_kubernetes_cluster' table provides insights into Kubernetes clusters within Azure Kubernetes Service (AKS). As a DevOps engineer, explore cluster-specific details through this table, including version, node count, and associated metadata. Utilize it to uncover information about clusters, such as those with specific configurations, the relationships between clusters, and the verification of cluster settings. The schema presents a range of attributes of the Kubernetes cluster for your analysis, like the cluster ID, creation date, attached network policies, and associated tags.

## Examples

### Basic Info
Explore which Azure Kubernetes clusters are available, by identifying their names, IDs, locations, types, and SKU details. This can help in managing resources and understanding the distribution of clusters across different locations and types.

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
Determine the areas in which clusters with a system-assigned identity are located. This query is useful to understand the distribution and arrangement of these clusters across different regions.

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
Determine the areas in which role-based access control (RBAC) is disabled on clusters. This is useful for identifying potential security vulnerabilities within your Azure Kubernetes clusters.

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
Discover clusters that are running on an outdated version, specifically older than 1.20.5. This is useful for identifying potential security risks and planning necessary updates to maintain optimal performance.

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