---
title: "Steampipe Table: azure_hybrid_kubernetes_connected_cluster - Query Azure Hybrid Kubernetes Connected Clusters using SQL"
description: "Allows users to query Hybrid Kubernetes Connected Clusters in Azure, specifically detailing the configuration, status, and properties of each connected cluster."
folder: "Hybrid Kubernetes"
---

# Table: azure_hybrid_kubernetes_connected_cluster - Query Azure Hybrid Kubernetes Connected Clusters using SQL

Azure Hybrid Kubernetes Connected Clusters is a feature offered by Microsoft Azure that allows users to manage and govern Kubernetes clusters across on-premises, edge, and multi-cloud environments from a single pane of glass. Its unified approach offers consistent visibility, governance, and control across different environments, making it easier to manage Kubernetes resources. It provides a comprehensive view of all Kubernetes applications, irrespective of where they are running.

## Table Usage Guide

The `azure_hybrid_kubernetes_connected_cluster` table provides insights into Hybrid Kubernetes Connected Clusters within Microsoft Azure. As a DevOps engineer, this table can be utilized to explore cluster-specific details, including configuration, status, and associated properties. Use it to uncover information about clusters, such as their health status, the Kubernetes version they're running, and their connectivity state with Azure.

## Examples

### Basic info
Explore which Azure Hybrid Kubernetes clusters are currently connected and their respective provisioning states. This can help in assessing the overall health and status of your hybrid cloud infrastructure.

```sql+postgres
select
  name,
  id,
  connectivity_status,
  provisioning_state,
  region
from
  azure_hybrid_kubernetes_connected_cluster;
```

```sql+sqlite
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
Explore which hybrid Kubernetes clusters in your Azure environment have expired connectivity. This is useful in maintaining an up-to-date and secure network by promptly addressing any expired clusters.

```sql+postgres
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

```sql+sqlite
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