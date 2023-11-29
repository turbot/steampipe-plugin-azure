---
title: "Steampipe Table: azure_hybrid_kubernetes_connected_cluster - Query Azure Kubernetes Service Connected Clusters using SQL"
description: "Allows users to query Azure Kubernetes Service Connected Clusters"
---

# Table: azure_hybrid_kubernetes_connected_cluster - Query Azure Kubernetes Service Connected Clusters using SQL

Azure Kubernetes Service (AKS) is a managed container orchestration service provided by Microsoft Azure. AKS simplifies the deployment, scaling, and operations of Kubernetes. The Connected Cluster feature allows users to bring their existing Kubernetes clusters running outside of Azure into the Azure Resource Model.

## Table Usage Guide

The 'azure_hybrid_kubernetes_connected_cluster' table provides insights into Connected Clusters within Azure Kubernetes Service (AKS). As a DevOps engineer, explore cluster-specific details through this table, including cluster versions, node counts, and associated metadata. Utilize it to uncover information about clusters, such as their provisioning states, the Kubernetes versions they are running, and their network profiles. The schema presents a range of attributes of the Connected Cluster for your analysis, like the cluster ID, creation date, provisioning state, and associated tags.

## Examples

### Basic info
Explore which Azure Hybrid Kubernetes clusters are provisioned and their respective connectivity statuses to understand their operational readiness across different regions. This is particularly useful in managing resources and ensuring optimal cluster performance.

```sql
select
  name,
  id,
  connectivity_status,
  provisioning_state,
  region
from
  azure_hybrid_kubernetes_connected_cluster;
```

### List expired clusters
Explore which hybrid Kubernetes clusters in your Azure environment have expired. This is useful in maintaining optimal resource allocation and ensuring all active clusters are in good health.

```sql
select
  name,
  id,
  type,
  provisioning_state,
  connectivity_status,
  region
from
  azure_hybrid_kubernetes_connected_cluster
where
  connectivity_status = 'Expired';
```