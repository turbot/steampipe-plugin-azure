---
title: "Steampipe Table: azure_synapse_workspace - Query Azure Synapse Workspaces using SQL"
description: "Allows users to query Azure Synapse Workspaces, providing insights into the analytics service that brings together enterprise data warehousing and Big Data analytics."
folder: "Synapse Analytics"
---

# Table: azure_synapse_workspace - Query Azure Synapse Workspaces using SQL

Azure Synapse Workspace is a feature within Microsoft Azure that integrates with big data and data warehouse technology for immediate insights. It offers a unified experience to ingest, prepare, manage, and serve data for immediate business intelligence and machine learning needs. Azure Synapse Workspace is designed to enable collaboration between data professionals and business decision-makers in a secure and compliant manner.

## Table Usage Guide

The `azure_synapse_workspace` table provides insights into Azure Synapse Workspaces within Microsoft Azure. As a data analyst or data scientist, explore workspace-specific details through this table, including managed private endpoints, firewall settings, and associated metadata. Utilize it to uncover information about workspaces, such as those with private endpoint connections, the status of managed private endpoints, and the verification of firewall rules.

## Examples

### Basic info
Explore the status and type of your Synapse workspaces in Azure to understand their current operation and provisioning state. This can help in managing and optimizing your resources effectively.

```sql+postgres
select
  id,
  name,
  type,
  provisioning_state
from
  azure_synapse_workspace;
```

```sql+sqlite
select
  id,
  name,
  type,
  provisioning_state
from
  azure_synapse_workspace;
```

### List synapse workspaces with public network access enabled
Identify instances where Synapse workspaces in Azure have public network access enabled. This can be useful for security audits to ensure that sensitive data is not exposed to the public internet.

```sql+postgres
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

```sql+sqlite
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
Discover the segments that utilize user-assigned identities within Synapse workspaces. This is beneficial for those wanting to understand which workspaces are configured with specific identity types, aiding in security and access management.

```sql+postgres
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

```sql+sqlite
select
  id,
  name,
  json_extract(identity, '$.type') as identity_type
from
  azure_synapse_workspace
where
  instr(json_extract(identity, '$.type'), 'UserAssigned') > 0;
```

### List private endpoint connection details for synapse workspaces
Explore the details of private endpoint connections for Synapse workspaces. This is beneficial for understanding the status, actions required, and provisioning state of these connections, which can aid in managing and troubleshooting your Azure Synapse Workspaces.

```sql+postgres
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

```sql+sqlite
select
  name as workspace_name,
  w.id as workspace_id,
  json_extract(connections.value, '$.id') as connection_id,
  json_extract(connections.value, '$.privateEndpointPropertyId') as connection_private_endpoint_property_id,
  json_extract(connections.value, '$.privateLinkServiceConnectionStateActionsRequired') as connection_actions_required,
  json_extract(connections.value, '$.privateLinkServiceConnectionStateDescription') as connection_description,
  json_extract(connections.value, '$.privateLinkServiceConnectionStateStatus') as connection_status,
  json_extract(connections.value, '$.provisioningState') as connection_provisioning_state
from
  azure_synapse_workspace as w,
  json_each(private_endpoint_connections) as connections;
```

### List encryption details for synapse workspaces
Explore encryption details for Synapse workspaces to understand the status and level of security measures in place. This can be particularly useful for security audits or for ensuring compliance with data protection regulations.

```sql+postgres
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

```sql+sqlite
select
  name as workspace_name,
  id as workspace_id,
  json_extract(encryption, '$.CmkKey.keyVaultUrl') as cmk_key_vault_url,
  json_extract(encryption, '$.CmkKey.name') as cmk_key_name,
  json_extract(encryption, '$.CmkStatus') as cmk_status,
  json_extract(encryption, '$.DoubleEncryptionEnabled') as double_encryption_enabled
from
  azure_synapse_workspace;
```