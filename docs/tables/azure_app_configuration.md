# Table: azure_app_configuration

Azure App Configuration provides a service to centrally manage application settings and feature flags. App Configuration is used to store all the settings for your application and secure their accesses in one place.

## Examples

### Basic info

```sql
select
  id,
  name,
  type,
  provisioning_state,
  creation_date
from
  azure_app_configuration;
```

### List public network access enabled app configurations

```sql
select
  id,
  name,
  type,
  provisioning_state,
  public_network_access
from
  azure_app_configuration
where
  public_network_access = 'Enabled';
```

### List app configurations with user assigned identities

```sql
select
  id,
  name,
  identity -> 'type' as identity_type,
  jsonb_pretty(identity -> 'userAssignedIdentities') as identity_user_assigned_identities
from
  azure_app_configuration
where
    exists (
      select
      from
        unnest(regexp_split_to_array(identity ->> 'type', ',')) elem
      where
        trim(elem) = 'UserAssigned'
  );
```

### List private endpoint connection details for app configurations

```sql
select
  name as app_config_name,
  id as app_config_id,
  connections ->> 'id' as connection_id,
  connections ->> 'privateEndpointPropertyId' as connection_private_endpoint_property_id,
  connections ->> 'privateLinkServiceConnectionStateActionsRequired' as connection_actions_required,
  connections ->> 'privateLinkServiceConnectionStateDescription' as connection_description,
  connections ->> 'privateLinkServiceConnectionStateStatus' as connection_status,
  connections ->> 'provisioningState' as connection_provisioning_state
from
  azure_app_configuration,
  jsonb_array_elements(private_endpoint_connections) as connections;
```

### List encryption details for app configurations

```sql
select
  name as app_config_name,
  id as app_config_id,
  encryption -> 'keyVaultProperties' ->> 'identityClientId' as key_vault_identity_client_id,
  encryption -> 'keyVaultProperties' ->> 'keyIdentifier' as key_vault_key_identifier
from
  azure_app_configuration;
```
