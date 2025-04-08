---
title: "Steampipe Table: azure_hybrid_compute_machine - Query Azure Hybrid Compute Machines using SQL"
description: "Allows users to query Azure Hybrid Compute Machines, providing insights into the configuration and status of hybrid machines in the Azure environment."
folder: "Hybrid"
---

# Table: azure_hybrid_compute_machine - Query Azure Hybrid Compute Machines using SQL

Azure Hybrid Compute Machines are part of the Azure Arc service, which extends Azure services and management to any infrastructure. It enables management of Windows and Linux machines hosted outside of Azure, on the corporate network, or other cloud provider. This feature allows consistent Azure management across environments, providing a single control plane with access to the same familiar cloud-native Azure management experiences.

## Table Usage Guide

The `azure_hybrid_compute_machine` table provides insights into Azure Hybrid Compute Machines within Azure Arc. As a system administrator, explore machine-specific details through this table, including machine properties, operating system details, and status information. Utilize it to uncover information about machines, such as their current provisioning state, the version of the installed agent, and the time of the last agent heartbeat.

## Examples

### Basic info
This query provides a way to gain insights into the status and location of your Azure hybrid compute machines. This can be useful for managing resources and ensuring optimal performance across different regions.

```sql+postgres
select
  name,
  id,
  status,
  provisioning_state,
  region
from
  azure_hybrid_compute_machine;
```

```sql+sqlite
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
Explore which machines in your Azure hybrid computing environment are disconnected. This is useful to identify potential issues in your network and ensure all systems are functioning properly.

```sql+postgres
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

```sql+sqlite
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