---
title: "Steampipe Table: azure_cognitive_account - Query Azure Cognitive Services Accounts using SQL"
description: "Allows users to query Azure Cognitive Services Accounts."
---

# Table: azure_cognitive_account - Query Azure Cognitive Services Accounts using SQL

Azure Cognitive Services is a suite of artificial intelligence (AI) services and cognitive APIs to help you build intelligent apps. It provides developers with APIs that help in building applications that can see, hear, speak, understand, and even begin to reason. The APIs are designed to be easy to use, while also providing a comprehensive set of capabilities.

## Table Usage Guide

The 'azure_cognitive_account' table provides insights into Cognitive Services Accounts within Azure Cognitive Services. As a DevOps engineer, explore account-specific details through this table, including the kind of cognitive service, the network rules set, and associated metadata. Utilize it to uncover information about accounts, such as those with specific capabilities, the network rules applied to them, and the status of the accounts. The schema presents a range of attributes of the Cognitive Services Account for your analysis, like the account name, creation date, endpoint, and associated tags.

## Examples

### Basic info
Explore which Azure cognitive accounts are currently being provisioned, by understanding their type and kind. This can help in managing resources and planning for capacity.

```sql
select
  name,
  id,
  kind,
  type,
  provisioning_state
from
  azure_cognitive_account;
```

### List accounts with enabled public network access
Determine the areas in which public network access is enabled within your Azure cognitive accounts. This can assist in identifying potential security risks and ensuring your data remains protected.

```sql
select
  name,
  id,
  kind,
  type,
  provisioning_state,
  public_network_access
from
  azure_cognitive_account
where
  public_network_access = 'Enabled';
```

### List private endpoint connection details for accounts
This example helps in exploring the details of private endpoint connections linked to cognitive accounts in Azure. It can assist in understanding the connections' status and type, which is essential for managing network accessibility and ensuring secure data communication.

```sql
select
  name,
  id,
  connections ->> 'ID' as connection_id,
  connections ->> 'Name' as connection_name,
  connections ->> 'PrivateEndpointID' as property_private_endpoint_id,
  jsonb_pretty(connections -> 'PrivateLinkServiceConnectionState') as property_private_link_service_connection_state,
  connections ->> 'Type' as connection_type
from
  azure_cognitive_account,
  jsonb_array_elements(private_endpoint_connections) as connections;
```

### List diagnostic setting details for accounts
This query allows you to analyze the diagnostic settings of your Azure Cognitive Services accounts. It's useful for understanding the log and metric settings of each account, which can help in monitoring and troubleshooting.

```sql
select
  name,
  id,
  settings ->> 'id' as settings_id,
  settings ->> 'name' as settings_name,
  jsonb_pretty(settings -> 'properties' -> 'logs') as settings_properties_logs,
  jsonb_pretty(settings -> 'properties' -> 'metrics') as settings_properties_metrics,
  settings -> 'properties' ->> 'workspaceId' as settings_properties_workspaceId,
  settings ->> 'type' as settings_type
from
  azure_cognitive_account,
  jsonb_array_elements(diagnostic_settings) as settings;
```