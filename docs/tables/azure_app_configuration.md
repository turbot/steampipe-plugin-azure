---
title: "Steampipe Table: azure_app_configuration - Query Azure App Configuration using SQL"
description: "Allows users to query App Configurations in Azure, providing insights into application settings and feature management."
---

# Table: azure_app_configuration - Query Azure App Configuration using SQL

Azure App Configuration is a service within Microsoft Azure that provides a way to centrally manage application settings and feature flags. It helps developers to separate configuration from code, making applications more modular and scalable. Azure App Configuration is fully managed, which allows developers to focus on code rather than managing and distributing configuration.

## Table Usage Guide

The `azure_app_configuration` table provides insights into application configurations within Microsoft Azure. As a developer or DevOps engineer, you can explore configuration-specific details through this table, including settings, feature management, and associated metadata. Utilize it to manage and monitor application settings, understand feature flags, and ensure the scalability and modularity of your applications.

## Examples

### Basic info
Explore which Azure App configurations are currently active and when they were created. This is useful for understanding the status and timeline of your app's setup and deployment.

```sql+postgres
select
  id,
  name,
  type,
  provisioning_state,
  creation_date
from
  azure_app_configuration;
```

```sql+sqlite
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
Explore which app configurations in Azure have public network access enabled. This can be beneficial in assessing potential security risks and ensuring appropriate network access settings are in place.

```sql+postgres
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

```sql+sqlite
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
This query is useful to identify and analyze the configurations of apps that have user-assigned identities within your Azure environment. It helps in managing and auditing access control, thereby enhancing the security of your applications.

```sql+postgres
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

```sql+sqlite
Error: SQLite does not support regexp_split_to_array function.
```

### List private endpoint connection details for app configurations
Explore the status and details of private connections for app configurations in Azure. This can help identify any required actions or understand the provisioning state for these connections.

```sql+postgres
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

```sql+sqlite
select
  name as app_config_name,
  c.id as app_config_id,
  json_extract(connections.value, '$.id') as connection_id,
  json_extract(connections.value, '$.privateEndpointPropertyId') as connection_private_endpoint_property_id,
  json_extract(connections.value, '$.privateLinkServiceConnectionStateActionsRequired') as connection_actions_required,
  json_extract(connections.value, '$.privateLinkServiceConnectionStateDescription') as connection_description,
  json_extract(connections.value, '$.privateLinkServiceConnectionStateStatus') as connection_status,
  json_extract(connections.value, '$.provisioningState') as connection_provisioning_state
from
  azure_app_configuration as c,
  json_each(private_endpoint_connections) as connections;
```

### List encryption details for app configurations
Explore encryption specifics for your applications, particularly focusing on identity client IDs and key identifiers. This is useful for assessing the security measures in place for your app configurations.

```sql+postgres
select
  name as app_config_name,
  id as app_config_id,
  encryption -> 'keyVaultProperties' ->> 'identityClientId' as key_vault_identity_client_id,
  encryption -> 'keyVaultProperties' ->> 'keyIdentifier' as key_vault_key_identifier
from
  azure_app_configuration;
```

```sql+sqlite
select
  name as app_config_name,
  id as app_config_id,
  json_extract(encryption, '$.keyVaultProperties.identityClientId') as key_vault_identity_client_id,
  json_extract(encryption, '$.keyVaultProperties.keyIdentifier') as key_vault_key_identifier
from
  azure_app_configuration;
```