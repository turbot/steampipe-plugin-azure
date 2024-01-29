---
title: "Steampipe Table: azure_app_service_web_app_slot - Query Azure App Service Web App Slots using SQL"
description: "Allows users to query Azure App Service Web App Slots, providing insights into the configuration, status, and metadata of each slot."
---

# Table: azure_app_service_web_app_slot - Query Azure App Service Web App Slots using SQL

Azure App Service Web App Slots are live apps with their own hostnames. They are part of Azure App Service and are designed to help developers handle app deployments. With slots, you can deploy your apps in a controlled manner and avoid downtime.

## Table Usage Guide

The `azure_app_service_web_app_slot` table provides insights into Azure App Service Web App Slots. As a developer or DevOps engineer, you can use this table to get detailed information about each slot, including its configuration, status, and metadata. This can be particularly useful when managing app deployments and ensuring smooth transitions between different versions of your app.

## Examples

### Basic info
Explore which web app slots are currently active within the Azure App Service. This allows you to assess their status, identify the ones that have been recently modified, and understand their configuration for better management and optimization.

```sql+postgres
select
  name,
  app_name,
  id,
  kind,
  state,
  type,
  last_modified_time_utc,
  repository_site_name,
  enabled
from
  azure_app_service_web_app_slot;
```

```sql+sqlite
select
  name,
  app_name,
  id,
  kind,
  state,
  type,
  last_modified_time_utc,
  repository_site_name,
  enabled
from
  azure_app_service_web_app_slot;
```

### List slots where the apps are enabled
Examine the active slots within Azure's app service to understand where applications are currently operational. This is useful for managing resources, ensuring optimal app performance, and identifying potential areas for scaling or re-allocation.

```sql+postgres
select
  name,
  app_name,
  state,
  type,
  reserved,
  server_farm_id,
  target_swap_slot
  enabled
from
  azure_app_service_web_app_slot
where
  enabled;
```

```sql+sqlite
select
  name,
  app_name,
  state,
  type,
  reserved,
  server_farm_id,
  target_swap_slot,
  enabled
from
  azure_app_service_web_app_slot
where
  enabled = 1;
```

### List slots that accept HTTP traffic (i.e only HTTPS is disabled)
Determine the areas in your Azure App Service where only HTTP traffic is allowed, which could potentially expose your web applications to security risks. This query is useful to identify these areas and implement necessary security measures to restrict traffic to HTTPS only.

```sql+postgres
select
  name,
  https_only,
  kind,
  region
from
  azure_app_service_web_app_slot
where
  not https_only;
```

```sql+sqlite
select
  name,
  https_only,
  kind,
  region
from
  azure_app_service_web_app_slot
where
  https_only = 0;
```

### Host names of each slot
Determine the areas in which your Azure App Service Web App Slots are being utilized. This can help you understand the distribution of your resources across different regions and resource groups, aiding in efficient resource management.

```sql+postgres
select
  name,
  host_names,
  kind,
  region,
  resource_group
from
  azure_app_service_web_app_slot;
```

```sql+sqlite
select
  name,
  host_names,
  kind,
  region,
  resource_group
from
  azure_app_service_web_app_slot;
```

### List enabled host names
Determine the areas in which host names are enabled to ensure the proper functioning of your Azure App Service Web App Slots. This allows you to manage and monitor your web applications effectively.

```sql+postgres
select
  name,
  id,
  type,
  kind,
  enabled_host_names
from
  azure_app_service_web_app_slot;
```

```sql+sqlite
select
  name,
  id,
  type,
  kind,
  enabled_host_names
from
  azure_app_service_web_app_slot;
```

### Get slot swap status of each slot
This query allows you to monitor the status of slot swaps in your Azure App Service Web App Slots. It's useful for keeping track of your deployment process and ensuring smooth transitions between different versions of your web applications.

```sql+postgres
select
  name,
  type,
  slot_swap_status ->> 'SlotSwapStatus' as slot_swap_status,
  slot_swap_status ->> 'SourceSlotName' as source_slot_name,
  slot_swap_status ->> 'DestinationSlotName' as destination_slot_name
from
  azure_app_service_web_app_slot;
```

```sql+sqlite
select
  name,
  type,
  json_extract(slot_swap_status, '$.SlotSwapStatus') as slot_swap_status,
  json_extract(slot_swap_status, '$.SourceSlotName') as source_slot_name,
  json_extract(slot_swap_status, '$.DestinationSlotName') as destination_slot_name
from
  azure_app_service_web_app_slot;
```

### Get site config details of each slot
Explore the configuration details of each slot in your Azure App Service Web App to understand the settings of individual workers and software versions. This can be useful for performance tuning and troubleshooting.

```sql+postgres
select
  name,
  id,
  site_config ->> 'NumberOfWorkers' as number_of_workers,
  site_config ->> 'DefaultDocuments' as DefaultDocuments,
  site_config ->> 'NetFrameworkVersion' as NetFrameworkVersion,
  site_config ->> 'PhpVersion' as PhpVersion,
  site_config ->> 'PythonVersion' as PythonVersion,
  site_config ->> 'NodeVersion' as NodeVersion,
  site_config ->> 'PowerShellVersion' as PowerShellVersion,
  site_config ->> 'LinuxFxVersion' as LinuxFxVersion,
  site_config ->> 'WindowsFxVersion' as WindowsFxVersion,
  site_config ->> 'RequestTracingEnabled' as RequestTracingEnabled,
  site_config ->> 'RequestTracingExpirationTime' as RequestTracingExpirationTime,
  site_config ->> 'RemoteDebuggingEnabled' as RemoteDebuggingEnabled,
  site_config ->> 'RemoteDebuggingVersion' as RemoteDebuggingVersion,
  site_config ->> 'HTTPLoggingEnabled' as HTTPLoggingEnabled,
  site_config ->> 'DetailedErrorLoggingEnabled' as DetailedErrorLoggingEnabled,
  site_config ->> 'PublishingUsername' as PublishingUsername,
  site_config ->> 'AppSettings' as AppSettings,
  site_config ->> 'ConnectionStrings' as ConnectionStrings,
  site_config ->> 'MachineKey' as MachineKey,
  site_config ->> 'HandlerMappings' as HandlerMappings,
  site_config ->> 'DocumentRoot' as DocumentRoot
from
  azure_app_service_web_app_slot;
```

```sql+sqlite
select
  name,
  id,
  json_extract(site_config, '$.NumberOfWorkers') as number_of_workers,
  json_extract(site_config, '$.DefaultDocuments') as DefaultDocuments,
  json_extract(site_config, '$.NetFrameworkVersion') as NetFrameworkVersion,
  json_extract(site_config, '$.PhpVersion') as PhpVersion,
  json_extract(site_config, '$.PythonVersion') as PythonVersion,
  json_extract(site_config, '$.NodeVersion') as NodeVersion,
  json_extract(site_config, '$.PowerShellVersion') as PowerShellVersion,
  json_extract(site_config, '$.LinuxFxVersion') as LinuxFxVersion,
  json_extract(site_config, '$.WindowsFxVersion') as WindowsFxVersion,
  json_extract(site_config, '$.RequestTracingEnabled') as RequestTracingEnabled,
  json_extract(site_config, '$.RequestTracingExpirationTime') as RequestTracingExpirationTime,
  json_extract(site_config, '$.RemoteDebuggingEnabled') as RemoteDebuggingEnabled,
  json_extract(site_config, '$.RemoteDebuggingVersion') as RemoteDebuggingVersion,
  json_extract(site_config, '$.HTTPLoggingEnabled') as HTTPLoggingEnabled,
  json_extract(site_config, '$.DetailedErrorLoggingEnabled') as DetailedErrorLoggingEnabled,
  json_extract(site_config, '$.PublishingUsername') as PublishingUsername,
  json_extract(site_config, '$.AppSettings') as AppSettings,
  json_extract(site_config, '$.ConnectionStrings') as ConnectionStrings,
  json_extract(site_config, '$.MachineKey') as MachineKey,
  json_extract(site_config, '$.HandlerMappings') as HandlerMappings,
  json_extract(site_config, '$.DocumentRoot') as DocumentRoot
from
  azure_app_service_web_app_slot;
```