---
title: "Steampipe Table: azure_synapse_workspace - Query Azure Synapse Analytics Workspaces using SQL"
description: "Allows users to query Azure Synapse Analytics Workspaces."
---

# Table: azure_synapse_workspace - Query Azure Synapse Analytics Workspaces using SQL

Azure Synapse Analytics is an integrated analytics service that accelerates time to insight across data warehouses and big data systems. It blends big data and data warehousing into an on-demand resource that brings together enterprise data warehousing and Big Data analytics. It gives you the freedom to query data on your terms, using either serverless or provisioned resources, at scale.

## Table Usage Guide

The 'azure_synapse_workspace' table provides insights into workspaces within Azure Synapse Analytics. As a data engineer or data scientist, explore workspace-specific details through this table, including managed private endpoints, managed private endpoint connections, and firewall settings. Utilize it to uncover information about workspaces, such as those with specific managed private endpoint settings, the firewall settings between workspaces, and the verification of managed private endpoint connections. The schema presents a range of attributes of the workspace for your analysis, like the workspace name, ID, type, and associated tags.

## Examples

### Basic info
Analyze the settings to understand the status and type of your Azure Synapse workspaces. This can be useful to manage and monitor your workspaces efficiently.

```sql
select
  id,
  name,
  type,
  provisioning_state
from
  azure_synapse_workspace;
```

### List synapse workspaces with public network access enabled
Discover the segments that have public network access enabled within Azure Synapse workspaces. This allows for a quick assessment of potential security risks and helps in maintaining secure configurations.

```sql
select
  id,
  name,
  type,
  provisioning_state,
  public_network_access
from
  azure_synapse_workspace
where
  public_network_access = 'Enabled';
```

### List synapse workspaces with user assigned identities
Determine the areas in which user-assigned identities are utilized within Azure Synapse workspaces. This is useful for managing access control and ensuring appropriate permissions are in place.

```sql
select
  id,
  name,
  identity -> 'type' as identity_type
from
  azure_synapse_workspace
where
    exists (
      select
      from
        unnest(regexp_split_to_array(identity ->> 'type', ',')) elem
      where
        trim(elem) = 'UserAssigned'
  );
```

### List private endpoint connection details for synapse workspaces
Explore the private endpoint connections of Synapse workspaces to understand the current state and any actions required. This is useful in managing and maintaining secure network connections in your data analytics environment.

```sql
select
  name as workspace_name,
  id as workspace_id,
  connections ->> 'id' as connection_id,
  connections ->> 'privateEndpointPropertyId' as connection_private_endpoint_property_id,
  connections ->> 'privateLinkServiceConnectionStateActionsRequired' as connection_actions_required,
  connections ->> 'privateLinkServiceConnectionStateDescription' as connection_description,
  connections ->> 'privateLinkServiceConnectionStateStatus' as connection_status,
  connections ->> 'provisioningState' as connection_provisioning_state
from
  azure_synapse_workspace,
  jsonb_array_elements(private_endpoint_connections) as connections;
```

### List encryption details for synapse workspaces
Explore the encryption details for Synapse workspaces to gain insights into the security measures in place, including the status of the customer-managed key (CMK) and whether double encryption is enabled. This can help assess the security posture and compliance of your data workspaces.

```sql
select
  name as workspace_name,
  id as workspace_id,
  encryption -> 'CmkKey' ->> 'keyVaultUrl' as cmk_key_vault_url,
  encryption -> 'CmkKey' ->> 'name' as cmk_key_name,
  encryption ->> 'CmkStatus' as cmk_status,
  encryption -> 'DoubleEncryptionEnabled' as double_encryption_enabled
from
  azure_synapse_workspace;
```