---
title: "Steampipe Table: azure_app_configuration - Query Azure App Configuration Stores using SQL"
description: "Allows users to query Azure App Configuration Stores"
---

# Table: azure_app_configuration - Query Azure App Configuration Stores using SQL

Azure App Configuration is a managed service that helps developers centralize their application and feature settings simply and securely. It provides a way to manage and distribute application settings, helping to improve the speed and reliability of application deployment. Azure App Configuration also allows you to automate the process of managing and updating these settings across multiple environments.

## Table Usage Guide

The 'azure_app_configuration' table provides insights into App Configuration Stores within Azure App Configuration. As a DevOps engineer, explore store-specific details through this table, including store names, resource groups, subscription IDs, and associated metadata. Utilize it to uncover information about stores, such as their provisioning states, creation times, and the number of failed requests. The schema presents a range of attributes of the App Configuration Store for your analysis, like the store name, creation date, provisioning state, and associated tags.

## Examples

### Basic info
Explore the status and creation dates of your Azure application configurations. This can help you understand the overall state of your applications, allowing for better management and timely updates.

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
Explore which app configurations have public network access enabled. This can be useful in identifying potential security risks and ensuring your app configurations adhere to best practices.

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
This query helps in identifying the application configurations within Azure that have been assigned user identities. It is useful in managing and tracking user access, contributing to improved security and compliance.

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
Explore the details of private endpoint connections for your app configurations. This can help you understand their current status, any required actions, and their provisioning state, which can be useful for troubleshooting or optimizing your app's performance.

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
Explore the encryption details of your app configurations to ensure secure data handling. This is particularly useful in maintaining data security standards and regulatory compliance.

```sql
select
  name as app_config_name,
  id as app_config_id,
  encryption -> 'keyVaultProperties' ->> 'identityClientId' as key_vault_identity_client_id,
  encryption -> 'keyVaultProperties' ->> 'keyIdentifier' as key_vault_key_identifier
from
  azure_app_configuration;
```