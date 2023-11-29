---
title: "Steampipe Table: azure_hybrid_compute_machine - Query Azure Hybrid Compute Machines using SQL"
description: "Allows users to query Azure Hybrid Compute Machines"
---

# Table: azure_hybrid_compute_machine - Query Azure Hybrid Compute Machines using SQL

Azure Hybrid Compute Machines are a part of the Azure Arc service that extends Azure management and services to any infrastructure. It allows you to manage and govern Windows and Linux machines hosted outside of Azure, on your corporate network, or other cloud provider. This service simplifies complex and distributed environments across on-premises, edge, and multi-cloud into a unified central point.

## Table Usage Guide

The 'azure_hybrid_compute_machine' table provides insights into Hybrid Compute Machines within Azure Arc. As a DevOps engineer, explore machine-specific details through this table, including machine properties, operating system details, and associated metadata. Utilize it to uncover information about machines, such as their status, location, and the version of the Azure Arc agent installed on them. The schema presents a range of attributes of the Hybrid Compute Machine for your analysis, like the machine's ID, name, location, and operating system.

## Examples

### Basic info
Explore which machines in your Azure hybrid environment are active and where they are located. This can assist in managing resources and understanding the distribution of your infrastructure.

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
Identify instances where machines in the Azure hybrid compute environment are disconnected. This can be useful in diagnosing network issues or managing system availability.

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