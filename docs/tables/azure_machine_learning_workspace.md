---
title: "Steampipe Table: azure_machine_learning_workspace - Query Azure Machine Learning Workspaces using SQL"
description: "Allows users to query Azure Machine Learning Workspaces."
---

# Table: azure_machine_learning_workspace - Query Azure Machine Learning Workspaces using SQL

Azure Machine Learning is a cloud-based service for creating and managing machine learning solutions. It's designed to help data scientists and developers to prepare data, develop experiments, and deploy models at cloud scale. The service supports a wide range of open-source machine learning frameworks like TensorFlow, PyTorch, and scikit-learn.

## Table Usage Guide

The 'azure_machine_learning_workspace' table provides insights into Machine Learning Workspaces within Azure Machine Learning. As a data scientist or developer, explore workspace-specific details through this table, including SKUs, identities, and associated metadata. Utilize it to uncover information about workspaces, such as their provisioning states, their associated application insights, and their linked storage accounts. The schema presents a range of attributes of the Machine Learning Workspace for your analysis, like the workspace name, creation time, and associated tags.

## Examples

### Basic info
Explore the status and types of your Azure Machine Learning workspaces to better understand your resource allocation and management. This can help you identify areas for optimization or reallocation to improve your machine learning workflows.

```sql
select
  name,
  id,
  type,
  provisioning_state
from
  azure_machine_learning_workspace;
```

### List system assigned identity type workspace
Gain insights into Azure Machine Learning Workspaces that are using system-assigned identities. This is beneficial for managing and auditing security and access controls within your Azure environment.

```sql
select
  name,
  id,
  type,
  identity ->> 'type' as identity_type
from
  azure_machine_learning_workspace
where
  identity ->> 'type' = 'SystemAssigned';
```

### List key vault used by workspaces with soft deletion disabled
Explore which workspaces are using key vaults that have soft deletion disabled. This can help identify potential areas of risk and ensure data protection measures are in place.

```sql
select
  m.name as workspace_name,
  m.id as workspace_id,
  v.soft_delete_enabled
from
  azure_machine_learning_workspace as m,
  azure_key_vault as v
where
  lower(m.key_vault) = lower(v.id) and not v.soft_delete_enabled;
```