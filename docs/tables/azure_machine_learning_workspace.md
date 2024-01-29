---
title: "Steampipe Table: azure_machine_learning_workspace - Query Azure Machine Learning Workspaces using SQL"
description: "Allows users to query Azure Machine Learning Workspaces, providing comprehensive information on configuration, status, and properties of each workspace."
---

# Table: azure_machine_learning_workspace - Query Azure Machine Learning Workspaces using SQL

Azure Machine Learning is a cloud-based environment that enables developers to build, train, and deploy machine learning models. Workspaces in Azure Machine Learning are the top-level resource for the service, providing a centralized place to work with all the artifacts you create. A workspace is tied to an Azure subscription and the resources are used to run the workspace and its experiments.

## Table Usage Guide

The `azure_machine_learning_workspace` table provides insights into Azure Machine Learning Workspaces. As a data scientist or machine learning engineer, you can explore workspace-specific details through this table, including configurations, status, and properties. Utilize it to uncover information about workspace, such as its associated resources, location, and SKU details, enabling effective management and optimization of your machine learning experiments.

## Examples

### Basic info
Explore which Azure Machine Learning Workspaces are currently provisioned and understand their types. This can be useful for managing resources and understanding the distribution of workspace types within your Azure environment.

```sql+postgres
select
  name,
  id,
  type,
  provisioning_state
from
  azure_machine_learning_workspace;
```

```sql+sqlite
select
  name,
  id,
  type,
  provisioning_state
from
  azure_machine_learning_workspace;
```

### List system assigned identity type workspace
Determine the areas in which system-assigned identities are used within Azure's machine learning workspace. This can help streamline system management by identifying where these automated identities are deployed.

```sql+postgres
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

```sql+sqlite
select
  name,
  id,
  type,
  json_extract(identity, '$.type') as identity_type
from
  azure_machine_learning_workspace
where
  json_extract(identity, '$.type') = 'SystemAssigned';
```

### List key vault used by workspaces with soft deletion disabled
Determine the areas in which the key vault used by workspaces has soft deletion disabled. This is beneficial in identifying potential vulnerabilities and ensuring data protection and recovery strategies are in place.

```sql+postgres
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

```sql+sqlite
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