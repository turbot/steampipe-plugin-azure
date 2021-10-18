# Table: azure_synapse_workspace

Azure Synapse is an enterprise analytics service that accelerates time to insight across data warehouses and big data systems. Azure Synapse brings together the best of SQL technologies used in enterprise data warehousing, Spark technologies used for big data, Pipelines for data integration and ETL/ELT, and deep integration with other Azure services such as Power BI, CosmosDB, and AzureML.

## Examples

### Basic info

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
